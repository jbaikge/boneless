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

func NewClass(classId string, parentId string, name string, created time.Time, updated time.Time, fields []*Field) (*Class, error) {
	if classId != "" && !id.IsValid(classId) {
		return nil, errors.New("invalid class id")
	}
	if parentId != "" && !id.IsValid(parentId) {
		return nil, errors.New("invalid parent id")
	}
	if name == "" {
		return nil, errors.New("empty class name")
	}
	if created.IsZero() {
		created = time.Now()
	}
	if updated.IsZero() {
		updated = created
	}
	if classId == "" {
		classId = id.NewWithTime(created)
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
