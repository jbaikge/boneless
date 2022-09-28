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

const classPrefix = "class#"

func dynamoClassIds(id string) (pk string, sk string) {
	pk = classPrefix + id
	sk = "class"
	return
}

type dynamoClass struct {
	PK       string
	SK       string
	ParentId string
	Name     string
	Created  time.Time
	Updated  time.Time
	Data     []boneless.Field
}

func newDynamoClass(c *boneless.Class) (dyn *dynamoClass) {
	pk, sk := dynamoClassIds(c.Id)
	dyn = &dynamoClass{
		PK:       pk,
		SK:       sk,
		ParentId: c.ParentId,
		Name:     c.Name,
		Created:  c.Created,
		Updated:  c.Updated,
		Data:     make([]boneless.Field, len(c.Fields)),
	}
	copy(dyn.Data, c.Fields)
	return
}

func (dyn *dynamoClass) ToClass() (c boneless.Class) {
	c = boneless.Class{
		Id:       dyn.PK[len(classPrefix):],
		ParentId: dyn.ParentId,
		Name:     dyn.Name,
		Created:  dyn.Created,
		Updated:  dyn.Updated,
		Fields:   make([]boneless.Field, len(dyn.Data)),
	}
	copy(c.Fields, dyn.Data)
	return
}

type dynamoClassByName []*dynamoClass

func (arr dynamoClassByName) Len() int           { return len(arr) }
func (arr dynamoClassByName) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoClassByName) Less(i, j int) bool { return arr[i].Name < arr[j].Name }

func (repo *DynamoDBRepository) CreateClass(ctx context.Context, class *boneless.Class) (err error) {
	return repo.putItem(ctx, newDynamoClass(class))
}

func (repo *DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	pk, sk := dynamoClassIds(id)
	return repo.deleteItem(ctx, pk, sk)
}

func (repo *DynamoDBRepository) GetClassById(ctx context.Context, id string) (class boneless.Class, err error) {
	pk, sk := dynamoClassIds(id)
	dbClass := new(dynamoClass)
	if err = repo.getItem(ctx, pk, sk, dbClass); err != nil {
		err = fmt.Errorf("getItem failed: %w", err)
		return
	}
	return dbClass.ToClass(), nil
}

func (repo *DynamoDBRepository) GetClassList(ctx context.Context, filter boneless.ClassFilter) (list []boneless.Class, r boneless.Range, err error) {
	dbClasses := make([]*dynamoClass, 0, 16)

	_, sk := dynamoClassIds("")
	skId, err := attributevalue.Marshal(sk)
	if err != nil {
		err = fmt.Errorf("marshalling sort key (%s): %w", sk, err)
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
		var response *dynamodb.ScanOutput
		response, err = paginator.NextPage(ctx)
		if err != nil {
			err = fmt.Errorf("paginator next page failed: %w", err)
			return
		}

		tmp := make([]*dynamoClass, 0, len(response.Items))
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &tmp); err != nil {
			err = fmt.Errorf("unmarshal failed: %w", err)
		}

		dbClasses = append(dbClasses, tmp...)
	}

	sort.Sort(dynamoClassByName(dbClasses))

	r.Size = len(dbClasses)

	// Convert dynamo classes to boneless classes, but only ones within range
	list = make([]boneless.Class, 0, filter.Range.End-filter.Range.Start+1)
	for i := filter.Range.Start; i < len(dbClasses) && i <= filter.Range.End; i++ {
		list = append(list, dbClasses[i].ToClass())
	}

	// If start = 0 and list is empty, there just aren't any records
	if filter.Range.Start > 0 && len(list) == 0 {
		err = ErrBadRange
		return
	}

	// Set new range bounds for pagination purposes
	r.Start = filter.Range.Start
	r.End = r.Start
	if length := len(list); length > 0 {
		r.End += length - 1
	}

	return
}

func (repo *DynamoDBRepository) UpdateClass(ctx context.Context, class *boneless.Class) (err error) {
	pk, sk := dynamoClassIds(class.Id)
	values := map[string]interface{}{
		"ParentId": class.ParentId,
		"Name":     class.Name,
		"Data":     class.Fields,
		"Updated":  class.Updated,
	}
	return repo.updateItem(ctx, pk, sk, values)
}
