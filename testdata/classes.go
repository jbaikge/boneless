package testdata

import "github.com/jbaikge/boneless/models"

func Classes() []models.Class {
	return []models.Class{
		{
			Id:   "page",
			Name: "Page",
			Fields: []models.Field{
				{Name: "title", Sort: true},
				{Name: "content"},
			},
		},
		{
			Id:   "blog",
			Name: "Blog",
			Fields: []models.Field{
				{Name: "title", Sort: true},
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "news",
			Name: "News",
			Fields: []models.Field{
				{Name: "title", Sort: true},
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "event",
			Name: "Event",
			Fields: []models.Field{
				{Name: "title", Sort: true},
				{Name: "start", Sort: true},
				{Name: "end"},
				{Name: "date_display"},
				{Name: "time_display"},
			},
		},
		{
			Id:   "session",
			Name: "Session",
			Fields: []models.Field{
				{Name: "title"},
				{Name: "start", Sort: true},
				{Name: "end"},
				{Name: "location"},
			},
		},
		{
			Id:   "speaker",
			Name: "Speaker",
			Fields: []models.Field{
				{Name: "name"},
				{Name: "prefix"},
				{Name: "first_name"},
				{Name: "last_name"},
				{Name: "suffix"},
				{Name: "title"},
				{Name: "sort_name", Sort: true},
			},
		},
	}
}
