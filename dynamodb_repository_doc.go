package gocms

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
		dyn.ParentId = xid.NilID().String()
	}
}

func (dyn dynamoDocument) ToDocument() (d Document) {
	d.Id = dyn.DocumentId
	d.ClassId = dyn.ClassId
	d.ParentId = dyn.ParentId
	d.TemplateId = dyn.TemplateId
	d.Title = dyn.Title
	d.Version = dyn.Version
	d.Created = dyn.Created
	d.Updated = dyn.Updated

	if id, err := xid.FromString(d.ParentId); err != nil || id == xid.NilID() {
		d.ParentId = ""
	}
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
	return
}

func (repo DynamoDBRepository) GetDocumentById(ctx context.Context, id string) (doc Document, err error) {
	return
}

func (repo DynamoDBRepository) GetDocumentList(ctx context.Context, filter DocumentFilter) (docs []Document, r Range, err error) {
	return
}

func (repo DynamoDBRepository) UpdateDocument(ctx context.Context, doc *Document) (err error) {
	return
}
