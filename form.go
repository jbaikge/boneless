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

}
