package bucket

import (
	"context"
	"fmt"
	"io"
	"maicare_go/util"
	"mime/multipart"

	"github.com/Backblaze/blazer/b2"
)

type B2Client struct {
	Client *b2.Client
	Bucket *b2.Bucket
}

func NewB2Client(config util.Config) (*B2Client, error) {
	ctx := context.Background()

	id := config.B2KeyID
	key := config.B2Key
	bucketName := config.B2Bucket

	client, err := b2.NewClient(ctx, id, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create B2 client: %v", err)
	}

	bucket, err := client.Bucket(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %v", err)
	}

	return &B2Client{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (b *B2Client) UploadToB2(ctx context.Context, file multipart.File, filename string) error {
	// Create a new writer for the B2 object
	obj := b.Bucket.Object(filename)
	writer := obj.NewWriter(ctx)

	writer.ConcurrentUploads = 4

	// Copy the file to B2
	_, err := io.Copy(writer, file)
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to copy file to B2: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}
