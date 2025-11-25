package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateSchedule(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateScheduleParams
		checks func(t *testing.T, schedule CreateScheduleRow, err error)
	}{
		{
			name: "successful schedule creation",
			setup: func(ctx context.Context, qtx *Queries) CreateScheduleParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				creator := createRandomEmployeeProfile(ctx, qtx)
				startTime := time.Now().Add(time.Hour)
				endTime := startTime.Add(time.Hour)
				return CreateScheduleParams{
					EmployeeID:          employee.ID,
					LocationID:          location.ID,
					LocationShiftID:     nil,
					Color:               randomStringPtrSchedule(7),
					IsCustom:            false,
					CreatedByEmployeeID: creator.ID,
					StartDatetime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndDatetime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, schedule CreateScheduleRow, err error) {
				require.NoError(t, err, "CreateSchedule should not return an error")
				require.NotZero(t, schedule.ID)
				require.Equal(t, false, schedule.IsCustom)
			},
		},
		{
			name: "schedule creation with shift",
			setup: func(ctx context.Context, qtx *Queries) CreateScheduleParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				shift := createRandomShift(ctx, qtx)
				creator := createRandomEmployeeProfile(ctx, qtx)
				startTime := time.Now().Add(time.Hour)
				endTime := startTime.Add(time.Hour)
				return CreateScheduleParams{
					EmployeeID:          employee.ID,
					LocationID:          location.ID,
					LocationShiftID:     &shift.ID,
					Color:               randomStringPtrSchedule(7),
					IsCustom:            true,
					CreatedByEmployeeID: creator.ID,
					StartDatetime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndDatetime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, schedule CreateScheduleRow, err error) {
				require.NoError(t, err, "CreateSchedule should not return an error")
				require.NotZero(t, schedule.ID)
				require.True(t, schedule.IsCustom)
				require.NotNil(t, schedule.LocationShiftID)
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
			schedule, err := qtx.CreateSchedule(ctx, params)
			tt.checks(t, schedule, err)
		})
	}
}

func TestGetScheduleById(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, schedule GetScheduleByIdRow, err error)
	}{
		{
			name: "get existing schedule",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				location := createRandomLocation(ctx, qtx)
				schedule := createRandomSchedule(ctx, qtx, location.ID)
				return schedule.ID
			},
			checks: func(t *testing.T, schedule GetScheduleByIdRow, err error) {
				require.NoError(t, err, "GetScheduleById should not return an error")
				require.NotZero(t, schedule.ID)
			},
		},
		{
			name: "get non-existent schedule",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, schedule GetScheduleByIdRow, err error) {
				require.Error(t, err, "GetScheduleById should error for non-existent schedule")
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

			scheduleID := tt.setup(ctx, qtx)
			schedule, err := qtx.GetScheduleById(ctx, scheduleID)
			tt.checks(t, schedule, err)
		})
	}
}

