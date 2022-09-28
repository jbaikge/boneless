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

const (
	sortPrefix   = "sort#"
	sortValueLen = 64
)

var _ dynamoDocumentInterface = &dynamoSort{}

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

func (repo *DynamoDBRepository) deleteSortDocuments(ctx context.Context, id string) (err error) {
	var key struct {
		PK string
		SK string
	}

	prefix, err := attributevalue.Marshal(sortPrefix)
	if err != nil {
		return
	}
	idValue, err := attributevalue.Marshal(id)
	if err != nil {
		return
	}
	params := &dynamodb.ScanInput{
		TableName:            &repo.resources.Table,
		ProjectionExpression: aws.String("PK,SK"),
		FilterExpression:     aws.String("begins_with(PK, :prefix) AND DocumentId = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":prefix": prefix,
			":id":     idValue,
		},
	}
	paginator := dynamodb.NewScanPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, item := range response.Items {
			if err = attributevalue.UnmarshalMap(item, &key); err != nil {
				return err
			}
			if err = repo.deleteItem(ctx, key.PK, key.SK); err != nil {
				return err
			}
		}
	}

	return
}

// Gets items using the sort indexes. This is preferred as it is much faster
// than manually sorting after a table scan.
func (repo *DynamoDBRepository) getSortDocuments(ctx context.Context, filter boneless.DocumentFilter) (list []boneless.Document, r boneless.Range, err error) {
	// Class ID and sort field are required to proceed.
	if filter.ClassId == "" || filter.Sort.Field == "" {
		err = ErrBadFilter
		return
	}

	// Fetch class to cross-reference sort field
	class, err := repo.GetClassById(ctx, filter.ClassId)
	if err != nil {
		return
	}

	// Verify sort field is valid
	keys := class.SortFields()
	if i := sort.SearchStrings(keys, filter.Sort.Field); i < len(keys) && keys[i] != filter.Sort.Field {
		err = ErrBadFilter
		return
	}

	// Get pre-marshalled pk out of key
	key, err := repo.marshalKey(dynamoSortIds(filter.ClassId, filter.Sort.Field, "", ""))
	if err != nil {
		return
	}
	params := &dynamodb.QueryInput{
		TableName:              &repo.resources.Table,
		ScanIndexForward:       aws.Bool(filter.Sort.Ascending()),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": key["PK"],
		},
	}

	// Add parent ID filter if necessary
	if filter.ParentId != "" {
		params.ExpressionAttributeValues[":parent_id"], err = attributevalue.Marshal(filter.ParentId)
		if err != nil {
			err = fmt.Errorf("marshal parent_id: %w", err)
			return
		}
		params.FilterExpression = aws.String("ParentId = :parent_id")
	}

	list = make([]boneless.Document, 0, filter.Range.SliceLen())
	seen := 0
	var response *dynamodb.QueryOutput
	paginator := dynamodb.NewQueryPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err = paginator.NextPage(ctx)
		if err != nil {
			err = fmt.Errorf("retrieving next page: %w", err)
			return
		}

		// Annoyingly, need to iterate over the entire query response to get the
		// final size.
		r.Size += len(response.Items)

		for _, item := range response.Items {
			// Break out if there is no reason to process items
			if seen > filter.Range.End {
				break
			}
			// Skip if not within slice range
			if seen < filter.Range.Start {
				continue
			}
			seen++
			dbSort := new(dynamoSort)
			if err = attributevalue.UnmarshalMap(item, dbSort); err != nil {
				err = fmt.Errorf("unmarshal item: %w", err)
				return
			}
			list = append(list, dbSort.ToDocument())
		}
	}

	r.Start = filter.Range.Start
	r.End = filter.Range.Start
	if length := len(list); length > 0 {
		r.End += length - 1
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
	for _, key := range class.SortFields() {
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
