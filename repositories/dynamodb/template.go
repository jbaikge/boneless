package dynamodb

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jbaikge/boneless"
)

const templatePrefix = "template#"

func dynamoTemplateIds(id string, version int) (pk string, sk string) {
	pk = templatePrefix + id
	sk = templatePrefix + fmt.Sprintf("v%06d", version)
	return
}

type dynamoTemplate struct {
	PK      string
	SK      string
	Name    string
	Version int
	Created time.Time
	Updated time.Time
}

func newDynamoTemplate(template *boneless.Template) (dyn *dynamoTemplate) {
	pk, sk := dynamoTemplateIds(template.Id, template.Version)
	dyn = &dynamoTemplate{
		PK:      pk,
		SK:      sk,
		Name:    template.Name,
		Version: template.Version,
		Created: template.Created,
		Updated: template.Updated,
	}
	return
}

func (dyn *dynamoTemplate) ToTemplate() (template boneless.Template) {
	template = boneless.Template{
		Id:      dyn.PK[len(templatePrefix):],
		Name:    dyn.Name,
		Version: dyn.Version,
		Created: dyn.Created,
		Updated: dyn.Updated,
	}
	return
}

// Sort by name
type dynamoTemplateByName []*dynamoTemplate

func (arr dynamoTemplateByName) Len() int           { return len(arr) }
func (arr dynamoTemplateByName) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr dynamoTemplateByName) Less(i, j int) bool { return arr[i].Name < arr[j].Name }

func (repo *DynamoDBRepository) CreateTemplate(ctx context.Context, template *boneless.Template) (err error) {
	template.Version = 1
	dbTemplate := newDynamoTemplate(template)
	for _, version := range []int{0, 1} {
		_, dbTemplate.SK = dynamoTemplateIds(template.Id, version)
		if err = repo.putItem(ctx, dbTemplate); err != nil {
			return
		}
	}
	return repo.putTemplateBody(ctx, template)
}

func (repo *DynamoDBRepository) DeleteTemplate(ctx context.Context, id string) (err error) {
	return
}

func (repo *DynamoDBRepository) GetTemplateById(ctx context.Context, id string) (template boneless.Template, err error) {
	pk, sk := dynamoTemplateIds(id, 0)
	dbTemplate := new(dynamoTemplate)
	if err = repo.getItem(ctx, pk, sk, dbTemplate); err != nil {
		return
	}
	template = dbTemplate.ToTemplate()
	if err = repo.getTemplateBody(ctx, &template); err != nil {
		return
	}
	return
}

func (repo *DynamoDBRepository) GetTemplateList(ctx context.Context, filter boneless.TemplateFilter) (templates []boneless.Template, r boneless.Range, err error) {
	return
}

func (repo *DynamoDBRepository) UpdateTemplate(ctx context.Context, template *boneless.Template) (err error) {
	return
}

func (repo *DynamoDBRepository) templateKey(template *boneless.Template) string {
	return fmt.Sprintf("templates/%s/v%06d.html", template.Id, template.Version)
}

func (repo *DynamoDBRepository) getTemplateBody(ctx context.Context, template *boneless.Template) (err error) {
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

func (repo *DynamoDBRepository) putTemplateBody(ctx context.Context, template *boneless.Template) (err error) {
	params := &s3.PutObjectInput{
		Bucket:      &repo.resources.Bucket,
		Key:         aws.String(repo.templateKey(template)),
		Body:        strings.NewReader(template.Body),
		ContentType: aws.String("text/html"),
	}
	_, err = repo.s3.PutObject(ctx, params)
	return
}
