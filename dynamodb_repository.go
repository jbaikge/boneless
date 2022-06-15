package gocms

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
	return
}

func (r DynamoDBRepository) GetAllClasses(ctx context.Context) (classes []Class, err error) {
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
			"Id": keyId,
		},
	}

	response, err := r.client.GetItem(ctx, params)
	if err != nil {
		return
	}

	err = attributevalue.Unmarshal(response.Item, &class)
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
