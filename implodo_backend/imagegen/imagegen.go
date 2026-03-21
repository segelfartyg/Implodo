// Package imagegen defines the interface for the Nano Banana image generation service
// and provides a stub implementation until the real SDK is wired in.
package imagegen

import "context"

// Client sends an input image to Nano Banana and receives a generated image back.
type Client interface {
	// Generate takes the raw bytes and MIME type of the input image and returns
	// the raw bytes of the generated image.
	Generate(ctx context.Context, imageData []byte, mimeType string) ([]byte, error)
}

// Stub is a placeholder that returns an error until the real SDK is integrated.
type Stub struct{}

func NewStub() *Stub { return &Stub{} }

func (s *Stub) Generate(_ context.Context, _ []byte, _ string) ([]byte, error) {
	return nil, errNotImplemented("imagegen: not yet integrated")
}
