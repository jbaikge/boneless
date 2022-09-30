package dynamodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jbaikge/boneless/models"
)

func (repo *DynamoDBRepository) CreateFile(ctx context.Context, f *models.File) (location string, err error) {
	path := fmt.Sprintf("%s/%s", time.Now().Format("2006/01/02"), f.Filename)
	params := &s3.PutObjectInput{
		Bucket:      &repo.resources.StaticBucket,
		Key:         &path,
		Body:        f.Data,
		ContentType: &f.ContentType,
	}
	if _, err = repo.s3.PutObject(ctx, params); err != nil {
		return
	}

	return fmt.Sprintf("https://%s/%s", repo.resources.StaticDomain, path), nil
}

func (repo *DynamoDBRepository) CreateUploadUrl(ctx context.Context, request models.FileUploadRequest) (response models.FileUploadResponse, err error) {
	key := strings.TrimLeft(request.Key, "/")
	params := &s3.PutObjectInput{
		Bucket:      &repo.resources.StaticBucket,
		Key:         &key,
		ContentType: &request.ContentType,
	}

	expires, err := time.ParseDuration(request.Expires)
	if err != nil {
		err = fmt.Errorf("bad duration, %s: %w", request.Expires, err)
		return
	}
	addDuration := func(po *s3.PresignOptions) {
		po.Expires = expires
	}

	signed, err := s3.NewPresignClient(repo.s3).PresignPutObject(ctx, params, addDuration)
	if err != nil {
		return
	}

	response.URL = signed.URL
	response.Method = signed.Method
	response.Headers = signed.SignedHeader
	response.Location = fmt.Sprintf("https://%s/%s", repo.resources.StaticDomain, key)

	return
}
