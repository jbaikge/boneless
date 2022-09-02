package gocms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	dynamoClassPrefix     = "class#"
	dynamoClassSortF      = dynamoClassPrefix + "v%04d"
	dynamoDocPrefix       = "doc#"
	dynamoDocSortF        = dynamoDocPrefix + "v%04d"
	dynamoPathPrefix      = "path#"
	dynamoPathSortKey     = "path"
	dynamoSortPartitionF  = "sort#%s#%s"
	dynamoSortSortF       = "%s#%s"
	dynamoSortValueLength = 64
	dynamoTemplatePrefix  = "template#"
	dynamoTemplateSortF   = dynamoTemplatePrefix + "v%04d"
)

var (
	ErrBadRange = errors.New("invalid range")
	ErrNotExist = errors.New("item does not exist")
)

type dynamoItem interface {
	PartitionKey() string
	SortKey() string
	UpdateValues() map[string]interface{}
}

// Class Types

type dynamoClass struct {
	PK       string
	SK       string
	ParentId string
	Name     string
	Created  time.Time
	Updated  time.Time
	Fields   []Field
}

func (dyn *dynamoClass) FromClass(c *Class) {
	dyn.PK = dynamoClassPrefix + c.Id
	dyn.SK = fmt.Sprintf(dynamoClassSortF, 0)
	dyn.ParentId = c.ParentId
	dyn.Name = c.Name
	dyn.Created = c.Created
	dyn.Updated = c.Updated
	dyn.Fields = make([]Field, len(c.Fields))
	copy(dyn.Fields, c.Fields)
}

func (dyn dynamoClass) ToClass() (c Class) {
	c.Id = dyn.PK[len(dynamoClassPrefix):]
	c.ParentId = dyn.ParentId
	c.Name = dyn.Name
	c.Created = dyn.Created
	c.Updated = dyn.Updated
	c.Fields = make([]Field, len(dyn.Fields))
	copy(c.Fields, dyn.Fields)
	return
}

func (dyn dynamoClass) PartitionKey() string {
	return dyn.PK
}

func (dyn dynamoClass) SortKey() string {
	return dyn.SK
}

func (dyn dynamoClass) UpdateValues() map[string]interface{} {
	return map[string]interface{}{
		"ParentId": dyn.ParentId,
		"Name":     dyn.Name,
		"Fields":   dyn.Fields,
		"Updated":  dyn.Updated,
	}
}

type dynamoClassByName []*dynamoClass

func (arr dynamoClassByName) Len() int           { return len(arr) }
func (arr dynamoClassByName) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoClassByName) Less(i, j int) bool { return arr[i].Name < arr[j].Name }

// Document Types

type dynamoDocumentInterface interface {
	ToDocument() Document
	GetCreated() time.Time
	GetUpdated() time.Time
}

// Sort by created time
type dynamoDocumentByCreated []dynamoDocumentInterface

func (arr dynamoDocumentByCreated) Len() int      { return len(arr) }
func (arr dynamoDocumentByCreated) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoDocumentByCreated) Less(i, j int) bool {
	return arr[i].GetCreated().Before(arr[j].GetCreated())
}

// Sort by updated time
type dynamoDocumentByUpdated []dynamoDocumentInterface

func (arr dynamoDocumentByUpdated) Len() int      { return len(arr) }
func (arr dynamoDocumentByUpdated) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoDocumentByUpdated) Less(i, j int) bool {
	return arr[i].GetUpdated().Before(arr[j].GetUpdated())
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
}

func (dyn *dynamoDocument) FromDocument(doc *Document) {
	dyn.SetSK(doc.Version)
	dyn.PK = dynamoDocPrefix + doc.Id
	dyn.ClassId = doc.ClassId
	dyn.ParentId = doc.ParentId
	dyn.TemplateId = doc.TemplateId
	dyn.Version = doc.Version
	dyn.Path = doc.Path
	dyn.Created = doc.Created
	dyn.Updated = doc.Updated
}

func (dyn dynamoDocument) ToDocument() (doc Document) {
	doc.Id = dyn.PK[len(dynamoDocPrefix):]
	doc.ClassId = dyn.ClassId
	doc.ParentId = dyn.ParentId
	doc.TemplateId = dyn.TemplateId
	doc.Version = dyn.Version
	doc.Path = dyn.Path
	doc.Created = dyn.Created
	doc.Updated = dyn.Updated
	return
}

func (dyn *dynamoDocument) SetSK(version int) {
	dyn.SK = fmt.Sprintf(dynamoDocSortF, version)
}

func (dyn dynamoDocument) PartitionKey() string {
	return dyn.PK
}

