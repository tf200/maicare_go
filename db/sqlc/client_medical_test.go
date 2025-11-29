package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

// TestCreateClientDiagnosis tests the CreateClientDiagnosis function
func TestCreateClientDiagnosis(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateClientDiagnosisParams
		checks func(t *testing.T, diagnosis ClientDiagnosis)
	}{
		{
			name: "successful creation with minimal fields",
			setup: func(ctx context.Context, qtx *Queries) CreateClientDiagnosisParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientDiagnosisParams{
					ClientID:            client.ID,
					DiagnosisCode:       "ICD-10-001",
					Description:         "Test diagnosis description",
					Status:              "ACTIVE",
					Title:               nil,
					Severity:            nil,
					DiagnosingClinician: nil,
					Notes:               nil,
				}
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis) {
				require.NotZero(t, diagnosis.ID, "diagnosis ID should be set")
				require.Equal(t, "ICD-10-001", diagnosis.DiagnosisCode)
				require.Equal(t, "Test diagnosis description", diagnosis.Description)
				require.Equal(t, "ACTIVE", diagnosis.Status)
				require.Nil(t, diagnosis.Title)
				require.Nil(t, diagnosis.Severity)
				require.Nil(t, diagnosis.DiagnosingClinician)
				require.Nil(t, diagnosis.Notes)
			},
		},
		{
			name: "successful creation with all fields",
			setup: func(ctx context.Context, qtx *Queries) CreateClientDiagnosisParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientDiagnosisParams{
					ClientID:            client.ID,
					Title:               util.StringPtr("Hypertension"),
					DiagnosisCode:       "ICD-10-I10",
					Description:         "High blood pressure condition",
					Severity:            util.StringPtr("MODERATE"),
					Status:              "ACTIVE",
					DiagnosingClinician: util.StringPtr("Dr. Smith"),
					Notes:               util.StringPtr("Patient requires regular monitoring"),
				}
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis) {
				require.NotZero(t, diagnosis.ID)
				require.NotNil(t, diagnosis.Title)
				require.Equal(t, "Hypertension", *diagnosis.Title)
				require.NotNil(t, diagnosis.Severity)
				require.Equal(t, "MODERATE", *diagnosis.Severity)
				require.NotNil(t, diagnosis.DiagnosingClinician)
				require.Equal(t, "Dr. Smith", *diagnosis.DiagnosingClinician)
				require.NotNil(t, diagnosis.Notes)
				require.Equal(t, "Patient requires regular monitoring", *diagnosis.Notes)
			},
		},
		{
			name: "creation with empty optional fields",
			setup: func(ctx context.Context, qtx *Queries) CreateClientDiagnosisParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateClientDiagnosisParams{
					ClientID:            client.ID,
					DiagnosisCode:       "ICD-10-002",
					Description:         "Another diagnosis",
					Status:              "INACTIVE",
					Title:               util.StringPtr(""),
					Severity:            util.StringPtr(""),
					DiagnosingClinician: util.StringPtr(""),
					Notes:               util.StringPtr(""),
				}
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis) {
				require.NotZero(t, diagnosis.ID)
				require.Equal(t, "ICD-10-002", diagnosis.DiagnosisCode)
				require.Equal(t, "INACTIVE", diagnosis.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			diagnosis, err := qtx.CreateClientDiagnosis(ctx, params)
			require.NoError(t, err, "CreateClientDiagnosis() should not error")

			tt.checks(t, diagnosis)
		})
	}
}

