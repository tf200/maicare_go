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

func TestAssignEmployee(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AssignEmployeeParams
		checks func(t *testing.T, row AssignEmployeeRow, params AssignEmployeeParams, err error)
	}{
		{
			name: "successful assignment",
			setup: func(ctx context.Context, qtx *Queries) AssignEmployeeParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				client := createRandomClientDetails(ctx, qtx)
				return AssignEmployeeParams{
					ClientID:   client.ID,
					EmployeeID: employee.ID,
					StartDate:  pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Role:       "Manager",
				}
			},
			checks: func(t *testing.T, row AssignEmployeeRow, params AssignEmployeeParams, err error) {
				require.NoError(t, err, "AssignEmployee should not return an error")
				require.NotZero(t, row.ID, "Assigned employee ID should not be zero")
				require.Equal(t, row.ClientID, params.ClientID, "ClientID should match")
				require.Equal(t, row.EmployeeID, params.EmployeeID, "EmployeeID should match")
				require.Equal(t, row.Role, params.Role, "Role should match")
				require.True(t, row.StartDate.Valid, "StartDate should be valid")
				require.Equal(t, row.StartDate.Time, params.StartDate.Time, "StartDate should match")
			},
		},
		{
			name: "assignment with different role",
			setup: func(ctx context.Context, qtx *Queries) AssignEmployeeParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				client := createRandomClientDetails(ctx, qtx)
				return AssignEmployeeParams{
					ClientID:   client.ID,
					EmployeeID: employee.ID,
					StartDate:  pgtype.Date{Time: time.Now(), Valid: true},
					Role:       "Supervisor",
				}
			},
			checks: func(t *testing.T, row AssignEmployeeRow, params AssignEmployeeParams, err error) {
				require.NoError(t, err, "AssignEmployee should not return an error")
				require.Equal(t, "Supervisor", row.Role, "Role should match")
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
			row, err := qtx.AssignEmployee(ctx, params)
			tt.checks(t, row, params, err)
		})
	}
}

func TestAssignSender(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AssignSenderParams
		checks func(t *testing.T, row ClientDetail, params AssignSenderParams, err error)
	}{
		{
			name: "successful sender assignment",
			setup: func(ctx context.Context, qtx *Queries) AssignSenderParams {
				client := createRandomClientDetails(ctx, qtx)
				sender := createRandomSenders(ctx, qtx)
				return AssignSenderParams{
					ID:       client.ID,
					SenderID: &sender.ID,
				}
			},
			checks: func(t *testing.T, row ClientDetail, params AssignSenderParams, err error) {
				require.NoError(t, err, "AssignSender should not return an error")
				require.NotNil(t, row.SenderID, "SenderID should not be nil")
				require.Equal(t, params.SenderID, row.SenderID, "SenderID should match")
			},
		},
		{
			name: "successful unassignment of sender",
			setup: func(ctx context.Context, qtx *Queries) AssignSenderParams {
				client := createRandomClientDetails(ctx, qtx)
				return AssignSenderParams{
					ID:       client.ID,
					SenderID: nil,
				}
			},
			checks: func(t *testing.T, row ClientDetail, params AssignSenderParams, err error) {
				require.NoError(t, err, "AssignSender should not return an error")
				require.Nil(t, row.SenderID, "SenderID should be nil")
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
			row, err := qtx.AssignSender(ctx, params)
			tt.checks(t, row, params, err)
		})
	}
}

func TestCreateEmergencyContact(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateEmemrgencyContactParams
		checks func(t *testing.T, contact ClientEmergencyContact, err error)
	}{
		{
			name: "successful creation with minimal fields",
			setup: func(ctx context.Context, qtx *Queries) CreateEmemrgencyContactParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateEmemrgencyContactParams{
					ClientID:         client.ID,
					FirstName:        util.StringPtr(util.RandomString(5)),
					LastName:         util.StringPtr(util.RandomString(5)),
					Email:            util.StringPtr(util.RandomEmail()),
					PhoneNumber:      util.StringPtr(util.RandomString(8)),
					MedicalReports:   true,
					IncidentsReports: false,
					GoalsReports:     false,
				}
			},
			checks: func(t *testing.T, contact ClientEmergencyContact, err error) {
				require.NoError(t, err, "CreateEmemrgencyContact should not return an error")
				require.NotZero(t, contact.ID)
				require.NotNil(t, contact.FirstName)
				require.NotNil(t, contact.Email)
			},
		},
		{
			name: "successful creation with all fields",
			setup: func(ctx context.Context, qtx *Queries) CreateEmemrgencyContactParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateEmemrgencyContactParams{
					ClientID:         client.ID,
					FirstName:        util.StringPtr(util.RandomString(5)),
					LastName:         util.StringPtr(util.RandomString(5)),
					Email:            util.StringPtr(util.RandomEmail()),
					PhoneNumber:      util.StringPtr(util.RandomString(9)),
					Address:          util.StringPtr(util.RandomString(10)),
					Relationship:     util.StringPtr("Family"),
					MedicalReports:   true,
					IncidentsReports: true,
					GoalsReports:     true,
				}
			},
			checks: func(t *testing.T, contact ClientEmergencyContact, err error) {
				require.NoError(t, err, "CreateEmemrgencyContact should not return an error")
				require.NotZero(t, contact.ID)
				require.NotNil(t, contact.Address)
				require.True(t, contact.GoalsReports)
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
			contact, err := qtx.CreateEmemrgencyContact(ctx, params)
			tt.checks(t, contact, err)
		})
	}
}

