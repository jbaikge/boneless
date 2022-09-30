package models

import (
	"time"
)

type Form struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Created time.Time   `json:"created"`
	Updated time.Time   `json:"updated"`
	Schema  interface{} `json:"schema"`
}

type FormFilterSort struct {
	Field     string
	Direction string
}

type FormFilter struct {
	Sort  FormFilterSort
	Range Range
}
