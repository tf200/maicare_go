package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateIntakeForm(t *testing.T) {
	tests := []struct {
		name   string
		params func(ctx context.Context, qtx *Queries) CreateIntakeFormParams
		checks func(t *testing.T, intakeForm IntakeForm)
	}{
		{
			name: "Successful Creation",
			params: func(ctx context.Context, qtx *Queries) CreateIntakeFormParams {
				regForm := createRandomRegistrationForm(ctx, qtx)
				maturityMatrix := createRandomMaturityMatrix(ctx, t, qtx)
				return CreateIntakeFormParams{
					RegistrationFormID:    regForm.ID,
					DateOfIntake:          pgtype.Timestamptz{Valid: true, Time: time.Now()},
					CareType:              IntakeCareTypeEnumAmbulatorySupport,
					IntakeParticipants:    []IntakeParticipantsEnum{IntakeParticipantsEnumClient, IntakeParticipantsEnumReferrer},
					FamilySituation:       util.StringPtr("Stable family situation"),
					PsychologicalState:    util.StringPtr("Good psychological state"),
					SelfSufficiency:       4,
					MaturityMatrixID:      &maturityMatrix.ID,
					Goals:                 util.StringPtr("Improve self-sufficiency"),
					RiskAssessment:        util.StringPtr("Low risk"),
					IntakeConclusion:      IntakeConclusionEnumFurtherInvestigation,
					IntakeConclusionNotes: util.StringPtr("Needs further assessment"),
					Signature:             util.StringPtr("Client Signature"),
				}
			},
			checks: func(t *testing.T, intakeForm IntakeForm) {
				require.NotZero(t, intakeForm.ID)
				require.Equal(t, IntakeCareTypeEnumAmbulatorySupport, intakeForm.CareType)
				require.Equal(t, 4, intakeForm.SelfSufficiency)
				require.Equal(t, IntakeConclusionEnumFurtherInvestigation, intakeForm.IntakeConclusion)
				require.NotNil(t, intakeForm.FamilySituation)
				require.Equal(t, "Stable family situation", *intakeForm.FamilySituation)
				require.True(t, intakeForm.CreatedAt.Valid)
			},
		},
		{
			name: "Successful Creation with Minimal Fields",
			params: func(ctx context.Context, qtx *Queries) CreateIntakeFormParams {
				regForm := createRandomRegistrationForm(ctx, qtx)
				return CreateIntakeFormParams{
					RegistrationFormID: regForm.ID,
					DateOfIntake:       pgtype.Timestamptz{Valid: true, Time: time.Now()},
					CareType:           IntakeCareTypeEnumProtectedLiving,
					IntakeParticipants: []IntakeParticipantsEnum{IntakeParticipantsEnumClient},
					SelfSufficiency:    2,
					IntakeConclusion:   IntakeConclusionEnumSuitable,
				}
			},
			checks: func(t *testing.T, intakeForm IntakeForm) {
				require.NotZero(t, intakeForm.ID)
				require.Equal(t, IntakeCareTypeEnumProtectedLiving, intakeForm.CareType)
				require.Equal(t, 2, intakeForm.SelfSufficiency)
				require.Equal(t, IntakeConclusionEnumSuitable, intakeForm.IntakeConclusion)
				require.Nil(t, intakeForm.FamilySituation)
				require.Nil(t, intakeForm.Goals)
			},
		},
		{
			name: "Successful Creation with All Optional Fields Nil",
			params: func(ctx context.Context, qtx *Queries) CreateIntakeFormParams {
				regForm := createRandomRegistrationForm(ctx, qtx)
				return CreateIntakeFormParams{
					RegistrationFormID:    regForm.ID,
					DateOfIntake:          pgtype.Timestamptz{Valid: true, Time: time.Now()},
					CareType:              IntakeCareTypeEnumTrainingCenter,
					IntakeParticipants:    []IntakeParticipantsEnum{IntakeParticipantsEnumReferrer},
					SelfSufficiency:       5,
					IntakeConclusion:      IntakeConclusionEnumUnsuitable,
					FamilySituation:       nil,
					PsychologicalState:    nil,
					MaturityMatrixID:      nil,
					Goals:                 nil,
					RiskAssessment:        nil,
					IntakeConclusionNotes: nil,
					Signature:             nil,
				}
			},
			checks: func(t *testing.T, intakeForm IntakeForm) {
				require.NotZero(t, intakeForm.ID)
				require.Equal(t, IntakeCareTypeEnumTrainingCenter, intakeForm.CareType)
				require.Equal(t, 5, intakeForm.SelfSufficiency)
				require.Equal(t, IntakeConclusionEnumUnsuitable, intakeForm.IntakeConclusion)
				require.Nil(t, intakeForm.FamilySituation)
				require.Nil(t, intakeForm.PsychologicalState)
				require.Nil(t, intakeForm.MaturityMatrixID)
			},
		},
		{
			name: "Successful Creation with Maturity Matrix",
			params: func(ctx context.Context, qtx *Queries) CreateIntakeFormParams {
				regForm := createRandomRegistrationForm(ctx, qtx)
				maturityMatrix := createRandomMaturityMatrix(ctx, t, qtx)
				return CreateIntakeFormParams{
					RegistrationFormID:    regForm.ID,
					DateOfIntake:          pgtype.Timestamptz{Valid: true, Time: time.Now()},
					CareType:              IntakeCareTypeEnumSupportedIndependentLiving,
					IntakeParticipants:    []IntakeParticipantsEnum{IntakeParticipantsEnumClient, IntakeParticipantsEnumParentsGuardians},
					FamilySituation:       util.StringPtr("Complex family dynamics"),
					PsychologicalState:    util.StringPtr("Requires monitoring"),
					SelfSufficiency:       3,
					MaturityMatrixID:      &maturityMatrix.ID,
					Goals:                 util.StringPtr("Build independence"),
					RiskAssessment:        util.StringPtr("Medium risk"),
					IntakeConclusion:      IntakeConclusionEnumFurtherInvestigation,
					IntakeConclusionNotes: util.StringPtr("Monitor progress closely"),
					Signature:             util.StringPtr("Authorized Signature"),
				}
			},
			checks: func(t *testing.T, intakeForm IntakeForm) {
				require.NotZero(t, intakeForm.ID)
				require.Equal(t, IntakeCareTypeEnumSupportedIndependentLiving, intakeForm.CareType)
				require.Equal(t, 3, intakeForm.SelfSufficiency)
				require.Equal(t, IntakeConclusionEnumFurtherInvestigation, intakeForm.IntakeConclusion)
				require.NotNil(t, intakeForm.MaturityMatrixID)
				require.Equal(t, "Complex family dynamics", *intakeForm.FamilySituation)
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

			params := tt.params(ctx, qtx)
			intakeForm, err := qtx.CreateIntakeForm(ctx, params)
			require.NoError(t, err)

			tt.checks(t, intakeForm)
		})
	}
}