func (dyn dynamoDocument) SortKey() string {
	return dyn.SK
}

func (dyn dynamoDocument) UpdateValues() map[string]interface{} {
	return map[string]interface{}{
		"ClassId":    dyn.ClassId,
		"ParentId":   dyn.ParentId,
		"TemplateId": dyn.TemplateId,
		"Version":    dyn.Version,
		"Path":       dyn.Path,
		"Updated":    dyn.Updated,
	}
}

func (dyn dynamoDocument) GetCreated() time.Time { return dyn.Created }
func (dyn dynamoDocument) GetUpdated() time.Time { return dyn.Updated }

// Path Type

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
}

func (dyn *dynamoPath) FromDocument(doc *Document) {
	dyn.PK = dynamoPathPrefix + doc.Path
	dyn.SK = dynamoPathSortKey
	dyn.DocumentId = doc.Id
	dyn.ClassId = doc.ClassId
	dyn.ParentId = doc.ParentId
	dyn.TemplateId = doc.TemplateId
	dyn.Version = doc.Version
	dyn.Created = doc.Created
	dyn.Updated = doc.Updated
}

func (dyn dynamoPath) ToDocument() (doc Document) {
	doc.Path = dyn.PK[len(dynamoPathPrefix):]
	doc.Id = dyn.DocumentId
	doc.ClassId = dyn.ClassId
	doc.ParentId = dyn.ParentId
	doc.TemplateId = dyn.TemplateId
	doc.Version = dyn.Version
	doc.Created = dyn.Created
	doc.Updated = dyn.Updated
	return
}

// Sort Type

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
}

func (dyn *dynamoSort) FromDocument(doc *Document, key string) (ok bool) {
	value, ok := doc.Values[key]
	if !ok {
		return
	}
	dyn.PK = fmt.Sprintf(dynamoSortPartitionF, doc.ClassId, key)
	dyn.SK = fmt.Sprintf(dynamoSortSortF, dyn.Truncate(value), doc.Id)
	dyn.DocumentId = doc.Id
	dyn.ClassId = doc.ClassId
	dyn.ParentId = doc.ParentId
	dyn.TemplateId = doc.TemplateId
	dyn.Version = doc.Version
	dyn.Path = doc.Path
	dyn.Created = doc.Created
	dyn.Updated = doc.Updated
	return true
}

func (dyn dynamoSort) ToDocument() (doc Document) {
	doc.Id = dyn.DocumentId
	doc.ClassId = dyn.ClassId
	doc.ParentId = dyn.ParentId
	doc.TemplateId = dyn.TemplateId
	doc.Version = dyn.Version
	doc.Path = dyn.Path
	doc.Created = dyn.Created
	doc.Updated = dyn.Updated
	return
}

func (dyn dynamoSort) Truncate(v interface{}) string {
	if t, ok := v.(time.Time); ok {
		v = t.UTC().Format(time.RFC3339)
	}
	return fmt.Sprintf("%.*s", dynamoSortValueLength, fmt.Sprintf("%v", v))
}

func (dyn dynamoSort) GetCreated() time.Time { return dyn.Created }
func (dyn dynamoSort) GetUpdated() time.Time { return dyn.Updated }

// Template Type

type dynamoTemplate struct {
	PK      string
	SK      string
	Name    string
	Version int
	Created time.Time
	Updated time.Time
}

func (dyn *dynamoTemplate) FromTemplate(template *Template) {
	dyn.SetSK(template.Version)
	dyn.PK = dynamoTemplatePrefix + template.Id
	dyn.Name = template.Name
	dyn.Version = template.Version
	dyn.Created = template.Created
	dyn.Updated = template.Updated
}

func (dyn dynamoTemplate) ToTemplate() (template Template) {
	template.Id = dyn.PK[len(dynamoTemplatePrefix):]
	template.Name = dyn.Name
	template.Version = dyn.Version
	template.Created = dyn.Created
	template.Updated = dyn.Updated
	return
}

func (dyn dynamoTemplate) PartitionKey() string {
	return dyn.PK
}

func (dyn *dynamoTemplate) SetSK(version int) {
	dyn.SK = fmt.Sprintf(dynamoTemplateSortF, version)
}

func (dyn dynamoTemplate) SortKey() string {
	return dyn.SK
}

func (dyn dynamoTemplate) UpdateValues() map[string]interface{} {
	return map[string]interface{}{
		"Name":    dyn.Name,
		"Version": dyn.Version,
		"Updated": dyn.Updated,
	}
}

// Sort by name
type dynamoTemplateByName []dynamoTemplate

