package models

import (
	"time"
)

type Template struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Version int       `json:"version"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type TemplateFilter struct {
	Field       string
	SortReverse bool
	Range       Range
}
