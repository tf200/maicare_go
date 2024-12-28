package api

import (
	"context"
	"testing"

	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomLocation(t *testing.T) *db.Location {
	arg := db.CreateLocationParams{
		Name:     util.RandomString(5),
		Address:  util.RandomString(8),
		Capacity: util.Int32Ptr(52),
	}

	location, err := testStore.CreateLocation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, location)

	// Check if the returned location matches the input
	require.Equal(t, arg.Name, location.Name)
	require.Equal(t, arg.Address, location.Address)
	require.Equal(t, arg.Capacity, location.Capacity)

	// Verify ID is generated
	require.NotZero(t, location.ID)
	return &location
}
