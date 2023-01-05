package class

import (
	"errors"
	"time"

	"github.com/jbaikge/boneless/internal/common/id"
)

type Class struct {
	id       string
	parentId string
	name     string
	created  time.Time
	updated  time.Time
	fields   []*Field
}

func NewClass(name string, parentId string, fields []*Field) *Class {
	created := time.Now()
	c := &Class{
		id:       id.NewWithTime(created),
		parentId: parentId,
		name:     name,
		created:  created,
		updated:  created,
		fields:   make([]*Field, len(fields)),
	}
	copy(c.fields, fields)
	return c
}

func Unmarshal(classId string, parentId string, name string, created time.Time, updated time.Time, fields []*Field) (*Class, error) {
	if !id.IsValid(classId) {
		return nil, errors.New("invalid class id")
	}
	if parentId != "" && !id.IsValid(parentId) {
		return nil, errors.New("invalid parent id")
	}
	if name == "" {
		return nil, errors.New("empty class name")
	}
	if created.IsZero() {
		return nil, errors.New("created timestamp is zero")
	}
	if updated.IsZero() {
		return nil, errors.New("updated timestamp is zero")
	}

	c := &Class{
		id:       classId,
		parentId: parentId,
		name:     name,
		created:  created,
		updated:  updated,
		fields:   make([]*Field, len(fields)),
	}
	copy(c.fields, fields)
	return c, nil
}

func (c Class) ID() string {
	return c.id
}

func (c Class) Name() string {
	return c.name
}