func (arr dynamoTemplateByName) Len() int           { return len(arr) }
func (arr dynamoTemplateByName) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoTemplateByName) Less(i, j int) bool { return arr[i].Name < arr[j].Name }

// Sort by created time
type dynamoTemplateByCreated []dynamoTemplate

func (arr dynamoTemplateByCreated) Len() int      { return len(arr) }
func (arr dynamoTemplateByCreated) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoTemplateByCreated) Less(i, j int) bool {
	return arr[i].Created.Before(arr[j].Created)
}

// Repository

type DynamoDBResources struct {
	Bucket string
	Table  string
}

func (res *DynamoDBResources) FromEnv() {
	res.Bucket = os.Getenv("REPOSITORY_BUCKET")
	res.Table = os.Getenv("REPOSITORY_TABLE")
}

type DynamoDBRepository struct {
	db        *dynamodb.Client
	s3        *s3.Client
	resources DynamoDBResources
	stats     RepositoryStats
}

// Ref: https://dynobase.dev/dynamodb-golang-query-examples/
func NewDynamoDBRepository(config aws.Config, resources DynamoDBResources) Repository {
	return &DynamoDBRepository{
		db:        dynamodb.NewFromConfig(config),
		s3:        s3.NewFromConfig(config),
		resources: resources,
	}
}

func (repo DynamoDBRepository) Stats() RepositoryStats {
	return repo.stats
}

// Class Methods

func (repo DynamoDBRepository) CreateClass(ctx context.Context, class *Class) (err error) {
	dbClass := new(dynamoClass)
	dbClass.FromClass(class)
	return repo.putItem(ctx, dbClass)
}

func (repo DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	return repo.deleteItem(ctx, dynamoClassPrefix+id, fmt.Sprintf(dynamoClassSortF, 0))
}

func (repo DynamoDBRepository) GetClassById(ctx context.Context, id string) (class Class, err error) {
	dbClass := new(dynamoClass)
	if err = repo.getItem(ctx, dynamoClassPrefix+id, fmt.Sprintf(dynamoClassSortF, 0), dbClass); err != nil {
		return
	}
	return dbClass.ToClass(), nil
}

func (repo DynamoDBRepository) GetClassList(ctx context.Context, filter ClassFilter) (list []Class, r Range, err error) {
	tmp := make([]*dynamoClass, 0, 128)

	skId, err := attributevalue.Marshal(fmt.Sprintf(dynamoClassSortF, 0))
	if err != nil {
		return
	}
	params := &dynamodb.ScanInput{
		TableName:        &repo.resources.Table,
		FilterExpression: aws.String("SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": skId,
		},
	}
	paginator := dynamodb.NewScanPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return list, r, err
		}

		// TODO goroutine
		dbClasses := make([]*dynamoClass, 0, len(response.Items))
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &dbClasses); err != nil {
			return list, r, err
		}
		tmp = append(tmp, dbClasses...)
	}

	sort.Sort(dynamoClassByName(tmp))

	r.Size = len(tmp)
	list = make([]Class, 0, filter.Range.End-filter.Range.Start+1)
	for i := filter.Range.Start; i < len(tmp) && i <= filter.Range.End; i++ {
		list = append(list, tmp[i].ToClass())
	}

	// If start = 0  and list is empty, then there just aren't any records
	if filter.Range.Start > 0 && len(list) == 0 {
		return list, r, ErrBadRange
	}

	// Kind of a weird situation here where equal start and end actually signify
	// one item, but size can be zero.
	r.Start = filter.Range.Start
	r.End = r.Start
	if len(list) > 0 {
		r.End += len(list) - 1
	}

	return
}

func (repo DynamoDBRepository) UpdateClass(ctx context.Context, class *Class) (err error) {
	dbClass := new(dynamoClass)
	dbClass.FromClass(class)
	return repo.updateItem(ctx, dbClass)
}

// Document Methods

// Document creation inserts two records: one with version zero and one with
// version one
func (repo DynamoDBRepository) CreateDocument(ctx context.Context, doc *Document) (err error) {
	if doc.ClassId == "" {
		return fmt.Errorf("classId is required")
	}

	doc.Version = 1

	dbDoc := new(dynamoDocument)
	dbDoc.FromDocument(doc)

	for _, version := range []int{0, 1} {
		dbDoc.SetSK(version)
		if err = repo.putItem(ctx, dbDoc); err != nil {
			return
		}
	}

	// Possible to have documents with no path as they are child documents
	if doc.Path != "" {
		dbPath := new(dynamoPath)
		dbPath.FromDocument(doc)
		if err = repo.putItem(ctx, dbPath); err != nil {
			return
		}
	}

	// Add sort records
	class, err := repo.GetClassById(ctx, doc.ClassId)
	if err != nil {
		return fmt.Errorf("could not retrieve class (%s): %w", doc.ClassId, err)
	}
	for _, key := range class.SortFields() {
		dbSort := new(dynamoSort)
		if ok := dbSort.FromDocument(doc, key); !ok {
			continue
		}
		if err = repo.putItem(ctx, dbSort); err != nil {
			return
		}
	}

	if err = repo.putValues(ctx, doc); err != nil {
		return
	}

	return
}

