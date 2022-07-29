package gocms

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/xid"
)

type dynamoDocument struct {
	DocumentId string
	ClassId    string
	ParentId   string
	TemplateId string
	VersionId  string
	Title      string
	Version    int
	Created    time.Time
	Updated    time.Time
}

func (dyn *dynamoDocument) FromDocument(d *Document) {
	dyn.DocumentId = d.Id
	dyn.ClassId = d.ClassId
	dyn.ParentId = d.ParentId
	dyn.TemplateId = d.TemplateId
	dyn.Title = d.Title
	dyn.Version = d.Version
	dyn.Created = d.Created
	dyn.Updated = d.Updated

	if dyn.ParentId == "" {
		dyn.ParentId = emptyParentId
	}
}

func (dyn dynamoDocument) ToDocument() (d Document) {
	d.Id = dyn.DocumentId
	d.ClassId = dyn.ClassId
	if dyn.ParentId != emptyParentId {
		d.ParentId = dyn.ParentId
	}
	d.TemplateId = dyn.TemplateId
	d.Title = dyn.Title
	d.Version = dyn.Version
	d.Created = dyn.Created
	d.Updated = dyn.Updated
	return
}

// Initial document inserts 2 records: one with version zero and one with
// version one. Both will have the same VersionId
func (repo DynamoDBRepository) CreateDocument(ctx context.Context, doc *Document) (err error) {
	dbDoc := new(dynamoDocument)
	dbDoc.FromDocument(doc)

	dbDoc.VersionId = xid.New().String()
	for _, version := range []int{0, 1} {
		dbDoc.Version = version
		item, err := attributevalue.MarshalMap(dbDoc)
		if err != nil {
			return err
		}

		params := &dynamodb.PutItemInput{
			TableName: &repo.resources.Tables.Document,
			Item:      item,
		}

		if _, err = repo.client.PutItem(ctx, params); err != nil {
			return err
		}
	}

	return
}

func (repo DynamoDBRepository) DeleteDocument(ctx context.Context, id string) (err error) {
	version, err := repo.getDocumentVersion(ctx, id)
	if err != nil {
		return
	}

	keyId, err := attributevalue.Marshal(id)
	if err != nil {
		return
	}

	for i := 0; i <= version; i++ {
		keyVersion, err := attributevalue.Marshal(i)
		if err != nil {
			return err
		}

		params := &dynamodb.DeleteItemInput{
			TableName: &repo.resources.Tables.Document,
			Key: map[string]types.AttributeValue{
				"DocumentId": keyId,
				"Version":    keyVersion,
			},
		}
		if _, err := repo.client.DeleteItem(ctx, params); err != nil {
			return err
		}
	}
	return
}

func (repo DynamoDBRepository) GetDocumentById(ctx context.Context, id string) (doc Document, err error) {
	keyId, err := attributevalue.Marshal(id)
	if err != nil {
		return
	}

	version, _ := attributevalue.Marshal(0)

	params := &dynamodb.GetItemInput{
		TableName: &repo.resources.Tables.Document,
		Key: map[string]types.AttributeValue{
			"DocumentId": keyId,
			"Version":    version,
		},
	}

	response, err := repo.client.GetItem(ctx, params)
	if err != nil {
		return doc, fmt.Errorf("bad response from GetItem: %w", err)
	}

	// Check for no-item-found condition
	if len(response.Item) == 0 {
		return doc, ErrNotFound
	}

	dbDoc := new(dynamoDocument)
	if err = attributevalue.UnmarshalMap(response.Item, dbDoc); err != nil {
		return doc, fmt.Errorf("unmarshal error: %w", err)
	}
	doc = dbDoc.ToDocument()

	return
}

func (repo DynamoDBRepository) GetDocumentList(ctx context.Context, filter DocumentFilter) (docs []Document, r Range, err error) {
	return
}

func (repo DynamoDBRepository) UpdateDocument(ctx context.Context, doc *Document) (err error) {
	dbDoc := new(dynamoDocument)
	dbDoc.FromDocument(doc)

	version, err := repo.getDocumentVersion(ctx, doc.Id)
	if err != nil {
		return fmt.Errorf("could not determine next version: %w", err)
	}

	dbDoc.VersionId = xid.New().String()

	for _, version := range []int{0, version + 1} {
		dbDoc.Version = version

		item, err := attributevalue.MarshalMap(dbDoc)
		if err != nil {
			return fmt.Errorf("could not marshal doc: %w", err)
		}

		params := &dynamodb.PutItemInput{
			TableName: &repo.resources.Tables.Document,
			Item:      item,
		}

		if _, err = repo.client.PutItem(ctx, params); err != nil {
			return fmt.Errorf("could not put doc with version %d: %w", version, err)
		}
	}

	return
}

func (repo DynamoDBRepository) getDocumentVersion(ctx context.Context, id string) (next int, err error) {
	params := &dynamodb.QueryInput{
		TableName:              &repo.resources.Tables.Document,
		Limit:                  aws.Int32(1),
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String("DocumentId = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: id},
		},
		ProjectionExpression: aws.String("Version"),
	}

	result, err := repo.client.Query(ctx, params)
	if err != nil {
		return
	}

	// Oh, this would be bad
	if len(result.Items) == 0 {
		err = fmt.Errorf("no documents found for id: %s", id)
		return
	}

	var row struct{ Version int }
	attributevalue.UnmarshalMap(result.Items[0], &row)

	// Bump up one since we want the next version
	next = row.Version
	return
}
