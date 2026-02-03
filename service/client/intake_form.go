package clientp

import (
	"context"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateIntakeForm(ctx context.Context, req *CreateIntakeFormRequest) (*CreateIntakeFormResponse, error) {
	// 1. Validate registration form exists and is processed
	regForm, err := s.Store.GetRegistrationForm(ctx, req.RegistrationFormID)
	if err != nil {
		return nil, fmt.Errorf("registration form not found: %w", err)
	}

	if regForm.FormStatus != db.FormStatusEnumProcessed {
		return nil, errors.New("intake can only be created for processed registration forms")
	}

	arg := db.CreateIntakeFormParams{
		RegistrationFormID:       req.RegistrationFormID,
		DateOfIntake:             pgtype.Timestamptz{Time: req.DateOfIntake, Valid: true},
		CareType:                 req.CareType,
		IntakeParticipants:       req.IntakeParticipants,
		FamilySituation:          req.FamilySituation,
		PsychologicalState:       req.PsychologicalState,
		SelfSufficiency:          req.SelfSufficiency,
		SenderID:                 req.SenderID,
		AssignedLocationID:       req.AssignedLocationID,
		RiskAssessment:           req.RiskAssessment,
		IntakeConclusion:         req.IntakeConclusion,
		IntakeConclusionNotes:    req.IntakeConclusionNotes,
		EvaluationIntervalsWeeks: req.EvaluationIntervalsWeeks,
		Signature:                req.Signature,
	}

	intakeForm, err := s.Store.CreateIntakeForm(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIntakeForm", "Failed to create intake form", zap.Error(err))
		return nil, err
	}
	res := &CreateIntakeFormResponse{
		ID:                       intakeForm.ID,
		RegistrationFormID:       intakeForm.RegistrationFormID,
		DateOfIntake:             intakeForm.DateOfIntake.Time,
		CareType:                 intakeForm.CareType,
		IntakeParticipants:       intakeForm.IntakeParticipants,
		FamilySituation:          intakeForm.FamilySituation,
		PsychologicalState:       intakeForm.PsychologicalState,
		SelfSufficiency:          intakeForm.SelfSufficiency,
		SenderID:                 intakeForm.SenderID,
		AssignedLocationID:       intakeForm.AssignedLocationID,
		RiskAssessment:           intakeForm.RiskAssessment,
		IntakeConclusion:         intakeForm.IntakeConclusion,
		IntakeConclusionNotes:    intakeForm.IntakeConclusionNotes,
		EvaluationIntervalsWeeks: intakeForm.EvaluationIntervalsWeeks,
		Signature:                intakeForm.Signature,
		CreatedAt:                intakeForm.CreatedAt,
		UpdatedAt:                intakeForm.UpdatedAt,
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
			ID:                       intakeForm.ID,
			RegistrationFormID:       intakeForm.RegistrationFormID,
			DateOfIntake:             intakeForm.DateOfIntake.Time,
			CareType:                 intakeForm.CareType,
			IntakeParticipants:       intakeForm.IntakeParticipants,
			FamilySituation:          intakeForm.FamilySituation,
			PsychologicalState:       intakeForm.PsychologicalState,
			SelfSufficiency:          intakeForm.SelfSufficiency,
			SenderID:                 intakeForm.SenderID,
			AssignedLocationID:       intakeForm.AssignedLocationID,
			RiskAssessment:           intakeForm.RiskAssessment,
			IntakeConclusion:         intakeForm.IntakeConclusion,
			IntakeConclusionNotes:    intakeForm.IntakeConclusionNotes,
			EvaluationIntervalsWeeks: intakeForm.EvaluationIntervalsWeeks,
			Signature:                intakeForm.Signature,
			CreatedAt:                intakeForm.CreatedAt,
			UpdatedAt:                intakeForm.UpdatedAt,
			ClientFirstName:          intakeForm.ClientFirstName,
			ClientLastName:           intakeForm.ClientLastName,
			ClientBsnNumber:          intakeForm.ClientBsnNumber,
		})
	}

	totalCount := intakeForms[0].TotalCount

	paginatedRes := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &paginatedRes, nil
}
