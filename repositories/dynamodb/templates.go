package dynamodb

import (
	"context"
	"fmt"

	"github.com/jbaikge/boneless"
)

const templatePrefix = "template#"

func dynamoTemplateIds(id string, version int) (pk string, sk string) {
	pk = templatePrefix + id
	sk = templatePrefix + fmt.Sprintf("v%06d", version)
	return
}

func (repo *DynamoDBRepository) CreateTemplate(context.Context, *boneless.Template) (err error) {
	return
}

func (repo *DynamoDBRepository) DeleteTemplate(context.Context, string) (err error) {
	return
}

func (repo *DynamoDBRepository) GetTemplateById(context.Context, string) (template boneless.Template, err error) {
	return
}

func (repo *DynamoDBRepository) GetTemplateList(context.Context, boneless.TemplateFilter) (templates []boneless.Template, r boneless.Range, err error) {
	return
}

func (repo *DynamoDBRepository) UpdateTemplate(context.Context, *boneless.Template) (err error) {
	return
}
