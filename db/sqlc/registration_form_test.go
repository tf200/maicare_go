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

func TestCreateRegistrationForm(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateRegistrationFormParams
		checks func(t *testing.T, form RegistrationForm, err error)
	}{
		{
			name: "successful registration form creation",
			setup: func(ctx context.Context, qtx *Queries) CreateRegistrationFormParams {
				return createRandomRegistrationFormParams()
			},
			checks: func(t *testing.T, form RegistrationForm, err error) {
				require.NoError(t, err, "CreateRegistrationForm should not return an error")
				require.NotZero(t, form.ID)
				require.Equal(t, "DRAFT", string(form.FormStatus))
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
			form, err := qtx.CreateRegistrationForm(ctx, params)
			tt.checks(t, form, err)
		})
	}
}

func TestGetRegistrationForm(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, form RegistrationForm, err error)
	}{
		{
			name: "get existing registration form",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				form := createRandomRegistrationForm(ctx, qtx)
				return form.ID
			},
			checks: func(t *testing.T, form RegistrationForm, err error) {
				require.NoError(t, err, "GetRegistrationForm should not return an error")
				require.NotZero(t, form.ID)
			},
		},
		{
			name: "get non-existent registration form",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 999999
			},
			checks: func(t *testing.T, form RegistrationForm, err error) {
				require.Error(t, err, "GetRegistrationForm should error for non-existent form")
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

			formID := tt.setup(ctx, qtx)
			form, err := qtx.GetRegistrationForm(ctx, formID)
			tt.checks(t, form, err)
		})
	}
}

func TestUpdateRegistrationForm(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateRegistrationFormParams
		checks func(t *testing.T, form RegistrationForm, err error)
	}{
		{
			name: "update existing registration form",
			setup: func(ctx context.Context, qtx *Queries) UpdateRegistrationFormParams {
				form := createRandomRegistrationForm(ctx, qtx)
				return UpdateRegistrationFormParams{
					ID:                         form.ID,
					ClientFirstName:            randomStringPtr(10),
					ClientLastName:             randomStringPtr(10),
					ClientBsnNumber:            randomStringPtr(9),
					ClientGender:               randomClientGender(),
					ClientNationality:          randomStringPtr(10),
					ClientPhoneNumber:          randomStringPtr(10),
					ClientEmail:                randomStringPtr(15),
					ClientStreet:               randomStringPtr(10),
					ClientHouseNumber:          randomStringPtr(5),
					ClientPostalCode:           randomStringPtr(7),
					ClientCity:                 randomStringPtr(10),
					ReferrerFirstName:          randomStringPtr(10),
					ReferrerLastName:           randomStringPtr(10),
					ReferrerOrganization:       randomStringPtr(15),
					ReferrerJobTitle:           randomStringPtr(10),
					ReferrerPhoneNumber:        randomStringPtr(10),
					ReferrerEmail:              randomStringPtr(15),
					Guardian1FirstName:         randomStringPtr(10),
					Guardian1LastName:          randomStringPtr(10),
					Guardian1Relationship:      randomStringPtr(10),
					Guardian1PhoneNumber:       randomStringPtr(10),
					Guardian1Email:             randomStringPtr(15),
					Guardian2FirstName:         randomStringPtr(10),
					Guardian2LastName:          randomStringPtr(10),
					Guardian2Relationship:      randomStringPtr(10),
					Guardian2PhoneNumber:       randomStringPtr(10),
					Guardian2Email:             randomStringPtr(15),
					EducationInstitution:       randomStringPtr(20),
					EducationMentorName:        randomStringPtr(10),
					EducationMentorPhone:       randomStringPtr(10),
					EducationMentorEmail:       randomStringPtr(15),
					EducationCurrentlyEnrolled: randomBoolPtr(),
					EducationAdditionalNotes:   randomStringPtr(20),
					EducationLevel:             randomClientEducationLevel(),
					WorkCurrentEmployer:        randomStringPtr(15),
					WorkEmployerPhone:          randomStringPtr(10),
					WorkEmployerEmail:          randomStringPtr(15),
					WorkCurrentPosition:        randomStringPtr(10),
					WorkCurrentlyEmployed:      randomBoolPtr(),
					WorkStartDate: pgtype.Date{
						Time:  time.Now(),
						Valid: true,
					},
					WorkAdditionalNotes:           randomStringPtr(20),
					CareProtectedLiving:           randomBoolPtr(),
					CareAssistedIndependentLiving: randomBoolPtr(),
					CareRoomTrainingCenter:        randomBoolPtr(),
					CareAmbulatoryGuidance:        randomBoolPtr(),
					RiskAggressiveBehavior:        randomBoolPtr(),
					RiskSuicidalSelfharm:          randomBoolPtr(),
					RiskSubstanceAbuse:            randomBoolPtr(),
					RiskPsychiatricIssues:         randomBoolPtr(),
					RiskCriminalHistory:           randomBoolPtr(),
					RiskFlightBehavior:            randomBoolPtr(),
					RiskWeaponPossession:          randomBoolPtr(),
					RiskSexualBehavior:            randomBoolPtr(),
					RiskDayNightRhythm:            randomBoolPtr(),
					RiskOther:                     randomBoolPtr(),
					RiskOtherDescription:          randomStringPtr(20),
					RiskAdditionalNotes:           randomStringPtr(20),
					DocumentReferral:              randomUUIDPtr(),
					DocumentEducationReport:       randomUUIDPtr(),
					DocumentPsychiatricReport:     randomUUIDPtr(),
					DocumentDiagnosis:             randomUUIDPtr(),
					DocumentSafetyPlan:            randomUUIDPtr(),
					DocumentIDCopy:                randomUUIDPtr(),
					ApplicationDate: pgtype.Date{
						Time:  time.Now(),
						Valid: true,
					},
					ReferrerSignature: randomBoolPtr(),
				}
			},
			checks: func(t *testing.T, form RegistrationForm, err error) {
				require.NoError(t, err, "UpdateRegistrationForm should not return an error")
				require.NotZero(t, form.ID)
			},
		},
		{
			name: "update non-existent registration form",
			setup: func(ctx context.Context, qtx *Queries) UpdateRegistrationFormParams {
				return UpdateRegistrationFormParams{
					ID:              999999,
					ClientFirstName: randomStringPtr(10),
				}
			},
			checks: func(t *testing.T, form RegistrationForm, err error) {
				require.Error(t, err, "UpdateRegistrationForm should error for non-existent form")
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
			form, err := qtx.UpdateRegistrationForm(ctx, params)
			tt.checks(t, form, err)
		})
	}
}

