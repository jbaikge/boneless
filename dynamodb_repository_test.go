package gocms

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/zeebo/assert"
)

func TestDynamoDBRepository(t *testing.T) {
	table := t.Name() + time.Now().Format("-20060102-150405")

	endpointResolverFunc := func(service string, region string, options ...interface{}) (endpoint aws.Endpoint, err error) {
		endpoint = aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:8000",
			SigningRegion: "local", // local to be compatible with dy command
		}
		return
	}

	endpointResolver := aws.EndpointResolverWithOptionsFunc(endpointResolverFunc)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)
	assert.NoError(t, err)

	// TODO Move this block elsewhere, maybe into the makefile or somewhere it
	// can match what is described in the deployment stack.
	client := dynamodb.NewFromConfig(cfg)
	_, err = client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: &table,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	assert.NoError(t, err)

	repo := NewDynamoDBRepository(cfg, table)

	// Use of UnixMicro trims off the m variable in the Time struct to make
	// reflect.DeepEqual function properly
	now := time.UnixMicro(time.Now().UnixMicro())
	class1 := Class{
		Id:      "class#1",
		Name:    "My Class",
		Slug:    "my_class",
		Created: now,
		Updated: now,
	}

	t.Run("InsertClass", func(t *testing.T) {
		assert.NoError(t, repo.InsertClass(context.Background(), &class1))
	})

	t.Run("GetClassById", func(t *testing.T) {
		check, err := repo.GetClassById(context.Background(), class1.Id)
		assert.NoError(t, err)
		assert.DeepEqual(t, class1, check)
	})

	t.Run("UpdateClass", func(t *testing.T) {
		class1.Slug = "my_new_class"
		assert.NoError(t, repo.UpdateClass(context.Background(), &class1))
		check, err := repo.GetClassById(context.Background(), class1.Id)
		assert.NoError(t, err)
		assert.DeepEqual(t, class1, check)
	})

	t.Run("GetAllClasses", func(t *testing.T) {
		count := 10
		for i := 2; i <= count; i++ {
			class := Class{
				Id:      fmt.Sprintf("class#%d", i),
				Name:    fmt.Sprintf("Class %d", i),
				Slug:    fmt.Sprintf("class_%d", i),
				Created: now,
				Updated: now,
			}
			assert.NoError(t, repo.InsertClass(context.Background(), &class))
		}
		classes, err := repo.GetAllClasses(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, count, len(classes))
	})

	t.Run("DeleteClass", func(t *testing.T) {
		assert.NoError(t, repo.DeleteClass(context.Background(), class1.Id))
		_, err := repo.GetClassById(context.Background(), class1.Id)
		assert.Error(t, err)
	})
}
