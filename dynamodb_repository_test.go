package gocms

import (
	"context"
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

	t.Run("InsertClass", func(t *testing.T) {
		now := time.Now()
		class := Class{
			Id:      "class#1",
			Name:    "My Class",
			Slug:    "my_class",
			Created: now,
			Updated: now,
		}
		assert.NoError(t, repo.InsertClass(context.Background(), &class))
	})
}
