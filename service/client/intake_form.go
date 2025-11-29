package clientp

import (
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateIntakeForm(ctx context.Context, req *CreateIntakeFormRequest) (*CreateIntakeFormResponse, error) {
	arg := db.CreateIntakeFormParams{
		RegistrationFormID:    req.RegistrationFormID,
		DateOfIntake:          pgtype.Timestamptz{Time: req.DateOfIntake, Valid: true},
		CareType:              req.CareType,
		IntakeParticipants:    req.IntakeParticipants,
		FamilySituation:       req.FamilySituation,
		PsychologicalState:    req.PsychologicalState,
		SelfSufficiency:       req.SelfSufficiency,
		MaturityMatrixID:      req.MaturityMatrixID,
		Goals:                 req.Goals,
		RiskAssessment:        req.RiskAssessment,
		IntakeConclusion:      req.IntakeConclusion,
		IntakeConclusionNotes: req.IntakeConclusionNotes,
		Signature:             req.Signature,
	}

	intakeForm, err := s.Store.CreateIntakeForm(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIntakeForm", "Failed to create intake form", zap.Error(err))
		return nil, err
	}
	res := &CreateIntakeFormResponse{
		ID:                    intakeForm.ID,
		RegistrationFormID:    intakeForm.RegistrationFormID,
		DateOfIntake:          intakeForm.DateOfIntake.Time,
		CareType:              intakeForm.CareType,
		IntakeParticipants:    intakeForm.IntakeParticipants,
		FamilySituation:       intakeForm.FamilySituation,
		PsychologicalState:    intakeForm.PsychologicalState,
		SelfSufficiency:       intakeForm.SelfSufficiency,
		MaturityMatrixID:      intakeForm.MaturityMatrixID,
		Goals:                 intakeForm.Goals,
		RiskAssessment:        intakeForm.RiskAssessment,
		IntakeConclusion:      intakeForm.IntakeConclusion,
		IntakeConclusionNotes: intakeForm.IntakeConclusionNotes,
		Signature:             intakeForm.Signature,
		CreatedAt:             intakeForm.CreatedAt,
		UpdatedAt:             intakeForm.UpdatedAt,
	}
	return res, nil
}

func (s *clientService) ListIntakeForms(ctx *gin.Context, req *ListIntakeFormsRequest) (*pagination.Response[ListIntakeFormsResponse], error) {
	params := req.GetParams()
	intakeForms, err := s.Store.ListIntakeForms(ctx, db.ListIntakeFormsParams{
		Limit:     params.Limit,
		Offset:    params.Offset,
		Search:    util.DerefString(req.Search),
		SortBy:    "created_at",
		SortOrder: util.DerefString(req.SortOrder),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListIntakeForms", "Failed to list intake forms", zap.Error(err))
		return nil, err
	}

	items := []ListIntakeFormsResponse{}
	for _, intakeForm := range intakeForms {
		items = append(items, ListIntakeFormsResponse{
			ID:                    intakeForm.ID,
			RegistrationFormID:    intakeForm.RegistrationFormID,
			DateOfIntake:          intakeForm.DateOfIntake.Time,
			CareType:              intakeForm.CareType,
			IntakeParticipants:    intakeForm.IntakeParticipants,
			FamilySituation:       intakeForm.FamilySituation,
			PsychologicalState:    intakeForm.PsychologicalState,
			SelfSufficiency:       intakeForm.SelfSufficiency,
			MaturityMatrixID:      intakeForm.MaturityMatrixID,
			Goals:                 intakeForm.Goals,
			RiskAssessment:        intakeForm.RiskAssessment,
			IntakeConclusion:      intakeForm.IntakeConclusion,
			IntakeConclusionNotes: intakeForm.IntakeConclusionNotes,
			Signature:             intakeForm.Signature,
			CreatedAt:             intakeForm.CreatedAt,
			UpdatedAt:             intakeForm.UpdatedAt,
		})
	}

	totalCount := intakeForms[0].TotalCount

	paginatedRes := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &paginatedRes, nil
}
