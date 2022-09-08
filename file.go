package gocms

import (
	"context"
	"io"
	"net/http"
	"time"
)

type File struct {
	Location    string    `json:"location"`
	ContentType string    `json:"content_type,omitempty"`
	Data        io.Reader `json:"data,omitempty"`
}

type FileUploadRequest struct {
	Key         string        `json:"key"`
	ContentType string        `json:"content_type"`
	Expires     time.Duration `json:"expires"`
}

type FileUploadResponse struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
}

type FileRepository interface {
	CreateFile(context.Context, *File) error
	CreateUploadUrl(context.Context, FileUploadRequest) (FileUploadResponse, error)
}
