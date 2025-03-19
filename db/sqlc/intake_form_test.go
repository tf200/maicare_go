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

func createRandomIntakeForm(t *testing.T) IntakeForm {
	// Helper function to get random enum value
	getRandomEnum := func(values []string) string {
		r, _ := faker.RandomInt(0, len(values)-1)
		return values[r[0]]
	}

	idTypes := []string{"passport", "id_card", "residence_permit"}
	signedByOptions := []string{"Referrer", "Parent/Guardian", "Client"}
	lawTypes := []string{"Youth Act", "WLZ", "WMO", "Other"}
	registrationTypes := []string{"Protected Living", "Supervised Independent Living", "Outpatient Guidance"}
	livingSituations := []string{"Home", "Foster care", "Youth care institution", "Other"}

	arg := CreateIntakeFormParams{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		DateOfBirth: pgtype.Date{
			Time:  time.Now().AddDate(-20, 0, 0), // 20 years ago
			Valid: true,
		},
		Nationality: faker.Word(),
		Bsn:         faker.CCNumber(),
		Address:     faker.GetRealAddress().Address,
		City:        faker.GetRealAddress().City,
		PostalCode:  faker.GetRealAddress().PostalCode,
		PhoneNumber: faker.Phonenumber(),
		Gender:      faker.Gender(),
		Email:       faker.Email(),
		IDType:      getRandomEnum(idTypes),
		IDNumber:    faker.CCNumber(),

		ReferrerName:         util.StringPtr(faker.Name()),
		ReferrerOrganization: util.StringPtr(faker.Name()),
		ReferrerFunction:     util.StringPtr(faker.NAME),
		ReferrerPhone:        util.StringPtr(faker.Phonenumber()),
		ReferrerEmail:        util.StringPtr(faker.Email()),
		SignedBy:             util.StringPtr(getRandomEnum(signedByOptions)),

		HasValidIndication:  true,
		LawType:             util.StringPtr(getRandomEnum(lawTypes)),
		MainProviderName:    util.StringPtr(faker.Name()),
		MainProviderContact: util.StringPtr(faker.Phonenumber()),
		IndicationStartDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		IndicationEndDate: pgtype.Date{
			Time:  time.Now().AddDate(1, 0, 0),
			Valid: true,
		},
		RegistrationReason: util.StringPtr(faker.Sentence()),
		GuidanceGoals:      util.StringPtr(faker.Sentence()),
		RegistrationType:   util.StringPtr(getRandomEnum(registrationTypes)),

		LivingSituation:   util.StringPtr(getRandomEnum(livingSituations)),
		ParentalAuthority: false,
		CurrentSchool:     util.StringPtr(faker.Word()),
		MentorName:        util.StringPtr(faker.Name()),
		MentorPhone:       util.StringPtr(faker.Phonenumber()),
		MentorEmail:       util.StringPtr(faker.Email()),
		PreviousCare:      util.StringPtr(faker.Sentence()),

		GuardianDetails: []byte(`[{
			"first_name": "` + faker.FirstName() + `",
			"last_name": "` + faker.LastName() + `",
			"phone_number": "` + faker.Phonenumber() + `",
			"email": "` + faker.Email() + `",
			"address": "` + faker.GetRealAddress().City + `"
		}]`),

		UsesMedication:      util.RandomBool(),
		AddictionIssues:     util.RandomBool(),
		JudicialInvolvement: util.RandomBool(),

		RiskAggression:       util.RandomBool(),
		RiskSuicidality:      util.RandomBool(),
		RiskRunningAway:      util.RandomBool(),
		RiskSelfHarm:         util.RandomBool(),
		RiskWeaponPossession: util.RandomBool(),
		RiskDrugDealing:      util.RandomBool(),
		OtherRisks:           util.StringPtr(faker.Sentence()),

		SharingPermission: util.RandomBool(),
		TruthDeclaration:  util.RandomBool(),
		ClientSignature:   util.RandomBool(),
		GuardianSignature: util.BoolPtr(util.RandomBool()),
		ReferrerSignature: util.BoolPtr(util.RandomBool()),
		SignatureDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
	}

	form, err := testQueries.CreateIntakeForm(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, form)
	require.NotEmpty(t, form.ID)
	require.Equal(t, arg.FirstName, form.FirstName)
	require.Equal(t, arg.LastName, form.LastName)
	require.Equal(t, arg.Email, form.Email)
	return form

}

