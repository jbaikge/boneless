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
		Id:     "from_class",
		Name:   t.Name(),
		Fields: []Field{{Name: "field_1"}},
	}

	dc := new(dynamoClass)
	dc.FromClass(&class)

	assert.Equal(t, dynamoClassPrefix+class.Id, dc.PK)
	assert.Equal(t, "class_v0", dc.SK)
	assert.Equal(t, t.Name(), dc.Name)
	assert.DeepEqual(t, class.Fields, dc.Fields)
}

func TestDynamoToClass(t *testing.T) {
	id := "to_class"
	dc := dynamoClass{
		PK:     dynamoClassPrefix + id,
		SK:     "class_v0",
		Name:   t.Name(),
		Fields: []Field{{Name: "field_1"}},
	}
	class := dc.ToClass()

	assert.Equal(t, id, class.Id)
	assert.Equal(t, t.Name(), class.Name)
	assert.DeepEqual(t, dc.Fields, class.Fields)
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

	t.Run("GetClassByIdSuccess", func(t *testing.T) {
		class := Class{
			Id:   "get_success",
			Name: t.Name(),
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		check, err := repo.GetClassById(ctx, class.Id)
		assert.NoError(t, err)
		assert.Equal(t, class.Id, check.Id)
		assert.Equal(t, class.Name, check.Name)
	})

	t.Run("GetClassByIdFail", func(t *testing.T) {
		_, err := repo.GetClassById(ctx, "get_fail")
		assert.Equal(t, ErrNotExist, err)
	})

	t.Run("UpdateClass", func(t *testing.T) {
		class := Class{
			Id:          "update_class",
			Name:        t.Name(),
			TableLabels: "Field 1; Field 2",
			TableFields: "field_1; field_2",
			Created:     time.Now(),
			Updated:     time.Now(),
			Fields:      []Field{{Name: "field_2"}, {Name: "field_1"}},
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		class.Name = t.Name() + "-Updated"
		class.TableLabels += "; Field 3"
		class.TableFields += "; field_3"
		class.Updated = time.UnixMicro(time.Now().UnixMicro())
		class.Fields = append(class.Fields, Field{Name: "field_3"})
		assert.NoError(t, repo.UpdateClass(ctx, &class))

		check, err := repo.GetClassById(ctx, class.Id)
		assert.NoError(t, err)
		assert.Equal(t, class.Name, check.Name)
		assert.Equal(t, class.TableFields, check.TableFields)
		assert.Equal(t, class.TableLabels, check.TableLabels)
		assert.Equal(t, class.Updated, check.Updated)
		assert.DeepEqual(t, class.Fields, check.Fields)
	})

	t.Run("DeleteClass", func(t *testing.T) {
		class := Class{
			Id:   "delete_class",
			Name: t.Name(),
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))
		assert.NoError(t, repo.DeleteClass(ctx, class.Id))
		_, err := repo.GetClassById(ctx, class.Id)
		assert.Equal(t, ErrNotExist, err)
	})
}
