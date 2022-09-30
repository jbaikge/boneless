package models

import (
	"strings"
	"time"
)

type Document struct {
	Id         string                 `json:"id"`
	ClassId    string                 `json:"class_id"`
	ParentId   string                 `json:"parent_id"`
	TemplateId string                 `json:"template_id"`
	Path       string                 `json:"path"`
	Version    int                    `json:"version"`
	Created    time.Time              `json:"created"`
	Updated    time.Time              `json:"updated"`
	Values     map[string]interface{} `json:"values"`
}

type DocumentFilterSort struct {
	Field     string
	Direction string
}

func (dfs DocumentFilterSort) Ascending() bool {
	return strings.ToUpper(dfs.Direction) == "ASC" || dfs.Direction == ""
}

func (dfs DocumentFilterSort) Descending() bool {
	return strings.ToUpper(dfs.Direction) == "DESC"
}

type DocumentFilter struct {
	ClassId  string
	ParentId string
	Sort     DocumentFilterSort
	Range    Range
}
