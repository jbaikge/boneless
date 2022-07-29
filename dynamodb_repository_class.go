package gocms

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jbaikge/gocms/internal/slicer"
)

const classIdPrefix = "class#"

type dynamoClass struct {
	ClassId     string
	Slug        string
	Name        string
	TableLabels string
	TableFields string
	Created     time.Time
	Updated     time.Time
	Fields      []Field
}

func (d *dynamoClass) FromClass(c *Class) {
	d.ClassId = classIdPrefix + c.Id
	d.Name = c.Name
	d.TableLabels = c.TableLabels
	d.TableFields = c.TableFields
	d.Created = c.Created
	d.Updated = c.Updated
	d.Fields = make([]Field, len(c.Fields))
	copy(d.Fields, c.Fields)
}

func (d *dynamoClass) ToClass() (c Class) {
	c.Id = d.ClassId[len(classIdPrefix):]
	c.Name = d.Name
	c.TableLabels = d.TableLabels
	c.TableFields = d.TableFields
	c.Created = d.Created
	c.Updated = d.Updated
	c.Fields = make([]Field, len(d.Fields))
	copy(c.Fields, d.Fields)
	return
}

func (repo DynamoDBRepository) CreateClass(ctx context.Context, class *Class) (err error) {
	dbClass := new(dynamoClass)
	dbClass.FromClass(class)

	item, err := attributevalue.MarshalMap(dbClass)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		TableName: &repo.tables.Class,
		Item:      item,
	}

	_, err = repo.client.PutItem(ctx, params)
	return
}

func (repo DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	prefixedId := classIdPrefix + id
	keyId, err := attributevalue.Marshal(prefixedId)
	if err != nil {
		return
	}

	params := &dynamodb.DeleteItemInput{
		TableName: &repo.tables.Class,
		Key: map[string]types.AttributeValue{
			"ClassId": keyId,
		},
	}

	_, err = repo.client.DeleteItem(ctx, params)
	return
}

func (repo DynamoDBRepository) GetClassById(ctx context.Context, id string) (class Class, err error) {
	prefixedId := classIdPrefix + id
	keyId, err := attributevalue.Marshal(prefixedId)
	if err != nil {
		return
	}

	params := &dynamodb.GetItemInput{
		TableName: &repo.tables.Class,
		Key: map[string]types.AttributeValue{
			"ClassId": keyId,
		},
	}

	response, err := repo.client.GetItem(ctx, params)
	if err != nil {
		return
	}

	// Check for no-item-found condition
	if len(response.Item) == 0 {
		err = ErrNotFound
		return
	}

	dbClass := new(dynamoClass)
	if err = attributevalue.UnmarshalMap(response.Item, dbClass); err != nil {
		return
	}
	class = dbClass.ToClass()

	return
}

func (repo DynamoDBRepository) GetClassList(ctx context.Context, filter ClassFilter) (classes []Class, r Range, err error) {
	classes = make([]Class, 0, filter.Range.End-filter.Range.Start)
	slicer := slicer.NewSlicer(filter.Range.Start, filter.Range.End)

	// Ref: https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go
	// filterPK := expression.Name("ClassId").BeginsWith(classIdPrefix)
	// expr, err := expression.NewBuilder().WithFilter(filterPK).Build()
	// if err != nil {
	// 	return
	// }

	params := &dynamodb.ScanInput{
		TableName: &repo.tables.Class,
		// ExpressionAttributeNames:  expr.Names(),
		// ExpressionAttributeValues: expr.Values(),
		// FilterExpression:          expr.Filter(),
	}
	paginator := dynamodb.NewScanPaginator(repo.client, params)
	for paginator.HasMorePages() {
		page, pageErr := paginator.NextPage(ctx)
		if pageErr != nil {
			err = pageErr
			return
		}

		slicer.Add(int(page.Count))
		start := slicer.Start()
		end := slicer.End()
		if start == 0 && end == 0 {
			continue
		}

		dbClasses := make([]dynamoClass, 0, end-start)
		if err = attributevalue.UnmarshalListOfMaps(page.Items[start:end], &dbClasses); err != nil {
			return
		}

		for _, dbClass := range dbClasses {
			classes = append(classes, dbClass.ToClass())
		}
	}

	r.Size = slicer.Total()

	// No data returned, just return an empty slice
	if r.Size == 0 {
		return
	}

	if filter.Range.Start >= r.Size {
		err = ErrBadRange
		return
	}

	r.Start = filter.Range.Start
	r.End = r.Start + len(classes) - 1

	return
}

func (repo DynamoDBRepository) UpdateClass(ctx context.Context, class *Class) (err error) {
	dbClass := new(dynamoClass)
	dbClass.FromClass(class)
	item, err := attributevalue.MarshalMap(dbClass)
	if err != nil {
		return
	}

	params := &dynamodb.PutItemInput{
		TableName: &repo.tables.Class,
		Item:      item,
	}

	_, err = repo.client.PutItem(ctx, params)
	return
}
