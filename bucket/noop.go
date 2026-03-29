package bucket

import (
	"context"
	"errors"
	"mime/multipart"
	"time"
)

var errBucketDisabled = errors.New("object storage disabled")

type NoopObjectStorageClient struct{}

func NewNoopObjectStorageClient() ObjectStorageInterface {
	return &NoopObjectStorageClient{}
}

func (n *NoopObjectStorageClient) Upload(ctx context.Context, file multipart.File, filename string, contentType string) (string, int64, error) {
	return "", 0, errBucketDisabled
}

func (n *NoopObjectStorageClient) GeneratePresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	return "", errBucketDisabled
}

func (n *NoopObjectStorageClient) GeneratePresignedUploadURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	return "", errBucketDisabled
}

func (n *NoopObjectStorageClient) GetFileInfo(ctx context.Context, objectKey string) (int64, error) {
	return 0, errBucketDisabled
}

func (n *NoopObjectStorageClient) GetFileInfos(ctx context.Context, objectKeys []string) (map[string]int64, error) {
	return nil, errBucketDisabled
}

func (n *NoopObjectStorageClient) Delete(ctx context.Context, objectKey string) error {
	return errBucketDisabled
}
