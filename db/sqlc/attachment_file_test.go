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

func TestCreateAttachment(t *testing.T) {
	createRandomAttachmentFile(t)
}
