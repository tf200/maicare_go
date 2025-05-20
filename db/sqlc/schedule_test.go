package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateSchedule(t *testing.T) {
	location := CreateRandomLocation(t)
	employee, _ := createRandomEmployee(t)

	arg := CreateScheduleParams{
		EmployeeID:    employee.ID,
		LocationID:    location.ID,
		StartDatetime: pgtype.Timestamp{Time: time.Now(), Valid: true},
		EndDatetime:   pgtype.Timestamp{Time: time.Now().Add(24 * time.Hour), Valid: true},
	}
	schedule, err := testQueries.CreateSchedule(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, schedule)
	require.Equal(t, arg.EmployeeID, schedule.EmployeeID)
}
