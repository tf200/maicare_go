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

// type B2Client struct {
// 	Client *b2.Client
// 	Bucket *b2.Bucket
// }

// func NewB2Client(config util.Config) (*B2Client, error) {
// 	ctx := context.Background()

// 	id := config.B2KeyID
// 	key := config.B2Key
// 	bucketName := config.B2Bucket

// 	client, err := b2.NewClient(ctx, id, key)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create B2 client: %v", err)
// 	}

// 	bucket, err := client.Bucket(ctx, bucketName)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get bucket: %v", err)
// 	}

// 	return &B2Client{
// 		Client: client,
// 		Bucket: bucket,
// 	}, nil
// }

// func (b *B2Client) UploadToB2(ctx context.Context, file multipart.File, filename string) error {
// 	// Create a new writer for the B2 object
// 	obj := b.Bucket.Object(filename)
// 	writer := obj.NewWriter(ctx)

// 	writer.ConcurrentUploads = 4

// 	// Copy the file to B2
// 	_, err := io.Copy(writer, file)
// 	if err != nil {
// 		writer.Close()
// 		return fmt.Errorf("failed to copy file to B2: %v", err)
// 	}

// 	if err := writer.Close(); err != nil {
// 		return fmt.Errorf("failed to close writer: %v", err)
// 	}

// 	return nil
// }

// func (b *B2Client) DeleteFromB2(ctx context.Context, filename string) error {
// 	obj := b.Bucket.Object(filename)
// 	if err := obj.Delete(ctx); err != nil {
// 		return fmt.Errorf("failed to delete file from B2: %v", err)
// 	}
// 	return nil
// }

// func (b *B2Client) DeleteFromB2URL(ctx context.Context, fileURL string) error {
// 	// Parse the URL
// 	parsedURL, err := url.Parse(fileURL)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse URL: %v", err)
// 	}

// 	// Extract filename from path
// 	// Path format: /file/bucket-name/filename
// 	pathParts := strings.Split(parsedURL.Path, "/")
// 	if len(pathParts) < 4 {
// 		return fmt.Errorf("invalid B2 file URL format")
// 	}
// 	filename := pathParts[len(pathParts)-1]

// 	// Delete the file
// 	return b.DeleteFromB2(ctx, filename)
// }

// new object storage client
type ObjectStorageClient struct {
	Client *minio.Client
	Bucket string
}

func NewObjectStorageClient(ctx context.Context, config util.Config) (*ObjectStorageClient, error) {
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
