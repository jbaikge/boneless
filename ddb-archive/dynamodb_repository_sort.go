package gocms

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoSort struct {
	ClassField string
	Value      string
	DocumentId string
}

func (repo DynamoDBRepository) newSortRecord(docId, classId, field, value string) (record dynamoSort) {
	record.ClassField = fmt.Sprintf("%s#%s", classId, field)
	record.Value = fmt.Sprintf("%s#%s", value, docId)
	record.DocumentId = docId
	return
}

func (repo DynamoDBRepository) putSortDocument(ctx context.Context, doc *Document) (err error) {
	table := &repo.resources.Tables.Sort
	class, err := repo.GetClassById(ctx, doc.ClassId)
	if err != nil {
		return fmt.Errorf("trouble getting related class: %w", err)
	}

	// Scan table and remove old values
	queryParams := &dynamodb.QueryInput{
		TableName:              table,
		IndexName:              aws.String("GSI-Document"),
		KeyConditionExpression: aws.String("DocumentId = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: doc.Id},
		},
		ProjectionExpression: aws.String("ClassField,Value"),
	}
	response, err := repo.client.Query(ctx, queryParams)
	if err != nil {
		return fmt.Errorf("error retrieving sort records: %w", err)
	}

	for _, item := range response.Items {
		deleteParams := &dynamodb.DeleteItemInput{
			TableName: table,
			Key:       item,
		}
		if _, err = repo.client.DeleteItem(ctx, deleteParams); err != nil {
			return fmt.Errorf("error deleting record with primary key %v: %w", item, err)
		}
	}

	for _, field := range class.SortFields() {
		value, ok := doc.Values[field]
		if !ok {
			continue
		}
		strValue := fmt.Sprintf("%v", value)

		sortRecord := repo.newSortRecord(doc.Id, doc.ClassId, field, strValue)
		item, err := attributevalue.MarshalMap(sortRecord)
		if err != nil {
			return fmt.Errorf("error marshalling sort record: %w", err)
		}

		params := &dynamodb.PutItemInput{
			TableName: table,
			Item:      item,
		}
		if _, err = repo.client.PutItem(ctx, params); err != nil {
			return fmt.Errorf("error during PutItem: %w", err)
		}
	}
	return
}