func TestUpdateSchedule(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateScheduleParams
		checks func(t *testing.T, schedule UpdateScheduleRow, err error)
	}{
		{
			name: "update existing schedule",
			setup: func(ctx context.Context, qtx *Queries) UpdateScheduleParams {
				location := createRandomLocation(ctx, qtx)
				schedule := createRandomSchedule(ctx, qtx, location.ID)
				newEmployee := createRandomEmployeeProfile(ctx, qtx)
				newLocation := createRandomLocation(ctx, qtx)
				startTime := time.Now().Add(2 * time.Hour)
				endTime := startTime.Add(time.Hour)
				return UpdateScheduleParams{
					ID:              schedule.ID,
					EmployeeID:      newEmployee.ID,
					LocationID:      newLocation.ID,
					LocationShiftID: nil,
					Color:           randomStringPtrSchedule(7),
					StartDatetime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndDatetime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, schedule UpdateScheduleRow, err error) {
				require.NoError(t, err, "UpdateSchedule should not return an error")
				require.NotZero(t, schedule.ID)
			},
		},
		{
			name: "update non-existent schedule",
			setup: func(ctx context.Context, qtx *Queries) UpdateScheduleParams {
				newEmployee := createRandomEmployeeProfile(ctx, qtx)
				newLocation := createRandomLocation(ctx, qtx)
				startTime := time.Now().Add(2 * time.Hour)
				endTime := startTime.Add(time.Hour)
				return UpdateScheduleParams{
					ID:              uuid.New(),
					EmployeeID:      newEmployee.ID,
					LocationID:      newLocation.ID,
					LocationShiftID: nil,
					Color:           randomStringPtrSchedule(7),
					StartDatetime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndDatetime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, schedule UpdateScheduleRow, err error) {
				require.Error(t, err, "UpdateSchedule should error for non-existent schedule")
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
			schedule, err := qtx.UpdateSchedule(ctx, params)
			tt.checks(t, schedule, err)
		})
	}
}

func TestDeleteSchedule(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete existing schedule",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				location := createRandomLocation(ctx, qtx)
				schedule := createRandomSchedule(ctx, qtx, location.ID)
				return schedule.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteSchedule should not error")
			},
		},
		{
			name: "delete non-existent schedule",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteSchedule should not error for non-existent schedule")
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

			scheduleID := tt.setup(ctx, qtx)
			err = qtx.DeleteSchedule(ctx, scheduleID)
			tt.checks(t, err)
		})
	}
}

func TestGetDailySchedulesByLocation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) GetDailySchedulesByLocationParams
		checks func(t *testing.T, schedules []GetDailySchedulesByLocationRow, err error)
	}{
		{
			name: "get schedules for day with schedules",
			setup: func(ctx context.Context, qtx *Queries) GetDailySchedulesByLocationParams {
				location := createRandomLocation(ctx, qtx)
				// Create a schedule for today
				createRandomSchedule(ctx, qtx, location.ID)
				now := time.Now()
				return GetDailySchedulesByLocationParams{
					Year:       int32(now.Year()),
					Month:      int32(now.Month()),
					Day:        int32(now.Day()),
					LocationID: location.ID,
				}
			},
			checks: func(t *testing.T, schedules []GetDailySchedulesByLocationRow, err error) {
				require.NoError(t, err, "GetDailySchedulesByLocation should not error")
				require.GreaterOrEqual(t, len(schedules), 1)
			},
		},
		{
			name: "get schedules for day without schedules",
			setup: func(ctx context.Context, qtx *Queries) GetDailySchedulesByLocationParams {
				location := createRandomLocation(ctx, qtx)
				// Use a past date
				past := time.Now().AddDate(0, 0, -10)
				return GetDailySchedulesByLocationParams{
					Year:       int32(past.Year()),
					Month:      int32(past.Month()),
					Day:        int32(past.Day()),
					LocationID: location.ID,
				}
			},
			checks: func(t *testing.T, schedules []GetDailySchedulesByLocationRow, err error) {
				require.NoError(t, err, "GetDailySchedulesByLocation should not error")
				require.Len(t, schedules, 0)
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
			schedules, err := qtx.GetDailySchedulesByLocation(ctx, params)
			tt.checks(t, schedules, err)
		})
	}
}

