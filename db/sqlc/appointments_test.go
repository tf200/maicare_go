package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomAppointment(t *testing.T, employeeID int64) Appointment {
	arg := CreateAppointmentParams{
		CreatorEmployeeID:  employeeID,
		StartTime:          pgtype.Timestamp{Time: time.Now(), Valid: true},
		EndTime:            pgtype.Timestamp{Time: time.Now().Add(1 * time.Hour), Valid: true},
		Location:           util.StringPtr("Test Location"),
		Description:        util.StringPtr("Test Description"),
		Status:             "Scheduled",
		RecurrenceType:     util.StringPtr("NONE"),
		RecurrenceInterval: util.Int32Ptr(0),
		RecurrenceEndDate:  pgtype.Date{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
	}

	appointment, err := testQueries.CreateAppointment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointment)
	return appointment
}

func TestCreateAppointment(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	createRandomAppointment(t, employee.ID)
}

func TestAddAppointmentParticipant(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	employee2, _ := createRandomEmployee(t)
	appointment := createRandomAppointment(t, employee.ID)

	arg := AddAppointmentParticipantParams{
		AppointmentID: appointment.ID,
		EmployeeID:    employee2.ID,
	}

	err := testQueries.AddAppointmentParticipant(context.Background(), arg)
	require.NoError(t, err)
}

func TestAddAppointmentClient(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	client := createRandomClientDetails(t)
	appointment := createRandomAppointment(t, employee.ID)

	arg := AddAppointmentClientParams{
		AppointmentID: appointment.ID,
		ClientID:      client.ID,
	}

	err := testQueries.AddAppointmentClient(context.Background(), arg)
	require.NoError(t, err)
}

func TestListAppointmentsForEmployeeInRange(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	appointment1 := createRandomAppointment(t, employee.ID)
	appointment2 := createRandomAppointment(t, employee.ID)

	startTime := appointment1.StartTime.Time.Add(-1 * time.Hour)
	endTime := appointment2.EndTime.Time.Add(1 * time.Hour)

	arg := ListAppointmentsForEmployeeInRangeParams{
		EmployeeID: employee.ID,
		StartDate:  pgtype.Timestamp{Time: startTime, Valid: true},
		EndDate:    pgtype.Timestamp{Time: endTime, Valid: true},
	}

	appointments, err := testQueries.ListAppointmentsForEmployeeInRange(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointments)
}

func TestListAppointmentsForClientInRange(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	appointment1 := createRandomAppointment(t, employee.ID)
	appointment2 := createRandomAppointment(t, employee.ID)

	err := testQueries.AddAppointmentClient(context.Background(), AddAppointmentClientParams{
		AppointmentID: appointment1.ID,
		ClientID:      client.ID,
	})
	require.NoError(t, err)
	err = testQueries.AddAppointmentClient(context.Background(), AddAppointmentClientParams{
		AppointmentID: appointment2.ID,
		ClientID:      client.ID,
	})
	require.NoError(t, err)

	startTime := appointment1.StartTime.Time.Add(-1 * time.Hour)
	endTime := appointment2.EndTime.Time.Add(1 * time.Hour)

	arg := ListAppointmentsForClientInRangeParams{
		ClientID:  client.ID,
		StartDate: pgtype.Timestamp{Time: startTime, Valid: true},
		EndDate:   pgtype.Timestamp{Time: endTime, Valid: true},
	}

	appointments, err := testQueries.ListAppointmentsForClientInRange(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointments)
}
