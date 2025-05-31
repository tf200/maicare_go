package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomSchedule(t *testing.T, employeeID int64) Schedule {
	location := CreateRandomLocation(t)

	arg := CreateScheduleParams{
		EmployeeID:    employeeID,
		LocationID:    location.ID,
		StartDatetime: pgtype.Timestamp{Time: time.Now(), Valid: true},
		EndDatetime:   pgtype.Timestamp{Time: time.Now().Add(24 * time.Hour), Valid: true},
	}
	schedule, err := testQueries.CreateSchedule(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, schedule)
	require.Equal(t, arg.EmployeeID, schedule.EmployeeID)
	return schedule
}

func TestCreateSchedule(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	createRandomSchedule(t, employee.ID)
}

func TestGetMonthlySchedulesByLocation(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	schedule := createRandomSchedule(t, employee.ID)
	year := schedule.StartDatetime.Time.Year()
	month := int32(schedule.StartDatetime.Time.Month())
	schedules, err := testQueries.GetMonthlySchedulesByLocation(context.Background(), GetMonthlySchedulesByLocationParams{
		LocationID: schedule.LocationID,
		Year:       int32(year),
		Month:      month,
	})

	require.NoError(t, err)
	require.NotEmpty(t, schedules)
	for _, s := range schedules {
		require.Equal(t, schedule.LocationID, s.LocationID)
		require.Equal(t, schedule.EmployeeID, s.EmployeeID)
		require.NotEmpty(t, s.StartDatetime.Time)
		require.NotEmpty(t, s.EndDatetime.Time)
	}

}

func TestGetDailySchedulesByLocation(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	schedule := createRandomSchedule(t, employee.ID)
	year, month, day := schedule.StartDatetime.Time.Date()
	schedules, err := testQueries.GetDailySchedulesByLocation(context.Background(), GetDailySchedulesByLocationParams{
		LocationID: schedule.LocationID,
		Year:       int32(year),
		Month:      int32(month),
		Day:        int32(day),
	})
	

	require.NoError(t, err)
	require.NotEmpty(t, schedules)
	for _, s := range schedules {
		require.Equal(t, schedule.LocationID, s.LocationID)
		require.Equal(t, schedule.EmployeeID, s.EmployeeID)
		require.NotEmpty(t, s.StartDatetime.Time)
		require.NotEmpty(t, s.EndDatetime.Time)
	}
}
