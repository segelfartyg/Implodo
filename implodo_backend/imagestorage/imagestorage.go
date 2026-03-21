// Package imagestorage defines the interface for blob storage and provides
// a stub and a GCS-backed implementation.
package imagestorage

import "context"

// UploadResult holds the public URL and provider-assigned key for a stored object.
type UploadResult struct {
	URL string
	Key string
}

// GetResult holds the raw bytes and content type of a retrieved object.
type GetResult struct {
	Data        []byte
	ContentType string
}

// Client stores and retrieves objects in blob storage.
type Client interface {
	// Upload writes data under the given key and returns the public URL.
	// The key is typically a path like "<google_id>/<filename>".
	Upload(ctx context.Context, key string, data []byte, contentType string) (UploadResult, error)

	// Get retrieves the object at key and returns its bytes and content type.
	Get(ctx context.Context, key string) (GetResult, error)

	// List returns all objects whose keys start with prefix, as UploadResults.
	List(ctx context.Context, prefix string) ([]UploadResult, error)
}
