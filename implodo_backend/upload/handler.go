package upload

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"implodo_backend/auth"
	"implodo_backend/imagegen"
	"implodo_backend/imagestorage"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 20 << 20 // 20 MB

var allowedMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

type Handler struct {
	gen   imagegen.Client
	store imagestorage.Client
}

func NewHandler(gen imagegen.Client, store imagestorage.Client) *Handler {
	return &Handler{gen: gen, store: store}
}

// ListImages godoc
//
//	GET /api/images
//
// Lists all image pairs uploaded by the authenticated user, newest first.
// Returns: [{ id, original_key, generated_key }]
func (h *Handler) ListImages(c *gin.Context) {
	claims := auth.GetClaims(c)
	prefix := sanitise(claims.GoogleID) + "/"

	keys, err := h.store.List(c.Request.Context(), prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list images"})
		return
	}

	// Group objects into original/generated pairs keyed by their shared hex ID.
	type pair struct {
		ID           string `json:"id"`
		OriginalURL  string `json:"original_url"`
		GeneratedURL string `json:"generated_url"`
	}
	byID := map[string]*pair{}
	var order []string // preserve insertion order for stable output

	for _, obj := range keys {
		base := obj.Key[len(prefix):] // strip "<googleID>/"
		if id, _, ok := strings.Cut(base, "_original"); ok {
			if byID[id] == nil {
				byID[id] = &pair{ID: id}
				order = append(order, id)
			}
			byID[id].OriginalURL = obj.URL
		} else if id, _, ok := strings.Cut(base, "_generated"); ok {
			if byID[id] == nil {
				byID[id] = &pair{ID: id}
				order = append(order, id)
			}
			byID[id].GeneratedURL = obj.URL
		}
	}

	// Only return complete pairs, newest first (order is reversed since GCS
	// returns lexicographically; hex IDs are random so we just reverse the list).
	result := make([]*pair, 0, len(order))
	for i := len(order) - 1; i >= 0; i-- {
		p := byID[order[i]]
		if p.OriginalURL != "" && p.GeneratedURL != "" {
			result = append(result, p)
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetImage godoc
//
//	GET /api/images/*key
//
// Fetches the image at the given storage key from GCS and streams it to the client.
// Requires a valid JWT — the middleware on the /api group handles that.
func (h *Handler) GetImage(c *gin.Context) {
	key := c.Param("key")
	if key == "" || key == "/" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}
	// Trim the leading slash that gin includes in wildcard params.
	key = strings.TrimPrefix(key, "/")

	result, err := h.store.Get(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	contentType := result.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Data(http.StatusOK, contentType, result.Data)
}

// UploadImage godoc
//
//	POST /api/upload
//
// Accepts a multipart/form-data request with an "image" field.
//  1. Reads and validates the image.
//  2. Sends it to Nano Banana to generate a second image.
//  3. Uploads both images to blob storage.
//  4. Returns the public URLs of both.
func (h *Handler) UploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large (max 20 MB)"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image field is required"})
		return
	}
	defer file.Close()

	// Read the full file into memory so we can re-use it.
	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read image"})
		return
	}

	// Detect MIME type from the first 512 bytes.
	mediaType, ext, ok := detectMIME(imageData, header.Header.Get("Content-Type"))
	if !ok {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "only JPEG, PNG, GIF, and WebP images are accepted"})
		return
	}

	claims := auth.GetClaims(c)
	ctx := context.Background()

	// --- Step 1: generate a second image via Nano Banana ---
	generatedData, err := h.gen.Generate(ctx, imageData, mediaType)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("image generation failed: %v", err)})
		return
	}

	// --- Step 2: upload both images to blob storage ---
	prefix := sanitise(claims.GoogleID) + "/" + randomHex(16)

	originalKey := prefix + "_original" + ext
	originalResult, err := h.store.Upload(ctx, originalKey, imageData, mediaType)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("storage upload failed: %v", err)})
		return
	}

	generatedKey := prefix + "_generated" + ext
	generatedResult, err := h.store.Upload(ctx, generatedKey, generatedData, mediaType)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("storage upload failed: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"original_key":  originalResult.Key,
		"generated_key": generatedResult.Key,
	})
}

// detectMIME returns the media type and file extension for imageData.
// Falls back to the declared content type header if byte sniffing is inconclusive.
func detectMIME(data []byte, declared string) (mediaType, ext string, ok bool) {
	sniff := http.DetectContentType(data[:min(512, len(data))])
	mt, _, _ := mime.ParseMediaType(sniff)
	if e, found := allowedMIME[mt]; found {
		return mt, e, true
	}
	mt2, _, _ := mime.ParseMediaType(declared)
	if e, found := allowedMIME[mt2]; found {
		return mt2, e, true
	}
	return "", "", false
}

func sanitise(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func randomHex(n int) string {
	buf := make([]byte, n)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
