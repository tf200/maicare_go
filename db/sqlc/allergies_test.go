package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateallergy(t *testing.T) {
	allergyType, err := testQueries.Createallergy(context.Background(), "Peanut")
	require.NoError(t, err)
	require.NotEmpty(t, allergyType)
	require.Equal(t, "Peanut", allergyType.Name)
}