func TestDeleteRegistrationForm(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete existing registration form",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				form := createRandomRegistrationForm(ctx, qtx)
				return form.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteRegistrationForm should not error")
			},
		},
		{
			name: "delete non-existent registration form",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 999999
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteRegistrationForm should not error for non-existent form")
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

			formID := tt.setup(ctx, qtx)
			err = qtx.DeleteRegistrationForm(ctx, formID)
			tt.checks(t, err)
		})
	}
}

func TestListRegistrationForms(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListRegistrationFormsParams
		checks func(t *testing.T, forms []RegistrationForm, err error)
	}{
		{
			name: "list registration forms without filters",
			setup: func(ctx context.Context, qtx *Queries) ListRegistrationFormsParams {
				// Create a form
				createRandomRegistrationForm(ctx, qtx)
				return ListRegistrationFormsParams{
					Limit:  10,
					Offset: 0,
					Status: NullFormStatusEnum{
						FormStatusEnum: "",
						Valid:          false,
					},
				}
			},
			checks: func(t *testing.T, forms []RegistrationForm, err error) {
				require.NoError(t, err, "ListRegistrationForms should not error")
				require.GreaterOrEqual(t, len(forms), 1)
			},
		},
		{
			name: "list registration forms with status filter",
			setup: func(ctx context.Context, qtx *Queries) ListRegistrationFormsParams {
				// Create a form
				createRandomRegistrationForm(ctx, qtx)
				return ListRegistrationFormsParams{
					Limit:  10,
					Offset: 0,
					Status: NullFormStatusEnum{
						FormStatusEnum: "DRAFT",
						Valid:          true,
					},
				}
			},
			checks: func(t *testing.T, forms []RegistrationForm, err error) {
				require.NoError(t, err, "ListRegistrationForms should not error")
				require.GreaterOrEqual(t, len(forms), 1)
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
			forms, err := qtx.ListRegistrationForms(ctx, params)
			tt.checks(t, forms, err)
		})
	}
}

func TestCountRegistrationForms(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CountRegistrationFormsParams
		checks func(t *testing.T, count int64, err error)
	}{
		{
			name: "count registration forms without filters",
			setup: func(ctx context.Context, qtx *Queries) CountRegistrationFormsParams {
				// Create a form
				createRandomRegistrationForm(ctx, qtx)
				return CountRegistrationFormsParams{
					Status: NullFormStatusEnum{
						FormStatusEnum: "",
						Valid:          false,
					},
				}
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountRegistrationForms should not error")
				require.GreaterOrEqual(t, count, 1)
			},
		},
		{
			name: "count registration forms with risk filter",
			setup: func(ctx context.Context, qtx *Queries) CountRegistrationFormsParams {
				// Create a form
				createRandomRegistrationForm(ctx, qtx)
				risk := true
				return CountRegistrationFormsParams{
					Status: NullFormStatusEnum{
						FormStatusEnum: "",
						Valid:          false,
					},
					RiskAggressiveBehavior: &risk,
				}
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountRegistrationForms should not error")
				// Count might be 0 or 1 depending on random data
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
			count, err := qtx.CountRegistrationForms(ctx, params)
			tt.checks(t, count, err)
		})
	}
}

