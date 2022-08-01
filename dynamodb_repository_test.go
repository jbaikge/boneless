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

var dynamoTablePrefix = time.Now().Format("20060102-150405-")

func testDynamoConfig(tableName string) (cfg aws.Config, err error) {
	endpointResolverFunc := func(service string, region string, options ...interface{}) (endpoint aws.Endpoint, err error) {
		endpoint = aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:8000",
			SigningRegion: "local",
		}
		return
	}
	endpointResolver := aws.EndpointResolverWithOptionsFunc(endpointResolverFunc)

	cfg, err = config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)
	if err != nil {
		return
	}

	// Build table before returning
	createTable := &dynamodb.CreateTableInput{
		TableName:   &tableName,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
	}

	client := dynamodb.NewFromConfig(cfg)
	if _, err = client.CreateTable(context.Background(), createTable); err != nil {
		return
	}

	return
}

func TestDynamoFromClass(t *testing.T) {
	class := Class{
		Id:   "from_class",
		Name: t.Name(),
	}

	dc := new(dynamoClass)
	dc.FromClass(&class)

	assert.Equal(t, dynamoClassPrefix+class.Id, dc.PK)
	assert.Equal(t, "class_v0", dc.SK)
	assert.Equal(t, t.Name(), dc.Name)
}

func TestDynamoToClass(t *testing.T) {
	id := "to_class"
	dc := dynamoClass{
		PK:   dynamoClassPrefix + id,
		SK:   "class_v0",
		Name: t.Name(),
	}
	class := dc.ToClass()

	assert.Equal(t, id, class.Id)
	assert.Equal(t, t.Name(), class.Name)
}

func TestDynamoDBRepository(t *testing.T) {
	resources := DynamoDBResources{
		Table: dynamoTablePrefix + "Test",
	}
	cfg, err := testDynamoConfig(resources.Table)
	assert.NoError(t, err)

	repo := NewDynamoDBRepository(cfg, resources)
	ctx := context.Background()

	t.Run("CreateClass", func(t *testing.T) {
		class := Class{
			Id:   "class_create",
			Name: t.Name(),
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))
	})
}
