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

func TestCreateAppointment(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateAppointmentParams
		checks func(t *testing.T, appointment ScheduledAppointment, err error)
	}{
		{
			name: "successful appointment creation",
			setup: func(ctx context.Context, qtx *Queries) CreateAppointmentParams {
				employee := createRandomEmployee(ctx, qtx)
				startTime := time.Now().Add(time.Hour)
				endTime := startTime.Add(time.Hour)
				return CreateAppointmentParams{
					CreatorEmployeeID: &employee.ID,
					StartTime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndTime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
					Location:    randomStringPtr(10),
					Color:       randomStringPtr(7),
					Description: randomStringPtr(20),
				}
			},
			checks: func(t *testing.T, appointment ScheduledAppointment, err error) {
				require.NoError(t, err, "CreateAppointment should not return an error")
				require.NotZero(t, appointment.ID)
				require.Equal(t, "SCHEDULED", string(appointment.Status))
				require.False(t, appointment.IsConfirmed)
			},
		},
		{
			name: "appointment creation without creator",
			setup: func(ctx context.Context, qtx *Queries) CreateAppointmentParams {
				startTime := time.Now().Add(time.Hour)
				endTime := startTime.Add(time.Hour)
				return CreateAppointmentParams{
					CreatorEmployeeID: nil,
					StartTime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndTime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
					Location:    randomStringPtr(10),
					Color:       randomStringPtr(7),
					Description: randomStringPtr(20),
				}
			},
			checks: func(t *testing.T, appointment ScheduledAppointment, err error) {
				require.NoError(t, err, "CreateAppointment should not return an error")
				require.Nil(t, appointment.CreatorEmployeeID)
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
			appointment, err := qtx.CreateAppointment(ctx, params)
			tt.checks(t, appointment, err)
		})
	}
}

func TestCreateAppointmentTemplate(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateAppointmentTemplateParams
		checks func(t *testing.T, template AppointmentTemplate, err error)
	}{
		{
			name: "successful template creation",
			setup: func(ctx context.Context, qtx *Queries) CreateAppointmentTemplateParams {
				employee := createRandomEmployee(ctx, qtx)
				startTime := time.Now().Add(time.Hour)
				endTime := startTime.Add(time.Hour)
				return CreateAppointmentTemplateParams{
					CreatorEmployeeID: employee.ID,
					StartTime: pgtype.Timestamp{
						Time:  startTime,
						Valid: true,
					},
					EndTime: pgtype.Timestamp{
						Time:  endTime,
						Valid: true,
					},
					Location:           randomStringPtr(10),
					Description:        randomStringPtr(20),
					Color:              randomStringPtr(7),
					RecurrenceType:     RecurrenceTypeEnum("DAILY"),
					RecurrenceInterval: randomInt32Ptr(1, 30),
					RecurrenceEndDate: pgtype.Date{
						Time:  time.Now().Add(30 * 24 * time.Hour),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, template AppointmentTemplate, err error) {
				require.NoError(t, err, "CreateAppointmentTemplate should not return an error")
				require.NotZero(t, template.ID)
				require.Equal(t, RecurrenceTypeEnum("DAILY"), template.RecurrenceType)
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
			template, err := qtx.CreateAppointmentTemplate(ctx, params)
			tt.checks(t, template, err)
		})
	}
}

func TestGetScheduledAppointmentByID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, appointment GetScheduledAppointmentByIDRow, err error)
	}{
		{
			name: "get existing appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				appointment := createRandomAppointment(ctx, qtx)
				return appointment.ID
			},
			checks: func(t *testing.T, appointment GetScheduledAppointmentByIDRow, err error) {
				require.NoError(t, err, "GetScheduledAppointmentByID should not return an error")
				require.NotZero(t, appointment.ID)
			},
		},
		{
			name: "get non-existent appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, appointment GetScheduledAppointmentByIDRow, err error) {
				require.Error(t, err, "GetScheduledAppointmentByID should error for non-existent appointment")
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

			appointmentID := tt.setup(ctx, qtx)
			appointment, err := qtx.GetScheduledAppointmentByID(ctx, appointmentID)
			tt.checks(t, appointment, err)
		})
	}
}

func TestGetAppointmentTemplate(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, template AppointmentTemplate, err error)
	}{
		{
			name: "get existing template",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				template := createRandomAppointmentTemplate(ctx, qtx)
				return template.ID
			},
			checks: func(t *testing.T, template AppointmentTemplate, err error) {
				require.NoError(t, err, "GetAppointmentTemplate should not return an error")
				require.NotZero(t, template.ID)
			},
		},
		{
			name: "get non-existent template",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, template AppointmentTemplate, err error) {
				require.Error(t, err, "GetAppointmentTemplate should error for non-existent template")
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

			templateID := tt.setup(ctx, qtx)
			template, err := qtx.GetAppointmentTemplate(ctx, templateID)
			tt.checks(t, template, err)
		})
	}
}

