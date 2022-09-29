package dynamodb

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jbaikge/boneless"
)

const documentPrefix = "doc#"

var _ dynamoDocumentInterface = &dynamoDocument{}

func dynamoDocumentIds(id string, version int) (pk string, sk string) {
	pk = documentPrefix + id
	sk = documentPrefix + fmt.Sprintf("v%06d", version)
	return
}

type dynamoDocumentInterface interface {
	ToDocument() boneless.Document
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

// Sorters

type dynamoDocumentByCreated []*dynamoDocument

func (arr dynamoDocumentByCreated) Len() int           { return len(arr) }
func (arr dynamoDocumentByCreated) Less(i, j int) bool { return arr[i].Created.Before(arr[j].Created) }
func (arr dynamoDocumentByCreated) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }

type dynamoDocumentByUpdated []*dynamoDocument

func (arr dynamoDocumentByUpdated) Len() int           { return len(arr) }
func (arr dynamoDocumentByUpdated) Less(i, j int) bool { return arr[i].Updated.Before(arr[j].Updated) }
func (arr dynamoDocumentByUpdated) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }

type dynamoDocumentByValue struct {
	Key  string
	Docs []*dynamoDocument
}

func (sorter dynamoDocumentByValue) Len() int { return len(sorter.Docs) }
func (sorter dynamoDocumentByValue) Less(i, j int) bool {
	iVal, iFound := sorter.Docs[i].Data[sorter.Key]
	if !iFound {
		iVal = ""
	}

	jVal, jFound := sorter.Docs[j].Data[sorter.Key]
	if !jFound {
		jVal = ""
	}

	switch v := iVal.(type) {
	case string:
		return v < jVal.(string)
	case int:
		return v < jVal.(int)
	default:
		return fmt.Sprint(iVal) < fmt.Sprint(jVal)
	}
}
func (sorter dynamoDocumentByValue) Swap(i, j int) {
	sorter.Docs[i], sorter.Docs[j] = sorter.Docs[j], sorter.Docs[i]
}

// API Methods

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

func (repo *DynamoDBRepository) GetDocumentList(ctx context.Context, filter boneless.DocumentFilter) (list []boneless.Document, r boneless.Range, err error) {
	list, r, err = repo.getSortDocuments(ctx, filter)

	// Success!
	if err == nil {
		return
	}

	// A bad filter means we didn't have enough valid information to pull
	// sorted documents. Reset the error and process below.
	if err == ErrBadFilter {
		err = nil
	}

	// Something serious happened and we need to let the user know
	if err != nil {
		return
	}

	// Pass through and perform an expensive scan and sort

	key, err := repo.marshalKey(dynamoDocumentIds("", 0))
	if err != nil {
		err = fmt.Errorf("marshal key: %w", err)
		return
	}

	filterExpression := "SK = :sk"
	params := &dynamodb.ScanInput{
		TableName: &repo.resources.Table,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": key["SK"],
		},
	}

	if filter.ParentId != "" {
		filterExpression += " AND ParentId = :parent_id"
		params.ExpressionAttributeValues[":parent_id"], err = attributevalue.Marshal(filter.ParentId)
		if err != nil {
			err = fmt.Errorf("marshal parent_id: %w", err)
			return
		}
	}

	if filter.ClassId != "" {
		filterExpression += " AND ClassId = :class_id"
		params.ExpressionAttributeValues[":class_id"], err = attributevalue.Marshal(filter.ClassId)
		if err != nil {
			err = fmt.Errorf("marshal class_id: %w", err)
			return
		}
	}

	params.FilterExpression = &filterExpression

	// Pull the data out of the database
	var response *dynamodb.ScanOutput
	dbDocs := make([]*dynamoDocument, 0, 64)
	paginator := dynamodb.NewScanPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err = paginator.NextPage(ctx)
		if err != nil {
			err = fmt.Errorf("unable to next page: %w", err)
			return
		}
		tmp := make([]*dynamoDocument, 0, len(response.Items))
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &tmp); err != nil {
			err = fmt.Errorf("unmarshal list of maps: %w", err)
			return
		}
		dbDocs = append(dbDocs, tmp...)
	}

	// Crank up the sorter
	var sorter sort.Interface
	switch filter.Sort.Field {
	case "":
		sorter = sort.Reverse(dynamoDocumentByCreated(dbDocs))
	case "created":
		sorter = dynamoDocumentByCreated(dbDocs)
	case "updated":
		sorter = dynamoDocumentByUpdated(dbDocs)
	default:
		sorter = dynamoDocumentByValue{
			Docs: dbDocs,
			Key:  filter.Sort.Field,
		}
	}

	// Reverse the sorter if explicitly requested or the sort field is blank
	if filter.Sort.Descending() {
		sorter = sort.Reverse(sorter)
	}

	// Sort documents
	sort.Sort(sorter)

	r.Size = len(dbDocs)

	// Pull out the requested slice
	list = make([]boneless.Document, 0, r.SliceLen())
	for i := filter.Range.Start; i < len(dbDocs) && i <= filter.Range.End; i++ {
		list = append(list, dbDocs[i].ToDocument())
	}

	r.Start = filter.Range.Start
	r.End = filter.Range.Start
	if length := len(list); length > 0 {
		r.End += length - 1
	}

	return
}

func (repo *DynamoDBRepository) UpdateDocument(ctx context.Context, doc *boneless.Document) (err error) {
	// Fetch the current version of the document in the database
	oldDoc := new(dynamoDocument)
	pk, sk := dynamoDocumentIds(doc.Id, 0)
	if err = repo.getItem(ctx, pk, sk, oldDoc); err != nil {
		return
	}

	// Check for path conflict before continuing.
	if oldDoc.Path != doc.Path && repo.hasPathDocument(ctx, doc) {
		return fmt.Errorf("document already exists for path (%s)", doc.Path)
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

	// Pull out sort items
	if err = repo.deleteSortDocuments(ctx, doc.Id); err != nil {
		return fmt.Errorf("delete sort documents: %w", err)
	}

	// Replace sort items with new ones
	if err = repo.putSortDocuments(ctx, doc); err != nil {
		return fmt.Errorf("put sort documents: %w", err)
	}

	if oldDoc.Path != doc.Path {
		if err = repo.deletePathDocument(ctx, oldDoc.Path); err != nil {
			return fmt.Errorf("delete path document: %w", err)
		}
		if err = repo.putPathDocument(ctx, doc); err != nil {
			return fmt.Errorf("put path document: %w", err)
		}
	}

	return
}
