package dynamodb

import (
	"context"
	"strings"
	"testing"

	"github.com/jbaikge/boneless/models"
	"github.com/zeebo/assert"
)

func TestSortUpdate(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + strings.ToLower(t.Name()),
		Table:  dynamoPrefix + t.Name(),
	}
	repo, err := newRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	class := models.Class{
		Id:   "class",
		Name: "Class",
		Fields: []models.Field{
			{Name: "sort_field_1", Sort: true},
		},
	}
	assert.NoError(t, repo.CreateClass(ctx, &class))

	docs := []models.Document{
		{
			Id:      "doc1",
			ClassId: "class",
			Values: map[string]interface{}{
				class.Fields[0].Name: "abc",
			},
		},
		{
			Id:      "doc2",
			ClassId: "class",
			Values: map[string]interface{}{
				class.Fields[0].Name: "xyz",
			},
		},
	}
	for _, doc := range docs {
		assert.NoError(t, repo.CreateDocument(ctx, &doc))
	}

	filter := models.DocumentFilter{
		Sort: models.DocumentFilterSort{
			Field:     class.Fields[0].Name,
			Direction: "ASC",
		},
		Range: models.Range{
			End: 9,
		},
	}
	t.Run("InitialOrder", func(t *testing.T) {
		list, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, len(docs), r.Size)
		assert.DeepEqual(t, []string{"doc1", "doc2"}, []string{list[0].Id, list[1].Id})
	})

	// Push doc2 in front of doc1 and completely change the values to ensure
	// the update did the right thing.
	t.Run("Reorder", func(t *testing.T) {
		docs[0].Values[class.Fields[0].Name] = "tuv"
		docs[1].Values[class.Fields[0].Name] = "def"
		for _, doc := range docs {
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		}
	})

	t.Run("NewOrder", func(t *testing.T) {
		list, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, len(docs), r.Size)
		assert.DeepEqual(t, []string{"doc2", "doc1"}, []string{list[0].Id, list[1].Id})
	})
}
