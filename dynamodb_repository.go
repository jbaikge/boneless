package gocms

import (
	"context"
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
)

const (
	dynamoClassPrefix = "class#"
	dynamoClassSortF  = dynamoClassPrefix + "v%04d"
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

type dynamoClass struct {
	PK          string
	SK          string
	Name        string
	TableFields string
	TableLabels string
	Created     time.Time
	Updated     time.Time
	Fields      []Field
}

func (dyn *dynamoClass) FromClass(c *Class) {
	dyn.PK = dynamoClassPrefix + c.Id
	dyn.SK = fmt.Sprintf(dynamoClassSortF, 0)
	dyn.Name = c.Name
	dyn.TableFields = c.TableFields
	dyn.TableLabels = c.TableLabels
	dyn.Created = c.Created
	dyn.Updated = c.Updated
	dyn.Fields = make([]Field, len(c.Fields))
	copy(dyn.Fields, c.Fields)
}

func (dyn dynamoClass) ToClass() (c Class) {
	c.Id = dyn.PK[len(dynamoClassPrefix):]
	c.Name = dyn.Name
	c.TableFields = dyn.TableFields
	c.TableLabels = dyn.TableLabels
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
		"Name":        dyn.Name,
		"TableFields": dyn.TableFields,
		"TableLabels": dyn.TableLabels,
		"Fields":      dyn.Fields,
		"Updated":     dyn.Updated,
	}
}

type dynamoClassByName []*dynamoClass

func (arr dynamoClassByName) Len() int           { return len(arr) }
func (arr dynamoClassByName) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoClassByName) Less(i, j int) bool { return arr[i].Name < arr[j].Name }

type DynamoDBResources struct {
	Bucket string
	Table  string
}

func (res *DynamoDBResources) FromEnv() {
	res.Bucket = os.Getenv("REPOSITORY_BUCKET")
	res.Table = os.Getenv("REPOSITORY_TABLE")
}

type DynamoDBRepository struct {
	client    *dynamodb.Client
	resources DynamoDBResources
}

// Ref: https://dynobase.dev/dynamodb-golang-query-examples/
func NewDynamoDBRepository(config aws.Config, resources DynamoDBResources) Repository {
	return &DynamoDBRepository{
		client:    dynamodb.NewFromConfig(config),
		resources: resources,
	}
}

// Class Methods
func (repo DynamoDBRepository) CreateClass(ctx context.Context, class *Class) (err error) {
	dc := new(dynamoClass)
	dc.FromClass(class)
	return repo.putItem(ctx, dc)
}

func (repo DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	return repo.deleteItem(ctx, dynamoClassPrefix+id, fmt.Sprintf(dynamoClassSortF, 0))
}

func (repo DynamoDBRepository) GetClassById(ctx context.Context, id string) (class Class, err error) {
	dc := new(dynamoClass)
	if err = repo.getItem(ctx, dynamoClassPrefix+id, fmt.Sprintf(dynamoClassSortF, 0), dc); err != nil {
		return
	}
	return dc.ToClass(), nil
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
	paginator := dynamodb.NewScanPaginator(repo.client, params)
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(ctx)
		if err != nil {
			return list, r, err
		}

		// TODO goroutine
		dcs := make([]*dynamoClass, 0, len(response.Items))
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &dcs); err != nil {
			return list, r, err
		}
		tmp = append(tmp, dcs...)
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
	dc := new(dynamoClass)
	dc.FromClass(class)
	return repo.updateItem(ctx, dc)
}

// Document Methods
func (repo DynamoDBRepository) CreateDocument(ctx context.Context, doc *Document) (err error) {
	return
}

func (repo DynamoDBRepository) DeleteDocument(ctx context.Context, id string) (err error) {
	return
}

func (repo DynamoDBRepository) GetDocumentById(ctx context.Context, id string) (doc Document, err error) {
	return
}

func (repo DynamoDBRepository) GetDocumentList(ctx context.Context, filter DocumentFilter) (list []Document, r Range, err error) {
	return
}

func (repo DynamoDBRepository) UpdateDocument(ctx context.Context, doc *Document) (err error) {
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
	_, err = repo.client.DeleteItem(ctx, params)

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
	response, err := repo.client.GetItem(ctx, params)

	if len(response.Item) == 0 {
		return ErrNotExist
	}

	err = attributevalue.UnmarshalMap(response.Item, dst)

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
	_, err = repo.client.PutItem(ctx, params)

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

	_, err = repo.client.UpdateItem(ctx, params)

	return
}
