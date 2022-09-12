package gocms

import (
	"context"
	"io"
	"net/http"
)

type File struct {
	Location    string    `json:"location"`
	ContentType string    `json:"content_type,omitempty"`
	Filename    string    `json:"filename,omitempty"`
	Data        io.Reader `json:"data,omitempty"`
}

type FileUploadRequest struct {
	Key         string `json:"key"`
	ContentType string `json:"content_type"`
	Expires     string `json:"expires"`
}

type FileUploadResponse struct {
	URL      string      `json:"url"`
	Method   string      `json:"method"`
	Headers  http.Header `json:"headers"`
	Location string      `json:"location"`
}

type FileRepository interface {
	CreateFile(context.Context, *File) error
	CreateUploadUrl(context.Context, FileUploadRequest) (FileUploadResponse, error)
}

type FileService struct {
	repo FileRepository
}

func NewFileService(repo FileRepository) FileService {
	return FileService{
		repo: repo,
	}
}

func (s FileService) UploadUrl(ctx context.Context, request FileUploadRequest) (FileUploadResponse, error) {
	return s.repo.CreateUploadUrl(ctx, request)
}