func TestGetEmergencyContact(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, contact ClientEmergencyContact, err error)
	}{
		{
			name: "get existing contact",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				contact := createRandomEmergencyContact(ctx, qtx, client.ID)
				return contact.ID
			},
			checks: func(t *testing.T, contact ClientEmergencyContact, err error) {
				require.NoError(t, err, "GetEmergencyContact should not error")
				require.NotZero(t, contact.ID)
			},
		},
		{
			name:  "get non-existent contact",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID { return uuid.New() },
			checks: func(t *testing.T, contact ClientEmergencyContact, err error) {
				require.Error(t, err, "GetEmergencyContact should error for non-existent ID")
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

			id := tt.setup(ctx, qtx)
			contact, err := qtx.GetEmergencyContact(ctx, id)
			tt.checks(t, contact, err)
		})
	}
}

func TestUpdateEmergencyContact(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateEmergencyContactParams
		checks func(t *testing.T, contact ClientEmergencyContact, params UpdateEmergencyContactParams, err error)
	}{
		{
			name: "update existing contact",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmergencyContactParams {
				client := createRandomClientDetails(ctx, qtx)
				contact := createRandomEmergencyContact(ctx, qtx, client.ID)
				newEmail := util.RandomEmail()
				return UpdateEmergencyContactParams{
					ID:    contact.ID,
					Email: util.StringPtr(newEmail),
				}
			},
			checks: func(t *testing.T, contact ClientEmergencyContact, params UpdateEmergencyContactParams, err error) {
				require.NoError(t, err, "UpdateEmergencyContact should not error")
				require.NotNil(t, contact.Email)
				require.Equal(t, *params.Email, *contact.Email)
			},
		},
		{
			name: "update non-existent contact",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmergencyContactParams {
				return UpdateEmergencyContactParams{
					ID:    uuid.New(),
					Email: util.StringPtr(util.RandomEmail()),
				}
			},
			checks: func(t *testing.T, contact ClientEmergencyContact, params UpdateEmergencyContactParams, err error) {
				require.Error(t, err, "UpdateEmergencyContact should error for non-existent ID")
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
			contact, err := qtx.UpdateEmergencyContact(ctx, params)
			tt.checks(t, contact, params, err)
		})
	}
}

func TestDeleteEmergencyContact(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, contact ClientEmergencyContact, err error)
	}{
		{
			name: "delete existing contact",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				contact := createRandomEmergencyContact(ctx, qtx, client.ID)
				return contact.ID
			},
			checks: func(t *testing.T, contact ClientEmergencyContact, err error) {
				require.NoError(t, err, "DeleteEmergencyContact should not error")
				require.NotZero(t, contact.ID)
			},
		},
		{
			name:  "delete non-existent contact",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID { return uuid.New() },
			checks: func(t *testing.T, contact ClientEmergencyContact, err error) {
				require.Error(t, err, "DeleteEmergencyContact should error for non-existent ID")
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

			id := tt.setup(ctx, qtx)
			contact, err := qtx.DeleteEmergencyContact(ctx, id)
			tt.checks(t, contact, err)
		})
	}
}

