package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jbaikge/boneless/models"
)

type DocumentRepository interface {
	CreateDocument(context.Context, *models.Document) error
	DeleteDocument(context.Context, string) error
	GetDocumentById(context.Context, string) (models.Document, error)
	GetDocumentByPath(context.Context, string) (models.Document, error)
	GetDocumentList(context.Context, models.DocumentFilter) ([]models.Document, models.Range, error)
	UpdateDocument(context.Context, *models.Document) error
}

type DocumentService struct {
	repo DocumentRepository
}

func NewDocumentService(repo DocumentRepository) DocumentService {
	return DocumentService{
		repo: repo,
	}
}

func (s DocumentService) ById(ctx context.Context, id string) (models.Document, error) {
	if !idProvider.IsValid(id) {
		return models.Document{}, fmt.Errorf("invalid document ID: %s", id)
	}
	return s.repo.GetDocumentById(ctx, id)
}

func (s DocumentService) ByPath(ctx context.Context, path string) (models.Document, error) {
	return s.repo.GetDocumentByPath(ctx, path)
}

func (s DocumentService) Create(ctx context.Context, doc *models.Document) (err error) {
	if doc.Id != "" {
		return fmt.Errorf("document already has an ID")
	}

	now := time.Now()
	doc.Id = idProvider.NewWithTime(now)
	doc.Created = now
	doc.Updated = now

	return s.repo.CreateDocument(ctx, doc)
}

func (s DocumentService) Delete(ctx context.Context, id string) (err error) {
	if !idProvider.IsValid(id) {
		return fmt.Errorf("invalid document ID: %s", id)
	}
	return s.repo.DeleteDocument(ctx, id)
}

func (s DocumentService) List(ctx context.Context, filter models.DocumentFilter) ([]models.Document, models.Range, error) {
	return s.repo.GetDocumentList(ctx, filter)
}

func (s DocumentService) Update(ctx context.Context, doc *models.Document) (err error) {
	if doc.Id == "" {
		return fmt.Errorf("document has no ID")
	}

	doc.Updated = time.Now()

	return s.repo.UpdateDocument(ctx, doc)
}
