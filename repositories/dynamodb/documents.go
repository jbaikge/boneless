package dynamodb

import (
	"context"
	"fmt"

	"github.com/jbaikge/boneless"
)

const documentPrefix = "doc#"

func dynamoDocumentIds(id string, version int) (pk string, sk string) {
	pk = documentPrefix + id
	sk = documentPrefix + fmt.Sprintf("v%06d", version)
	return
}

func (repo *DynamoDBRepository) CreateDocument(context.Context, *boneless.Document) (err error) {
	return
}

func (repo *DynamoDBRepository) DeleteDocument(context.Context, string) (err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentById(context.Context, string) (doc boneless.Document, err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentByPath(context.Context, string) (doc boneless.Document, err error) {
	return
}

func (repo *DynamoDBRepository) GetDocumentList(context.Context, boneless.DocumentFilter) (docs []boneless.Document, r boneless.Range, err error) {
	return
}

func (repo *DynamoDBRepository) UpdateDocument(context.Context, *boneless.Document) (err error) {
	return
}
