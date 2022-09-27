package dynamodb

import (
	"context"

	"github.com/jbaikge/boneless"
)

const formPrefix = "form#"

func dynamoFormIds(id string) (pk string, sk string) {
	pk = formPrefix + id
	return
}

func (repo *DynamoDBRepository) CreateForm(context.Context, *boneless.Form) (err error) {
	return
}

func (repo *DynamoDBRepository) DeleteForm(context.Context, string) (err error) {
	return
}

func (repo *DynamoDBRepository) GetFormById(context.Context, string) (form boneless.Form, err error) {
	return
}

func (repo *DynamoDBRepository) GetFormList(context.Context, boneless.FormFilter) (forms []boneless.Form, r boneless.Range, err error) {
	return
}

func (repo *DynamoDBRepository) UpdateForm(context.Context, *boneless.Form) (err error) {
	return
}