func (repo DynamoDBRepository) DeleteDocument(ctx context.Context, id string) (err error) {
	dbDoc := new(dynamoDocument)
	docId := dynamoDocPrefix + id
	if err = repo.getItem(ctx, docId, fmt.Sprintf(dynamoDocSortF, 0), dbDoc); err != nil {
		return
	}

	// Delete all versions of the document
	objects := make([]s3types.ObjectIdentifier, 0, dbDoc.Version)
	doc := dbDoc.ToDocument()
	for i := 0; i <= dbDoc.Version; i++ {
		if err = repo.deleteItem(ctx, docId, fmt.Sprintf(dynamoDocSortF, i)); err != nil {
			return
		}
		if i > 0 {
			doc.Version = i
			objects = append(objects, s3types.ObjectIdentifier{
				Key: aws.String(repo.valuesKey(&doc)),
			})
		}
	}

	// Delete path item
	if err = repo.deleteItem(ctx, dynamoPathPrefix+dbDoc.Path, dynamoPathSortKey); err != nil {
		return
	}

	// Delete sort items
	type PK_SK struct {
		PK string
		SK string
	}

	prefix, err := attributevalue.Marshal("sort#")
	if err != nil {
		return
	}
	docIdValue, err := attributevalue.Marshal(id)
	if err != nil {
		return
	}
	params := &dynamodb.ScanInput{
		TableName:            &repo.resources.Table,
		ProjectionExpression: aws.String("PK,SK"),
		FilterExpression:     aws.String("begins_with(PK, :prefix) AND DocumentId = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":prefix": prefix,
			":id":     docIdValue,
		},
	}
	paginator := dynamodb.NewScanPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, item := range response.Items {
			var pk_sk PK_SK
			if err = attributevalue.UnmarshalMap(item, &pk_sk); err != nil {
				return err
			}
			if err = repo.deleteItem(ctx, pk_sk.PK, pk_sk.SK); err != nil {
				return err
			}
		}
	}

	// Delete all values from S3 - this is at the end of the func to make sure
	// everything is removed from the database first in case the S3 delete fails
	delete := &s3.DeleteObjectsInput{
		Bucket: &repo.resources.Bucket,
		Delete: &s3types.Delete{
			Objects: objects,
		},
	}
	_, err = repo.s3.DeleteObjects(ctx, delete)

	return
}

// Always fetches the latest version (v0) of the document with requested id
func (repo DynamoDBRepository) GetDocumentById(ctx context.Context, id string) (doc Document, err error) {
	dbDoc := new(dynamoDocument)
	if err = repo.getItem(ctx, dynamoDocPrefix+id, fmt.Sprintf(dynamoDocSortF, 0), dbDoc); err != nil {
		return
	}
	doc = dbDoc.ToDocument()
	err = repo.getValues(ctx, &doc)
	return
}

func (repo DynamoDBRepository) GetDocumentByPath(ctx context.Context, path string) (doc Document, err error) {
	dbPath := new(dynamoPath)
	if err = repo.getItem(ctx, dynamoPathPrefix+path, dynamoPathSortKey, dbPath); err != nil {
		return
	}
	doc = dbPath.ToDocument()
	err = repo.getValues(ctx, &doc)
	return
}

