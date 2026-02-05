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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (s *clientService) GetIntakeForm(ctx context.Context, intakeFormID uuid.UUID) (*GetIntakeFormResponse, error) {
	intakeForm, err := s.Store.GetIntakeFormDetails(ctx, intakeFormID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetIntakeForm", "Failed to get intake form", zap.Error(err))
		return nil, err
	}

	assessments, err := s.Store.GetIntakeMaturityAssessments(ctx, intakeFormID)
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
	}

	return res, nil
}

func (s *clientService) CreateIntakeFormGoals(ctx context.Context, intakeFormID uuid.UUID, req *CreateIntakeFormGoalsRequest) (*CreateIntakeFormGoalsResponse, error) {
	if len(req.Assessments) == 0 {
		return &CreateIntakeFormGoalsResponse{Assessments: []ListIntakeMaturityAssessmentsResponse{}}, nil
	}

	itemsJSON, err := json.Marshal(req.Assessments)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIntakeFormGoals", "Failed to marshal intake goals", zap.Error(err))
		return nil, err
	}

	rows, err := s.Store.CreateIntakeMaturityAssessmentsBatch(ctx, db.CreateIntakeMaturityAssessmentsBatchParams{
		IntakeFormID: intakeFormID,
		Items:        itemsJSON,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIntakeFormGoals", "Failed to create intake goals", zap.Error(err))
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
