package bucket

import (
	"bytes"
	"context"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	// Create a random text file to upload
	content := []byte("This is a test file for upload.")

	var f multipart.File = &InMemoryFile{bytes.NewReader(content)}
	defer f.Close()

	file, size, err := testBucketClient.Upload(context.Background(), f, "test_upload.txt", "text/plain")
	require.NoError(t, err)
	require.NotEmpty(t, file)
	require.Equal(t, int64(len(content)), size)

	err = testBucketClient.Delete(context.Background(), file)
	require.NoError(t, err)
}