// TestGetClientDiagnosis tests the GetClientDiagnosis function
func TestGetClientDiagnosis(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, diagnosis ClientDiagnosis, err error)
	}{
		{
			name: "get existing diagnosis",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				return diagnosis.ID
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.NoError(t, err, "GetClientDiagnosis() should not error")
				require.NotZero(t, diagnosis.ID)
				require.NotZero(t, diagnosis.ClientID)
				require.NotEmpty(t, diagnosis.DiagnosisCode)
				require.NotEmpty(t, diagnosis.Description)
				require.NotEmpty(t, diagnosis.Status)
			},
		},
		{
			name: "get non-existent diagnosis",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New() // Non-existent ID
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.Error(t, err, "GetClientDiagnosis() should error for non-existent diagnosis")
			},
		},
		{
			name: "get diagnosis with all fields populated",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				client := createRandomClientDetails(ctx, qtx)
				diagnosis, err := qtx.CreateClientDiagnosis(ctx, CreateClientDiagnosisParams{
					ClientID:            client.ID,
					Title:               util.StringPtr("Diabetes"),
					DiagnosisCode:       "ICD-10-E11",
					Description:         "Type 2 diabetes mellitus",
					Severity:            util.StringPtr("HIGH"),
					Status:              "ACTIVE",
					DiagnosingClinician: util.StringPtr("Dr. Johnson"),
					Notes:               util.StringPtr("Requires insulin therapy"),
				})
				require.NoError(t, err)
				return diagnosis.ID
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.NoError(t, err)
				require.NotNil(t, diagnosis.Title)
				require.Equal(t, "Diabetes", *diagnosis.Title)
				require.NotNil(t, diagnosis.Severity)
				require.Equal(t, "HIGH", *diagnosis.Severity)
				require.NotNil(t, diagnosis.DiagnosingClinician)
				require.NotNil(t, diagnosis.Notes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			diagnosis, err := qtx.GetClientDiagnosis(ctx, id)
			tt.checks(t, diagnosis, err)
		})
	}
}

// TestUpdateClientDiagnosis tests the UpdateClientDiagnosis function
func TestUpdateClientDiagnosis(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateClientDiagnosisParams
		checks func(t *testing.T, diagnosis ClientDiagnosis, err error)
	}{
		{
			name: "successful update with new title",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDiagnosisParams {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				return UpdateClientDiagnosisParams{
					ID:    diagnosis.ID,
					Title: util.StringPtr("Updated Title"),
				}
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.NoError(t, err, "UpdateClientDiagnosis() should not error")
				require.Equal(t, "Updated Title", *diagnosis.Title)
			},
		},
		{
			name: "update multiple fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDiagnosisParams {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				return UpdateClientDiagnosisParams{
					ID:          diagnosis.ID,
					Severity:    util.StringPtr("CRITICAL"),
					Status:      util.StringPtr("INACTIVE"),
					Description: util.StringPtr("Updated description"),
				}
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.NoError(t, err)
				require.Equal(t, "CRITICAL", *diagnosis.Severity)
				require.Equal(t, "INACTIVE", diagnosis.Status)
				require.Equal(t, "Updated description", diagnosis.Description)
			},
		},
		{
			name: "update non-existent diagnosis",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDiagnosisParams {
				return UpdateClientDiagnosisParams{
					ID:    uuid.New(),
					Title: util.StringPtr("Non-existent"),
				}
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.Error(t, err, "UpdateClientDiagnosis() should error for non-existent diagnosis")
			},
		},
		{
			name: "update with null values",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientDiagnosisParams {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				return UpdateClientDiagnosisParams{
					ID:       diagnosis.ID,
					Severity: nil,
					Notes:    util.StringPtr(""),
				}
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			diagnosis, err := qtx.UpdateClientDiagnosis(ctx, params)
			tt.checks(t, diagnosis, err)
		})
	}
}

// TestDeleteClientDiagnosis tests the DeleteClientDiagnosis function
func TestDeleteClientDiagnosis(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, diagnosis ClientDiagnosis, err error)
	}{
		{
			name: "successful delete existing diagnosis",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				return diagnosis.ID
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.NoError(t, err, "DeleteClientDiagnosis() should not error")
				require.NotZero(t, diagnosis.ID)
			},
		},
		{
			name: "delete non-existent diagnosis",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, diagnosis ClientDiagnosis, err error) {
				require.Error(t, err, "DeleteClientDiagnosis() should error for non-existent diagnosis")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			diagnosis, err := qtx.DeleteClientDiagnosis(ctx, id)
			tt.checks(t, diagnosis, err)
		})
	}
}

