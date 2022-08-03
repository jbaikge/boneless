package gocms

import "time"

func testClasses() []Class {
	return []Class{
		{
			Id:   "page",
			Name: "Page",
			Fields: []Field{
				{Name: "content"},
			},
		},
		{
			Id:   "blog",
			Name: "Blog",
			Fields: []Field{
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "news",
			Name: "News",
			Fields: []Field{
				{Name: "published", Sort: true},
				{Name: "excerpt"},
				{Name: "author"},
			},
		},
		{
			Id:   "event",
			Name: "Event",
			Fields: []Field{
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
				"content": "Homepage content goes here",
			},
		},
		{
			Id:      "blog-1",
			ClassId: "blog",
			Name:    "Blog 1",
			Path:    "/blogs/blog-1",
			Values: map[string]interface{}{
				"published": time.Unix(1659530000, 0),
				"excerpt":   "Blog excerpt",
				"content":   "Blog content",
				"author":    "author-1",
			},
		},
		{
			Id:      "blog-2",
			ClassId: "blog",
			Name:    "Blog 2",
			Path:    "/blogs/blog-2",
			Values: map[string]interface{}{
				"published": time.Unix(1659550000, 0),
				"excerpt":   "Blog excerpt",
				"content":   "Blog content",
				"author":    "author-3",
			},
		},
		{
			Id:      "blog-3",
			ClassId: "blog",
			Name:    "Blog 3",
			Path:    "/blogs/blog-1",
			Values: map[string]interface{}{
				"published": time.Unix(1659570000, 0),
				"excerpt":   "Blog excerpt",
				"content":   "Blog content",
				"author":    "author-2",
			},
		},
		{
			Id:      "news-1",
			ClassId: "news",
			Name:    "News 1",
			Path:    "/news/news-1",
			Values: map[string]interface{}{
				"published": time.Unix(1659540000, 0),
				"excerpt":   "News excerpt",
				"content":   "News content",
				"author":    "author-1",
			},
		},
		{
			Id:      "news-2",
			ClassId: "news",
			Name:    "News 2",
			Path:    "/news/news-2",
			Values: map[string]interface{}{
				"published": time.Unix(1659560000, 0),
				"excerpt":   "News excerpt",
				"content":   "News content",
				"author":    "author-3",
			},
		},
		{
			Id:      "news-3",
			ClassId: "news",
			Name:    "News 3",
			Path:    "/news/news-1",
			Values: map[string]interface{}{
				"published": time.Unix(1659580000, 0),
				"excerpt":   "News excerpt",
				"content":   "News content",
				"author":    "author-2",
			},
		},
		{
			Id:      "event-1",
			ClassId: "event",
			Name:    "First Event",
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
			Name:     "Session 1",
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
			Name:     "Session 2",
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
			Name:     "Session 3",
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
			Name:     "Felicity Grantham",
			Values: map[string]interface{}{
				"first_name": "Felicity",
				"last_name":  "Grantham",
				"sort_name":  "Grantham, Felicity",
			},
		},
		{
			Id:       "speaker-2",
			ClassId:  "speaker",
			ParentId: "session-2",
			Name:     "Sibby Begg",
			Values: map[string]interface{}{
				"first_name": "Sibby",
				"last_name":  "Begg",
				"sort_name":  "Begg, Sibby",
			},
		},
		{
			Id:       "speaker-3",
			ClassId:  "speaker",
			ParentId: "session-2",
			Name:     "Gordon Pont",
			Values: map[string]interface{}{
				"first_name": "Gordon",
				"last_name":  "Pont",
				"sort_name":  "Pont, Gordon",
			},
		},
		{
			Id:       "speaker-4",
			ClassId:  "speaker",
			ParentId: "session-1",
			Name:     "Alon Keohane",
			Values: map[string]interface{}{
				"first_name": "Alon",
				"last_name":  "Keohane",
				"sort_name":  "Keohane, Alon",
			},
		},
		{
			Id:       "speaker-5",
			ClassId:  "speaker",
			ParentId: "session-3",
			Name:     "Darlene Blackmore",
			Values: map[string]interface{}{
				"first_name": "Darlene",
				"last_name":  "Blackmore",
				"sort_name":  "Blackmore, Darlene",
			},
		},
		{
			Id:       "speaker-1",
			ClassId:  "speaker",
			ParentId: "session-1",
			Name:     "Wylie Bussey",
			Values: map[string]interface{}{
				"first_name": "Wylie",
				"last_name":  "Bussey",
				"sort_name":  "Bussey, Wylie",
			},
		},
	}
}
