package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomShift(t *testing.T, locationID int64) LocationShift {
	arg := CreateShiftParams{
		LocationID: locationID,
		ShiftName:  "Ochtenddienst",
		StartTime:  pgtype.Time{Microseconds: (7*3600 + 30*60) * 1000000, Valid: true},  // 07:30:00
		EndTime:    pgtype.Time{Microseconds: (15*3600 + 30*60) * 1000000, Valid: true}, // 15:30:00
	}

	arg2 := CreateShiftParams{
		LocationID: locationID,
		ShiftName:  "Avonddienst",
		StartTime:  pgtype.Time{Microseconds: (15 * 3600) * 1000000, Valid: true}, // 15:00:00
		EndTime:    pgtype.Time{Microseconds: (23 * 3600) * 1000000, Valid: true}, // 23:00:00
	}

	arg3 := CreateShiftParams{
		LocationID: locationID,
		ShiftName:  "Slaapdienst of Waakdienst",
		StartTime:  pgtype.Time{Microseconds: (23 * 3600) * 1000000, Valid: true},      // 23:00:00
		EndTime:    pgtype.Time{Microseconds: (7*3600 + 30*60) * 1000000, Valid: true}, // 07:30:00
	}

	shift, err := testQueries.CreateShift(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, shift)
	shift2, err2 := testQueries.CreateShift(context.Background(), arg2)
	require.NoError(t, err2)
	require.NotEmpty(t, shift2)

	shift3, err3 := testQueries.CreateShift(context.Background(), arg3)
	require.NoError(t, err3)
	require.NotEmpty(t, shift3)

	return shift
}

func TestCreateShift(t *testing.T) {
	location := CreateRandomLocation(t)
	shift := createRandomShift(t, location.ID)

	require.NotEmpty(t, shift)
	require.Equal(t, location.ID, shift.LocationID)
}
