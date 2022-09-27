package testdata

import "github.com/jbaikge/boneless"

func Classes() []boneless.Class {
	return []boneless.Class{
		{
			Id:   "page",
			Name: "Page",
			Fields: []boneless.Field{
				{Name: "title", Sort: true},
				{Name: "content"},
			},
		},
		{
			Id:   "blog",
			Name: "Blog",
			Fields: []boneless.Field{
				{Name: "title", Sort: true},
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "news",
			Name: "News",
			Fields: []boneless.Field{
				{Name: "title", Sort: true},
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "event",
			Name: "Event",
			Fields: []boneless.Field{
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
			Fields: []boneless.Field{
				{Name: "title"},
				{Name: "start", Sort: true},
				{Name: "end"},
				{Name: "location"},
			},
		},
		{
			Id:   "speaker",
			Name: "Speaker",
			Fields: []boneless.Field{
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
