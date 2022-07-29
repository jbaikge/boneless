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

func (dyn dynamoDocument) FromDocument(d *Document) {
	dyn.DocumentId = d.Id
	dyn.ClassId = d.ClassId
	dyn.ParentId = d.ParentId
	dyn.TemplateId = d.TemplateId
	dyn.Title = d.Title
	dyn.Version = d.Version
	dyn.Created = d.Created
	dyn.Updated = d.Updated
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