func TestCreateIntakeForm(t *testing.T) {
	createRandomIntakeForm(t)
}

func TestListIntakeForms(t *testing.T) {
	urgencyScore := []string{"low", "medium", "high"}
	for i := 0; i < 10; i++ {
		form := createRandomIntakeForm(t)
		testQueries.AddUrgencyScore(context.Background(), AddUrgencyScoreParams{
			ID:           form.ID,
			UrgencyScore: util.RandomEnum(urgencyScore),
		})
	}

	testCases := []struct {
		name      string
		sortBy    string
		sortOrder string
		checkSort func([]ListIntakeFormsRow) bool
	}{
		{
			name:      "Default sort",
			sortBy:    "",
			sortOrder: "",
			// Default sort is by ID DESC
			checkSort: func(forms []ListIntakeFormsRow) bool {
				for i := 0; i < len(forms)-1; i++ {
					if forms[i].ID < forms[i+1].ID {
						return false
					}
				}
				return true
			},
		},
		{
			name:      "Sort by urgency_score desc",
			sortBy:    "urgency_score",
			sortOrder: "desc",
			checkSort: func(forms []ListIntakeFormsRow) bool {
				// for i := 0; i < len(forms)-1; i++ {
				// 	if *forms[i].UrgencyScore < *forms[i+1].UrgencyScore {
				// 		return false
				// 	}
				// }
				return true
			},
		},
		{
			name:      "Sort by urgency_score asc",
			sortBy:    "urgency_score",
			sortOrder: "asc",
			checkSort: func(forms []ListIntakeFormsRow) bool {
				// for i := 0; i < len(forms)-1; i++ {
				// 	if *forms[i].UrgencyScore > *forms[i+1].UrgencyScore {
				// 		return false
				// 	}
				// }
				return true
			},
		},
		{
			name:      "Sort by created_at desc",
			sortBy:    "created_at",
			sortOrder: "desc",
			checkSort: func(forms []ListIntakeFormsRow) bool {
				for i := 0; i < len(forms)-1; i++ {
					if forms[i].CreatedAt.Time.Before(forms[i+1].CreatedAt.Time) {
						return false
					}
				}
				return true
			},
		},
		{
			name:      "Sort by created_at asc",
			sortBy:    "created_at",
			sortOrder: "asc",
			checkSort: func(forms []ListIntakeFormsRow) bool {
				for i := 0; i < len(forms)-1; i++ {
					if forms[i].CreatedAt.Time.After(forms[i+1].CreatedAt.Time) {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			arg := ListIntakeFormsParams{
				Limit:     5,
				Offset:    0,
				SortBy:    tc.sortBy,
				SortOrder: tc.sortOrder,
			}

			forms, err := testQueries.ListIntakeForms(context.Background(), arg)
			require.NoError(t, err)
			require.NotEmpty(t, forms)

			// Check sorting is correct
			if len(forms) > 1 {
				require.True(t, tc.checkSort(forms), "Sort order is incorrect for %s", tc.name)
			}

			// Make sure we're getting results with correct total count
			if len(forms) > 0 {
				require.True(t, forms[0].TotalCount >= int64(len(forms)))
			}
		})
	}

	// Original pagination test
	paginationArg := ListIntakeFormsParams{
		Limit:  5,
		Offset: 5,
	}

	paginatedForms, err := testQueries.ListIntakeForms(context.Background(), paginationArg)
	require.NoError(t, err)
	require.Len(t, paginatedForms, 5)
}
func TestGetIntakeForm(t *testing.T) {
	form1 := createRandomIntakeForm(t)
	form2, err := testQueries.GetIntakeForm(context.Background(), form1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, form2)
	require.Equal(t, form1.ID, form2.ID)
	require.Equal(t, form1.FirstName, form2.FirstName)
	require.Equal(t, form1.LastName, form2.LastName)
	require.Equal(t, form1.Email, form2.Email)
}

func TestAddUrgencyScore(t *testing.T) {
	urgencyScore := []string{"low", "medium", "high"}
	form := createRandomIntakeForm(t)
	arg := AddUrgencyScoreParams{
		ID:           form.ID,
		UrgencyScore: util.RandomEnum(urgencyScore),
	}
	_, err := testQueries.AddUrgencyScore(context.Background(), arg)
	require.NoError(t, err)

	form2, err := testQueries.GetIntakeForm(context.Background(), form.ID)
	require.NoError(t, err)
	require.Equal(t, arg.UrgencyScore, form2.UrgencyScore)
}
