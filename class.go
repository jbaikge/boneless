package gocms

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

type Class struct {
	Id          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	TableLabels string    `json:"table_labels"`
	TableFields string    `json:"table_fields"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Fields      []Field   `json:"fields"`
}

type ClassRepository interface {
	DeleteClass(string) error
	GetAllClasses() ([]Class, error)
	GetClassById(string) (Class, error)
	InsertClass(*Class) error
	UpdateClass(*Class) error
}

type ClassService struct {
	repo ClassRepository
}

func NewClassService(repo ClassRepository) ClassService {
	return ClassService{
		repo: repo,
	}
}

func (s ClassService) All() ([]Class, error) {
	return s.repo.GetAllClasses()
}

func (s ClassService) Insert(class *Class) (err error) {
	// TODO validate internal fields

	if class.Id != "" {
		return fmt.Errorf("class already has an ID")
	}

	// TODO check for existing class with same slug

	now := time.Now()
	class.Id = "class#" + xid.NewWithTime(now).String()
	class.Created = now
	class.Updated = now

	return s.repo.InsertClass(class)
}

func (s ClassService) Update(class *Class) (err error) {
	if class.Id == "" {
		return fmt.Errorf("class has no ID")
	}

	// TODO check for slug overwrite

	class.Updated = time.Now()

	return s.repo.UpdateClass(class)
}
