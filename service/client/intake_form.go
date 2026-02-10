package clientp

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateIntakeForm(ctx context.Context, req *CreateIntakeFormRequest) (*CreateIntakeFormResponse, error) {
	arg := db.CreateIntakeFormParams{
		RegistrationFormID:    req.RegistrationFormID,
		DateOfIntake:          pgtype.Timestamptz{Time: req.DateOfIntake, Valid: true},
		CareType:              intakeCareTypeFromEnum(req.CareType),
		IntakeParticipants:    req.IntakeParticipants,
		FamilySituation:       req.FamilySituation,
		PsychologicalState:    req.PsychologicalState,
		SelfSufficiency:       &req.SelfSufficiency,
		MaturityMatrixID:      req.MaturityMatrixID,
		Goals:                 req.Goals,
		RiskAssessment:        req.RiskAssessment,
		IntakeConclusion:      intakeConclusionFromEnum(req.IntakeConclusion),
		IntakeConclusionNotes: req.IntakeConclusionNotes,
		Signature:             req.Signature,
		Status:                db.IntakeStatusEnumInProgress,
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
		CareType:              intakeCareTypeValue(intakeForm.CareType),
		IntakeParticipants:    intakeForm.IntakeParticipants,
		FamilySituation:       intakeForm.FamilySituation,
		PsychologicalState:    intakeForm.PsychologicalState,
		SelfSufficiency:       util.DerefInt32(intakeForm.SelfSufficiency),
		MaturityMatrixID:      intakeForm.MaturityMatrixID,
		Goals:                 intakeForm.Goals,
		RiskAssessment:        intakeForm.RiskAssessment,
		IntakeConclusion:      intakeConclusionValue(intakeForm.IntakeConclusion),
		IntakeConclusionNotes: intakeForm.IntakeConclusionNotes,
		Signature:             intakeForm.Signature,
		Status:                intakeForm.Status,
		UrgencyLevel:          urgencyLevelPtrFromNull(intakeForm.UrgencyLevel),
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
			CareType:              intakeCareTypeValue(intakeForm.CareType),
			IntakeParticipants:    intakeForm.IntakeParticipants,
			FamilySituation:       intakeForm.FamilySituation,
			PsychologicalState:    intakeForm.PsychologicalState,
			SelfSufficiency:       util.DerefInt32(intakeForm.SelfSufficiency),
			MaturityMatrixID:      intakeForm.MaturityMatrixID,
			Goals:                 intakeForm.Goals,
			RiskAssessment:        intakeForm.RiskAssessment,
			IntakeConclusion:      intakeConclusionValue(intakeForm.IntakeConclusion),
			IntakeConclusionNotes: intakeForm.IntakeConclusionNotes,
			Signature:             intakeForm.Signature,
			Status:                intakeForm.Status,
			UrgencyLevel:          urgencyLevelPtrFromNull(intakeForm.UrgencyLevel),
			CreatedAt:             intakeForm.CreatedAt,
			UpdatedAt:             intakeForm.UpdatedAt,
		})
	}

	totalCount := intakeForms[0].TotalCount

	paginatedRes := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &paginatedRes, nil
}

func (s *clientService) CompleteIntakeForm(ctx context.Context, req *CompleteIntakeFormRequest, intakeID uuid.UUID) (*CompleteIntakeFormResponse, error) {
	if req.Outcome == "accepted" && req.ClientDetails == nil {
		return nil, fmt.Errorf("client_details is required when outcome is accepted")
	}

	var conclusion db.IntakeConclusionEnum
	switch req.Outcome {
	case "accepted":
		conclusion = db.IntakeConclusionEnumSuitable
	case "rejected":
		conclusion = db.IntakeConclusionEnumUnsuitable
	case "further_investigation":
		conclusion = db.IntakeConclusionEnumFurtherInvestigation
	default:
		return nil, fmt.Errorf("invalid outcome")
	}

	urgency := urgencyLevelFromPtr(req.UrgencyLevel)
	updatedIntake, err := s.Store.UpdateIntakeFormOutcome(ctx, db.UpdateIntakeFormOutcomeParams{
		ID:                    intakeID,
		Status:                db.IntakeStatusEnumCompleted,
		IntakeConclusion:      intakeConclusionFromEnum(conclusion),
		UrgencyLevel:          urgency,
		IntakeConclusionNotes: req.ReportSummary,
		RiskAssessment:        req.RiskAssessment,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CompleteIntakeForm", "Failed to complete intake form", zap.Error(err), zap.String("IntakeID", intakeID.String()))
		return nil, fmt.Errorf("failed to complete intake form: %w", err)
	}

	var clientID *uuid.UUID
	if req.Outcome == "accepted" {
		clientReq := *req.ClientDetails
		clientReq.IntakeFormID = &updatedIntake.ID
		createdClient, err := s.CreateClientDetails(clientReq, ctx)
		if err != nil {
			return nil, err
		}
		clientID = &createdClient.ID
	} else {
		status := "in_review"
		if req.Outcome == "rejected" {
			status = "rejected"
		}
		_, err = s.Store.UpdateRegistrationFormStatusOnly(ctx, db.UpdateRegistrationFormStatusOnlyParams{
			ID:         updatedIntake.RegistrationFormID,
			FormStatus: db.FormStatusEnum(status),
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CompleteIntakeForm", "Failed to update registration status", zap.Error(err), zap.String("RegistrationFormID", updatedIntake.RegistrationFormID.String()))
		}
	}

	return &CompleteIntakeFormResponse{
		IntakeID:           updatedIntake.ID,
		RegistrationFormID: updatedIntake.RegistrationFormID,
		ClientID:           clientID,
		Outcome:            req.Outcome,
		Status:             updatedIntake.Status,
	}, nil
}

func intakeCareTypeFromEnum(value db.IntakeCareTypeEnum) db.NullIntakeCareTypeEnum {
	return db.NullIntakeCareTypeEnum{IntakeCareTypeEnum: value, Valid: true}
}

func intakeCareTypeValue(value db.NullIntakeCareTypeEnum) db.IntakeCareTypeEnum {
	if value.Valid {
		return value.IntakeCareTypeEnum
	}
	return ""
}

func intakeConclusionFromEnum(value db.IntakeConclusionEnum) db.NullIntakeConclusionEnum {
	return db.NullIntakeConclusionEnum{IntakeConclusionEnum: value, Valid: true}
}

func intakeConclusionValue(value db.NullIntakeConclusionEnum) db.IntakeConclusionEnum {
	if value.Valid {
		return value.IntakeConclusionEnum
	}
	return ""
}

func urgencyLevelFromPtr(value *string) db.NullUrgencyLevelEnum {
	if value == nil {
		return db.NullUrgencyLevelEnum{Valid: false}
	}
	return db.NullUrgencyLevelEnum{UrgencyLevelEnum: db.UrgencyLevelEnum(*value), Valid: true}
}

func urgencyLevelPtrFromNull(value db.NullUrgencyLevelEnum) *db.UrgencyLevelEnum {
	if !value.Valid {
		return nil
	}
	level := value.UrgencyLevelEnum
	return &level
}
