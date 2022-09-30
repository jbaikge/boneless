package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jbaikge/boneless/models"
)

type TemplateRepository interface {
	CreateTemplate(context.Context, *models.Template) error
	DeleteTemplate(context.Context, string) error
	GetTemplateById(context.Context, string) (models.Template, error)
	GetTemplateList(context.Context, models.TemplateFilter) ([]models.Template, models.Range, error)
	UpdateTemplate(context.Context, *models.Template) error
}

type TemplateService struct {
	repo TemplateRepository
}

func NewTemplateService(repo TemplateRepository) TemplateService {
	return TemplateService{
		repo: repo,
	}
}

func (s TemplateService) ById(ctx context.Context, id string) (models.Template, error) {
	if !idProvider.IsValid(id) {
		return models.Template{}, fmt.Errorf("invalid template ID: %s", id)
	}
	return s.repo.GetTemplateById(ctx, id)
}

func (s TemplateService) Create(ctx context.Context, template *models.Template) (err error) {
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

func (s TemplateService) List(ctx context.Context, filter models.TemplateFilter) ([]models.Template, models.Range, error) {
	return s.repo.GetTemplateList(ctx, filter)
}

func (s TemplateService) Update(ctx context.Context, template *models.Template) (err error) {
	if template.Id == "" {
		return fmt.Errorf("template has no ID")
	}

	template.Updated = time.Now()

	return s.repo.UpdateTemplate(ctx, template)
}
