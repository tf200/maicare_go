package db

import (
	"context"
	"maicare_go/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomClientAllergy(t *testing.T, clientID int64) ClientAllergy {

	arg := CreateClientAllergyParams{
		ClientID:      clientID,
		AllergyTypeID: 1,
		Severity:      "Mild",
		Reaction:      "test reaction",
		Notes:         util.StringPtr("test note"),
	}

	allergy, err := testQueries.CreateClientAllergy(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, allergy)
	require.Equal(t, arg.ClientID, allergy.ClientID)
	return allergy
}

func TestCreateClientAllergy(t *testing.T) {
	client := createRandomClientDetails(t)
	createRandomClientAllergy(t, client.ID)
}

func TestListClientAllergies(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 20; i++ {
		_ = createRandomClientAllergy(t, client.ID)
	}
	testCases := []struct {
		name  string
		arg   ListClientAllergiesParams
		check func(t *testing.T, allergies []ListClientAllergiesRow)
	}{
		{
			name: "base case",
			arg: ListClientAllergiesParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   0,
			},
			check: func(t *testing.T, allergies []ListClientAllergiesRow) {
				require.NotEmpty(t, allergies)
				require.Len(t, allergies, 5)
				require.Equal(t, int64(20), allergies[0].TotalAllergies)
			},
		},
		{
			name: "with offset",
			arg: ListClientAllergiesParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   5,
			},
			check: func(t *testing.T, allergies []ListClientAllergiesRow) {
				require.NotEmpty(t, allergies)
				require.Len(t, allergies, 5)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allergies, err := testQueries.ListClientAllergies(context.Background(), tc.arg)
			require.NoError(t, err)
			tc.check(t, allergies)
		})
	}
}

func TestGetClientAllergy(t *testing.T) {
	client := createRandomClientDetails(t)
	allergy1 := createRandomClientAllergy(t, client.ID)

	allergy2, err := testQueries.GetClientAllergy(context.Background(), allergy1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, allergy2)
	require.Equal(t, allergy1.ID, allergy2.ID)
}

func TestUpdateClientAllergyy(t *testing.T) {
	client := createRandomClientDetails(t)
	allergy1 := createRandomClientAllergy(t, client.ID)

	arg := UpdateClientAllergyParams{
		ID:       allergy1.ID,
		Severity: util.StringPtr("Severe"),
	}

	allergy2, err := testQueries.UpdateClientAllergy(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, allergy2)
	require.Equal(t, allergy1.ID, allergy2.ID)
	require.NotEqual(t, allergy1.Severity, allergy2.Severity)

}

func TestDeleteClientAllergy(t *testing.T) {
	client := createRandomClientDetails(t)
	allergy1 := createRandomClientAllergy(t, client.ID)

	_, err := testQueries.DeleteClientAllergy(context.Background(), allergy1.ID)
	require.NoError(t, err)

}

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

func createRandomClientMedication(t *testing.T, clientID int64, employeeID int64) ClientMedication {

	arg := CreateClientMedicationParams{
		ClientID:         clientID,
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
	require.Equal(t, arg.ClientID, medication.ClientID)
	return medication
}

func TestCreateClientMedication(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	createRandomClientMedication(t, client.ID, employee.ID)
}

func TestListClientMedications(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	for i := 0; i < 20; i++ {
		_ = createRandomClientMedication(t, client.ID, employee.ID)
	}

	testCases := []struct {
		name  string
		arg   ListClientMedicationsParams
		check func(t *testing.T, medications []ListClientMedicationsRow)
	}{
		{
			name: "base case",
			arg: ListClientMedicationsParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   0,
			},
			check: func(t *testing.T, medications []ListClientMedicationsRow) {
				require.NotEmpty(t, medications)
				require.Len(t, medications, 5)
				require.Equal(t, int64(20), medications[0].TotalMedications)
			},
		},
		{
			name: "with offset",
			arg: ListClientMedicationsParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   5,
			},
			check: func(t *testing.T, medications []ListClientMedicationsRow) {
				require.NotEmpty(t, medications)
				require.Len(t, medications, 5)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			medications, err := testQueries.ListClientMedications(context.Background(), tc.arg)
			require.NoError(t, err)
			tc.check(t, medications)
		})
	}
}

func TestGetClientMedication(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	medication1 := createRandomClientMedication(t, client.ID, employee.ID)

	medication2, err := testQueries.GetClientMedication(context.Background(), medication1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, medication2)
	require.Equal(t, medication1.ID, medication2.ID)
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

	_, err := testQueries.DeleteClientMedication(context.Background(), medication1.ID)
	require.NoError(t, err)

}
