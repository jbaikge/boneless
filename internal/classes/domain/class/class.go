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
	if created.IsZero() {
		return nil, errors.New("created timestamp is zero")
	}
	if updated.IsZero() {
		return nil, errors.New("updated timestamp is zero")
	}

	c := &Class{
		id:      classId,
		created: created,
		fields:  make([]*Field, len(fields)),
	}
	copy(c.fields, fields)
	c.UpdateName(name)
	c.UpdateParentID(parentId)
	c.updated = updated
	return c, nil
}

func (c Class) ID() string {
	return c.id
}

func (c Class) ParentID() string {
	return c.parentId
}

func (c *Class) UpdateParentID(parentId string) error {
	if parentId != "" && !id.IsValid(parentId) {
		return errors.New("invalid parent id")
	}
	c.parentId = parentId
	c.modified()
	return nil
}

func (c Class) Name() string {
	return c.name
}

func (c *Class) UpdateName(name string) error {
	if name == "" {
		return errors.New("empty class name")
	}
	c.name = name
	c.modified()
	return nil
}

func (c Class) Created() time.Time {
	return c.created
}

func (c Class) Updated() time.Time {
	return c.updated
}

func (c Class) Fields() []*Field {
	return c.fields
}

func (c *Class) UpdateFields(fields []*Field) error {
	c.fields = fields
	c.modified()
	return nil
}

func (c *Class) modified() {
	c.updated = time.Now()
}
