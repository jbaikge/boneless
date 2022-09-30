package models

import (
	"io"
	"net/http"
)

type File struct {
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
