package boneless

import (
	"context"
	"fmt"
	"time"
)

type Form struct {
	Id         string        `json:"id"`
	Created    time.Time     `json:"created"`
	Updated    time.Time     `json:"updated"`
	Components []interface{} `json:"components"`
}

type FormFilterSort struct {
	Field     string
	Direction string
}

type FormFilter struct {
	Sort  FormFilterSort
	Range Range
}

type FormRepository interface {
	CreateForm(context.Context, *Form) error
	DeleteForm(context.Context, string) error
	GetFormById(context.Context, string) (Form, error)
	GetFormList(context.Context, FormFilter) ([]Form, Range, error)
	UpdateForm(context.Context, *Form) error
}

type FormService struct {
	repo FormRepository
}

func NewFormService(repo FormRepository) FormService {
	return FormService{
		repo: repo,
	}
}

func (s FormService) ById(ctx context.Context, id string) (Form, error) {
	if !idProvider.IsValid(id) {
		return Form{}, fmt.Errorf("invalid form ID: %s", id)
	}
	return s.repo.GetFormById(ctx, id)
}

func (s FormService) Create(ctx context.Context, form *Form) (err error) {
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

func (s FormService) List(ctx context.Context, filter FormFilter) ([]Form, Range, error) {
	return s.repo.GetFormList(ctx, filter)
}

func (s FormService) Update(ctx context.Context, form *Form) (err error) {
	if !idProvider.IsValid(form.Id) {
		return fmt.Errorf("invalid form ID: %s", form.Id)
	}
	form.Updated = time.Now()
	return s.repo.UpdateForm(ctx, form)
}