func (repo DynamoDBRepository) GetDocumentList(ctx context.Context, filter DocumentFilter) (list []Document, r Range, err error) {
	tmp := make([]dynamoDocumentInterface, 0, 128)

	var sortField string
	if filter.ClassId != "" && filter.Sort.Field != "" {
		class, err := repo.GetClassById(ctx, filter.ClassId)
		if err != nil {
			return list, r, fmt.Errorf("could not retrieve class [%s]: %w", filter.ClassId, err)
		}
		sortable := class.SortFields()
		if i := sort.SearchStrings(sortable, filter.Sort.Field); i < len(sortable) && sortable[i] == filter.Sort.Field {
			sortField = filter.Sort.Field
		}
	}

	// Handle a class-field search, ordered by field values
	switch {
	case sortField != "":
		pk, err := attributevalue.Marshal(fmt.Sprintf(dynamoSortPartitionF, filter.ClassId, sortField))
		if err != nil {
			return list, r, err
		}
		params := &dynamodb.QueryInput{
			TableName:              &repo.resources.Table,
			ScanIndexForward:       aws.Bool(filter.Sort.Ascending()),
			Limit:                  aws.Int32(int32(filter.Range.End + 1)),
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": pk,
			},
		}
		if filter.ParentId != "" {
			parentId, err := attributevalue.Marshal(filter.ParentId)
			if err != nil {
				return list, r, err
			}
			params.FilterExpression = aws.String("ParentId = :parent_id")
			params.ExpressionAttributeValues[":parent_id"] = parentId
		}
		paginator := dynamodb.NewQueryPaginator(repo.db, params)
		for paginator.HasMorePages() {
			response, err := paginator.NextPage(ctx)
			if err != nil {
				return list, r, err
			}
			dbSorts := make([]dynamoSort, 0, len(response.Items))
			if err = attributevalue.UnmarshalListOfMaps(response.Items, &dbSorts); err != nil {
				return list, r, err
			}
			// Cannot use append(tmp, dbSorts...)
			for i := range dbSorts {
				tmp = append(tmp, dbSorts[i])
			}
		}
	default:
		sk, err := attributevalue.Marshal(fmt.Sprintf(dynamoDocSortF, 0))
		if err != nil {
			return list, r, err
		}
		filterExpression := "SK = :sk"
		params := &dynamodb.ScanInput{
			TableName: &repo.resources.Table,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":sk": sk,
			},
		}
		if filter.ParentId != "" {
			parentId, err := attributevalue.Marshal(filter.ParentId)
			if err != nil {
				return list, r, err
			}
			filterExpression += " AND ParentId = :parent_id"
			params.ExpressionAttributeValues[":parent_id"] = parentId
		}
		if filter.ClassId != "" {
			classId, err := attributevalue.Marshal(filter.ClassId)
			if err != nil {
				return list, r, err
			}
			filterExpression += " AND ClassId = :class_id"
			params.ExpressionAttributeValues[":class_id"] = classId
		}
		params.FilterExpression = &filterExpression

		paginator := dynamodb.NewScanPaginator(repo.db, params)
		for paginator.HasMorePages() {
			response, err := paginator.NextPage(ctx)
			if err != nil {
				return list, r, err
			}
			dbDocs := make([]dynamoDocument, 0, len(response.Items))
			if err = attributevalue.UnmarshalListOfMaps(response.Items, &dbDocs); err != nil {
				return list, r, err
			}
			// Cannot use append(tmp, dbDocs...)
			for i := range dbDocs {
				tmp = append(tmp, dbDocs[i])
			}
		}
	}

	// Sort results in memory if the field is a certain option
	var sorter sort.Interface
	switch strings.ToLower(filter.Sort.Field) {
	case "id", "created", "":
		sorter = dynamoDocumentByCreated(tmp)
	case "updated":
		sorter = dynamoDocumentByUpdated(tmp)
	}
	if sorter != nil {
		if filter.Sort.Descending() || filter.Sort.Field == "" {
			sorter = sort.Reverse(sorter)
		}
		sort.Sort(sorter)
	}

	// Pull out range segment
	r.Size = len(tmp)
	list = make([]Document, 0, filter.Range.End-filter.Range.Start+1)
	for i := filter.Range.Start; i < len(tmp) && i <= filter.Range.End; i++ {
		doc := tmp[i].ToDocument()
		if err = repo.getValues(ctx, &doc); err != nil {
			return
		}
		list = append(list, doc)
	}

	// If start = 0  and list is empty, then there just aren't any records
	if filter.Range.Start > 0 && len(list) == 0 {
		return list, r, ErrBadRange
	}

	// Kind of a weird situation here where equal start and end actually signify
	// one item, but size can be zero.
	r.Start = filter.Range.Start
	r.End = r.Start
	if len(list) > 0 {
		r.End += len(list) - 1
	}

	return
}

// Updates are pretty expensive as all the various copies need to be updated as well
func (repo DynamoDBRepository) UpdateDocument(ctx context.Context, doc *Document) (err error) {
	// Fetch the current document in the database
	oldDoc, err := repo.GetDocumentById(ctx, doc.Id)
	if err != nil {
		return
	}

	// Update the version value to the next one
	doc.Version = oldDoc.Version + 1

	// Push in new document version
	dbDoc := new(dynamoDocument)
	dbDoc.FromDocument(doc)
	if err = repo.putItem(ctx, dbDoc); err != nil {
		return
	}
	// Update version zero
	dbDoc.SetSK(0)
	if err = repo.updateItem(ctx, dbDoc); err != nil {
		return
	}

	if err = repo.updatePathDocument(ctx, &oldDoc, doc); err != nil {
		return
	}

	if err = repo.updateSortDocuments(ctx, doc); err != nil {
		return
	}

	if err = repo.putValues(ctx, doc); err != nil {
		return
	}

	return
}

