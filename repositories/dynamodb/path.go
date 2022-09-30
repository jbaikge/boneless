package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jbaikge/boneless/models"
)

const pathPrefix = "path#"

var _ dynamoDocumentInterface = &dynamoPath{}

func dynamoPathIds(path string) (pk string, sk string) {
	pk = pathPrefix + path
	sk = "path"
	return
}

type dynamoPath struct {
	PK         string
	SK         string
	DocumentId string
	ClassId    string
	ParentId   string
	TemplateId string
	Version    int
	Created    time.Time
	Updated    time.Time
	Data       map[string]interface{}
}

func newDynamoPath(doc *models.Document) (dyn *dynamoPath) {
	pk, sk := dynamoPathIds(doc.Path)
	dyn = &dynamoPath{
		PK:         pk,
		SK:         sk,
		DocumentId: doc.Id,
		ClassId:    doc.ClassId,
		ParentId:   doc.ParentId,
		TemplateId: doc.TemplateId,
		Version:    doc.Version,
		Created:    doc.Created,
		Updated:    doc.Updated,
		Data:       make(map[string]interface{}),
	}
	for k, v := range doc.Values {
		dyn.Data[k] = v
	}
	return
}

func (dyn dynamoPath) ToDocument() (doc models.Document) {
	doc = models.Document{
		Id:         dyn.DocumentId,
		Path:       dyn.PK[len(pathPrefix):],
		ClassId:    dyn.ClassId,
		ParentId:   dyn.ParentId,
		TemplateId: dyn.TemplateId,
		Version:    dyn.Version,
		Created:    dyn.Created,
		Updated:    dyn.Updated,
		Values:     make(map[string]interface{}),
	}
	for k, v := range dyn.Data {
		doc.Values[k] = v
	}
	return
}

func (repo *DynamoDBRepository) GetDocumentByPath(ctx context.Context, path string) (doc models.Document, err error) {
	pk, sk := dynamoPathIds(path)
	dbPath := new(dynamoPath)
	if err = repo.getItem(ctx, pk, sk, dbPath); err != nil {
		return
	}
	return dbPath.ToDocument(), nil
}

func (repo *DynamoDBRepository) deletePathDocument(ctx context.Context, path string) (err error) {
	if path == "" {
		return
	}

	pk, sk := dynamoPathIds(path)
	return repo.deleteItem(ctx, pk, sk)
}

func (repo *DynamoDBRepository) hasPathDocument(ctx context.Context, doc *models.Document) (exists bool) {
	pk, sk := dynamoPathIds(doc.Path)
	dbPath := new(dynamoPath)
	return !errors.Is(repo.getItem(ctx, pk, sk, dbPath), ErrNotExist)
}

func (repo *DynamoDBRepository) putPathDocument(ctx context.Context, doc *models.Document) (err error) {
	if doc.Path == "" {
		return
	}

	if repo.hasPathDocument(ctx, doc) {
		return fmt.Errorf("document already exists for path (%s)", doc.Path)
	}

	return repo.putItem(ctx, newDynamoPath(doc))
}