func TestUpdateAppointment(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateAppointmentParams
		checks func(t *testing.T, appointment ScheduledAppointment, err error)
	}{
		{
			name: "update existing appointment",
			setup: func(ctx context.Context, qtx *Queries) UpdateAppointmentParams {
				appointment := createRandomAppointment(ctx, qtx)
				newStartTime := time.Now().Add(2 * time.Hour)
				newEndTime := newStartTime.Add(time.Hour)
				return UpdateAppointmentParams{
					ID: appointment.ID,
					StartTime: pgtype.Timestamp{
						Time:  newStartTime,
						Valid: true,
					},
					EndTime: pgtype.Timestamp{
						Time:  newEndTime,
						Valid: true,
					},
					Location:    randomStringPtr(10),
					Description: randomStringPtr(20),
					Color:       randomStringPtr(7),
				}
			},
			checks: func(t *testing.T, appointment ScheduledAppointment, err error) {
				require.NoError(t, err, "UpdateAppointment should not return an error")
				require.NotZero(t, appointment.ID)
			},
		},
		{
			name: "update non-existent appointment",
			setup: func(ctx context.Context, qtx *Queries) UpdateAppointmentParams {
				newStartTime := time.Now().Add(2 * time.Hour)
				newEndTime := newStartTime.Add(time.Hour)
				return UpdateAppointmentParams{
					ID: uuid.New(),
					StartTime: pgtype.Timestamp{
						Time:  newStartTime,
						Valid: true,
					},
					EndTime: pgtype.Timestamp{
						Time:  newEndTime,
						Valid: true,
					},
					Location:    randomStringPtr(10),
					Description: randomStringPtr(20),
					Color:       randomStringPtr(7),
				}
			},
			checks: func(t *testing.T, appointment ScheduledAppointment, err error) {
				require.Error(t, err, "UpdateAppointment should error for non-existent appointment")
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
			appointment, err := qtx.UpdateAppointment(ctx, params)
			tt.checks(t, appointment, err)
		})
	}
}

func TestDeleteAppointment(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete existing appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				appointment := createRandomAppointment(ctx, qtx)
				return appointment.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteAppointment should not error")
			},
		},
		{
			name: "delete non-existent appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteAppointment should not error for non-existent appointment")
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

			appointmentID := tt.setup(ctx, qtx)
			err = qtx.DeleteAppointment(ctx, appointmentID)
			tt.checks(t, err)
		})
	}
}

func TestConfirmAppointment(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ConfirmAppointmentParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "confirm existing appointment",
			setup: func(ctx context.Context, qtx *Queries) ConfirmAppointmentParams {
				appointment := createRandomAppointment(ctx, qtx)
				employee := createRandomEmployee(ctx, qtx)
				return ConfirmAppointmentParams{
					ID:         appointment.ID,
					EmployeeID: &employee.ID,
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "ConfirmAppointment should not error")
			},
		},
		{
			name: "confirm appointment without employee",
			setup: func(ctx context.Context, qtx *Queries) ConfirmAppointmentParams {
				appointment := createRandomAppointment(ctx, qtx)
				return ConfirmAppointmentParams{
					ID:         appointment.ID,
					EmployeeID: nil,
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "ConfirmAppointment should not error")
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
			err = qtx.ConfirmAppointment(ctx, params)
			tt.checks(t, err)
		})
	}
}

func TestBulkAddAppointmentClients(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) BulkAddAppointmentClientsParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "bulk add clients to appointment",
			setup: func(ctx context.Context, qtx *Queries) BulkAddAppointmentClientsParams {
				appointment := createRandomAppointment(ctx, qtx)
				client1 := createRandomClientDetails(ctx, qtx)
				client2 := createRandomClientDetails(ctx, qtx)
				return BulkAddAppointmentClientsParams{
					AppointmentID: appointment.ID,
					ClientIds:     []uuid.UUID{client1.ID, client2.ID},
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "BulkAddAppointmentClients should not error")
			},
		},
		{
			name: "bulk add empty clients list",
			setup: func(ctx context.Context, qtx *Queries) BulkAddAppointmentClientsParams {
				appointment := createRandomAppointment(ctx, qtx)
				return BulkAddAppointmentClientsParams{
					AppointmentID: appointment.ID,
					ClientIds:     []uuid.UUID{},
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "BulkAddAppointmentClients should not error with empty list")
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
			err = qtx.BulkAddAppointmentClients(ctx, params)
			tt.checks(t, err)
		})
	}
}

