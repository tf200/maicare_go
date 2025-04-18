package db

import (
	"context"
	"maicare_go/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomClientDiagnosis(t *testing.T, clientID int64) ClientDiagnosis {

	arg := CreateClientDiagnosisParams{
		ClientID:            clientID,
		Title:               util.StringPtr("test title"),
		DiagnosisCode:       "test code",
		Description:         "test description",
		Severity:            util.StringPtr("Mild"),
		Status:              "Active",
		DiagnosingClinician: util.StringPtr("Dr. Test"),
		Notes:               util.StringPtr("test note"),
	}

	diagnosis, err := testQueries.CreateClientDiagnosis(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, diagnosis)
	require.Equal(t, arg.ClientID, diagnosis.ClientID)
	return diagnosis
}

func TestCreateClientDiagnosis(t *testing.T) {
	client := createRandomClientDetails(t)
	createRandomClientDiagnosis(t, client.ID)
}

func TestListClientDiagnoses(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 20; i++ {
		_ = createRandomClientDiagnosis(t, client.ID)
	}
	testCases := []struct {
		name  string
		arg   ListClientDiagnosesParams
		check func(t *testing.T, diagnoses []ListClientDiagnosesRow)
	}{
		{
			name: "base case",
			arg: ListClientDiagnosesParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   0,
			},
			check: func(t *testing.T, diagnoses []ListClientDiagnosesRow) {
				require.NotEmpty(t, diagnoses)
				require.Len(t, diagnoses, 5)
				require.Equal(t, int64(20), diagnoses[0].TotalDiagnoses)
			},
		},
		{
			name: "with offset",
			arg: ListClientDiagnosesParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   5,
			},
			check: func(t *testing.T, diagnoses []ListClientDiagnosesRow) {
				require.NotEmpty(t, diagnoses)
				require.Len(t, diagnoses, 5)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diagnoses, err := testQueries.ListClientDiagnoses(context.Background(), tc.arg)
			require.NoError(t, err)
			tc.check(t, diagnoses)
		})
	}
}

func TestGetClientDiagnosis(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis1 := createRandomClientDiagnosis(t, client.ID)

	diagnosis2, err := testQueries.GetClientDiagnosis(context.Background(), diagnosis1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, diagnosis2)
	require.Equal(t, diagnosis1.ID, diagnosis2.ID)
}

func TestUpdateClientDiagnosis(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis1 := createRandomClientDiagnosis(t, client.ID)

	arg := UpdateClientDiagnosisParams{
		ID:       diagnosis1.ID,
		Severity: util.StringPtr("Severe"),
	}

	diagnosis2, err := testQueries.UpdateClientDiagnosis(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, diagnosis2)
	require.Equal(t, diagnosis1.ID, diagnosis2.ID)
	require.NotEqual(t, diagnosis1.Severity, diagnosis2.Severity)

}

func TestDeleteClientDiagnosis(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis1 := createRandomClientDiagnosis(t, client.ID)

	_, err := testQueries.DeleteClientDiagnosis(context.Background(), diagnosis1.ID)
	require.NoError(t, err)

}

func createRandomClientMedication(t *testing.T, diagnosisID int64, employeeID int64) ClientMedication {

	arg := CreateClientMedicationParams{
		DiagnosisID:      &diagnosisID,
		Name:             "test name",
		Dosage:           "test dosage",
		StartDate:        pgtype.Date{Time: util.RandomTIme(), Valid: true},
		EndDate:          pgtype.Date{Time: util.RandomTIme(), Valid: true},
		Notes:            util.StringPtr("test note"),
		SelfAdministered: true,
		AdministeredByID: util.IntPtr(employeeID),
		IsCritical:       true,
	}

	medication, err := testQueries.CreateClientMedication(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, medication)
	require.Equal(t, arg.DiagnosisID, medication.DiagnosisID)
	return medication
}

func TestCreateClientMedication(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	createRandomClientMedication(t, client.ID, employee.ID)
}

func TestGetMedication(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis := createRandomClientDiagnosis(t, client.ID)
	employee, _ := createRandomEmployee(t)
	medication1 := createRandomClientMedication(t, diagnosis.ID, employee.ID)

	medication2, err := testQueries.GetMedication(context.Background(), medication1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, medication2)
	require.Equal(t, medication1.ID, medication2.ID)
}

func TestGetMedicationsByDiagnosisID(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis := createRandomClientDiagnosis(t, client.ID)
	employee, _ := createRandomEmployee(t)
	medication1 := createRandomClientMedication(t, diagnosis.ID, employee.ID)

	arg := ListMedicationsByDiagnosisIDParams{
		DiagnosisID: &diagnosis.ID,
		Limit:       5,
		Offset:      0,
	}

	medication2, err := testQueries.ListMedicationsByDiagnosisID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, medication2)
	require.Equal(t, medication1.ID, medication2[0].ID)
	require.Equal(t, medication1.Name, medication2[0].Name)
}

func TestUpdateClientMedication(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	medication1 := createRandomClientMedication(t, client.ID, employee.ID)

	arg := UpdateClientMedicationParams{
		ID:         medication1.ID,
		IsCritical: util.BoolPtr(false),
	}

	medication2, err := testQueries.UpdateClientMedication(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, medication2)
	require.Equal(t, medication1.ID, medication2.ID)
	require.NotEqual(t, medication1.IsCritical, medication2.IsCritical)

}

func TestDeleteClientMedication(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	medication1 := createRandomClientMedication(t, client.ID, employee.ID)

	err := testQueries.DeleteClientMedication(context.Background(), medication1.ID)
	require.NoError(t, err)

}