// TestListClientDiagnoses tests the ListClientDiagnoses function
func TestListClientDiagnoses(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListClientDiagnosesParams
		checks func(t *testing.T, diagnoses []ListClientDiagnosesRow, err error)
	}{
		{
			name: "list diagnoses for client with multiple diagnoses",
			setup: func(ctx context.Context, qtx *Queries) ListClientDiagnosesParams {
				client := createRandomClientDetails(ctx, qtx)
				clientID := client.ID
				// Create multiple diagnoses for the same client
				for i := 0; i < 3; i++ {
					_, err := qtx.CreateClientDiagnosis(ctx, CreateClientDiagnosisParams{
						ClientID:      clientID,
						DiagnosisCode: "ICD-10-00" + string(rune('1'+i)),
						Description:   "Diagnosis " + string(rune('1'+i)),
						Status:        "ACTIVE",
					})
					require.NoError(t, err)
				}
				return ListClientDiagnosesParams{
					ClientID: clientID,
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, diagnoses []ListClientDiagnosesRow, err error) {
				require.NoError(t, err, "ListClientDiagnoses() should not error")
				require.Len(t, diagnoses, 3, "should return 3 diagnoses")
				require.Equal(t, int64(3), diagnoses[0].TotalDiagnoses, "total_diagnoses should be 3")
			},
		},
		{
			name: "list diagnoses with pagination",
			setup: func(ctx context.Context, qtx *Queries) ListClientDiagnosesParams {
				client := createRandomClientDetails(ctx, qtx)
				clientID := client.ID
				for i := 0; i < 5; i++ {
					_, err := qtx.CreateClientDiagnosis(ctx, CreateClientDiagnosisParams{
						ClientID:      clientID,
						DiagnosisCode: "ICD-10-00" + string(rune('1'+i)),
						Description:   "Diagnosis " + string(rune('1'+i)),
						Status:        "ACTIVE",
					})
					require.NoError(t, err)
				}
				return ListClientDiagnosesParams{
					ClientID: clientID,
					Limit:    2,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, diagnoses []ListClientDiagnosesRow, err error) {
				require.NoError(t, err)
				require.Len(t, diagnoses, 2, "should return 2 diagnoses with limit 2")
				require.Equal(t, int64(5), diagnoses[0].TotalDiagnoses, "total_diagnoses should still be 5")
			},
		},
		{
			name: "list diagnoses for client with no diagnoses",
			setup: func(ctx context.Context, qtx *Queries) ListClientDiagnosesParams {
				return ListClientDiagnosesParams{
					ClientID: uuid.New(),
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, diagnoses []ListClientDiagnosesRow, err error) {
				require.NoError(t, err)
				require.Empty(t, diagnoses, "should return empty slice for client with no diagnoses")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			diagnoses, err := qtx.ListClientDiagnoses(ctx, params)
			tt.checks(t, diagnoses, err)
		})
	}
}

// TestCreateClientMedication tests the CreateClientMedication function
func TestCreateClientMedication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateClientMedicationParams
		checks func(t *testing.T, medication ClientMedication)
	}{
		{
			name: "successful creation with minimal fields",
			setup: func(ctx context.Context, qtx *Queries) CreateClientMedicationParams {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				return CreateClientMedicationParams{
					DiagnosisID:      &diagnosis.ID,
					Name:             "Aspirin",
					Dosage:           "500mg",
					StartDate:        pgtype.Date{Valid: true},
					EndDate:          pgtype.Date{Valid: true},
					SelfAdministered: true,
					IsCritical:       false,
					Notes:            nil,
					AdministeredByID: nil,
				}
			},
			checks: func(t *testing.T, medication ClientMedication) {
				require.NotZero(t, medication.ID)
				require.Equal(t, "Aspirin", medication.Name)
				require.Equal(t, "500mg", medication.Dosage)
				require.True(t, medication.SelfAdministered)
				require.False(t, medication.IsCritical)
				require.Nil(t, medication.Notes)
				require.Nil(t, medication.AdministeredByID)
			},
		},
		{
			name: "successful creation with all fields",
			setup: func(ctx context.Context, qtx *Queries) CreateClientMedicationParams {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				employee := createRandomEmployeeProfile(ctx, qtx)
				return CreateClientMedicationParams{
					DiagnosisID:      &diagnosis.ID,
					Name:             "Insulin",
					Dosage:           "10 units",
					StartDate:        pgtype.Date{Valid: true},
					EndDate:          pgtype.Date{Valid: true},
					Notes:            util.StringPtr("Inject before meals"),
					SelfAdministered: false,
					AdministeredByID: &employee.ID,
					IsCritical:       true,
				}
			},
			checks: func(t *testing.T, medication ClientMedication) {
				require.NotZero(t, medication.ID)
				require.Equal(t, "Insulin", medication.Name)
				require.Equal(t, "10 units", medication.Dosage)
				require.False(t, medication.SelfAdministered)
				require.True(t, medication.IsCritical)
				require.NotNil(t, medication.Notes)
				require.Equal(t, "Inject before meals", *medication.Notes)
				require.NotNil(t, medication.AdministeredByID)
			},
		},
		{
			name: "creation without diagnosis ID",
			setup: func(ctx context.Context, qtx *Queries) CreateClientMedicationParams {
				return CreateClientMedicationParams{
					DiagnosisID:      nil,
					Name:             "Paracetamol",
					Dosage:           "1000mg",
					StartDate:        pgtype.Date{Valid: true},
					EndDate:          pgtype.Date{Valid: true},
					SelfAdministered: true,
					IsCritical:       false,
				}
			},
			checks: func(t *testing.T, medication ClientMedication) {
				require.NotZero(t, medication.ID)
				require.Nil(t, medication.DiagnosisID)
				require.Equal(t, "Paracetamol", medication.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			medication, err := qtx.CreateClientMedication(ctx, params)
			require.NoError(t, err, "CreateClientMedication() should not error")

			tt.checks(t, medication)
		})
	}
}

// TestGetMedication tests the GetMedication function
func TestGetMedication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, medication GetMedicationRow, err error)
	}{
		{
			name: "get medication with administered_by relationship",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				employee := createRandomEmployeeProfile(ctx, qtx)
				medication, err := qtx.CreateClientMedication(ctx, CreateClientMedicationParams{
					DiagnosisID:      &diagnosis.ID,
					Name:             "Medication with admin",
					Dosage:           "500mg",
					StartDate:        pgtype.Date{Valid: true},
					EndDate:          pgtype.Date{Valid: true},
					SelfAdministered: false,
					AdministeredByID: &employee.ID,
					IsCritical:       false,
				})
				require.NoError(t, err)
				return medication.ID
			},
			checks: func(t *testing.T, medication GetMedicationRow, err error) {
				require.NoError(t, err, "GetMedication() should not error")
				require.NotZero(t, medication.ID)
				require.Equal(t, "Medication with admin", medication.Name)
				require.NotNil(t, medication.AdministeredByID)
				require.NotEmpty(t, medication.AdministeredByFirstName)
				require.NotEmpty(t, medication.AdministeredByLastName)
			},
		},
		{
			name: "get non-existent medication",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, medication GetMedicationRow, err error) {
				require.Error(t, err, "GetMedication() should error for non-existent medication")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			medication, err := qtx.GetMedication(ctx, id)
			tt.checks(t, medication, err)
		})
	}
}

