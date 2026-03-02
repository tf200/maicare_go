package clientp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var ErrIntakeGoalsUpdateBlockedByActiveClient = errors.New("intake goals cannot be updated when an active client exists for this intake form")
var ErrIntakeFormUpdateBlockedByActiveClient = errors.New("intake form cannot be updated when an active client exists for this intake form")
var ErrIntakeFormUpdateConflict = errors.New("intake form was updated by another request")
var ErrNoIntakeFormFieldsToUpdate = errors.New("no intake form fields provided for update")
var ErrInvalidIntakeFormClearField = errors.New("invalid intake form clear field")
var ErrInvalidIntakeFormUpdate = errors.New("invalid intake form update")

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
		UpdatedAt:                intakeForm.UpdatedAt,
	}
	return res, nil
}

func (s *clientService) ListIntakeForms(ctx *gin.Context, req *ListIntakeFormsRequest) (*pagination.Response[ListIntakeFormsResponse], error) {
	params := req.GetParams()
	status := db.NullIntakeConclusionEnum{}
	if req.Status != nil {
		status = db.NullIntakeConclusionEnum{
			IntakeConclusionEnum: *req.Status,
			Valid:                true,
		}
	}
	intakeForms, err := s.Store.ListIntakeForms(ctx, db.ListIntakeFormsParams{
		Limit:     params.Limit,
		Offset:    params.Offset,
		Search:    util.DerefString(req.Search),
		Status:    status,
		SortBy:    "created_at",
		SortOrder: util.DerefString(req.SortOrder),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListIntakeForms", "Failed to list intake forms", zap.Error(err))
		return nil, err
	}

	items := []ListIntakeFormsResponse{}
	for _, intakeForm := range intakeForms {
		var assignedLocationAddress *AssignedLocationAddress
		if intakeForm.AssignedLocationID != nil {
			assignedLocationAddress = &AssignedLocationAddress{
				Street:              intakeForm.AssignedLocationStreet,
				HouseNumber:         intakeForm.AssignedLocationHouseNumber,
				HouseNumberAddition: intakeForm.AssignedLocationHouseNumberAddition,
				PostalCode:          intakeForm.AssignedLocationPostalCode,
				City:                intakeForm.AssignedLocationCity,
			}
		}
		items = append(items, ListIntakeFormsResponse{
			ID:                      intakeForm.ID,
			RegistrationFormID:      intakeForm.RegistrationFormID,
			DateOfIntake:            intakeForm.DateOfIntake.Time,
			ClientFirstName:         intakeForm.ClientFirstName,
			ClientLastName:          intakeForm.ClientLastName,
			ClientBsnNumber:         intakeForm.ClientBsnNumber,
			IntakeStatus:            intakeForm.IntakeConclusion,
			GoalAssessmentDone:      intakeForm.GoalAssessmentDone,
			CareType:                intakeForm.CareType,
			AssignedLocationID:      intakeForm.AssignedLocationID,
			AssignedLocationAddress: assignedLocationAddress,
		})
	}

	totalCount := int64(0)
	if len(intakeForms) > 0 {
		totalCount = intakeForms[0].TotalCount
	}

	paginatedRes := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &paginatedRes, nil
}

func (s *clientService) GetIntakeFormTotals(ctx context.Context) (*GetIntakeFormTotalsResponse, error) {
	totals, err := s.Store.GetIntakeFormTotals(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetIntakeFormTotals", "Failed to get intake form totals", zap.Error(err))
		return nil, err
	}

	return &GetIntakeFormTotalsResponse{
		FurtherInvestigationTotal: totals.FurtherInvestigationTotal,
		WithoutGoalsTotal:         totals.WithoutGoalsTotal,
	}, nil
}

func (s *clientService) GetIntakeForm(ctx context.Context, intakeFormID uuid.UUID) (*GetIntakeFormResponse, error) {
	intakeForm, err := s.Store.GetIntakeFormDetails(ctx, intakeFormID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetIntakeForm", "Failed to get intake form", zap.Error(err))
		return nil, err
	}

	assessments, err := s.Store.GetIntakeTopicsAssessments(ctx, intakeFormID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetIntakeForm", "Failed to get intake assessments", zap.Error(err))
		return nil, err
	}

	intakeGoalsAssigned := make([]IntakeGoalTopic, 0, len(assessments))
	for _, assessment := range assessments {
		goals, err := jsonToGoals(assessment.ProposedGoals)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetIntakeForm", "Failed to parse proposed goals", zap.Error(err))
			return nil, err
		}
		intakeGoalsAssigned = append(intakeGoalsAssigned, IntakeGoalTopic{
			AssessmentID:  assessment.ID,
			TopicID:       assessment.TopicID,
			TopicName:     assessment.TopicName,
			CurrentLevel:  assessment.CurrentLevel,
			ProposedGoals: goals,
			Notes:         assessment.Notes,
		})
	}

	var location *IntakeFormLocationDetails
	if intakeForm.LocationName != nil || intakeForm.LocationStreet != nil || intakeForm.LocationHouseNumber != nil || intakeForm.LocationPostalCode != nil || intakeForm.LocationCity != nil {
		location = &IntakeFormLocationDetails{
			Name:                intakeForm.LocationName,
			Street:              intakeForm.LocationStreet,
			HouseNumber:         intakeForm.LocationHouseNumber,
			HouseNumberAddition: intakeForm.LocationHouseNumberAddition,
			PostalCode:          intakeForm.LocationPostalCode,
			City:                intakeForm.LocationCity,
		}
	}

	res := &GetIntakeFormResponse{
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
		CreatedAt:                intakeForm.CreatedAt.Time,
		UpdatedAt:                intakeForm.UpdatedAt.Time,
		ClientFirstName:          intakeForm.ClientFirstName,
		ClientLastName:           intakeForm.ClientLastName,
		ClientBsnNumber:          intakeForm.ClientBsnNumber,
		DesiredGoals:             intakeForm.ClientGoals,
		SenderName:               intakeForm.SenderName,
		Location:                 location,
		IntakeGoalsAssigned:      intakeGoalsAssigned,
		HasClient:                intakeForm.HasClient,
	}

	return res, nil
}

