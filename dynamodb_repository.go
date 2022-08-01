package gocms

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	dynamoClassPrefix = "class#"
)

var (
	ErrNotExist = errors.New("item does not exist")
)

type dynamoClass struct {
	PK          string
	SK          string
	Name        string
	TableFields string
	TableLabels string
	Created     time.Time
	Updated     time.Time
	Fields      []Field
}

func (dyn *dynamoClass) FromClass(c *Class) {
	dyn.PK = dynamoClassPrefix + c.Id
	dyn.SK = "class_v0"
	dyn.Name = c.Name
	dyn.TableFields = c.TableFields
	dyn.TableLabels = c.TableLabels
	dyn.Created = c.Created
	dyn.Updated = c.Updated
	dyn.Fields = make([]Field, len(c.Fields))
	copy(dyn.Fields, c.Fields)
}

func (dyn dynamoClass) ToClass() (c Class) {
	c.Id = dyn.PK[len(dynamoClassPrefix):]
	c.Name = dyn.Name
	c.TableFields = dyn.TableFields
	c.TableLabels = dyn.TableLabels
	c.Created = dyn.Created
	c.Updated = dyn.Updated
	c.Fields = make([]Field, len(dyn.Fields))
	copy(c.Fields, dyn.Fields)
	return
}

type DynamoDBResources struct {
	Bucket string
	Table  string
}

func (res *DynamoDBResources) FromEnv() {
	res.Bucket = os.Getenv("REPOSITORY_BUCKET")
	res.Table = os.Getenv("REPOSITORY_TABLE")
}

type DynamoDBRepository struct {
	client    *dynamodb.Client
	resources DynamoDBResources
}

// Ref: https://dynobase.dev/dynamodb-golang-query-examples/
func NewDynamoDBRepository(config aws.Config, resources DynamoDBResources) Repository {
	return &DynamoDBRepository{
		client:    dynamodb.NewFromConfig(config),
		resources: resources,
	}
}

// Class Methods
func (repo DynamoDBRepository) CreateClass(ctx context.Context, class *Class) (err error) {
	dc := new(dynamoClass)
	dc.FromClass(class)
	item, err := attributevalue.MarshalMap(dc)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: &repo.resources.Table,
	}
	_, err = repo.client.PutItem(ctx, params)

	return
}

func (repo DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	pkId, err := attributevalue.Marshal(dynamoClassPrefix + id)
	if err != nil {
		return
	}

	skId, err := attributevalue.Marshal("class_v0")
	if err != nil {
		return
	}

	params := &dynamodb.DeleteItemInput{
		TableName: &repo.resources.Table,
		Key: map[string]types.AttributeValue{
			"PK": pkId,
			"SK": skId,
		},
	}
	_, err = repo.client.DeleteItem(ctx, params)

	return
}

func (repo DynamoDBRepository) GetClassById(ctx context.Context, id string) (class Class, err error) {
	pkId, err := attributevalue.Marshal(dynamoClassPrefix + id)
	if err != nil {
		return
	}

	skId, err := attributevalue.Marshal("class_v0")
	if err != nil {
		return
	}

	params := &dynamodb.GetItemInput{
		TableName: &repo.resources.Table,
		Key: map[string]types.AttributeValue{
			"PK": pkId,
			"SK": skId,
		},
	}
	response, err := repo.client.GetItem(ctx, params)

	if len(response.Item) == 0 {
		return class, ErrNotExist
	}

	dc := new(dynamoClass)
	if err = attributevalue.UnmarshalMap(response.Item, dc); err != nil {
		return
	}
	class = dc.ToClass()

	return
}

func (repo DynamoDBRepository) GetClassList(ctx context.Context, filter ClassFilter) (list []Class, r Range, err error) {
	return
}

func (repo DynamoDBRepository) UpdateClass(ctx context.Context, class *Class) (err error) {
	pkId, err := attributevalue.Marshal(dynamoClassPrefix + class.Id)
	if err != nil {
		return
	}

	skId, err := attributevalue.Marshal("class_v0")
	if err != nil {
		return
	}

	rawValues := map[string]interface{}{
		"Name":        class.Name,
		"TableFields": class.TableFields,
		"TableLabels": class.TableLabels,
		"Fields":      class.Fields,
		"Updated":     class.Updated,
	}

	sets := make([]string, 0, len(rawValues))
	values := make(map[string]types.AttributeValue)
	names := make(map[string]string)
	for key, value := range rawValues {
		index := len(sets)
		placeholder := ":" + key
		if values[placeholder], err = attributevalue.Marshal(value); err != nil {
			return fmt.Errorf("failed to marshal %s: %w", key, err)
		}
		sets = append(sets, fmt.Sprintf("#param_%d = %s", index, placeholder))
		names[fmt.Sprintf("#param_%d", index)] = key
	}
	updateExpression := "SET " + strings.Join(sets, ", ")

	params := &dynamodb.UpdateItemInput{
		TableName: &repo.resources.Table,
		Key: map[string]types.AttributeValue{
			"PK": pkId,
			"SK": skId,
		},
		UpdateExpression:          &updateExpression,
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
	}

	_, err = repo.client.UpdateItem(ctx, params)

	return
}

// Document Methods
func (repo DynamoDBRepository) CreateDocument(ctx context.Context, doc *Document) (err error) {
	return
}

func (repo DynamoDBRepository) DeleteDocument(ctx context.Context, id string) (err error) {
	return
}

func (repo DynamoDBRepository) GetDocumentById(ctx context.Context, id string) (doc Document, err error) {
	return
}

func (repo DynamoDBRepository) GetDocumentList(ctx context.Context, filter DocumentFilter) (list []Document, r Range, err error) {
	return
}

func (repo DynamoDBRepository) UpdateDocument(ctx context.Context, doc *Document) (err error) {
	return
}