func (repo DynamoDBRepository) updatePathDocument(ctx context.Context, oldDoc *Document, doc *Document) (err error) {
	oldPath := new(dynamoPath)

	// Attempt to find old path using the old document's information
	if oldDoc.Path != "" {
		err = repo.getItem(ctx, dynamoPathPrefix+oldDoc.Path, dynamoPathSortKey, oldPath)
		if err != nil && err != ErrNotExist {
			return
		}
	}

	// Uh-oh need to do a table scan because somehow the path updated on
	// the document without the path partition key getting updated
	if oldPath.PK == "" {
		sk, err := attributevalue.Marshal(dynamoPathSortKey)
		if err != nil {
			return err
		}
		id, err := attributevalue.Marshal(oldDoc.Id)
		if err != nil {
			return err
		}
		params := &dynamodb.ScanInput{
			TableName:        &repo.resources.Table,
			FilterExpression: aws.String("SK = :sk AND DocumentId = :id"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":sk": sk,
				":id": id,
			},
		}
		pagination := dynamodb.NewScanPaginator(repo.db, params)
		for pagination.HasMorePages() {
			response, err := pagination.NextPage(ctx)
			if err != nil {
				return err
			}
			if len(response.Items) == 0 {
				continue
			}
			if err = attributevalue.UnmarshalMap(response.Items[0], oldPath); err != nil {
				return err
			}
		}
	}

	// Remove old path
	if oldPath.PK != "" {
		if err = repo.deleteItem(ctx, oldPath.PK, oldPath.SK); err != nil {
			return
		}
	}

	// Add new path
	if doc.Path != "" {
		newPath := new(dynamoPath)
		newPath.FromDocument(doc)
		if err = repo.putItem(ctx, newPath); err != nil {
			return
		}
	}

	return
}

func (repo DynamoDBRepository) updateSortDocuments(ctx context.Context, doc *Document) (err error) {
	class, err := repo.GetClassById(ctx, doc.ClassId)
	if err != nil {
		return
	}
	for _, key := range class.SortFields() {
		pk, err := attributevalue.Marshal(fmt.Sprintf(dynamoSortPartitionF, class.Id, key))
		if err != nil {
			return err
		}
		query := &dynamodb.QueryInput{
			TableName:              &repo.resources.Table,
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": pk,
			},
			ProjectionExpression: aws.String("PK,SK"),
		}
		paginator := dynamodb.NewQueryPaginator(repo.db, query)
		for paginator.HasMorePages() {
			response, err := paginator.NextPage(ctx)
			if err != nil {
				return err
			}
			for _, item := range response.Items {
				var sk string
				if err = attributevalue.Unmarshal(item["SK"], &sk); err != nil {
					return fmt.Errorf("error unmarshalling sortkey: %w", err)
				}
				if !strings.HasSuffix(sk, "#"+doc.Id) {
					continue
				}
				delete := &dynamodb.DeleteItemInput{
					TableName: &repo.resources.Table,
					Key: map[string]types.AttributeValue{
						"PK": item["PK"],
						"SK": item["SK"],
					},
				}
				if _, err = repo.db.DeleteItem(ctx, delete); err != nil {
					return err
				}
			}
		}

		// Add sort record
		dbSort := new(dynamoSort)
		if ok := dbSort.FromDocument(doc, key); ok {
			if err = repo.putItem(ctx, dbSort); err != nil {
				return err
			}
		}
	}
	return
}

// S3 Document interaction

func (repo DynamoDBRepository) valuesKey(doc *Document) string {
	return fmt.Sprintf("documents/%s/%s/v%04d.json", doc.ClassId, doc.Id, doc.Version)
}

func (repo DynamoDBRepository) getValues(ctx context.Context, doc *Document) (err error) {
	params := &s3.GetObjectInput{
		Bucket: &repo.resources.Bucket,
		Key:    aws.String(repo.valuesKey(doc)),
	}
	response, err := repo.s3.GetObject(ctx, params)
	if err != nil {
		return
	}
	if doc.Values == nil {
		doc.Values = make(map[string]interface{})
	}
	return json.NewDecoder(response.Body).Decode(&doc.Values)
}

