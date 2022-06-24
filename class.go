package gocms

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type Class struct {
	Id          string    `json:"id" dynamodbav:"PrimaryKey"`
	SortKey     string    `json:"-" dynamodbav:"SortKey"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	TableLabels string    `json:"table_labels"`
	TableFields string    `json:"table_fields"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Fields      []Field   `json:"fields"`
}

type ClassRepository interface {
	DeleteClass(context.Context, string) error
	GetAllClasses(context.Context) ([]Class, error)
	GetClassById(context.Context, string) (Class, error)
	InsertClass(context.Context, *Class) error
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

func (s ClassService) All(ctx context.Context) ([]Class, error) {
	return s.repo.GetAllClasses(ctx)
}

func (s ClassService) ById(ctx context.Context, id string) (Class, error) {
	if _, err := xid.FromString(id); err != nil {
		return Class{}, err
	}
	return s.repo.GetClassById(ctx, id)
}

func (s ClassService) Insert(ctx context.Context, class *Class) (err error) {
	// TODO validate internal fields

	if class.Id != "" {
		return fmt.Errorf("class already has an ID")
	}

	// TODO check for existing class with same slug

	now := time.Now()
	class.Id = "class#" + xid.NewWithTime(now).String()
	class.Created = now
	class.Updated = now

	return s.repo.InsertClass(ctx, class)
}

func (s ClassService) Update(ctx context.Context, class *Class) (err error) {
	if class.Id == "" {
		return fmt.Errorf("class has no ID")
	}

	// TODO check for slug overwrite

	class.Updated = time.Now()

	return s.repo.UpdateClass(ctx, class)
}
