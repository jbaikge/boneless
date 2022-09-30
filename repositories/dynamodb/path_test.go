package dynamodb

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jbaikge/boneless/models"
	"github.com/zeebo/assert"
)

func TestPathUpdate(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + strings.ToLower(t.Name()),
		Table:  dynamoPrefix + t.Name(),
	}
	repo, err := newRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	class := models.Class{
		Id:     "class",
		Name:   "Class",
		Fields: []models.Field{},
	}
	assert.NoError(t, repo.CreateClass(ctx, &class))

	docs := []models.Document{
		{
			Id:      "doc1",
			ClassId: "class",
			Path:    "/doc/1",
		},
		{
			Id:      "doc2",
			ClassId: "class",
			Path:    "/doc/2",
		},
	}
	for _, doc := range docs {
		assert.NoError(t, repo.CreateDocument(ctx, &doc))
	}

	t.Run("DeletePath", func(t *testing.T) {
		oldPath := docs[0].Path
		docs[0].Path = ""
		for _, doc := range docs {
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		}

		_, err := repo.GetDocumentByPath(ctx, oldPath)
		assert.True(t, errors.Is(err, ErrNotExist))
	})

	t.Run("OverwritePath", func(t *testing.T) {
		docs[0].Path = docs[1].Path
		assert.Error(t, repo.UpdateDocument(ctx, &docs[0]))
		assert.NoError(t, repo.UpdateDocument(ctx, &docs[1]))
	})
}
