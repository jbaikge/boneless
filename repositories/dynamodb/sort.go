package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/jbaikge/boneless"
)

const (
	sortPrefix   = "sort#"
	sortValueLen = 64
)

func dynamoSortIds(classId string, key string, docId string, value interface{}) (pk string, sk string) {
	pk = sortPrefix + classId + "#" + key
	if t, ok := value.(time.Time); ok {
		value = t.UTC().Format(time.RFC3339)
	}
	sk = fmt.Sprintf("%.*s#%s", sortValueLen, fmt.Sprintf("%v", value), docId)
	return
}

type dynamoSort struct {
	PK         string
	SK         string
	DocumentId string
	ClassId    string
	ParentId   string
	TemplateId string
	Version    int
	Path       string
	Created    time.Time
	Updated    time.Time
	Data       map[string]interface{}
}

// func newDynamoSort(doc *boneless.Document, key string) (dyn *dynamoSort, ok bool) {
// 	value, ok := doc.Values[key]
// 	if !ok {
// 		return
// 	}
// 	dyn = newDynamoSortBase(doc)
// 	dyn.PK, dyn.SK = dynamoSortIds(doc.ClassId, key, doc.Id, value)
// 	return
// }

func newDynamoSortBase(doc *boneless.Document) (dyn *dynamoSort) {
	dyn = &dynamoSort{
		DocumentId: doc.Id,
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

func (dyn dynamoSort) ToDocument() (doc boneless.Document) {
	doc = boneless.Document{
		Id:         dyn.DocumentId,
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

func (repo *DynamoDBRepository) putSortDocuments(ctx context.Context, doc *boneless.Document) (err error) {
	if doc.ClassId == "" {
		return fmt.Errorf("no class ID")
	}

	class, err := repo.GetClassById(ctx, doc.ClassId)
	if err != nil {
		return fmt.Errorf("get class failed: %w", err)
	}

	dbSort := newDynamoSortBase(doc)
	for _, field := range class.Fields {
		if !field.Sort {
			continue
		}

		key := field.Name
		value, ok := doc.Values[key]
		if !ok {
			continue
		}

		dbSort.PK, dbSort.SK = dynamoSortIds(class.Id, key, doc.Id, value)
		if err = repo.putItem(ctx, dbSort); err != nil {
			return fmt.Errorf("put sort document failed: %w", err)
		}
	}

	return
}