func (repo DynamoDBRepository) putValues(ctx context.Context, doc *Document) (err error) {
	buffer := new(bytes.Buffer)
	if err = json.NewEncoder(buffer).Encode(doc.Values); err != nil {
		return
	}

	params := &s3.PutObjectInput{
		Bucket:      &repo.resources.Bucket,
		Key:         aws.String(repo.valuesKey(doc)),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String("application/json"),
	}
	_, err = repo.s3.PutObject(ctx, params)
	return
}

// Template Methods

func (repo DynamoDBRepository) CreateTemplate(ctx context.Context, template *Template) (err error) {
	dbTemplate := new(dynamoTemplate)
	template.Version = 1
	for _, version := range []int{0, 1} {
		dbTemplate.FromTemplate(template)
		dbTemplate.SetSK(version)
		if err = repo.putItem(ctx, dbTemplate); err != nil {
			return
		}
	}
	if err = repo.putTemplateBody(ctx, template); err != nil {
		return
	}
	return
}

func (repo DynamoDBRepository) DeleteTemplate(ctx context.Context, id string) (err error) {
	dbTemplate := new(dynamoTemplate)
	if err = repo.getItem(ctx, dynamoTemplatePrefix+id, fmt.Sprintf(dynamoTemplateSortF, 0), dbTemplate); err != nil {
		return
	}

	// Delete all versions of the template
	objects := make([]s3types.ObjectIdentifier, 0, dbTemplate.Version)
	for i := 0; i <= dbTemplate.Version; i++ {
		dbTemplate.SetSK(i)
		if err = repo.deleteItem(ctx, dbTemplate.PK, dbTemplate.SK); err != nil {
			return
		}
		if i > 0 {
			template := Template{Id: id, Version: i}
			objects = append(objects, s3types.ObjectIdentifier{
				Key: aws.String(repo.templateKey(&template)),
			})
		}
	}

	// Setup S3 delete
	delete := &s3.DeleteObjectsInput{
		Bucket: &repo.resources.Bucket,
		Delete: &s3types.Delete{
			Objects: objects,
		},
	}
	_, err = repo.s3.DeleteObjects(ctx, delete)

	return
}

func (repo DynamoDBRepository) GetTemplateById(ctx context.Context, id string) (template Template, err error) {
	dbTemplate := new(dynamoTemplate)
	if err = repo.getItem(ctx, dynamoTemplatePrefix+id, fmt.Sprintf(dynamoTemplateSortF, 0), dbTemplate); err != nil {
		return
	}
	template = dbTemplate.ToTemplate()
	if err = repo.getTemplateBody(ctx, &template); err != nil {
		return
	}
	return
}

func (repo DynamoDBRepository) GetTemplateList(ctx context.Context, filter TemplateFilter) (list []Template, r Range, err error) {
	tmp := make([]dynamoTemplate, 0, 128)

	// Filter out template-specific SortKeys
	sk, err := attributevalue.Marshal(fmt.Sprintf(dynamoTemplateSortF, 0))
	if err != nil {
		return list, r, err
	}
	params := &dynamodb.ScanInput{
		TableName:        &repo.resources.Table,
		FilterExpression: aws.String("SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": sk,
		},
	}
	paginator := dynamodb.NewScanPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return list, r, err
		}
		dbTemplates := make([]dynamoTemplate, 0, len(response.Items))
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &dbTemplates); err != nil {
			return list, r, err
		}
		tmp = append(tmp, dbTemplates...)
	}

	// Sort by name or created
	var sorter sort.Interface
	switch strings.ToLower(filter.Field) {
	case "created":
		sorter = dynamoTemplateByCreated(tmp)
	default:
		sorter = dynamoTemplateByName(tmp)
	}
	if filter.SortReverse {
		sorter = sort.Reverse(sorter)
	}
	sort.Sort(sorter)

	r.Size = len(tmp)
	list = make([]Template, 0, filter.Range.End-filter.Range.Start+1)
	for i := filter.Range.Start; i < len(tmp) && i <= filter.Range.End; i++ {
		template := tmp[i].ToTemplate()
		if err = repo.getTemplateBody(ctx, &template); err != nil {
			return
		}
		list = append(list, template)
	}

	if filter.Range.Start > 0 && len(list) == 0 {
		return list, r, ErrBadRange
	}

	r.Start = filter.Range.Start
	r.End = r.Start
	if len(list) > 0 {
		r.End += len(list) - 1
	}

	return
}

