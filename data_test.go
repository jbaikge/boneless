package gocms

import "time"

func testClasses() []Class {
	return []Class{
		{
			Id:   "page",
			Name: "Page",
			Fields: []Field{
				{Name: "title", Sort: true},
				{Name: "content"},
			},
		},
		{
			Id:   "blog",
			Name: "Blog",
			Fields: []Field{
				{Name: "title", Sort: true},
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "news",
			Name: "News",
			Fields: []Field{
				{Name: "title", Sort: true},
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "event",
			Name: "Event",
			Fields: []Field{
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
			Fields: []Field{
				{Name: "title"},
				{Name: "start", Sort: true},
				{Name: "end"},
				{Name: "location"},
			},
		},
		{
			Id:   "speaker",
			Name: "Speaker",
			Fields: []Field{
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

func testDocuments() []Document {
	return []Document{
		{
			Id:      "page-1",
			ClassId: "page",
			Path:    "/",
			Values: map[string]interface{}{
				"title":   "Homepage",
				"content": "Homepage content goes here",
			},
		},
		{
			Id:      "blog-1",
			ClassId: "blog",
			Path:    "/blogs/blog-1",
			Values: map[string]interface{}{
				"title":     "Blog 1",
				"published": time.Unix(1659530000, 0),
				"excerpt":   "Blog excerpt",
				"content":   "Blog content",
				"author":    "author-1",
			},
		},
		{
			Id:      "blog-2",
			ClassId: "blog",
			Path:    "/blogs/blog-2",
			Values: map[string]interface{}{
				"title":     "Blog 2",
				"published": time.Unix(1659550000, 0),
				"excerpt":   "Blog excerpt",
				"content":   "Blog content",
				"author":    "author-3",
			},
		},
		{
			Id:      "blog-3",
			ClassId: "blog",
			Path:    "/blogs/blog-1",
			Values: map[string]interface{}{
				"title":     "Blog 3",
				"published": time.Unix(1659570000, 0),
				"excerpt":   "Blog excerpt",
				"content":   "Blog content",
				"author":    "author-2",
			},
		},
		{
			Id:      "news-1",
			ClassId: "news",
			Path:    "/news/news-1",
			Values: map[string]interface{}{
				"title":     "News 1",
				"published": time.Unix(1659540000, 0),
				"excerpt":   "News excerpt",
				"content":   "News content",
				"author":    "author-1",
			},
		},
		{
			Id:      "news-2",
			ClassId: "news",
			Path:    "/news/news-2",
			Values: map[string]interface{}{
				"title":     "News 2",
				"published": time.Unix(1659560000, 0),
				"excerpt":   "News excerpt",
				"content":   "News content",
				"author":    "author-3",
			},
		},
		{
			Id:      "news-3",
			ClassId: "news",
			Path:    "/news/news-1",
			Values: map[string]interface{}{
				"title":     "News 3",
				"published": time.Unix(1659580000, 0),
				"excerpt":   "News excerpt",
				"content":   "News content",
				"author":    "author-2",
			},
		},
		{
			Id:      "event-1",
			ClassId: "event",
			Path:    "/events/event-1",
			Values: map[string]interface{}{
				"start":        time.Unix(1659600000, 0),
				"end":          time.Unix(1659603600, 0),
				"date_display": "Hopefully in the Future",
				"time_display": "About an hour",
			},
		},
		{
			Id:       "session-1",
			ClassId:  "session",
			ParentId: "event-1",
			Values: map[string]interface{}{
				"title":    "First Session",
				"start":    "09:30",
				"end":      "09:45",
				"location": "Hall B",
			},
		},
		{
			Id:       "session-2",
			ClassId:  "session",
			ParentId: "event-1",
			Values: map[string]interface{}{
				"title":    "First Session",
				"start":    "09:45",
				"end":      "10:00",
				"location": "Hall B",
			},
		},
		{
			Id:       "session-3",
			ClassId:  "session",
			ParentId: "event-1",
			Values: map[string]interface{}{
				"title":    "First Session",
				"start":    "10:00",
				"end":      "10:30",
				"location": "Hall B",
			},
		},
		{
			Id:       "speaker-1",
			ClassId:  "speaker",
			ParentId: "session-1",
			Path:     "/speakers/speaker-1",
			Values: map[string]interface{}{
				"name":       "Felicity Grantham",
				"first_name": "Felicity",
				"last_name":  "Grantham",
				"sort_name":  "Grantham, Felicity",
			},
		},
		{
			Id:       "speaker-2",
			ClassId:  "speaker",
			ParentId: "session-2",
			Path:     "/speakers/speaker-2",
			Values: map[string]interface{}{
				"name":       "Sibby Begg",
				"first_name": "Sibby",
				"last_name":  "Begg",
				"sort_name":  "Begg, Sibby",
			},
		},
		{
			Id:       "speaker-3",
			ClassId:  "speaker",
			ParentId: "session-2",
			Path:     "/speakers/speaker-3",
			Values: map[string]interface{}{
				"name":       "Gordon Pont",
				"first_name": "Gordon",
				"last_name":  "Pont",
				"sort_name":  "Pont, Gordon",
			},
		},
		{
			Id:       "speaker-4",
			ClassId:  "speaker",
			ParentId: "session-1",
			Path:     "/speakers/speaker-4",
			Values: map[string]interface{}{
				"name":       "Alon Keohane",
				"first_name": "Alon",
				"last_name":  "Keohane",
				"sort_name":  "Keohane, Alon",
			},
		},
		{
			Id:       "speaker-5",
			ClassId:  "speaker",
			ParentId: "session-3",
			Path:     "/speakers/speaker-5",
			Values: map[string]interface{}{
				"name":       "Darlene Blackmore",
				"first_name": "Darlene",
				"last_name":  "Blackmore",
				"sort_name":  "Blackmore, Darlene",
			},
		},
		{
			Id:       "speaker-6",
			ClassId:  "speaker",
			ParentId: "session-1",
			Path:     "/speakers/speaker-6",
			Values: map[string]interface{}{
				"name":       "Wylie Bussey",
				"first_name": "Wylie",
				"last_name":  "Bussey",
				"sort_name":  "Bussey, Wylie",
			},
		},
	}
}
