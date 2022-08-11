package gocms

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type Class struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	TableLabels string    `json:"table_labels"`
	TableFields string    `json:"table_fields"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Fields      []Field   `json:"fields"`
}

func (c Class) SortFields() (fields []string) {
	fields = make([]string, 0, len(c.Fields))
	for _, field := range c.Fields {
		if field.Sort {
			fields = append(fields, field.Name)
		}
	}
	return
}

type ClassFilter struct {
	Range Range
}

type ClassRepository interface {
	CreateClass(context.Context, *Class) error
	DeleteClass(context.Context, string) error
	GetClassById(context.Context, string) (Class, error)
	GetClassList(context.Context, ClassFilter) ([]Class, Range, error)
	UpdateClass(context.Context, *Class) error
}

type ClassService struct {
	repo ClassRepository
}

func NewClassService(repo ClassRepository) ClassService {
	return ClassService{
		repo: repo,
	}
}

func (s ClassService) ById(ctx context.Context, id string) (Class, error) {
	if !idProvider.IsValid(id) {
		return Class{}, fmt.Errorf("invalid class ID: %s", id)
	}
	return s.repo.GetClassById(ctx, id)
}

func (s ClassService) Create(ctx context.Context, class *Class) (err error) {
	if class.Id != "" {
		return fmt.Errorf("class already has an ID")
	}

	// TODO validate internal fields

	now := time.Now()
	class.Id = xid.NewWithTime(now).String()
	class.Created = now
	class.Updated = now

	return s.repo.CreateClass(ctx, class)
}

func (s ClassService) Delete(ctx context.Context, id string) (err error) {
	return s.repo.DeleteClass(ctx, id)
}

func (s ClassService) List(ctx context.Context, filter ClassFilter) ([]Class, Range, error) {
	return s.repo.GetClassList(ctx, filter)
}

func (s ClassService) Update(ctx context.Context, class *Class) (err error) {
	if class.Id == "" {
		return fmt.Errorf("class has no ID")
	}

	class.Updated = time.Now()

	return s.repo.UpdateClass(ctx, class)
}