// TestUpdateClientMedication tests the UpdateClientMedication function
func TestUpdateClientMedication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateClientMedicationParams
		checks func(t *testing.T, medication ClientMedication, err error)
	}{
		{
			name: "successful update medication name",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientMedicationParams {
				medication := createRandomMedication(ctx, qtx)
				return UpdateClientMedicationParams{
					ID:   medication.ID,
					Name: util.StringPtr("Updated Medicine Name"),
				}
			},
			checks: func(t *testing.T, medication ClientMedication, err error) {
				require.NoError(t, err, "UpdateClientMedication() should not error")
				require.Equal(t, "Updated Medicine Name", medication.Name)
			},
		},
		{
			name: "update medication dosage and notes",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientMedicationParams {
				medication := createRandomMedication(ctx, qtx)
				return UpdateClientMedicationParams{
					ID:     medication.ID,
					Dosage: util.StringPtr("1000mg"),
					Notes:  util.StringPtr("Take with food"),
				}
			},
			checks: func(t *testing.T, medication ClientMedication, err error) {
				require.NoError(t, err)
				require.Equal(t, "1000mg", medication.Dosage)
				require.NotNil(t, medication.Notes)
				require.Equal(t, "Take with food", *medication.Notes)
			},
		},
		{
			name: "update medication critical flag",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientMedicationParams {
				medication := createRandomMedication(ctx, qtx)
				isCritical := true
				return UpdateClientMedicationParams{
					ID:         medication.ID,
					IsCritical: &isCritical,
				}
			},
			checks: func(t *testing.T, medication ClientMedication, err error) {
				require.NoError(t, err)
				require.True(t, medication.IsCritical)
			},
		},
		{
			name: "update non-existent medication",
			setup: func(ctx context.Context, qtx *Queries) UpdateClientMedicationParams {
				return UpdateClientMedicationParams{
					ID:   uuid.New(),
					Name: util.StringPtr("Non-existent"),
				}
			},
			checks: func(t *testing.T, medication ClientMedication, err error) {
				require.Error(t, err, "UpdateClientMedication() should error for non-existent medication")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			medication, err := qtx.UpdateClientMedication(ctx, params)
			tt.checks(t, medication, err)
		})
	}
}

