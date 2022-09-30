package models

import (
	"time"
)

type Class struct {
	Id       string    `json:"id"`
	ParentId string    `json:"parent_id"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Fields   []Field   `json:"fields"`
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
