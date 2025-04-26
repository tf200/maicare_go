package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomAppointment(t *testing.T, employeeID *int64) ScheduledAppointment {
	arg := CreateAppointmentParams{
		CreatorEmployeeID: employeeID,
		StartTime:         pgtype.Timestamp{Time: time.Now(), Valid: true},
		EndTime:           pgtype.Timestamp{Time: time.Now().Add(1 * time.Hour), Valid: true},
		Location:          util.StringPtr("Test Location"),
		Description:       util.StringPtr("Test Description"),
	}

	appointment, err := testQueries.CreateAppointment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointment)
	return appointment
}

func TestCreateAppointment(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	createRandomAppointment(t, &employee.ID)
}

func TestAddAppointmentParticipant(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	employee2, _ := createRandomEmployee(t)
	appointment := createRandomAppointment(t, &employee.ID)

	arg := BulkAddAppointmentParticipantsParams{
		AppointmentID: appointment.ID,
		EmployeeIds:   []int64{employee.ID, employee2.ID},
	}

	err := testQueries.BulkAddAppointmentParticipants(context.Background(), arg)
	require.NoError(t, err)
}

func TestAddAppointmentClient(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	client := createRandomClientDetails(t)
	appointment := createRandomAppointment(t, &employee.ID)

	arg := BulkAddAppointmentClientsParams{
		AppointmentID: appointment.ID,
		ClientIds:     []int64{client.ID},
	}

	err := testQueries.BulkAddAppointmentClients(context.Background(), arg)
	require.NoError(t, err)
}
