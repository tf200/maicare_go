package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateShift(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateShiftParams
		checks func(t *testing.T, shift LocationShift, err error)
	}{
		{
			name: "successful shift creation",
			setup: func(ctx context.Context, qtx *Queries) CreateShiftParams {
				location := createRandomLocation(ctx, qtx)
				return CreateShiftParams{
					LocationID: location.ID,
					ShiftName:  util.RandomString(10),
					StartTime:  randomPgTime(),
					EndTime:    randomPgTime(),
				}
			},
			checks: func(t *testing.T, shift LocationShift, err error) {
				require.NoError(t, err, "CreateShift should not return an error")
				require.NotZero(t, shift.ID)
				require.NotZero(t, shift.CreatedAt)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			shift, err := qtx.CreateShift(ctx, params)
			tt.checks(t, shift, err)
		})
	}
}

func TestGetShiftByID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, shift LocationShift, err error)
	}{
		{
			name: "get existing shift",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				shift := createRandomShift(ctx, qtx)
				return shift.ID
			},
			checks: func(t *testing.T, shift LocationShift, err error) {
				require.NoError(t, err, "GetShiftByID should not return an error")
				require.NotZero(t, shift.ID)
			},
		},
		{
			name: "get non-existent shift",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 999999 // Assuming this ID doesn't exist
			},
			checks: func(t *testing.T, shift LocationShift, err error) {
				require.Error(t, err, "GetShiftByID should error for non-existent shift")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			shiftID := tt.setup(ctx, qtx)
			shift, err := qtx.GetShiftByID(ctx, shiftID)
			tt.checks(t, shift, err)
		})
	}
}

func TestUpdateShift(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateShiftParams
		checks func(t *testing.T, shift LocationShift, err error)
	}{
		{
			name: "update existing shift",
			setup: func(ctx context.Context, qtx *Queries) UpdateShiftParams {
				shift := createRandomShift(ctx, qtx)
				return UpdateShiftParams{
					ID:        shift.ID,
					ShiftName: util.RandomString(10),
					StartTime: randomPgTime(),
					EndTime:   randomPgTime(),
				}
			},
			checks: func(t *testing.T, shift LocationShift, err error) {
				require.NoError(t, err, "UpdateShift should not return an error")
				require.NotZero(t, shift.ID)
			},
		},
		{
			name: "update non-existent shift",
			setup: func(ctx context.Context, qtx *Queries) UpdateShiftParams {
				return UpdateShiftParams{
					ID:        999999,
					ShiftName: util.RandomString(10),
					StartTime: randomPgTime(),
					EndTime:   randomPgTime(),
				}
			},
			checks: func(t *testing.T, shift LocationShift, err error) {
				require.Error(t, err, "UpdateShift should error for non-existent shift")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			shift, err := qtx.UpdateShift(ctx, params)
			tt.checks(t, shift, err)
		})
	}
}

func TestDeleteShift(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete existing shift",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				shift := createRandomShift(ctx, qtx)
				return shift.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteShift should not error")
			},
		},
		{
			name: "delete non-existent shift",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 999999
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteShift should not error for non-existent shift")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			shiftID := tt.setup(ctx, qtx)
			err = qtx.DeleteShift(ctx, shiftID)
			tt.checks(t, err)
		})
	}
}

func TestGetShiftsByLocationID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, shifts []LocationShift, err error)
	}{
		{
			name: "get shifts for location with shifts",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				location := createRandomLocation(ctx, qtx)
				// Create a shift for this location
				createRandomShift(ctx, qtx)
				return location.ID
			},
			checks: func(t *testing.T, shifts []LocationShift, err error) {
				require.NoError(t, err, "GetShiftsByLocationID should not error")
				require.GreaterOrEqual(t, len(shifts), 1)
			},
		},
		{
			name: "get shifts for location without shifts",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				location := createRandomLocation(ctx, qtx)
				return location.ID
			},
			checks: func(t *testing.T, shifts []LocationShift, err error) {
				require.NoError(t, err, "GetShiftsByLocationID should not error")
				require.Len(t, shifts, 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			locationID := tt.setup(ctx, qtx)
			shifts, err := qtx.GetShiftsByLocationID(ctx, locationID)
			tt.checks(t, shifts, err)
		})
	}
}

func TestCheckAllShiftsExist(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CheckAllShiftsExistParams
		checks func(t *testing.T, exists bool, err error)
	}{
		{
			name: "all shifts exist",
			setup: func(ctx context.Context, qtx *Queries) CheckAllShiftsExistParams {
				shift1 := createRandomShift(ctx, qtx)
				shift2 := createRandomShift(ctx, qtx)
				return CheckAllShiftsExistParams{
					ExpectedCount: 2,
					Ids:           []int32{int32(shift1.ID), int32(shift2.ID)},
				}
			},
			checks: func(t *testing.T, exists bool, err error) {
				require.NoError(t, err, "CheckAllShiftsExist should not error")
				require.True(t, exists)
			},
		},
		{
			name: "not all shifts exist",
			setup: func(ctx context.Context, qtx *Queries) CheckAllShiftsExistParams {
				shift := createRandomShift(ctx, qtx)
				return CheckAllShiftsExistParams{
					ExpectedCount: 2,
					Ids:           []int32{int32(shift.ID), 999999},
				}
			},
			checks: func(t *testing.T, exists bool, err error) {
				require.NoError(t, err, "CheckAllShiftsExist should not error")
				require.False(t, exists)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			exists, err := qtx.CheckAllShiftsExist(ctx, params)
			tt.checks(t, exists, err)
		})
	}
}

// Helpers
func createRandomShift(ctx context.Context, qtx *Queries) LocationShift {
	location := createRandomLocation(ctx, qtx)
	params := CreateShiftParams{
		LocationID: location.ID,
		ShiftName:  util.RandomString(10),
		StartTime:  randomPgTime(),
		EndTime:    randomPgTime(),
	}
	shift, err := qtx.CreateShift(ctx, params)
	if err != nil {
		panic(err)
	}
	return shift
}

func randomPgTime() pgtype.Time {
	hour := util.RandomInt(0, 23)
	minute := util.RandomInt(0, 59)
	second := util.RandomInt(0, 59)
	microseconds := int64(hour*3600+minute*60+second) * 1000000
	return pgtype.Time{Microseconds: microseconds, Valid: true}
}
