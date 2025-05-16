package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomRegistrationForm(t *testing.T) RegistrationForm {
	arg := CreateRegistrationFormParams{
		ClientFirstName:               faker.FirstName(),
		ClientLastName:                faker.LastName(),
		ClientBsnNumber:               "123456789",
		ClientGender:                  "male",
		ClientNationality:             "Dutch",
		ClientPhoneNumber:             faker.PhoneNumber,
		ClientEmail:                   faker.Email(),
		ClientStreet:                  "123 Main St",
		ClientHouseNumber:             "1A",
		ClientPostalCode:              "1234AB",
		ClientCity:                    "Amsterdam",
		ReferrerFirstName:             faker.FirstName(),
		ReferrerLastName:              faker.LastName(),
		ReferrerOrganization:          "Referrer Org",
		ReferrerJobTitle:              "Referrer Job",
		ReferrerPhoneNumber:           faker.PhoneNumber,
		ReferrerEmail:                 faker.Email(),
		Guardian1FirstName:            faker.FirstName(),
		Guardian1LastName:             faker.LastName(),
		Guardian1Relationship:         "Parent",
		Guardian1PhoneNumber:          faker.PhoneNumber,
		Guardian1Email:                faker.Email(),
		Guardian2FirstName:            faker.FirstName(),
		Guardian2LastName:             faker.LastName(),
		Guardian2Relationship:         "Parent",
		Guardian2PhoneNumber:          faker.PhoneNumber,
		Guardian2Email:                faker.Email(),
		EducationInstitution:          "Education Institution",
		EducationMentorName:           faker.Name(),
		EducationMentorPhone:          faker.PhoneNumber,
		EducationMentorEmail:          faker.Email(),
		EducationCurrentlyEnrolled:    true,
		EducationAdditionalNotes:      util.StringPtr("Additional notes"),
		CareProtectedLiving:           util.BoolPtr(true),
		CareAssistedIndependentLiving: util.BoolPtr(false),
		CareRoomTrainingCenter:        util.BoolPtr(true),
		CareAmbulatoryGuidance:        util.BoolPtr(false),
		RiskAggressiveBehavior:        util.BoolPtr(true),
		RiskSuicidalSelfharm:          util.BoolPtr(false),
		RiskSubstanceAbuse:            util.BoolPtr(true),
		RiskPsychiatricIssues:         util.BoolPtr(false),
		RiskCriminalHistory:           util.BoolPtr(true),
		RiskFlightBehavior:            util.BoolPtr(false),
		RiskWeaponPossession:          util.BoolPtr(true),
		RiskSexualBehavior:            util.BoolPtr(false),
		RiskDayNightRhythm:            util.BoolPtr(true),
		RiskOther:                     util.BoolPtr(false),
		RiskOtherDescription:          nil,
		RiskAdditionalNotes:           nil,
		DocumentReferral:              nil,
		DocumentEducationReport:       nil,
		DocumentPsychiatricReport:     nil,
		DocumentDiagnosis:             nil,
		DocumentSafetyPlan:            nil,
		DocumentIDCopy:                nil,
		ApplicationDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		ReferrerSignature: util.BoolPtr(true),
	}

	registrationForm, err := testQueries.CreateRegistrationForm(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, registrationForm)
	return registrationForm
}

func TestCreateRegistrationForm(t *testing.T) {
	createRandomRegistrationForm(t)
}

func TestListRegistrationForms(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomRegistrationForm(t)
	}

	arg := ListRegistrationFormsParams{
		Limit:  5,
		Offset: 0,
	}
	registrationForms, err := testQueries.ListRegistrationForms(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, registrationForms)
	require.Len(t, registrationForms, 5)
}

func TestGetRegistrationForm(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)

	registrationFormFromDB, err := testQueries.GetRegistrationForm(context.Background(), registrationForm.ID)
	require.NoError(t, err)
	require.NotEmpty(t, registrationFormFromDB)
	require.Equal(t, registrationForm.ID, registrationFormFromDB.ID)
}

func TestUpdateRegistrationForm(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)

	arg := UpdateRegistrationFormParams{
		ID:                registrationForm.ID,
		ClientFirstName:   util.StringPtr(faker.FirstName()),
		ClientLastName:    util.StringPtr(faker.LastName()),
		ClientBsnNumber:   util.StringPtr("987654321"),
		ClientGender:      util.StringPtr("female"),
		ClientNationality: util.StringPtr("Dutch"),
	}

	registrationFormFromDB, err := testQueries.UpdateRegistrationForm(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, registrationFormFromDB)
	require.Equal(t, registrationForm.ID, registrationFormFromDB.ID)
	require.Equal(t, arg.ClientFirstName, &registrationFormFromDB.ClientFirstName)
	require.Equal(t, arg.ClientLastName, &registrationFormFromDB.ClientLastName)
}

func TestDeleteRegistrationForm(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)

	err := testQueries.DeleteRegistrationForm(context.Background(), registrationForm.ID)
	require.NoError(t, err)

	registrationFormFromDB, err := testQueries.GetRegistrationForm(context.Background(), registrationForm.ID)
	require.Error(t, err)
	require.Empty(t, registrationFormFromDB)
}

func TestUpdateRegistrationFormStatus(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)
	employee, _ := createRandomEmployee(t)

	arg := UpdateRegistrationFormStatusParams{
		ID:                    registrationForm.ID,
		FormStatus:            "approved",
		ProcessedByEmployeeID: &employee.ID,
	}

	err := testQueries.UpdateRegistrationFormStatus(context.Background(), arg)
	require.NoError(t, err)
}
