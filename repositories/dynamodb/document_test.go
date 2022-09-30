package dynamodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jbaikge/boneless/models"
	"github.com/jbaikge/boneless/testdata"
	"github.com/zeebo/assert"
)

func TestDocumentList(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + "documentlist",
		Table:  dynamoPrefix + "DocumentList",
	}
	repo, err := newRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	for _, class := range testdata.Classes() {
		assert.NoError(t, repo.CreateClass(ctx, &class))
	}

	for _, document := range testdata.Documents() {
		assert.NoError(t, repo.CreateDocument(ctx, &document))
	}

	t.Run("ListPagesByTitle", func(t *testing.T) {
		filter := models.DocumentFilter{
			ClassId: "page",
			Sort:    models.DocumentFilterSort{Field: "title"},
			Range:   models.Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, models.Range{End: 1, Size: 2}, r)
		assert.Equal(t, 2, len(docs))
		assert.Equal(t, "page-2", docs[0].Id)
		assert.Equal(t, "page-1", docs[1].Id)
	})

	t.Run("ListSessionsByStart", func(t *testing.T) {
		filter := models.DocumentFilter{
			ClassId: "session",
			Sort:    models.DocumentFilterSort{Field: "start"},
			Range:   models.Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, models.Range{End: 4, Size: 5}, r)
		assert.Equal(t, 5, len(docs))
	})

	t.Run("ListSessionsByEvent", func(t *testing.T) {
		filter := models.DocumentFilter{
			ClassId:  "session",
			ParentId: "event-1",
			Sort:     models.DocumentFilterSort{Field: "start"},
			Range:    models.Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, models.Range{End: 2, Size: 3}, r)
		assert.Equal(t, 3, len(docs))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		// Should list all documents, sorted by descending creation date
		filter := models.DocumentFilter{
			Range: models.Range{End: 99},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, models.Range{End: 19, Size: 20}, r)
		assert.Equal(t, 20, len(docs))
		assert.Equal(t, "speaker-6", docs[0].Id)
		assert.Equal(t, "event-1", docs[19].Id)
	})

	t.Run("AllChildren", func(t *testing.T) {
		// Should give all children of a parent document, regardless of class
		// then sorted by descending creation date
		filter := models.DocumentFilter{
			ParentId: "event-1",
			Range:    models.Range{End: 99},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, models.Range{End: 2, Size: 3}, r)
		assert.Equal(t, 3, len(docs))
		assert.Equal(t, "session-3", docs[0].Id)
		assert.Equal(t, "session-1", docs[2].Id)
	})
}

func TestTableScan(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + "tablescan",
		Table:  dynamoPrefix + "TableScan",
	}
	repo, err := newRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	class := models.Class{
		Id:      "class",
		Name:    "Class",
		Created: time.Now(),
		Updated: time.Now(),
		Fields: []models.Field{
			{Name: "sort_field", Sort: true},
			{Name: "scan_field"},
			{Name: "empty_field"},
		},
	}
	assert.NoError(t, repo.CreateClass(ctx, &class))

	doc := models.Document{
		ClassId: "class",
		Values:  make(map[string]interface{}),
	}
	data := [][]string{
		{"doc1", "B", "C"},
		{"doc2", "D", "A"},
		{"doc3", "C", "D"},
		{"doc4", "A", "B"},
	}
	for _, set := range data {
		doc.Id = set[0]
		doc.Values["sort_field"] = set[1]
		doc.Values["scan_field"] = set[2]
		assert.NoError(t, repo.CreateDocument(ctx, &doc))
	}

	testTable := []struct {
		Name   string
		Filter models.DocumentFilter
		Expect []string
	}{
		{
			Name: "UseSort",
			Filter: models.DocumentFilter{
				Range: models.Range{End: 9},
				Sort:  models.DocumentFilterSort{Field: "sort_field"},
			},
			Expect: []string{"doc4", "doc1", "doc3", "doc2"},
		},
		{
			Name: "UseScan",
			Filter: models.DocumentFilter{
				Range: models.Range{End: 9},
				Sort:  models.DocumentFilterSort{Field: "scan_field"},
			},
			Expect: []string{"doc2", "doc4", "doc1", "doc3"},
		},
		{
			Name: "UseEmpty",
			Filter: models.DocumentFilter{
				Range: models.Range{End: 9},
				Sort:  models.DocumentFilterSort{Field: "empty_field"},
			},
			Expect: []string{"doc2", "doc4", "doc3", "doc1"},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			docs, r, err := repo.GetDocumentList(ctx, test.Filter)
			assert.NoError(t, err)
			assert.Equal(t, len(data), r.Size)
			for i, doc := range docs {
				assert.Equal(t, test.Expect[i], doc.Id)
			}
		})
	}
}

func TestValues(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + "values",
		Table:  dynamoPrefix + "Values",
	}
	repo, err := newRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	class := models.Class{
		Id:      "class",
		Name:    "Class",
		Created: time.Now(),
		Updated: time.Now(),
		Fields: []models.Field{
			{Name: "field1", Sort: true},
			{Name: "field2"},
			{Name: "field3", Sort: true},
			{Name: "field4"},
		},
	}
	assert.NoError(t, repo.CreateClass(ctx, &class))

	doc := models.Document{
		Id:      "doc",
		ClassId: "class",
		Values: map[string]interface{}{
			"field1": "abc",
			"field2": "def",
			"field3": "ghi",
			"field4": "jkl",
		},
	}
	assert.NoError(t, repo.CreateDocument(ctx, &doc))

	docCheck, err := repo.GetDocumentById(ctx, doc.Id)
	assert.NoError(t, err)

	assert.DeepEqual(t, doc.Values, docCheck.Values)
	// Can't use DeepEqual for ... reasons?
	// Need to cast everything to string for assert.Equal
	for key, expect := range doc.Values {
		assert.Equal(t, fmt.Sprint(expect), fmt.Sprint(docCheck.Values[key]))
	}

	extra := models.Document{
		Id:      "extra",
		ClassId: "class",
		Values: map[string]interface{}{
			"field0": "abc", // extra
			"field1": "def",
			// field2 intentionally left out
			"field3": "ghi",
			"field4": "jkl",
			"field5": "mno", // extra
		},
	}
	assert.NoError(t, repo.CreateDocument(ctx, &extra))

	extraCheck, err := repo.GetDocumentById(ctx, extra.Id)
	assert.NoError(t, err)

	// Should contain all the submitted values and match what was sent into the
	// repo.
	assert.DeepEqual(t, extra.Values, extraCheck.Values)
}