func (s *clientService) UpdateIntakeForm(ctx context.Context, intakeFormID uuid.UUID, req *UpdateIntakeFormRequest) (*UpdateIntakeFormResponse, error) {
	clearFieldSet := map[string]bool{}
	for _, field := range req.ClearFields {
		normalized := strings.ToLower(strings.TrimSpace(field))
		switch normalized {
		case "family_situation", "psychological_state", "sender_id", "assigned_location_id", "risk_assessment", "signature":
			clearFieldSet[normalized] = true
		default:
			return nil, fmt.Errorf("%w: %s", ErrInvalidIntakeFormClearField, field)
		}
	}

	hasChanges := req.DateOfIntake != nil || req.CareType != nil || req.IntakeParticipants != nil || req.FamilySituation != nil || req.PsychologicalState != nil || req.SelfSufficiency != nil || req.SenderID != nil || req.AssignedLocationID != nil || req.RiskAssessment != nil || req.EvaluationIntervalsWeeks != nil || req.Signature != nil || len(clearFieldSet) > 0
	if !hasChanges {
		return nil, ErrNoIntakeFormFieldsToUpdate
	}

	if req.SelfSufficiency != nil && (*req.SelfSufficiency < 0 || *req.SelfSufficiency > 5) {
		return nil, fmt.Errorf("%w: self_sufficiency must be between 0 and 5", ErrInvalidIntakeFormUpdate)
	}

	if req.EvaluationIntervalsWeeks != nil && *req.EvaluationIntervalsWeeks < 0 {
		return nil, fmt.Errorf("%w: evaluation_intervals_weeks must be greater than or equal to 0", ErrInvalidIntakeFormUpdate)
	}

	arg := db.UpdateIntakeFormParams{
		DateOfIntake: func() pgtype.Timestamptz {
			if req.DateOfIntake != nil {
				return pgtype.Timestamptz{Time: *req.DateOfIntake, Valid: true}
			}
			return pgtype.Timestamptz{Valid: false}
		}(),
		CareType: func() db.NullIntakeCareTypeEnum {
			if req.CareType != nil {
				return db.NullIntakeCareTypeEnum{IntakeCareTypeEnum: *req.CareType, Valid: true}
			}
			return db.NullIntakeCareTypeEnum{Valid: false}
		}(),
		IntakeParticipants: func() []db.IntakeParticipantsEnum {
			if req.IntakeParticipants != nil {
				return *req.IntakeParticipants
			}
			return nil
		}(),
		ClearFamilySituation:     clearFieldSet["family_situation"],
		FamilySituation:          req.FamilySituation,
		ClearPsychologicalState:  clearFieldSet["psychological_state"],
		PsychologicalState:       req.PsychologicalState,
		SelfSufficiency:          req.SelfSufficiency,
		ClearSenderID:            clearFieldSet["sender_id"],
		SenderID:                 req.SenderID,
		ClearAssignedLocationID:  clearFieldSet["assigned_location_id"],
		AssignedLocationID:       req.AssignedLocationID,
		ClearRiskAssessment:      clearFieldSet["risk_assessment"],
		RiskAssessment:           req.RiskAssessment,
		EvaluationIntervalsWeeks: req.EvaluationIntervalsWeeks,
		ClearSignature:           clearFieldSet["signature"],
		Signature:                req.Signature,
		ID:                       intakeFormID,
	}

	intakeForm, err := s.Store.UpdateIntakeForm(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			hasActiveClient, activeClientErr := s.Store.HasActiveClientByIntakeFormID(ctx, &intakeFormID)
			if activeClientErr == nil && hasActiveClient {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "UpdateIntakeForm", "Blocked intake form update due to active client", zap.String("intake_form_id", intakeFormID.String()))
				return nil, ErrIntakeFormUpdateBlockedByActiveClient
			}

			_, getErr := s.Store.GetIntakeForm(ctx, intakeFormID)
			if errors.Is(getErr, pgx.ErrNoRows) {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "UpdateIntakeForm", "Intake form not found for update", zap.String("intake_form_id", intakeFormID.String()))
				return nil, pgx.ErrNoRows
			}
			if getErr != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateIntakeForm", "Failed to resolve intake form update conflict", zap.Error(getErr))
				return nil, getErr
			}

			return nil, ErrIntakeFormUpdateConflict
		}

		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateIntakeForm", "Failed to update intake form", zap.Error(err))
		return nil, err
	}

	return &UpdateIntakeFormResponse{
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
		UpdatedAt:                intakeForm.UpdatedAt.Time,
	}, nil
}