func TestBulkAddAppointmentParticipants(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) BulkAddAppointmentParticipantsParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "bulk add participants to appointment",
			setup: func(ctx context.Context, qtx *Queries) BulkAddAppointmentParticipantsParams {
				appointment := createRandomAppointment(ctx, qtx)
				employee1 := createRandomEmployee(ctx, qtx)
				employee2 := createRandomEmployee(ctx, qtx)
				return BulkAddAppointmentParticipantsParams{
					AppointmentID: appointment.ID,
					EmployeeIds:   []uuid.UUID{employee1.ID, employee2.ID},
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "BulkAddAppointmentParticipants should not error")
			},
		},
		{
			name: "bulk add empty participants list",
			setup: func(ctx context.Context, qtx *Queries) BulkAddAppointmentParticipantsParams {
				appointment := createRandomAppointment(ctx, qtx)
				return BulkAddAppointmentParticipantsParams{
					AppointmentID: appointment.ID,
					EmployeeIds:   []uuid.UUID{},
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "BulkAddAppointmentParticipants should not error with empty list")
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
			err = qtx.BulkAddAppointmentParticipants(ctx, params)
			tt.checks(t, err)
		})
	}
}

func TestGetAppointmentClients(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) []uuid.UUID
		checks func(t *testing.T, clients []GetAppointmentClientsRow, err error)
	}{
		{
			name: "get clients for existing appointment",
			setup: func(ctx context.Context, qtx *Queries) []uuid.UUID {
				appointment := createRandomAppointment(ctx, qtx)
				client1 := createRandomClientDetails(ctx, qtx)
				client2 := createRandomClientDetails(ctx, qtx)
				qtx.BulkAddAppointmentClients(ctx, BulkAddAppointmentClientsParams{
					AppointmentID: appointment.ID,
					ClientIds:     []uuid.UUID{client1.ID, client2.ID},
				})
				return []uuid.UUID{appointment.ID}
			},
			checks: func(t *testing.T, clients []GetAppointmentClientsRow, err error) {
				require.NoError(t, err, "GetAppointmentClients should not error")
				require.Len(t, clients, 2)
			},
		},
		{
			name: "get clients for non-existent appointment",
			setup: func(ctx context.Context, qtx *Queries) []uuid.UUID {
				return []uuid.UUID{uuid.New()}
			},
			checks: func(t *testing.T, clients []GetAppointmentClientsRow, err error) {
				require.NoError(t, err, "GetAppointmentClients should not error for non-existent appointment")
				require.Len(t, clients, 0)
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

			appointmentIDs := tt.setup(ctx, qtx)
			clients, err := qtx.GetAppointmentClients(ctx, appointmentIDs)
			tt.checks(t, clients, err)
		})
	}
}

func TestGetAppointmentParticipants(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) []uuid.UUID
		checks func(t *testing.T, participants []GetAppointmentParticipantsRow, err error)
	}{
		{
			name: "get participants for existing appointment",
			setup: func(ctx context.Context, qtx *Queries) []uuid.UUID {
				appointment := createRandomAppointment(ctx, qtx)
				employee1 := createRandomEmployee(ctx, qtx)
				employee2 := createRandomEmployee(ctx, qtx)
				qtx.BulkAddAppointmentParticipants(ctx, BulkAddAppointmentParticipantsParams{
					AppointmentID: appointment.ID,
					EmployeeIds:   []uuid.UUID{employee1.ID, employee2.ID},
				})
				return []uuid.UUID{appointment.ID}
			},
			checks: func(t *testing.T, participants []GetAppointmentParticipantsRow, err error) {
				require.NoError(t, err, "GetAppointmentParticipants should not error")
				require.Len(t, participants, 2)
			},
		},
		{
			name: "get participants for non-existent appointment",
			setup: func(ctx context.Context, qtx *Queries) []uuid.UUID {
				return []uuid.UUID{uuid.New()}
			},
			checks: func(t *testing.T, participants []GetAppointmentParticipantsRow, err error) {
				require.NoError(t, err, "GetAppointmentParticipants should not error for non-existent appointment")
				require.Len(t, participants, 0)
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

			appointmentIDs := tt.setup(ctx, qtx)
			participants, err := qtx.GetAppointmentParticipants(ctx, appointmentIDs)
			tt.checks(t, participants, err)
		})
	}
}

