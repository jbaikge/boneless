package boneless

import (
	"context"
	"fmt"
	"time"
)

type Template struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Version int       `json:"version"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type TemplateFilter struct {
	Field       string
	SortReverse bool
	Range       Range
}

type TemplateRepository interface {
	CreateTemplate(context.Context, *Template) error
	DeleteTemplate(context.Context, string) error
	GetTemplateById(context.Context, string) (Template, error)
	GetTemplateList(context.Context, TemplateFilter) ([]Template, Range, error)
	UpdateTemplate(context.Context, *Template) error
}

type TemplateService struct {
	repo TemplateRepository
}

func NewTemplateService(repo TemplateRepository) TemplateService {
	return TemplateService{
		repo: repo,
	}
}

func (s TemplateService) ById(ctx context.Context, id string) (Template, error) {
	if !idProvider.IsValid(id) {
		return Template{}, fmt.Errorf("invalid template ID: %s", id)
	}
	return s.repo.GetTemplateById(ctx, id)
}

func (s TemplateService) Create(ctx context.Context, template *Template) (err error) {
	if template.Id != "" {
		return fmt.Errorf("template already has an ID")
	}

	now := time.Now()
	template.Id = idProvider.NewWithTime(now)
	template.Created = now
	template.Updated = now

	return s.repo.CreateTemplate(ctx, template)
}

func (s TemplateService) Delete(ctx context.Context, id string) (err error) {
	return s.repo.DeleteTemplate(ctx, id)
}

func (s TemplateService) List(ctx context.Context, filter TemplateFilter) ([]Template, Range, error) {
	return s.repo.GetTemplateList(ctx, filter)
}

func (s TemplateService) Update(ctx context.Context, template *Template) (err error) {
	if template.Id == "" {
		return fmt.Errorf("template has no ID")
	}

	template.Updated = time.Now()

	return s.repo.UpdateTemplate(ctx, template)
}
