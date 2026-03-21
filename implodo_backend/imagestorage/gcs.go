package imagestorage

import (
	"context"
	"fmt"
	"io"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// GCSClient is a Client backed by Google Cloud Storage.
type GCSClient struct {
	bucket *gcs.BucketHandle
	// baseURL is the public URL prefix, e.g. "https://storage.googleapis.com/<bucket>".
	baseURL string
}

// NewGCSClient creates a GCS-backed client for the given bucket.
// Credentials are resolved from the GOOGLE_APPLICATION_CREDENTIALS environment
// variable (point it at your service account JSON file).
func NewGCSClient(ctx context.Context, bucketName string) (*GCSClient, error) {
	client, err := gcs.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("imagestorage: create GCS client: %w", err)
	}
	return &GCSClient{
		bucket:  client.Bucket(bucketName),
		baseURL: fmt.Sprintf("https://storage.googleapis.com/%s", bucketName),
	}, nil
}

func (g *GCSClient) Upload(ctx context.Context, key string, data []byte, contentType string) (UploadResult, error) {
	obj := g.bucket.Object(key)
	w := obj.NewWriter(ctx)
	w.ContentType = contentType

	if _, err := w.Write(data); err != nil {
		w.Close()
		return UploadResult{}, fmt.Errorf("imagestorage: write object: %w", err)
	}
	if err := w.Close(); err != nil {
		return UploadResult{}, fmt.Errorf("imagestorage: close object writer: %w", err)
	}

	url := fmt.Sprintf("%s/%s", g.baseURL, key)
	return UploadResult{URL: url, Key: key}, nil
}

func (g *GCSClient) List(ctx context.Context, prefix string) ([]UploadResult, error) {
	var results []UploadResult
	it := g.bucket.Objects(ctx, &gcs.Query{Prefix: prefix})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("imagestorage: list objects: %w", err)
		}
		results = append(results, UploadResult{
			Key: attrs.Name,
			URL: fmt.Sprintf("%s/%s", g.baseURL, attrs.Name),
		})
	}
	return results, nil
}

func (g *GCSClient) Get(ctx context.Context, key string) (GetResult, error) {
	r, err := g.bucket.Object(key).NewReader(ctx)
	if err != nil {
		return GetResult{}, fmt.Errorf("imagestorage: open object %q: %w", key, err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return GetResult{}, fmt.Errorf("imagestorage: read object %q: %w", key, err)
	}

	return GetResult{Data: data, ContentType: r.Attrs.ContentType}, nil
}
