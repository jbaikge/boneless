package dynamodb

import (
	"context"
	"testing"
	"time"

	"github.com/jbaikge/boneless"
	"github.com/zeebo/assert"
)

func TestTemplates(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + "templates",
		Table:  dynamoPrefix + "Templates",
	}
	repo, err := newRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("BasicInOut", func(t *testing.T) {
		template := boneless.Template{
			Id:      "template-1",
			Name:    "Test Template",
			Created: time.Now(),
			Updated: time.Now(),
			Body:    "This is a test template\n\nThis is another line in the test template",
		}
		assert.NoError(t, repo.CreateTemplate(ctx, &template))

		check, err := repo.GetTemplateById(ctx, template.Id)
		assert.NoError(t, err)
		assert.Equal(t, template.Body, check.Body)
	})

	t.Run("Version2", func(t *testing.T) {
		template := boneless.Template{
			Id:      "two-versions",
			Name:    "Version 1",
			Created: time.Now(),
			Updated: time.Now(),
			Body:    "This is version one content",
		}
		assert.NoError(t, repo.CreateTemplate(ctx, &template))

		template.Name = "Version 2"
		template.Body = "This is version two content"
		assert.NoError(t, repo.UpdateTemplate(ctx, &template))

		check, err := repo.GetTemplateById(ctx, template.Id)
		assert.NoError(t, err)
		assert.Equal(t, template.Body, check.Body)
	})

	t.Run("DeleteMultiVersion", func(t *testing.T) {
		template := boneless.Template{
			Id:      "delete-me",
			Name:    "Delete Me v1",
			Created: time.Now(),
			Updated: time.Now(),
			Body:    "Delete this content v1",
		}
		assert.NoError(t, repo.CreateTemplate(ctx, &template))

		template.Name = "Delete Me v2"
		template.Body = template.Body[len(template.Body)-2:] + "2"
		assert.NoError(t, repo.UpdateTemplate(ctx, &template))

		assert.NoError(t, repo.DeleteTemplate(ctx, template.Id))

		_, err := repo.GetTemplateById(ctx, template.Id)
		assert.Equal(t, ErrNotExist, err)
	})
}
