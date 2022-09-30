package dynamodb

import (
	"context"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jbaikge/boneless/models"
)

const formPrefix = "form#"

func dynamoFormIds(id string) (pk string, sk string) {
	pk = formPrefix + id
	sk = "form"
	return
}

type dynamoForm struct {
	PK      string
	SK      string
	Name    string
	Created time.Time
	Updated time.Time
	Data    interface{}
}

func newDynamoForm(form *models.Form) (dyn *dynamoForm) {
	pk, sk := dynamoFormIds(form.Id)
	dyn = &dynamoForm{
		PK:      pk,
		SK:      sk,
		Name:    form.Name,
		Created: form.Created,
		Updated: form.Updated,
		Data:    form.Schema,
	}
	return
}

func (dyn *dynamoForm) ToForm() (form models.Form) {
	form = models.Form{
		Id:      dyn.PK[len(formPrefix):],
		Name:    dyn.Name,
		Created: dyn.Created,
		Updated: dyn.Updated,
		Schema:  dyn.Data,
	}
	return
}

type dynamoFormByName []*dynamoForm

func (arr dynamoFormByName) Len() int           { return len(arr) }
func (arr dynamoFormByName) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoFormByName) Less(i, j int) bool { return arr[i].Name < arr[j].Name }

func (repo *DynamoDBRepository) CreateForm(ctx context.Context, form *models.Form) (err error) {
	return repo.putItem(ctx, newDynamoForm(form))
}

func (repo *DynamoDBRepository) DeleteForm(ctx context.Context, id string) (err error) {
	pk, sk := dynamoFormIds(id)
	return repo.deleteItem(ctx, pk, sk)
}

func (repo *DynamoDBRepository) GetFormById(ctx context.Context, id string) (form models.Form, err error) {
	pk, sk := dynamoFormIds(id)
	dbForm := new(dynamoForm)
	if err = repo.getItem(ctx, pk, sk, dbForm); err != nil {
		return
	}
	return dbForm.ToForm(), nil
}

func (repo *DynamoDBRepository) GetFormList(ctx context.Context, filter models.FormFilter) (list []models.Form, r models.Range, err error) {
	var response *dynamodb.ScanOutput
	dbForms := make([]*dynamoForm, 0, 64)

	key, err := repo.marshalKey(dynamoFormIds(""))
	if err != nil {
		return
	}

	params := &dynamodb.ScanInput{
		TableName:        &repo.resources.Table,
		FilterExpression: aws.String("SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": key["SK"],
		},
	}
	paginator := dynamodb.NewScanPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err = paginator.NextPage(ctx)
		if err != nil {
			return
		}
		tmp := make([]*dynamoForm, 0, len(response.Items))
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &tmp); err != nil {
			return
		}
		dbForms = append(dbForms, tmp...)
	}

	sort.Sort(dynamoFormByName(dbForms))

	r.Size = len(dbForms)
	list = make([]models.Form, 0, filter.Range.SliceLen())
	for i := filter.Range.Start; i < len(dbForms) && i <= filter.Range.End; i++ {
		list = append(list, dbForms[i].ToForm())
	}

	if filter.Range.Start > 0 && len(list) == 0 {
		err = ErrBadRange
		return
	}

	r.Start = filter.Range.Start
	r.End = filter.Range.Start
	if length := len(list); length > 0 {
		r.End += length - 1
	}

	return
}

func (repo *DynamoDBRepository) UpdateForm(ctx context.Context, form *models.Form) (err error) {
	pk, sk := dynamoFormIds(form.Id)
	values := map[string]interface{}{
		"Name":    form.Name,
		"Updated": form.Updated,
		"Data":    form.Schema,
	}
	return repo.updateItem(ctx, pk, sk, values)
}
