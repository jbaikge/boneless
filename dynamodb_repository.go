package gocms

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	dynamoClassPrefix = "class#"
)

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

func (dyn *dynamoClass) FromClass(c Class) {
	dyn.PK = dynamoClassPrefix + c.Id
	dyn.SK = "class_v0"
	dyn.Name = c.Name
	dyn.TableFields = c.TableFields
	dyn.TableLabels = c.TableLabels
	dyn.Created = c.Created
	dyn.Updated = c.Updated
	dyn.Fields = make([]Field, 0, len(c.Fields))
	copy(dyn.Fields, c.Fields)
}

func (dyn dynamoClass) ToClass() (c Class) {
	c.Id = dyn.PK[len(dynamoClassPrefix):]
	c.Name = dyn.Name
	c.TableFields = dyn.TableFields
	c.TableLabels = dyn.TableLabels
	c.Created = dyn.Created
	c.Updated = dyn.Updated
	c.Fields = make([]Field, 0, len(dyn.Fields))
	copy(c.Fields, dyn.Fields)
	return
}

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
	return
}

func (repo DynamoDBRepository) DeleteClass(ctx context.Context, id string) (err error) {
	return
}

func (repo DynamoDBRepository) GetClassById(ctx context.Context, id string) (class Class, err error) {
	return
}

func (repo DynamoDBRepository) GetClassList(ctx context.Context, filter ClassFilter) (list []Class, r Range, err error) {
	return
}

func (repo DynamoDBRepository) UpdateClass(ctx context.Context, class *Class) (err error) {
	return
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
