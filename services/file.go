package services

import (
	"context"

	"github.com/jbaikge/boneless/models"
)

type FileRepository interface {
	CreateFile(context.Context, *models.File) (string, error)
	CreateUploadUrl(context.Context, models.FileUploadRequest) (models.FileUploadResponse, error)
}

type FileService struct {
	repo FileRepository
}

func NewFileService(repo FileRepository) FileService {
	return FileService{
		repo: repo,
	}
}

func (s FileService) CreateFile(ctx context.Context, file *models.File) (location string, err error) {
	return s.repo.CreateFile(ctx, file)
}

func (s FileService) UploadUrl(ctx context.Context, request models.FileUploadRequest) (models.FileUploadResponse, error) {
	return s.repo.CreateUploadUrl(ctx, request)
}