func (s *clientService) CreateIntakeFormGoals(ctx context.Context, intakeFormID uuid.UUID, req *CreateIntakeFormGoalsRequest) (*CreateIntakeFormGoalsResponse, error) {
	var itemsJSON []byte
	var err error
	if len(req.Assessments) > 0 {
		itemsJSON, err = json.Marshal(req.Assessments)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIntakeFormGoals", "Failed to marshal intake goals", zap.Error(err))
			return nil, err
		}
	}

	rows := make([]db.CreateIntakeTopicAssessmentsBatchRow, 0)
	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if _, err := q.LockIntakeFormByID(ctx, intakeFormID); err != nil {
			return err
		}

		hasActiveClient, err := q.HasActiveClientByIntakeFormID(ctx, &intakeFormID)
		if err != nil {
			return err
		}
		if hasActiveClient {
			return ErrIntakeGoalsUpdateBlockedByActiveClient
		}

		if err := q.DeleteIntakeTopicAssessmentsByIntakeForm(ctx, intakeFormID); err != nil {
			return err
		}

		if len(req.Assessments) == 0 {
			rows = []db.CreateIntakeTopicAssessmentsBatchRow{}
			return nil
		}

		rows, err = q.CreateIntakeTopicAssessmentsBatch(ctx, db.CreateIntakeTopicAssessmentsBatchParams{
			IntakeFormID: intakeFormID,
			Items:        itemsJSON,
		})
		return err
	})
	if err != nil {
		if errors.Is(err, ErrIntakeGoalsUpdateBlockedByActiveClient) {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "CreateIntakeFormGoals", "Blocked intake goals update due to active client", zap.String("intake_form_id", intakeFormID.String()))
			return nil, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "CreateIntakeFormGoals", "Intake form not found for goals update", zap.String("intake_form_id", intakeFormID.String()))
			return nil, err
		}
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIntakeFormGoals", "Failed to replace intake goals", zap.Error(err))
		return nil, err
	}

	assessments := make([]ListIntakeMaturityAssessmentsResponse, 0, len(rows))
	for _, row := range rows {
		goals, err := jsonToGoals(row.ProposedGoals)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIntakeFormGoals", "Failed to parse proposed goals", zap.Error(err))
			return nil, err
		}
		assessments = append(assessments, ListIntakeMaturityAssessmentsResponse{
			ID:            row.ID,
			IntakeFormID:  row.IntakeFormID,
			TopicID:       row.TopicID,
			TopicName:     row.TopicName,
			CurrentLevel:  row.CurrentLevel,
			ProposedGoals: goals,
			Notes:         row.Notes,
			CreatedAt:     row.CreatedAt,
		})
	}

	return &CreateIntakeFormGoalsResponse{Assessments: assessments}, nil
}

func (s *clientService) UpdateIntakeConclusion(ctx context.Context, intakeFormID uuid.UUID, req *UpdateIntakeConclusionRequest) (*UpdateIntakeConclusionResponse, error) {
	decision := strings.ToLower(strings.TrimSpace(req.Decision))

	var conclusion db.IntakeConclusionEnum
	switch decision {
	case "accept":
		conclusion = db.IntakeConclusionEnumSuitable
	case "refuse":
		conclusion = db.IntakeConclusionEnumUnsuitable
	default:
		return nil, fmt.Errorf("invalid decision: %s", req.Decision)
	}

	updated, err := s.Store.UpdateIntakeConclusion(ctx, db.UpdateIntakeConclusionParams{
		ID:                    intakeFormID,
		IntakeConclusion:      conclusion,
		IntakeConclusionNotes: req.IntakeConclusionNotes,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateIntakeConclusion", "Failed to update intake conclusion", zap.Error(err))
		return nil, err
	}

	return &UpdateIntakeConclusionResponse{
		ID:                    updated.ID,
		IntakeConclusion:      updated.IntakeConclusion,
		IntakeConclusionNotes: updated.IntakeConclusionNotes,
		UpdatedAt:             updated.UpdatedAt.Time,
	}, nil
}
