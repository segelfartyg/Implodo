package imagestorage

import (
	"context"
	"errors"
)

// Stub is a placeholder that returns an error until a real implementation is wired in.
type Stub struct{}

func NewStub() *Stub { return &Stub{} }

func (s *Stub) Upload(_ context.Context, _ string, _ []byte, _ string) (UploadResult, error) {
	return UploadResult{}, errors.New("imagestorage: no backend configured")
}

func (s *Stub) Get(_ context.Context, _ string) (GetResult, error) {
	return GetResult{}, errors.New("imagestorage: no backend configured")
}

func (s *Stub) List(_ context.Context, _ string) ([]UploadResult, error) {
	return nil, errors.New("imagestorage: no backend configured")
}
