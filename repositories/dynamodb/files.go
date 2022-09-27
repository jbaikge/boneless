package dynamodb

import (
	"context"

	"github.com/jbaikge/boneless"
)

func (repo *DynamoDBRepository) CreateFile(context.Context, *boneless.File) (location string, err error) {
	return
}

func (repo *DynamoDBRepository) CreateUploadUrl(context.Context, boneless.FileUploadRequest) (response boneless.FileUploadResponse, err error) {
	return
}
