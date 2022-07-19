package gocms

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const classIdPrefix = "class#"

var ErrNotFound = errors.New("item not found")

type dynamoClass struct {
	PrimaryKey  string
	SortKey     string
	Slug        string
	Name        string
	TableLabels string
	TableFields string
	Created     time.Time
	Updated     time.Time
	Fields      []Field
}

func (d *dynamoClass) FromClass(c *Class) {
	d.PrimaryKey = classIdPrefix + c.Id
	d.SortKey = classIdPrefix + c.Id
	d.Slug = c.Slug
	d.Name = c.Name
	d.TableLabels = c.TableLabels
	d.TableFields = c.TableFields
	d.Created = c.Created
	d.Updated = c.Updated
	d.Fields = make([]Field, len(c.Fields))
	copy(d.Fields, c.Fields)
}

func (d *dynamoClass) ToClass() (c Class) {
	c.Id = d.PrimaryKey[len(classIdPrefix):]
	c.Slug = d.Slug
	c.Name = d.Name
	c.TableLabels = d.TableLabels
	c.TableFields = d.TableFields
	c.Created = d.Created
	c.Updated = d.Updated
	c.Fields = make([]Field, len(d.Fields))
	copy(c.Fields, d.Fields)
	return
}

type DynamoDBRepository struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDBRepository(config aws.Config, table string) Repository {
	return &DynamoDBRepository{
		client: dynamodb.NewFromConfig(config),
		table:  table,
	}
}

func (repo DynamoDBRepository) CreateClass(ctx context.Context, class *Class) (err error) {
	dbClass := new(dynamoClass)
	dbClass.FromClass(class)

	item, err := attributevalue.MarshalMap(dbClass)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		TableName: &repo.table,
		Item:      item,
	}

	_, err = repo.client.PutItem(ctx, params)
	return
}

func (repo DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	prefixedId := classIdPrefix + id
	keyId, err := attributevalue.Marshal(prefixedId)
	if err != nil {
		return
	}

	params := &dynamodb.DeleteItemInput{
		TableName: &repo.table,
		Key: map[string]types.AttributeValue{
			"PrimaryKey": keyId,
			"SortKey":    keyId,
		},
	}

	_, err = repo.client.DeleteItem(ctx, params)
	return
}

func (repo DynamoDBRepository) GetClassById(ctx context.Context, id string) (class Class, err error) {
	prefixedId := classIdPrefix + id
	keyId, err := attributevalue.Marshal(prefixedId)
	if err != nil {
		return
	}

	params := &dynamodb.GetItemInput{
		TableName: &repo.table,
		Key: map[string]types.AttributeValue{
			"PrimaryKey": keyId,
			"SortKey":    keyId,
		},
	}

	response, err := repo.client.GetItem(ctx, params)
	if err != nil {
		return
	}

	// Check for no-item-found condition
	if len(response.Item) == 0 {
		err = ErrNotFound
		return
	}

	dbClass := new(dynamoClass)
	if err = attributevalue.UnmarshalMap(response.Item, dbClass); err != nil {
		return
	}
	class = dbClass.ToClass()

	return
}

func (repo DynamoDBRepository) GetClassList(ctx context.Context, filter ClassFilter) (classes []Class, r Range, err error) {
	// Ref: https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go
	filterPK := expression.Name("PrimaryKey").BeginsWith(classIdPrefix)
	expr, err := expression.NewBuilder().WithFilter(filterPK).Build()
	if err != nil {
		return
	}

	params := &dynamodb.ScanInput{
		TableName:                 &repo.table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}
	result, err := repo.client.Scan(ctx, params)
	if err != nil {
		return
	}

	dbClasses := make([]dynamoClass, 0, result.Count)
	if err = attributevalue.UnmarshalListOfMaps(result.Items, &dbClasses); err != nil {
		return
	}

	classes = make([]Class, len(dbClasses))
	for i, dbClass := range dbClasses {
		classes[i] = dbClass.ToClass()
	}
	return
}

func (repo DynamoDBRepository) UpdateClass(ctx context.Context, class *Class) (err error) {
	dbClass := new(dynamoClass)
	dbClass.FromClass(class)
	item, err := attributevalue.MarshalMap(dbClass)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		TableName: &repo.table,
		Item:      item,
	}

	_, err = repo.client.PutItem(ctx, params)
	return
}
