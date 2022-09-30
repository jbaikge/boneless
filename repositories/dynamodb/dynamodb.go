package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jbaikge/boneless/services"
)

var (
	ErrBadRange  = errors.New("invalid range")
	ErrNotExist  = errors.New("item does not exist")
	ErrBadFilter = errors.New("filter not valid")
)

type DynamoDBResources struct {
	Bucket       string
	StaticBucket string
	StaticDomain string
	Table        string
}

func (res *DynamoDBResources) FromEnv() {
	res.Bucket = os.Getenv("REPOSITORY_BUCKET")
	res.StaticBucket = os.Getenv("STATIC_BUCKET")
	res.StaticDomain = os.Getenv("STATIC_DOMAIN")
	res.Table = os.Getenv("REPOSITORY_TABLE")
}

type DynamoDBRepository struct {
	db        *dynamodb.Client
	s3        *s3.Client
	resources DynamoDBResources
}

func NewRepository(config aws.Config, resources DynamoDBResources) services.Repository {
	return &DynamoDBRepository{
		db:        dynamodb.NewFromConfig(config),
		s3:        s3.NewFromConfig(config),
		resources: resources,
	}
}

func (repo *DynamoDBRepository) marshalKey(pk string, sk string) (key map[string]dynamotypes.AttributeValue, err error) {
	pkId, err := attributevalue.Marshal(pk)
	if err != nil {
		err = fmt.Errorf("marshalling partition key (%s): %w", pk, err)
		return
	}

	skId, err := attributevalue.Marshal(sk)
	if err != nil {
		err = fmt.Errorf("marshalling sort key (%s): %w", sk, err)
		return
	}

	key = map[string]dynamotypes.AttributeValue{
		"PK": pkId,
		"SK": skId,
	}
	return
}

func (repo *DynamoDBRepository) deleteItem(ctx context.Context, pk string, sk string) (err error) {
	key, err := repo.marshalKey(pk, sk)
	if err != nil {
		return
	}

	params := &dynamodb.DeleteItemInput{
		TableName: &repo.resources.Table,
		Key:       key,
	}

	_, err = repo.db.DeleteItem(ctx, params)

	return
}

func (repo *DynamoDBRepository) getItem(ctx context.Context, pk string, sk string, dst interface{}) (err error) {
	key, err := repo.marshalKey(pk, sk)
	if err != nil {
		return
	}

	params := &dynamodb.GetItemInput{
		TableName: &repo.resources.Table,
		Key:       key,
	}
	response, err := repo.db.GetItem(ctx, params)
	if err != nil {
		return fmt.Errorf("repo.db.GetItem: %w", err)
	}

	if len(response.Item) == 0 {
		return ErrNotExist
	}

	err = attributevalue.UnmarshalMap(response.Item, dst)

	return
}

func (repo *DynamoDBRepository) putItem(ctx context.Context, item interface{}) (err error) {
	inputItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		Item:      inputItem,
		TableName: &repo.resources.Table,
	}
	_, err = repo.db.PutItem(ctx, params)

	return
}

func (repo *DynamoDBRepository) updateItem(ctx context.Context, pk string, sk string, rawValues map[string]interface{}) (err error) {
	key, err := repo.marshalKey(pk, sk)
	if err != nil {
		return
	}

	sets := make([]string, 0, len(rawValues))
	values := make(map[string]dynamotypes.AttributeValue)
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
		TableName:                 &repo.resources.Table,
		Key:                       key,
		UpdateExpression:          &updateExpression,
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
	}

	_, err = repo.db.UpdateItem(ctx, params)

	return
}
