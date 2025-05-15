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

func TestCreateRegistrationForm(t *testing.T) {
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
		DocumentReferral:              pgtype.UUID{},
		DocumentEducationReport:       pgtype.UUID{},
		DocumentPsychiatricReport:     pgtype.UUID{},
		DocumentDiagnosis:             pgtype.UUID{},
		DocumentSafetyPlan:            pgtype.UUID{},
		DocumentIDCopy:                pgtype.UUID{},
		ApplicationDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		ReferrerSignature: util.BoolPtr(true),
	}

	registrationForm, err := testQueries.CreateRegistrationForm(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, registrationForm)
}