func (repo DynamoDBRepository) UpdateTemplate(ctx context.Context, template *Template) (err error) {
	// Fetch old template to get latest version number
	oldTemplate, err := repo.GetTemplateById(ctx, template.Id)
	if err != nil {
		return
	}
	template.Version = oldTemplate.Version + 1

	// Add new version
	dbTemplate := new(dynamoTemplate)
	dbTemplate.FromTemplate(template)
	if err = repo.putItem(ctx, dbTemplate); err != nil {
		return
	}

	// Update main version (zero)
	dbTemplate.SetSK(0)
	if err = repo.updateItem(ctx, dbTemplate); err != nil {
		return
	}

	// Update content on S3
	if err = repo.putTemplateBody(ctx, template); err != nil {
		return
	}

	return
}

// S3 Template interaction

func (repo DynamoDBRepository) templateKey(template *Template) string {
	return fmt.Sprintf("templates/%s/v%04d.html", template.Id, template.Version)
}

func (repo DynamoDBRepository) getTemplateBody(ctx context.Context, template *Template) (err error) {
	params := &s3.GetObjectInput{
		Bucket: &repo.resources.Bucket,
		Key:    aws.String(repo.templateKey(template)),
	}
	response, err := repo.s3.GetObject(ctx, params)
	if err != nil {
		return
	}
	defer response.Body.Close()
	buffer := new(bytes.Buffer)
	if _, err = buffer.ReadFrom(response.Body); err != nil {
		return
	}
	template.Body = buffer.String()
	return
}

func (repo DynamoDBRepository) putTemplateBody(ctx context.Context, template *Template) (err error) {
	params := &s3.PutObjectInput{
		Bucket:      &repo.resources.Bucket,
		Key:         aws.String(repo.templateKey(template)),
		Body:        strings.NewReader(template.Body),
		ContentType: aws.String("text/html"),
	}
	_, err = repo.s3.PutObject(ctx, params)
	return
}

// Abstracted API calls to handle generic operations

func (repo DynamoDBRepository) deleteItem(ctx context.Context, pk string, sk string) (err error) {
	pkId, err := attributevalue.Marshal(pk)
	if err != nil {
		return
	}

	skId, err := attributevalue.Marshal(sk)
	if err != nil {
		return
	}

	params := &dynamodb.DeleteItemInput{
		TableName: &repo.resources.Table,
		Key: map[string]types.AttributeValue{
			"PK": pkId,
			"SK": skId,
		},
	}
	_, err = repo.db.DeleteItem(ctx, params)

	repo.stats.Deletes++

	return
}

func (repo DynamoDBRepository) getItem(ctx context.Context, pk string, sk string, dst interface{}) (err error) {
	pkId, err := attributevalue.Marshal(pk)
	if err != nil {
		return
	}

	skId, err := attributevalue.Marshal(sk)
	if err != nil {
		return
	}

	params := &dynamodb.GetItemInput{
		TableName: &repo.resources.Table,
		Key: map[string]types.AttributeValue{
			"PK": pkId,
			"SK": skId,
		},
	}
	response, err := repo.db.GetItem(ctx, params)
	if err != nil {
		return
	}

	if len(response.Item) == 0 {
		return ErrNotExist
	}

	err = attributevalue.UnmarshalMap(response.Item, dst)

	repo.stats.Fetches++

	return
}

func (repo DynamoDBRepository) putItem(ctx context.Context, item interface{}) (err error) {
	inputItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		Item:      inputItem,
		TableName: &repo.resources.Table,
	}
	_, err = repo.db.PutItem(ctx, params)

	repo.stats.Inserts++

	return
}

func (repo DynamoDBRepository) updateItem(ctx context.Context, item dynamoItem) (err error) {
	pkId, err := attributevalue.Marshal(item.PartitionKey())
	if err != nil {
		return
	}

	skId, err := attributevalue.Marshal(item.SortKey())
	if err != nil {
		return
	}

	rawValues := item.UpdateValues()
	sets := make([]string, 0, len(rawValues))
	values := make(map[string]types.AttributeValue)
	names := make(map[string]string)
	for key, value := range rawValues {
		index := len(sets)
		placeholder := ":" + key
		if values[placeholder], err = attributevalue.Marshal(value); err != nil {
			return fmt.Errorf("failed to marshal %s: %w", key, err)
		}
		sets = append(sets, fmt.Sprintf("#param_%d = %s", index, placeholder))
		names[fmt.Sprintf("#param_%d", index)] = key
	}
	updateExpression := "SET " + strings.Join(sets, ", ")

	params := &dynamodb.UpdateItemInput{
		TableName: &repo.resources.Table,
		Key: map[string]types.AttributeValue{
			"PK": pkId,
			"SK": skId,
		},
		UpdateExpression:          &updateExpression,
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
	}

	_, err = repo.db.UpdateItem(ctx, params)

	repo.stats.Updates++

	return
}