func TestGetEmployeeSchedules(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) GetEmployeeSchedulesParams
		checks func(t *testing.T, schedules []GetEmployeeSchedulesRow, err error)
	}{
		{
			name: "get schedules for employee in period",
			setup: func(ctx context.Context, qtx *Queries) GetEmployeeSchedulesParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				creator := createRandomEmployeeProfile(ctx, qtx)
				// Create a schedule for this employee
				startTime := time.Now().Add(-30 * time.Minute)
				endTime := startTime.Add(time.Hour)
				_, err := qtx.CreateSchedule(ctx, CreateScheduleParams{
					EmployeeID:          employee.ID,
					LocationID:          location.ID,
					LocationShiftID:     nil,
					Color:               randomStringPtrSchedule(7),
					IsCustom:            false,
					CreatedByEmployeeID: creator.ID,
					StartDatetime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndDatetime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
				})
				require.NoError(t, err)
				start := time.Now().Add(-time.Hour)
				end := time.Now().Add(time.Hour)
				return GetEmployeeSchedulesParams{
					EmployeeID: employee.ID,
					PeriodStart: pgtype.Timestamp{
						Time:  start,
						Valid: true,
					},
					PeriodEnd: pgtype.Timestamp{
						Time:  end,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, schedules []GetEmployeeSchedulesRow, err error) {
				require.NoError(t, err, "GetEmployeeSchedules should not error")
				require.GreaterOrEqual(t, len(schedules), 1)
			},
		},
		{
			name: "get schedules for employee outside period",
			setup: func(ctx context.Context, qtx *Queries) GetEmployeeSchedulesParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				start := time.Now().Add(24 * time.Hour)
				end := time.Now().Add(25 * time.Hour)
				return GetEmployeeSchedulesParams{
					EmployeeID: employee.ID,
					PeriodStart: pgtype.Timestamp{
						Time:  start,
						Valid: true,
					},
					PeriodEnd: pgtype.Timestamp{
						Time:  end,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, schedules []GetEmployeeSchedulesRow, err error) {
				require.NoError(t, err, "GetEmployeeSchedules should not error")
				require.Len(t, schedules, 0)
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
			schedules, err := qtx.GetEmployeeSchedules(ctx, params)
			tt.checks(t, schedules, err)
		})
	}
}

func TestGetMonthlySchedulesByLocation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) GetMonthlySchedulesByLocationParams
		checks func(t *testing.T, schedules []GetMonthlySchedulesByLocationRow, err error)
	}{
		{
			name: "get schedules for month with schedules",
			setup: func(ctx context.Context, qtx *Queries) GetMonthlySchedulesByLocationParams {
				location := createRandomLocation(ctx, qtx)
				// Create a schedule for this month
				createRandomSchedule(ctx, qtx, location.ID)
				now := time.Now()
				return GetMonthlySchedulesByLocationParams{
					Year:       int32(now.Year()),
					Month:      int32(now.Month()),
					LocationID: location.ID,
				}
			},
			checks: func(t *testing.T, schedules []GetMonthlySchedulesByLocationRow, err error) {
				require.NoError(t, err, "GetMonthlySchedulesByLocation should not error")
				require.GreaterOrEqual(t, len(schedules), 1)
			},
		},
		{
			name: "get schedules for month without schedules",
			setup: func(ctx context.Context, qtx *Queries) GetMonthlySchedulesByLocationParams {
				location := createRandomLocation(ctx, qtx)
				// Use a past month
				past := time.Now().AddDate(0, -1, 0)
				return GetMonthlySchedulesByLocationParams{
					Year:       int32(past.Year()),
					Month:      int32(past.Month()),
					LocationID: location.ID,
				}
			},
			checks: func(t *testing.T, schedules []GetMonthlySchedulesByLocationRow, err error) {
				require.NoError(t, err, "GetMonthlySchedulesByLocation should not error")
				require.Len(t, schedules, 0)
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
			schedules, err := qtx.GetMonthlySchedulesByLocation(ctx, params)
			tt.checks(t, schedules, err)
		})
	}
}

// Helpers
func createRandomSchedule(ctx context.Context, qtx *Queries, locationID int64) CreateScheduleRow {
	employee := createRandomEmployeeProfile(ctx, qtx)
	creator := createRandomEmployeeProfile(ctx, qtx)
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(time.Hour)
	params := CreateScheduleParams{
		EmployeeID:          employee.ID,
		LocationID:          locationID,
		LocationShiftID:     nil,
		Color:               randomStringPtrSchedule(7),
		IsCustom:            false,
		CreatedByEmployeeID: creator.ID,
		StartDatetime: pgtype.Timestamp{
			Time:  startTime,
			Valid: true,
		},
		EndDatetime: pgtype.Timestamp{
			Time:  endTime,
			Valid: true,
		},
	}
	schedule, err := qtx.CreateSchedule(ctx, params)
	if err != nil {
		panic(err)
	}
	return schedule
}

func randomStringPtrSchedule(n int) *string {
	s := util.RandomString(n)
	return &s
}
