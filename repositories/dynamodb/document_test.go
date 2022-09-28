package dynamodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jbaikge/boneless"
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
		filter := boneless.DocumentFilter{
			ClassId: "page",
			Sort:    boneless.DocumentFilterSort{Field: "title"},
			Range:   boneless.Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, boneless.Range{End: 1, Size: 2}, r)
		assert.Equal(t, 2, len(docs))
		assert.Equal(t, "page-2", docs[0].Id)
		assert.Equal(t, "page-1", docs[1].Id)
	})

	t.Run("ListSessionsByStart", func(t *testing.T) {
		filter := boneless.DocumentFilter{
			ClassId: "session",
			Sort:    boneless.DocumentFilterSort{Field: "start"},
			Range:   boneless.Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, boneless.Range{End: 4, Size: 5}, r)
		assert.Equal(t, 5, len(docs))
	})

	t.Run("ListSessionsByEvent", func(t *testing.T) {
		filter := boneless.DocumentFilter{
			ClassId:  "session",
			ParentId: "event-1",
			Sort:     boneless.DocumentFilterSort{Field: "start"},
			Range:    boneless.Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, boneless.Range{End: 2, Size: 3}, r)
		assert.Equal(t, 3, len(docs))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		// Should list all documents, sorted by descending creation date
		filter := boneless.DocumentFilter{
			Range: boneless.Range{End: 99},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, boneless.Range{End: 19, Size: 20}, r)
		assert.Equal(t, 20, len(docs))
		assert.Equal(t, "speaker-6", docs[0].Id)
		assert.Equal(t, "event-1", docs[19].Id)
	})

	t.Run("AllChildren", func(t *testing.T) {
		// Should give all children of a parent document, regardless of class
		// then sorted by descending creation date
		filter := boneless.DocumentFilter{
			ParentId: "event-1",
			Range:    boneless.Range{End: 99},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, boneless.Range{End: 2, Size: 3}, r)
		assert.Equal(t, 3, len(docs))
		assert.Equal(t, "session-3", docs[0].Id)
		assert.Equal(t, "session-1", docs[2].Id)
	})
}

func TestValues(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + "values",
		Table:  dynamoPrefix + "Values",
	}
	repo, err := newRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	class := boneless.Class{
		Id:      "class",
		Name:    "Class",
		Created: time.Now(),
		Updated: time.Now(),
		Fields: []boneless.Field{
			{Name: "field1", Sort: true},
			{Name: "field2"},
			{Name: "field3", Sort: true},
			{Name: "field4"},
		},
	}
	assert.NoError(t, repo.CreateClass(ctx, &class))

	doc := boneless.Document{
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

	extra := boneless.Document{
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