// TestDeleteClientMedication tests the DeleteClientMedication function
func TestDeleteClientMedication(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "successful delete existing medication",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				medication := createRandomMedication(ctx, qtx)
				return medication.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteClientMedication() should not error")
			},
		},
		{
			name: "delete non-existent medication",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteClientMedication() should not error even for non-existent medication")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			err = qtx.DeleteClientMedication(ctx, id)
			tt.checks(t, err)
		})
	}
}

// TestListMedicationsByDiagnosisID tests the ListMedicationsByDiagnosisID function
func TestListMedicationsByDiagnosisID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListMedicationsByDiagnosisIDParams
		checks func(t *testing.T, medications []ListMedicationsByDiagnosisIDRow, err error)
	}{
		{
			name: "list medications for diagnosis with multiple medications",
			setup: func(ctx context.Context, qtx *Queries) ListMedicationsByDiagnosisIDParams {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				for i := 0; i < 3; i++ {
					_, err := qtx.CreateClientMedication(ctx, CreateClientMedicationParams{
						DiagnosisID:      &diagnosis.ID,
						Name:             util.RandomString(5),
						Dosage:           "500mg",
						StartDate:        pgtype.Date{Valid: true},
						EndDate:          pgtype.Date{Valid: true},
						SelfAdministered: true,
						IsCritical:       false,
					})
					require.NoError(t, err)
				}
				return ListMedicationsByDiagnosisIDParams{
					DiagnosisID: &diagnosis.ID,
					Limit:       10,
					Offset:      0,
				}
			},
			checks: func(t *testing.T, medications []ListMedicationsByDiagnosisIDRow, err error) {
				require.NoError(t, err, "ListMedicationsByDiagnosisID() should not error")
				require.Len(t, medications, 3, "should return 3 medications")
				require.Equal(t, int64(3), medications[0].TotalMedications)
			},
		},
		{
			name: "list medications with pagination",
			setup: func(ctx context.Context, qtx *Queries) ListMedicationsByDiagnosisIDParams {
				diagnosis := createRandomDiagnosis(ctx, qtx)
				for i := 0; i < 5; i++ {
					_, err := qtx.CreateClientMedication(ctx, CreateClientMedicationParams{
						DiagnosisID:      &diagnosis.ID,
						Name:             util.RandomString(5),
						Dosage:           "500mg",
						StartDate:        pgtype.Date{Valid: true},
						EndDate:          pgtype.Date{Valid: true},
						SelfAdministered: true,
						IsCritical:       false,
					})
					require.NoError(t, err)
				}
				return ListMedicationsByDiagnosisIDParams{
					DiagnosisID: &diagnosis.ID,
					Limit:       2,
					Offset:      0,
				}
			},
			checks: func(t *testing.T, medications []ListMedicationsByDiagnosisIDRow, err error) {
				require.NoError(t, err)
				require.Len(t, medications, 2, "should return 2 medications with limit 2")
				require.Equal(t, int64(5), medications[0].TotalMedications)
			},
		},
		{
			name: "list medications for non-existent diagnosis",
			setup: func(ctx context.Context, qtx *Queries) ListMedicationsByDiagnosisIDParams {
				diagnosisID := uuid.New()
				return ListMedicationsByDiagnosisIDParams{
					DiagnosisID: &diagnosisID,
					Limit:       10,
					Offset:      0,
				}
			},
			checks: func(t *testing.T, medications []ListMedicationsByDiagnosisIDRow, err error) {
				require.NoError(t, err)
				require.Empty(t, medications, "should return empty slice for non-existent diagnosis")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			medications, err := qtx.ListMedicationsByDiagnosisID(ctx, params)
			tt.checks(t, medications, err)
		})
	}
}

