package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomAttachmentFile(t *testing.T) AttachmentFile {

	tagvalue := "test"
	arg := CreateAttachmentParams{
		Name: util.RandomString(5),
		File: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
		Size: 23,
		Tag:  &tagvalue,
	}
	attachment, err := testQueries.CreateAttachment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, attachment)

	require.Equal(t, arg.Name, attachment.Name)
	require.Equal(t, arg.File, attachment.File)
	require.Equal(t, arg.Size, attachment.Size)
	require.Equal(t, arg.Tag, attachment.Tag)

	require.NotZero(t, attachment.Uuid)
	require.NotZero(t, attachment.Created)
	return attachment
}

func createRandomAttachmentImage(t *testing.T) AttachmentFile {

	tagvalue := "test"
	arg := CreateAttachmentParams{
		Name: util.RandomString(5),
		File: util.GetRandomImageURL(),
		Size: 23,
		Tag:  &tagvalue,
	}
	attachment, err := testQueries.CreateAttachment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, attachment)

	require.Equal(t, arg.Name, attachment.Name)
	require.Equal(t, arg.File, attachment.File)
	require.Equal(t, arg.Size, attachment.Size)
	require.Equal(t, arg.Tag, attachment.Tag)

	require.NotZero(t, attachment.Uuid)
	require.NotZero(t, attachment.Created)
	return attachment
}

func TestCreateAttachment(t *testing.T) {
	createRandomAttachmentFile(t)
}

func TestGetAttachmentById(t *testing.T) {
	attachment1 := createRandomAttachmentFile(t)
	attachment2, err := testQueries.GetAttachmentById(context.Background(), attachment1.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, attachment2)

	require.Equal(t, attachment1.Name, attachment2.Name)
	require.Equal(t, attachment1.File, attachment2.File)
	require.Equal(t, attachment1.Size, attachment2.Size)
	require.Equal(t, attachment1.Tag, attachment2.Tag)
	require.Equal(t, attachment1.Uuid, attachment2.Uuid)
	require.Equal(t, attachment1.Created, attachment2.Created)

}

func TestDeleteAttachement(t *testing.T) {
	attachment1 := createRandomAttachmentFile(t)
	_, err := testQueries.DeleteAttachment(context.Background(), attachment1.Uuid)
	require.NoError(t, err)

	attachment2, err := testQueries.GetAttachmentById(context.Background(), attachment1.Uuid)
	require.Error(t, err)
	require.Empty(t, attachment2)
}
