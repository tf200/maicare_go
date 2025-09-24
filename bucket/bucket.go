package bucket

import (
	"context"
	"fmt"
	"maicare_go/util"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//go:generate mockgen -package bucketmocks -destination=../bucket/mocks/bucket_mock.go maicare_go/bucket ObjectStorageInterface
type ObjectStorageInterface interface {
	Upload(ctx context.Context, file multipart.File, filename string, contentType string) (string, int64, error)
	GeneratePresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, objectKey string) error
}

// new object storage client
type ObjectStorageClient struct {
	Client *minio.Client
	Bucket string
}

func NewObjectStorageClient(ctx context.Context, config util.Config) (ObjectStorageInterface, error) {
	minioClient, err := minio.New(config.B2Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.B2KeyID, config.B2Key, ""),
		Secure: true,
		Region: "eu-central-003",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	// Check if the bucket exists
	exists, err := minioClient.BucketExists(ctx, config.B2Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", config.B2Bucket)
	}
	return &ObjectStorageClient{
		Client: minioClient,
		Bucket: config.B2Bucket,
	}, nil
}

func (o *ObjectStorageClient) Upload(ctx context.Context, file multipart.File, filename string, contentType string) (string, int64, error) {
	// Upload the file to the bucket
	uploadInfo, err := o.Client.PutObject(ctx, o.Bucket, filename, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", 0, fmt.Errorf("failed to upload file: %v", err)
	}

	// Return just the object key, not a full URL
	return uploadInfo.Key, uploadInfo.Size, nil
}

func (o *ObjectStorageClient) GeneratePresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	// Generate a presigned URL for private bucket access
	presignedURL, err := o.Client.PresignedGetObject(ctx, o.Bucket, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.String(), nil
}

func (o *ObjectStorageClient) Delete(ctx context.Context, objectKey string) error {
	// Delete the object from the bucket
	err := o.Client.RemoveObject(ctx, o.Bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}
	return nil
}
