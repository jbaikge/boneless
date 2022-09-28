package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jbaikge/boneless"
)

const documentPrefix = "doc#"

func dynamoDocumentIds(id string, version int) (pk string, sk string) {
	pk = documentPrefix + id
	sk = documentPrefix + fmt.Sprintf("v%06d", version)
	return
}

type dynamoDocument struct {
	PK         string
	SK         string
	ClassId    string
	ParentId   string
	TemplateId string
	Version    int
	Path       string
	Created    time.Time
	Updated    time.Time
	Data       map[string]interface{}
}

func newDynamoDocument(doc *boneless.Document) (dyn *dynamoDocument) {
	pk, sk := dynamoDocumentIds(doc.Id, doc.Version)
	dyn = &dynamoDocument{
		PK:         pk,
		SK:         sk,
		ClassId:    doc.ClassId,
		ParentId:   doc.ParentId,
		TemplateId: doc.TemplateId,
		Version:    doc.Version,
		Path:       doc.Path,
		Created:    doc.Created,
		Updated:    doc.Updated,
		Data:       make(map[string]interface{}),
	}
	for k, v := range doc.Values {
		dyn.Data[k] = v
	}
	return
}

func (dyn *dynamoDocument) ToDocument() (doc boneless.Document) {
	doc = boneless.Document{
		Id:         dyn.PK[len(documentPrefix):],
		ClassId:    dyn.ClassId,
		ParentId:   dyn.ParentId,
		TemplateId: dyn.TemplateId,
		Version:    dyn.Version,
		Path:       dyn.Path,
		Created:    dyn.Created,
		Updated:    dyn.Updated,
		Values:     make(map[string]interface{}),
	}
	for k, v := range dyn.Data {
		doc.Values[k] = v
	}
	return
}

func (repo *DynamoDBRepository) CreateDocument(ctx context.Context, doc *boneless.Document) (err error) {
	if doc.ClassId == "" {
		return fmt.Errorf("class ID required")
	}

	doc.Version = 1
	dbDoc := newDynamoDocument(doc)

	// Insert two copies of the document: v1 and the latest (v0)
	for _, version := range []int{0, 1} {
		_, dbDoc.SK = dynamoDocumentIds(doc.Id, version)
		if err = repo.putItem(ctx, dbDoc); err != nil {
			return fmt.Errorf("put document failed: %w", err)
		}
	}

	if err = repo.putPathDocument(ctx, doc); err != nil {
		return fmt.Errorf("put path document failed: %w", err)
	}

	if err = repo.putSortDocuments(ctx, doc); err != nil {
		return fmt.Errorf("put sort documents failed: %w", err)
	}

	return
}

func (repo *DynamoDBRepository) DeleteDocument(ctx context.Context, id string) (err error) {
	// Fetch document information from current (v0) document
	key, err := repo.marshalKey(dynamoDocumentIds(id, 0))
	if err != nil {
		return fmt.Errorf("marshal key failed: %w", err)
	}
	docParams := &dynamodb.GetItemInput{
		TableName:            &repo.resources.Table,
		Key:                  key,
		ProjectionExpression: aws.String("#pk, #sk, #version, #path"),
		ExpressionAttributeNames: map[string]string{
			"#pk":      "PK",
			"#sk":      "SK",
			"#version": "Version",
			"#path":    "Path",
		},
	}
	docResponse, err := repo.db.GetItem(ctx, docParams)
	if err != nil {
		return fmt.Errorf("get item failed: %w", err)
	}
	if len(docResponse.Item) == 0 {
		return ErrNotExist
	}
	dbDoc := new(dynamoDocument)
	if err = attributevalue.UnmarshalMap(docResponse.Item, dbDoc); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	// Delete all versions of the document
	for version := 0; version <= dbDoc.Version; version++ {
		pk, sk := dynamoDocumentIds(id, version)
		if err = repo.deleteItem(ctx, pk, sk); err != nil {
			return fmt.Errorf("delete %s v%d failed: %w", id, version, err)
		}
	}

	// Delete path item
	if err = repo.deletePathDocument(ctx, dbDoc.Path); err != nil {
		return fmt.Errorf("delete path (%s) failed: %w", dbDoc.Path, err)
	}

	// Delete sort items
	if err = repo.deleteSortDocuments(ctx, id); err != nil {
		return fmt.Errorf("delete sorts failed: %w", err)
	}

	return
}

// Always fetches the latest version (v0)
func (repo *DynamoDBRepository) GetDocumentById(ctx context.Context, id string) (doc boneless.Document, err error) {
	pk, sk := dynamoDocumentIds(id, 0)
	dbDoc := new(dynamoDocument)
	if err = repo.getItem(ctx, pk, sk, dbDoc); err != nil {
		return
	}
	return dbDoc.ToDocument(), nil
}

func (repo *DynamoDBRepository) GetDocumentList(ctx context.Context, filter boneless.DocumentFilter) (docs []boneless.Document, r boneless.Range, err error) {
	return
}

func (repo *DynamoDBRepository) UpdateDocument(ctx context.Context, doc *boneless.Document) (err error) {
	// Fetch the current version of the document in the database
	oldDoc := new(dynamoDocument)
	pk, sk := dynamoDocumentIds(doc.Id, 0)
	if err = repo.getItem(ctx, pk, sk, oldDoc); err != nil {
		return
	}

	// Increment version based on the current version in the database
	doc.Version = oldDoc.Version + 1

	// Push in the new version
	dbDoc := newDynamoDocument(doc)
	if err = repo.putItem(ctx, dbDoc); err != nil {
		return
	}

	// Force data to an empty map, otherwise it is stored as a null - don't want
	// unmarshal issues later.
	data := doc.Values
	if data == nil {
		data = make(map[string]interface{})
	}
	values := map[string]interface{}{
		"ClassId":    doc.ClassId,
		"ParentId":   doc.ParentId,
		"TemplateId": doc.TemplateId,
		"Version":    doc.Version,
		"Path":       doc.Path,
		"Updated":    doc.Updated,
		"Data":       data,
	}

	// Update the "current" (v0) version of the document
	if err = repo.updateItem(ctx, pk, sk, values); err != nil {
		return
	}
	return
}
