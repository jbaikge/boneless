package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jbaikge/boneless/models"
	"github.com/rs/xid"
)

type ClassRepository interface {
	CreateClass(context.Context, *models.Class) error
	DeleteClass(context.Context, string) error
	GetClassById(context.Context, string) (models.Class, error)
	GetClassList(context.Context, models.ClassFilter) ([]models.Class, models.Range, error)
	UpdateClass(context.Context, *models.Class) error
}

type ClassService struct {
	repo ClassRepository
}

func NewClassService(repo ClassRepository) ClassService {
	return ClassService{
		repo: repo,
	}
}

func (s ClassService) All(ctx context.Context) (classes []models.Class, err error) {
	filter := models.ClassFilter{
		Range: models.Range{End: 1000},
	}
	classes, _, err = s.List(ctx, filter)
	return
}

func (s ClassService) ById(ctx context.Context, id string) (models.Class, error) {
	if !idProvider.IsValid(id) {
		return models.Class{}, fmt.Errorf("invalid class ID: %s", id)
	}
	return s.repo.GetClassById(ctx, id)
}

func (s ClassService) Create(ctx context.Context, class *models.Class) (err error) {
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

func (s ClassService) List(ctx context.Context, filter models.ClassFilter) ([]models.Class, models.Range, error) {
	return s.repo.GetClassList(ctx, filter)
}

func (s ClassService) Update(ctx context.Context, class *models.Class) (err error) {
	if class.Id == "" {
		return fmt.Errorf("class has no ID")
	}

	class.Updated = time.Now()

	return s.repo.UpdateClass(ctx, class)
}