func TestDeleteAppointmentClients(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete clients from existing appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				appointment := createRandomAppointment(ctx, qtx)
				client := createRandomClientDetails(ctx, qtx)
				qtx.BulkAddAppointmentClients(ctx, BulkAddAppointmentClientsParams{
					AppointmentID: appointment.ID,
					ClientIds:     []uuid.UUID{client.ID},
				})
				return appointment.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteAppointmentClients should not error")
			},
		},
		{
			name: "delete clients from non-existent appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteAppointmentClients should not error for non-existent appointment")
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

			appointmentID := tt.setup(ctx, qtx)
			err = qtx.DeleteAppointmentClients(ctx, appointmentID)
			tt.checks(t, err)
		})
	}
}

func TestDeleteAppointmentParticipants(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete participants from existing appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				appointment := createRandomAppointment(ctx, qtx)
				employee := createRandomEmployee(ctx, qtx)
				qtx.BulkAddAppointmentParticipants(ctx, BulkAddAppointmentParticipantsParams{
					AppointmentID: appointment.ID,
					EmployeeIds:   []uuid.UUID{employee.ID},
				})
				return appointment.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteAppointmentParticipants should not error")
			},
		},
		{
			name: "delete participants from non-existent appointment",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteAppointmentParticipants should not error for non-existent appointment")
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

			appointmentID := tt.setup(ctx, qtx)
			err = qtx.DeleteAppointmentParticipants(ctx, appointmentID)
			tt.checks(t, err)
		})
	}
}

// List tests - these are more complex, might need to be simplified or expanded
func TestListClientAppointmentsInRange(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListClientAppointmentsInRangeParams
		checks func(t *testing.T, appointments []ListClientAppointmentsInRangeRow, err error)
	}{
		{
			name: "list appointments for client in range",
			setup: func(ctx context.Context, qtx *Queries) ListClientAppointmentsInRangeParams {
				client := createRandomClientDetails(ctx, qtx)
				appointment := createRandomAppointment(ctx, qtx)
				qtx.BulkAddAppointmentClients(ctx, BulkAddAppointmentClientsParams{
					AppointmentID: appointment.ID,
					ClientIds:     []uuid.UUID{client.ID},
				})
				startDate := time.Now().Add(-time.Hour)
				endDate := time.Now().Add(2 * time.Hour)
				return ListClientAppointmentsInRangeParams{
					ClientID: client.ID,
					StartDate: pgtype.Timestamp{
						Time:  startDate,
						Valid: true,
					},
					EndDate: pgtype.Timestamp{
						Time:  endDate,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, appointments []ListClientAppointmentsInRangeRow, err error) {
				require.NoError(t, err, "ListClientAppointmentsInRange should not error")
				require.GreaterOrEqual(t, len(appointments), 1)
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
			appointments, err := qtx.ListClientAppointmentsInRange(ctx, params)
			tt.checks(t, appointments, err)
		})
	}
}

// Helpers
func createRandomAppointment(ctx context.Context, qtx *Queries) ScheduledAppointment {
	employee := createRandomEmployee(ctx, qtx)
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(time.Hour)
	params := CreateAppointmentParams{
		CreatorEmployeeID: &employee.ID,
		StartTime: pgtype.Timestamp{
			Time:  startTime,
			Valid: true,
		},
		EndTime: pgtype.Timestamp{
			Time:  endTime,
			Valid: true,
		},
		Location:    randomStringPtr(10),
		Color:       randomStringPtr(7),
		Description: randomStringPtr(20),
	}
	appointment, err := qtx.CreateAppointment(ctx, params)
	if err != nil {
		panic(err)
	}
	return appointment
}

func createRandomAppointmentTemplate(ctx context.Context, qtx *Queries) AppointmentTemplate {
	employee := createRandomEmployee(ctx, qtx)
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(time.Hour)
	params := CreateAppointmentTemplateParams{
		CreatorEmployeeID: employee.ID,
		StartTime: pgtype.Timestamp{
			Time:  startTime,
			Valid: true,
		},
		EndTime: pgtype.Timestamp{
			Time:  endTime,
			Valid: true,
		},
		Location:           randomStringPtr(10),
		Description:        randomStringPtr(20),
		Color:              randomStringPtr(7),
		RecurrenceType:     RecurrenceTypeEnum("DAILY"),
		RecurrenceInterval: randomInt32Ptr(1, 30),
		RecurrenceEndDate: pgtype.Date{
			Time:  time.Now().Add(30 * 24 * time.Hour),
			Valid: true,
		},
	}
	template, err := qtx.CreateAppointmentTemplate(ctx, params)
	if err != nil {
		panic(err)
	}
	return template
}

func randomStringPtr(n int) *string {
	s := util.RandomString(n)
	return &s
}

func randomInt32Ptr(min, max int64) *int32 {
	v := int32(util.RandomInt(min, max))
	return &v
}
