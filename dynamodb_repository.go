package gocms

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ErrNotFound = errors.New("item not found")

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

func (r DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	keyId, err := attributevalue.Marshal(id)
	if err != nil {
		return
	}

	params := &dynamodb.DeleteItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"PrimaryKey": keyId,
			"SortKey":    keyId,
		},
	}

	_, err = r.client.DeleteItem(ctx, params)
	return
}

func (r DynamoDBRepository) GetAllClasses(ctx context.Context) (classes []Class, err error) {
	// Ref: https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go
	filterPK := expression.Name("PrimaryKey").BeginsWith("class#")
	filterSK := expression.Key("SortKey").BeginsWith("class#")
	expr, err := expression.NewBuilder().WithFilter(filterPK).WithKeyCondition(filterSK).Build()
	if err != nil {
		return
	}

	params := &dynamodb.QueryInput{
		TableName:                 &r.table,
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	result, err := r.client.Query(ctx, params)
	if err != nil {
		return
	}

	classes = make([]Class, 0, result.Count)
	err = attributevalue.UnmarshalListOfMaps(result.Items, &classes)
	return
}

func (r DynamoDBRepository) GetClassById(ctx context.Context, id string) (class Class, err error) {
	keyId, err := attributevalue.Marshal(id)
	if err != nil {
		return
	}

	params := &dynamodb.GetItemInput{
		TableName: &r.table,
		Key: map[string]types.AttributeValue{
			"PrimaryKey": keyId,
			"SortKey":    keyId,
		},
	}

	response, err := r.client.GetItem(ctx, params)
	if err != nil {
		return
	}

	// Check for no-item-found condition
	if len(response.Item) == 0 {
		err = ErrNotFound
		return
	}

	err = attributevalue.UnmarshalMap(response.Item, &class)
	return
}

func (r DynamoDBRepository) InsertClass(ctx context.Context, class *Class) (err error) {
	item, err := attributevalue.MarshalMap(class)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		TableName: &r.table,
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, params)
	return
}

func (r DynamoDBRepository) UpdateClass(ctx context.Context, class *Class) (err error) {
	item, err := attributevalue.MarshalMap(class)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		TableName: &r.table,
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, params)
	return
}
