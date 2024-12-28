package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func CreateRandomLocation(t *testing.T) *Location {
	arg := CreateLocationParams{
		Name:     util.RandomString(5),
		Address:  util.RandomString(8),
		Capacity: util.Int32Ptr(25),
	}

	location, err := testQueries.CreateLocation(context.Background(), arg)
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

func TestCreateLocation(t *testing.T) {
	CreateRandomLocation(t)
}

func TestListLocations(t *testing.T) {
	for i := 0; i < 4; i++ {
		CreateRandomLocation(t)
	}

	locations, err := testQueries.ListLocations(context.Background())
	require.NoError(t, err)

	for _, location := range locations {
		require.NotEmpty(t, location)
	}
}
