package gocms

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	ErrBadRange = errors.New("invalid range")
	ErrNotFound = errors.New("item not found")
)

type DynamoDBTables struct {
	Class    string
	Document string
	Path     string
	Template string
}

type DynamoDBResources struct {
	S3Bucket string
	Tables   DynamoDBTables
}

func (res *DynamoDBResources) FromEnv() {
	res.S3Bucket = os.Getenv("DYNAMODB_S3_BUCKET")
	res.Tables.Class = os.Getenv("DYNAMODB_CLASS_TABLE")
	res.Tables.Document = os.Getenv("DYNAMODB_DOCUMENT_TABLE")
	res.Tables.Path = os.Getenv("DYNAMODB_PATH_TABLE")
	res.Tables.Template = os.Getenv("DYNAMODB_TEMPLATE_TABLE")
}

type DynamoDBRepository struct {
	client    *dynamodb.Client
	resources DynamoDBResources
}

func NewDynamoDBRepository(config aws.Config, resources DynamoDBResources) Repository {
	return &DynamoDBRepository{
		client:    dynamodb.NewFromConfig(config),
		resources: resources,
	}
}
