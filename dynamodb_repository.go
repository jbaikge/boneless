package gocms

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var _ Repository = new(DynamoDBRepository)

type DynamoDBRepository struct {
	client *dynamodb.Client
}

func NewDynamoDBRepository(config aws.Config) Repository {
	return &DynamoDBRepository{}
}

func (r DynamoDBRepository) DeleteClass(id string) (err error) {
	return
}

func (r DynamoDBRepository) GetAllClasses() (classes []Class, err error) {
	return
}

func (r DynamoDBRepository) GetClassById(id string) (class Class, err error) {
	return
}

func (r DynamoDBRepository) InsertClass(class *Class) (err error) {
	return
}

func (r DynamoDBRepository) UpdateClass(class *Class) (err error) {
	return
}