// Helpers
func createRandomRegistrationForm(ctx context.Context, qtx *Queries) RegistrationForm {
	params := createRandomRegistrationFormParams()
	form, err := qtx.CreateRegistrationForm(ctx, params)
	if err != nil {
		panic(err)
	}
	return form
}

func createRandomRegistrationFormParams() CreateRegistrationFormParams {
	return CreateRegistrationFormParams{
		ClientFirstName:            util.RandomString(10),
		ClientLastName:             util.RandomString(10),
		ClientBsnNumber:            util.RandomString(9),
		ClientGender:               ClientGenderEnum("MALE"),
		ClientNationality:          util.RandomString(10),
		ClientPhoneNumber:          util.RandomString(10),
		ClientEmail:                util.RandomEmail(),
		ClientStreet:               util.RandomString(10),
		ClientHouseNumber:          util.RandomString(5),
		ClientPostalCode:           util.RandomString(7),
		ClientCity:                 util.RandomString(10),
		ReferrerFirstName:          util.RandomString(10),
		ReferrerLastName:           util.RandomString(10),
		ReferrerOrganization:       util.RandomString(15),
		ReferrerJobTitle:           util.RandomString(10),
		ReferrerPhoneNumber:        util.RandomString(10),
		ReferrerEmail:              util.RandomEmail(),
		Guardian1FirstName:         util.RandomString(10),
		Guardian1LastName:          util.RandomString(10),
		Guardian1Relationship:      util.RandomString(10),
		Guardian1PhoneNumber:       util.RandomString(10),
		Guardian1Email:             util.RandomEmail(),
		Guardian2FirstName:         util.RandomString(10),
		Guardian2LastName:          util.RandomString(10),
		Guardian2Relationship:      util.RandomString(10),
		Guardian2PhoneNumber:       util.RandomString(10),
		Guardian2Email:             util.RandomEmail(),
		EducationInstitution:       randomStringPtr(20),
		EducationMentorName:        randomStringPtr(10),
		EducationMentorPhone:       randomStringPtr(10),
		EducationMentorEmail:       randomStringPtr(15),
		EducationCurrentlyEnrolled: util.RandomBool(),
		EducationAdditionalNotes:   randomStringPtr(20),
		EducationLevel:             randomClientEducationLevel(),
		WorkCurrentEmployer:        randomStringPtr(15),
		WorkEmployerPhone:          randomStringPtr(10),
		WorkEmployerEmail:          randomStringPtr(15),
		WorkCurrentPosition:        randomStringPtr(10),
		WorkCurrentlyEmployed:      util.RandomBool(),
		WorkStartDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		WorkAdditionalNotes:           randomStringPtr(20),
		CareProtectedLiving:           randomBoolPtr(),
		CareAssistedIndependentLiving: randomBoolPtr(),
		CareRoomTrainingCenter:        randomBoolPtr(),
		CareAmbulatoryGuidance:        randomBoolPtr(),
		RiskAggressiveBehavior:        randomBoolPtr(),
		RiskSuicidalSelfharm:          randomBoolPtr(),
		RiskSubstanceAbuse:            randomBoolPtr(),
		RiskPsychiatricIssues:         randomBoolPtr(),
		RiskCriminalHistory:           randomBoolPtr(),
		RiskFlightBehavior:            randomBoolPtr(),
		RiskWeaponPossession:          randomBoolPtr(),
		RiskSexualBehavior:            randomBoolPtr(),
		RiskDayNightRhythm:            randomBoolPtr(),
		RiskOther:                     randomBoolPtr(),
		RiskOtherDescription:          randomStringPtr(20),
		RiskAdditionalNotes:           randomStringPtr(20),
		DocumentReferral:              randomUUIDPtr(),
		DocumentEducationReport:       randomUUIDPtr(),
		DocumentPsychiatricReport:     randomUUIDPtr(),
		DocumentDiagnosis:             randomUUIDPtr(),
		DocumentSafetyPlan:            randomUUIDPtr(),
		DocumentIDCopy:                randomUUIDPtr(),
		ApplicationDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		ReferrerSignature: randomBoolPtr(),
	}
}

func randomBoolPtr() *bool {
	b := util.RandomBool()
	return &b
}

func randomUUIDPtr() *uuid.UUID {
	id := uuid.New()
	return &id
}

func randomClientGender() NullClientGenderEnum {
	return NullClientGenderEnum{
		ClientGenderEnum: "MALE",
		Valid:            true,
	}
}

func randomClientEducationLevel() NullClientEducationLevelEnum {
	return NullClientEducationLevelEnum{
		ClientEducationLevelEnum: "HIGH_SCHOOL",
		Valid:                    true,
	}
}
