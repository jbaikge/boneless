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
	return
}

func (repo *DynamoDBRepository) GetDocumentById(ctx context.Context, id string) (doc boneless.Document, err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentByPath(ctx context.Context, path string) (doc boneless.Document, err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentList(ctx context.Context, filter boneless.DocumentFilter) (docs []boneless.Document, r boneless.Range, err error) {
	return
}

func (repo *DynamoDBRepository) UpdateDocument(ctx context.Context, doc *boneless.Document) (err error) {
	return
}
