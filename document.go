package boneless

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Document struct {
	Id         string                 `json:"id"`
	ClassId    string                 `json:"class_id"`
	ParentId   string                 `json:"parent_id"`
	TemplateId string                 `json:"template_id"`
	Path       string                 `json:"path"`
	Version    int                    `json:"version"`
	Created    time.Time              `json:"created"`
	Updated    time.Time              `json:"updated"`
	Values     map[string]interface{} `json:"values"`
}

type DocumentFilterSort struct {
	Field     string
	Direction string
}

func (dfs DocumentFilterSort) Ascending() bool {
	return strings.ToUpper(dfs.Direction) == "ASC" || dfs.Direction == ""
}

func (dfs DocumentFilterSort) Descending() bool {
	return strings.ToUpper(dfs.Direction) == "DESC"
}

type DocumentFilter struct {
	ClassId  string
	ParentId string
	Sort     DocumentFilterSort
	Range    Range
}

type DocumentRepository interface {
	CreateDocument(context.Context, *Document) error
	DeleteDocument(context.Context, string) error
	GetDocumentById(context.Context, string) (Document, error)
	GetDocumentByPath(context.Context, string) (Document, error)
	GetDocumentList(context.Context, DocumentFilter) ([]Document, Range, error)
	UpdateDocument(context.Context, *Document) error
}

type DocumentService struct {
	repo DocumentRepository
}

func NewDocumentService(repo DocumentRepository) DocumentService {
	return DocumentService{
		repo: repo,
	}
}

func (s DocumentService) ById(ctx context.Context, id string) (Document, error) {
	if !idProvider.IsValid(id) {
		return Document{}, fmt.Errorf("invalid document ID: %s", id)
	}
	return s.repo.GetDocumentById(ctx, id)
}

func (s DocumentService) ByPath(ctx context.Context, path string) (Document, error) {
	return s.repo.GetDocumentByPath(ctx, path)
}

func (s DocumentService) Create(ctx context.Context, doc *Document) (err error) {
	if doc.Id != "" {
		return fmt.Errorf("document already has an ID")
	}

	now := time.Now()
	doc.Id = idProvider.NewWithTime(now)
	doc.Version = 1
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

func (s DocumentService) List(ctx context.Context, filter DocumentFilter) ([]Document, Range, error) {
	return s.repo.GetDocumentList(ctx, filter)
}

func (s DocumentService) Update(ctx context.Context, doc *Document) (err error) {
	if doc.Id == "" {
		return fmt.Errorf("document has no ID")
	}

	doc.Updated = time.Now()

	return s.repo.UpdateDocument(ctx, doc)
}
