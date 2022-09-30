package dynamodb

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jbaikge/boneless/models"
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

func newDynamoTemplate(template *models.Template) (dyn *dynamoTemplate) {
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

func (dyn *dynamoTemplate) ToTemplate() (template models.Template) {
	template = models.Template{
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

func (repo *DynamoDBRepository) CreateTemplate(ctx context.Context, template *models.Template) (err error) {
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
	pk, sk := dynamoTemplateIds(id, 0)
	dbTemplate := new(dynamoTemplate)
	if err = repo.getItem(ctx, pk, sk, dbTemplate); err != nil {
		return
	}

	// Delete past template versions
	objects := make([]s3types.ObjectIdentifier, 0, dbTemplate.Version)
	for version := 0; version <= dbTemplate.Version; version++ {
		_, delSk := dynamoTemplateIds(id, version)
		if err = repo.deleteItem(ctx, pk, delSk); err != nil {
			return
		}

		// Version zero has no HTML
		if version == 0 {
			continue
		}

		objects = append(objects, s3types.ObjectIdentifier{
			Key: aws.String(repo.templateKey(id, version)),
		})
	}

	// Delete all the S3 objects at once
	params := &s3.DeleteObjectsInput{
		Bucket: &repo.resources.Bucket,
		Delete: &s3types.Delete{
			Objects: objects,
		},
	}
	_, err = repo.s3.DeleteObjects(ctx, params)

	return
}

func (repo *DynamoDBRepository) GetTemplateById(ctx context.Context, id string) (template models.Template, err error) {
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

func (repo *DynamoDBRepository) GetTemplateList(ctx context.Context, filter models.TemplateFilter) (list []models.Template, r models.Range, err error) {
	var response *dynamodb.ScanOutput
	dbTemplates := make([]*dynamoTemplate, 0, 64)

	key, err := repo.marshalKey(dynamoTemplateIds("", 0))
	if err != nil {
		return
	}

	params := &dynamodb.ScanInput{
		TableName:        &repo.resources.Table,
		FilterExpression: aws.String("SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sK": key["SK"],
		},
	}
	paginator := dynamodb.NewScanPaginator(repo.db, params)
	for paginator.HasMorePages() {
		response, err = paginator.NextPage(ctx)
		if err != nil {
			return
		}
		tmp := make([]*dynamoTemplate, 0, len(response.Items))
		if err = attributevalue.UnmarshalListOfMaps(response.Items, &dbTemplates); err != nil {
			return
		}
		dbTemplates = append(dbTemplates, tmp...)
	}

	sort.Sort(dynamoTemplateByName(dbTemplates))

	r.Size = len(dbTemplates)
	list = make([]models.Template, 0, filter.Range.SliceLen())
	for i := filter.Range.Start; i < len(dbTemplates) && i <= filter.Range.End; i++ {
		template := dbTemplates[i].ToTemplate()
		if err = repo.getTemplateBody(ctx, &template); err != nil {
			return
		}
		list = append(list, template)
	}

	if filter.Range.Start > 0 && len(list) == 0 {
		err = ErrBadRange
		return
	}

	r.Start = filter.Range.Start
	r.End = filter.Range.End
	if length := len(list); length > 0 {
		r.End += length - 1
	}

	return
}

func (repo *DynamoDBRepository) UpdateTemplate(ctx context.Context, template *models.Template) (err error) {
	// Fetch current template
	oldTemplate := new(dynamoTemplate)
	pk, sk := dynamoTemplateIds(template.Id, 0)
	if err = repo.getItem(ctx, pk, sk, oldTemplate); err != nil {
		return
	}

	// Increment version based on current version in database
	template.Version = oldTemplate.Version

	// Add new template with new version
	dbTemplate := newDynamoTemplate(template)
	if err = repo.putItem(ctx, dbTemplate); err != nil {
		return
	}

	// Update values in v0
	values := map[string]interface{}{
		"Name":    template.Name,
		"Version": template.Version,
		"Updated": template.Updated,
	}
	if err = repo.updateItem(ctx, pk, sk, values); err != nil {
		return
	}

	// Add new version of body to S3
	return repo.putTemplateBody(ctx, template)
}

func (repo *DynamoDBRepository) templateKey(id string, version int) string {
	return fmt.Sprintf("templates/%s/v%06d.html", id, version)
}

func (repo *DynamoDBRepository) getTemplateBody(ctx context.Context, template *models.Template) (err error) {
	params := &s3.GetObjectInput{
		Bucket: &repo.resources.Bucket,
		Key:    aws.String(repo.templateKey(template.Id, template.Version)),
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

func (repo *DynamoDBRepository) putTemplateBody(ctx context.Context, template *models.Template) (err error) {
	params := &s3.PutObjectInput{
		Bucket:      &repo.resources.Bucket,
		Key:         aws.String(repo.templateKey(template.Id, template.Version)),
		Body:        strings.NewReader(template.Body),
		ContentType: aws.String("text/html"),
	}
	_, err = repo.s3.PutObject(ctx, params)
	return
}