func TestListEmergencyContacts(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListEmergencyContactsParams
		checks func(t *testing.T, contacts []ListEmergencyContactsRow, err error)
	}{
		{
			name: "list contacts with results",
			setup: func(ctx context.Context, qtx *Queries) ListEmergencyContactsParams {
				client := createRandomClientDetails(ctx, qtx)
				_ = createRandomEmergencyContact(ctx, qtx, client.ID)
				_ = createRandomEmergencyContact(ctx, qtx, client.ID)
				return ListEmergencyContactsParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
					Search:   "",
				}
			},
			checks: func(t *testing.T, contacts []ListEmergencyContactsRow, err error) {
				require.NoError(t, err, "ListEmergencyContacts should not error")
				require.GreaterOrEqual(t, len(contacts), 2)
			},
		},
		{
			name: "list contacts with search",
			setup: func(ctx context.Context, qtx *Queries) ListEmergencyContactsParams {
				client := createRandomClientDetails(ctx, qtx)
				contact := createRandomEmergencyContact(ctx, qtx, client.ID)
				return ListEmergencyContactsParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
					Search:   *contact.FirstName,
				}
			},
			checks: func(t *testing.T, contacts []ListEmergencyContactsRow, err error) {
				require.NoError(t, err, "ListEmergencyContacts should not error")
				require.GreaterOrEqual(t, len(contacts), 1)
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
			contacts, err := qtx.ListEmergencyContacts(ctx, params)
			tt.checks(t, contacts, err)
		})
	}
}

func TestGetAssignedEmployee(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, row GetAssignedEmployeeRow, err error)
	}{
		{
			name: "get existing assigned employee",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				assigned := createRandomAssignedEmployee(ctx, qtx, client.ID)
				return assigned.ID
			},
			checks: func(t *testing.T, row GetAssignedEmployeeRow, err error) {
				require.NoError(t, err, "GetAssignedEmployee should not error")
				require.NotZero(t, row.ID)
			},
		},
		{
			name:  "get non-existent assigned employee",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID { return uuid.New() },
			checks: func(t *testing.T, row GetAssignedEmployeeRow, err error) {
				require.Error(t, err, "GetAssignedEmployee should error for non-existent ID")
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

			id := tt.setup(ctx, qtx)
			row, err := qtx.GetAssignedEmployee(ctx, id)
			tt.checks(t, row, err)
		})
	}
}

func TestUpdateAssignedEmployee(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateAssignedEmployeeParams
		checks func(t *testing.T, assigned AssignedEmployee, err error)
	}{
		{
			name: "update existing assigned employee",
			setup: func(ctx context.Context, qtx *Queries) UpdateAssignedEmployeeParams {
				client := createRandomClientDetails(ctx, qtx)
				assigned := createRandomAssignedEmployee(ctx, qtx, client.ID)
				return UpdateAssignedEmployeeParams{
					ID:   assigned.ID,
					Role: util.StringPtr("Updated Role"),
				}
			},
			checks: func(t *testing.T, assigned AssignedEmployee, err error) {
				require.NoError(t, err, "UpdateAssignedEmployee should not error")
				require.Equal(t, "Updated Role", assigned.Role)
			},
		},
		{
			name: "update non-existent assigned employee",
			setup: func(ctx context.Context, qtx *Queries) UpdateAssignedEmployeeParams {
				return UpdateAssignedEmployeeParams{
					ID:   uuid.New(),
					Role: util.StringPtr("Role"),
				}
			},
			checks: func(t *testing.T, assigned AssignedEmployee, err error) {
				require.Error(t, err, "UpdateAssignedEmployee should error for non-existent ID")
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
			assigned, err := qtx.UpdateAssignedEmployee(ctx, params)
			tt.checks(t, assigned, err)
		})
	}
}

func TestDeleteAssignedEmployee(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, assigned AssignedEmployee, err error)
	}{
		{
			name: "delete existing assigned employee",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				assigned := createRandomAssignedEmployee(ctx, qtx, client.ID)
				return assigned.ID
			},
			checks: func(t *testing.T, assigned AssignedEmployee, err error) {
				require.NoError(t, err, "DeleteAssignedEmployee should not error")
				require.NotZero(t, assigned.ID)
			},
		},
		{
			name:  "delete non-existent assigned employee",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID { return uuid.New() },
			checks: func(t *testing.T, assigned AssignedEmployee, err error) {
				require.Error(t, err, "DeleteAssignedEmployee should error for non-existent ID")
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

			id := tt.setup(ctx, qtx)
			assigned, err := qtx.DeleteAssignedEmployee(ctx, id)
			tt.checks(t, assigned, err)
		})
	}
}

