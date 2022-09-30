package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jbaikge/boneless/models"
)

type FormRepository interface {
	CreateForm(context.Context, *models.Form) error
	DeleteForm(context.Context, string) error
	GetFormById(context.Context, string) (models.Form, error)
	GetFormList(context.Context, models.FormFilter) ([]models.Form, models.Range, error)
	UpdateForm(context.Context, *models.Form) error
}

type FormService struct {
	repo FormRepository
}

func NewFormService(repo FormRepository) FormService {
	return FormService{
		repo: repo,
	}
}

func (s FormService) ById(ctx context.Context, id string) (models.Form, error) {
	if !idProvider.IsValid(id) {
		return models.Form{}, fmt.Errorf("invalid form ID: %s", id)
	}
	return s.repo.GetFormById(ctx, id)
}

func (s FormService) Create(ctx context.Context, form *models.Form) (err error) {
	if form.Id != "" {
		return fmt.Errorf("form already has an ID")
	}

	now := time.Now()
	form.Id = idProvider.NewWithTime(now)
	form.Created = now
	form.Updated = now

	return s.repo.CreateForm(ctx, form)
}

func (s FormService) Delete(ctx context.Context, id string) (err error) {
	if !idProvider.IsValid(id) {
		return fmt.Errorf("invalid form ID: %s", id)
	}

	return s.repo.DeleteForm(ctx, id)
}

func (s FormService) List(ctx context.Context, filter models.FormFilter) ([]models.Form, models.Range, error) {
	return s.repo.GetFormList(ctx, filter)
}

func (s FormService) Update(ctx context.Context, form *models.Form) (err error) {
	if !idProvider.IsValid(form.Id) {
		return fmt.Errorf("invalid form ID: %s", form.Id)
	}
	form.Updated = time.Now()
	return s.repo.UpdateForm(ctx, form)
}
