package dynamodb

import (
	"context"
	"fmt"
	"time"

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

func (repo *DynamoDBRepository) CreateDocument(context.Context, *boneless.Document) (err error) {
	return
}

func (repo *DynamoDBRepository) DeleteDocument(context.Context, string) (err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentById(context.Context, string) (doc boneless.Document, err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentByPath(context.Context, string) (doc boneless.Document, err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentList(context.Context, boneless.DocumentFilter) (docs []boneless.Document, r boneless.Range, err error) {
	return
}

func (repo *DynamoDBRepository) UpdateDocument(context.Context, *boneless.Document) (err error) {
	return
}