func TestListIntakeForms(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) []IntakeForm
		params ListIntakeFormsParams
		checks func(t *testing.T, rows []ListIntakeFormsRow, created []IntakeForm)
	}{
		{
			name: "list all intake forms without search",
			setup: func(ctx context.Context, qtx *Queries) []IntakeForm {
				var forms []IntakeForm
				for i := 0; i < 3;  i++ {
					regForm := createRandomRegistrationForm(ctx, qtx)
					form, err := qtx.CreateIntakeForm(ctx, CreateIntakeFormParams{
						RegistrationFormID: regForm.ID,
						DateOfIntake:       pgtype.Timestamptz{Valid: true, Time: time.Now()},
						CareType:           IntakeCareTypeEnumAmbulatorySupport,
						IntakeParticipants: []IntakeParticipantsEnum{IntakeParticipantsEnumClient},
						SelfSufficiency:    3,
						IntakeConclusion:   IntakeConclusionEnumSuitable,
					})
					require.NoError(t, err)
					forms = append(forms, form)
				}
				return forms
			},
			params: ListIntakeFormsParams{
				Limit:     10,
				Offset:    0,
				Search:    "",
				SortBy:    "",
				SortOrder: "",
			},
			checks: func(t *testing.T, rows []ListIntakeFormsRow, created []IntakeForm) {
				require.Len(t, rows, 3)
				require.Equal(t, int64(3), rows[0].TotalCount)
				// Check that all created forms are in the list
				ids := make(map[uuid.UUID]bool)
				for _, row := range rows {
					ids[row.ID] = true
				}
				for _, form := range created {
					require.True(t, ids[form.ID])
				}
			},
		},
		{
			name: "list intake forms with limit and offset",
			setup: func(ctx context.Context, qtx *Queries) []IntakeForm {
				var forms []IntakeForm
				for i := 0; i < 5; i++ {
					regForm := createRandomRegistrationForm(ctx, qtx)
					form, err := qtx.CreateIntakeForm(ctx, CreateIntakeFormParams{
						RegistrationFormID: regForm.ID,
						DateOfIntake:       pgtype.Timestamptz{Valid: true, Time: time.Now()},
						CareType:           IntakeCareTypeEnumProtectedLiving,
						IntakeParticipants: []IntakeParticipantsEnum{IntakeParticipantsEnumClient},
						SelfSufficiency:    4,
						IntakeConclusion:   IntakeConclusionEnumFurtherInvestigation,
					})
					require.NoError(t, err)
					forms = append(forms, form)
				}
				return forms
			},
			params: ListIntakeFormsParams{
				Limit:     2,
				Offset:    1,
				Search:    "",
				SortBy:    "",
				SortOrder: "",
			},
			checks: func(t *testing.T, rows []ListIntakeFormsRow, created []IntakeForm) {
				require.Len(t, rows, 2)
				require.Equal(t, int64(5), rows[0].TotalCount)
			},
		},
		{
			name: "list intake forms sorted by created_at desc",
			setup: func(ctx context.Context, qtx *Queries) []IntakeForm {
				var forms []IntakeForm
				for i := 0; i < 3; i++ {
					regForm := createRandomRegistrationForm(ctx, qtx)
					form, err := qtx.CreateIntakeForm(ctx, CreateIntakeFormParams{
						RegistrationFormID: regForm.ID,
						DateOfIntake:       pgtype.Timestamptz{Valid: true, Time: time.Now()},
						CareType:           IntakeCareTypeEnumTrainingCenter,
						IntakeParticipants: []IntakeParticipantsEnum{IntakeParticipantsEnumReferrer},
						SelfSufficiency:    2,
						IntakeConclusion:   IntakeConclusionEnumUnsuitable,
					})
					require.NoError(t, err)
					forms = append(forms, form)
					// Small delay to ensure different created_at times
					time.Sleep(1 * time.Millisecond)
				}
				return forms
			},
			params: ListIntakeFormsParams{
				Limit:     10,
				Offset:    0,
				Search:    "",
				SortBy:    "created_at",
				SortOrder: "desc",
			},
			checks: func(t *testing.T, rows []ListIntakeFormsRow, created []IntakeForm) {
				require.Len(t, rows, 3)
				// Check that rows are sorted by created_at desc (most recent first)
				require.True(t, rows[0].CreatedAt.Time.After(rows[1].CreatedAt.Time) || rows[0].CreatedAt.Time.Equal(rows[1].CreatedAt.Time))
				require.True(t, rows[1].CreatedAt.Time.After(rows[2].CreatedAt.Time) || rows[1].CreatedAt.Time.Equal(rows[2].CreatedAt.Time))
			},
		},
		{
			name: "list intake forms with empty result",
			setup: func(ctx context.Context, qtx *Queries) []IntakeForm {
				return []IntakeForm{} // No forms created
			},
			params: ListIntakeFormsParams{
				Limit:     10,
				Offset:    0,
				Search:    "",
				SortBy:    "",
				SortOrder: "",
			},
			checks: func(t *testing.T, rows []ListIntakeFormsRow, created []IntakeForm) {
				require.Len(t, rows, 0)
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

			created := tt.setup(ctx, qtx)
			rows, err := qtx.ListIntakeForms(ctx, tt.params)
			require.NoError(t, err)

			tt.checks(t, rows, created)
		})
	}
}
