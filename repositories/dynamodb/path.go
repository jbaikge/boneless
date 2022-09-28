package dynamodb

import (
	"context"
	"time"

	"github.com/jbaikge/boneless"
)

const pathPrefix = "path#"

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

func newDynamoPath(doc *boneless.Document) (dyn *dynamoPath) {
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

func (dyn dynamoPath) ToDocument() (doc boneless.Document) {
	doc = boneless.Document{
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

func (repo *DynamoDBRepository) putPathDocument(ctx context.Context, doc *boneless.Document) (err error) {
	if doc.Path == "" {
		return
	}

	return repo.putItem(ctx, newDynamoPath(doc))
}
