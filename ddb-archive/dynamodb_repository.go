package gocms

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const emptyParentId = "#NULL#"

var (
	ErrBadRange = errors.New("invalid range")
	ErrNotFound = errors.New("item not found")
)

type DynamoDBBuckets struct {
	Document string
	Static   string
}

type DynamoDBTables struct {
	Class    string
	Document string
	Path     string
	Sort     string
	Template string
}

type DynamoDBResources struct {
	Buckets DynamoDBBuckets
	Tables  DynamoDBTables
}

func (res *DynamoDBResources) FromEnv() {
	res.Buckets.Document = os.Getenv("BUCKET_DOCUMENT")
	res.Buckets.Static = os.Getenv("BUCKET_STATIC")
	res.Tables.Class = os.Getenv("DYNAMODB_CLASS_TABLE")
	res.Tables.Document = os.Getenv("DYNAMODB_DOCUMENT_TABLE")
	res.Tables.Path = os.Getenv("DYNAMODB_PATH_TABLE")
	res.Tables.Sort = os.Getenv("DYNAMODB_SORT_TABLE")
	res.Tables.Template = os.Getenv("DYNAMODB_TEMPLATE_TABLE")
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
