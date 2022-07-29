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

func (tables *DynamoDBTables) FromEnv() {
	tables.Class = os.Getenv("DYNAMODB_CLASS_TABLE")
	tables.Document = os.Getenv("DYNAMODB_DOCUMENT_TABLE")
	tables.Path = os.Getenv("DYNAMODB_PATH_TABLE")
	tables.Template = os.Getenv("DYNAMODB_TEMPLATE_TABLE")
}

type DynamoDBRepository struct {
	client *dynamodb.Client
	tables DynamoDBTables
}

func NewDynamoDBRepository(config aws.Config, tables DynamoDBTables) Repository {
	return &DynamoDBRepository{
		client: dynamodb.NewFromConfig(config),
		tables: tables,
	}
}