func TestListAssignedEmployees(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListAssignedEmployeesParams
		checks func(t *testing.T, list []ListAssignedEmployeesRow, err error)
	}{
		{
			name: "list assigned employees with results",
			setup: func(ctx context.Context, qtx *Queries) ListAssignedEmployeesParams {
				client := createRandomClientDetails(ctx, qtx)
				_ = createRandomAssignedEmployee(ctx, qtx, client.ID)
				_ = createRandomAssignedEmployee(ctx, qtx, client.ID)
				return ListAssignedEmployeesParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, list []ListAssignedEmployeesRow, err error) {
				require.NoError(t, err, "ListAssignedEmployees should not error")
				require.GreaterOrEqual(t, len(list), 2)
			},
		},
		{
			name: "list assigned employees empty",
			setup: func(ctx context.Context, qtx *Queries) ListAssignedEmployeesParams {
				client := createRandomClientDetails(ctx, qtx)
				return ListAssignedEmployeesParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, list []ListAssignedEmployeesRow, err error) {
				require.NoError(t, err, "ListAssignedEmployees should not error")
				require.Empty(t, list)
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
			list, err := qtx.ListAssignedEmployees(ctx, params)
			tt.checks(t, list, err)
		})
	}
}

func TestGetClientRelatedEmails(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, emails []string, err error)
	}{
		{
			name: "get related emails with data",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				employee := createRandomEmployeeProfile(ctx, qtx)
				_, _ = qtx.AssignEmployee(ctx, AssignEmployeeParams{ClientID: client.ID, EmployeeID: employee.ID, StartDate: pgtype.Date{Time: time.Now(), Valid: true}, Role: "Primary"})
				_, _ = qtx.CreateEmemrgencyContact(ctx, CreateEmemrgencyContactParams{ClientID: client.ID, FirstName: util.StringPtr("Contact"), LastName: util.StringPtr("One"), Email: util.StringPtr(util.RandomEmail()), MedicalReports: true, IncidentsReports: false, GoalsReports: false})
				return client.ID
			},
			checks: func(t *testing.T, emails []string, err error) {
				require.NoError(t, err, "GetClientRelatedEmails should not error")
				require.GreaterOrEqual(t, len(emails), 2)
			},
		},
		{
			name: "get related emails empty",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				return client.ID
			},
			checks: func(t *testing.T, emails []string, err error) {
				require.NoError(t, err, "GetClientRelatedEmails should not error")
				require.Empty(t, emails)
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

			clientID := tt.setup(ctx, qtx)
			emails, err := qtx.GetClientRelatedEmails(ctx, clientID)
			tt.checks(t, emails, err)
		})
	}
}

func TestGetClientSender(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, sender Sender, err error)
	}{
		{
			name: "get client sender",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				sender := createRandomSenders(ctx, qtx)
				_, _ = qtx.AssignSender(ctx, AssignSenderParams{ID: client.ID, SenderID: &sender.ID})
				return client.ID
			},
			checks: func(t *testing.T, sender Sender, err error) {
				require.NoError(t, err, "GetClientSender should not error")
				require.NotZero(t, sender.ID)
			},
		},
		{
			name: "get client sender no sender",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				return client.ID
			},
			checks: func(t *testing.T, sender Sender, err error) {
				require.Error(t, err, "GetClientSender should error when no sender assigned")
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

			clientID := tt.setup(ctx, qtx)
			sender, err := qtx.GetClientSender(ctx, clientID)
			tt.checks(t, sender, err)
		})
	}
}

// Helpers
func createRandomEmergencyContact(ctx context.Context, qtx *Queries, clientID uuid.UUID) ClientEmergencyContact {

	params := CreateEmemrgencyContactParams{
		ClientID:         clientID,
		FirstName:        util.StringPtr(util.RandomString(5)),
		LastName:         util.StringPtr(util.RandomString(5)),
		Email:            util.StringPtr(util.RandomEmail()),
		PhoneNumber:      util.StringPtr(util.RandomString(9)),
		Address:          util.StringPtr(util.RandomString(10)),
		Relationship:     util.StringPtr("Friend"),
		MedicalReports:   util.RandomBool(),
		IncidentsReports: util.RandomBool(),
		GoalsReports:     util.RandomBool(),
	}
	contact, err := qtx.CreateEmemrgencyContact(ctx, params)
	if err != nil {
		panic(err)
	}
	return contact
}

func createRandomAssignedEmployee(ctx context.Context, qtx *Queries, clientID uuid.UUID) AssignEmployeeRow {
	employee := createRandomEmployeeProfile(ctx, qtx)
	params := AssignEmployeeParams{
		ClientID:   clientID,
		EmployeeID: employee.ID,
		StartDate:  pgtype.Date{Time: time.Now(), Valid: true},
		Role:       "Random Role",
	}
	row, err := qtx.AssignEmployee(ctx, params)
	if err != nil {
		panic(err)
	}
	return row
}
