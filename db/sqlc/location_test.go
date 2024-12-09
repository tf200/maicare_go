package db

import (
	"context"
	"maicare_go/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateRandomLocation(t *testing.T) *Location {
	arg := CreateLocationParams{
		Name:    util.RandomString(5),
		Address: util.RandomString(8),
		Capacity: pgtype.Int4{
			Int32: 25,
			Valid: true,
		},
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