// TestListMedicationsByDiagnosisIDs tests the ListMedicationsByDiagnosisIDs function
func TestListMedicationsByDiagnosisIDs(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) []uuid.UUID
		checks func(t *testing.T, medications []ClientMedication, err error)
	}{
		{
			name: "list medications for multiple diagnoses",
			setup: func(ctx context.Context, qtx *Queries) []uuid.UUID {
				diagnosis1 := createRandomDiagnosis(ctx, qtx)
				diagnosis2 := createRandomDiagnosis(ctx, qtx)

				_, err := qtx.CreateClientMedication(ctx, CreateClientMedicationParams{
					DiagnosisID:      &diagnosis1.ID,
					Name:             "Medicine 1",
					Dosage:           "500mg",
					StartDate:        pgtype.Date{Valid: true},
					EndDate:          pgtype.Date{Valid: true},
					SelfAdministered: true,
					IsCritical:       false,
				})
				require.NoError(t, err)

				_, err = qtx.CreateClientMedication(ctx, CreateClientMedicationParams{
					DiagnosisID:      &diagnosis2.ID,
					Name:             "Medicine 2",
					Dosage:           "1000mg",
					StartDate:        pgtype.Date{Valid: true},
					EndDate:          pgtype.Date{Valid: true},
					SelfAdministered: true,
					IsCritical:       false,
				})
				require.NoError(t, err)

				return []uuid.UUID{diagnosis1.ID, diagnosis2.ID}
			},
			checks: func(t *testing.T, medications []ClientMedication, err error) {
				require.NoError(t, err, "ListMedicationsByDiagnosisIDs() should not error")
				require.Len(t, medications, 2, "should return 2 medications")
			},
		},
		{
			name: "list medications with empty diagnosis IDs",
			setup: func(ctx context.Context, qtx *Queries) []uuid.UUID {
				return []uuid.UUID{}
			},
			checks: func(t *testing.T, medications []ClientMedication, err error) {
				require.NoError(t, err)
				require.Empty(t, medications, "should return empty slice for empty diagnosis IDs")
			},
		},
		{
			name: "list medications for non-existent diagnosis IDs",
			setup: func(ctx context.Context, qtx *Queries) []uuid.UUID {
				return []uuid.UUID{uuid.New(), uuid.New()}
			},
			checks: func(t *testing.T, medications []ClientMedication, err error) {
				require.NoError(t, err)
				require.Empty(t, medications, "should return empty slice for non-existent diagnosis IDs")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			diagnosisIDs := tt.setup(ctx, qtx)
			medications, err := qtx.ListMedicationsByDiagnosisIDs(ctx, diagnosisIDs)
			tt.checks(t, medications, err)
		})
	}
}

// Random data generators for tests

func createRandomDiagnosis(ctx context.Context, qtx *Queries) ClientDiagnosis {
	client := createRandomClientDetails(ctx, qtx)
	diagnosis, err := qtx.CreateClientDiagnosis(ctx, CreateClientDiagnosisParams{
		ClientID:            client.ID,
		DiagnosisCode:       util.RandomString(10),
		Description:         util.RandomString(20),
		Status:              "ACTIVE",
		Title:               util.StringPtr(util.RandomString(10)),
		Severity:            util.StringPtr("MODERATE"),
		DiagnosingClinician: util.StringPtr(util.RandomString(10)),
		Notes:               util.StringPtr(util.RandomString(30)),
	})
	if err != nil {
		panic("failed to create random diagnosis: " + err.Error())
	}
	return diagnosis
}

func createRandomMedication(ctx context.Context, qtx *Queries) ClientMedication {
	diagnosis := createRandomDiagnosis(ctx, qtx)
	medication, err := qtx.CreateClientMedication(ctx, CreateClientMedicationParams{
		DiagnosisID:      &diagnosis.ID,
		Name:             util.RandomString(10),
		Dosage:           util.RandomString(5) + "mg",
		StartDate:        pgtype.Date{Valid: true},
		EndDate:          pgtype.Date{Valid: true},
		SelfAdministered: util.RandomBool(),
		IsCritical:       util.RandomBool(),
		Notes:            util.StringPtr(util.RandomString(20)),
	})
	if err != nil {
		panic("failed to create random medication: " + err.Error())
	}
	return medication
}
