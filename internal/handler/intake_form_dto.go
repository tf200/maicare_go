package handler

import (
	"time"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	db "maicare_go/db/sqlc"

	"github.com/google/uuid"
)

// ==================== Request DTOs ====================

type createIntakeFormRequest struct {
	RegistrationFormID       uuid.UUID                   `json:"registration_form_id" binding:"required,uuid"`
	DateOfIntake             time.Time                   `json:"date_of_intake" binding:"required"`
	CareType                 db.IntakeCareTypeEnum       `json:"care_type" binding:"required"`
	IntakeParticipants       []db.IntakeParticipantsEnum `json:"intake_participants"`
	FamilySituation          *string                     `json:"family_situation"`
	PsychologicalState       *string                     `json:"psychological_state"`
	SelfSufficiency          int32                       `json:"self_sufficiency"`
	SenderID                 *uuid.UUID                  `json:"sender_id"`
	AssignedLocationID       *uuid.UUID                  `json:"assigned_location_id"`
	RiskAssessment           *string                     `json:"risk_assessment"`
	IntakeConclusion         db.IntakeConclusionEnum     `json:"intake_conclusion" binding:"required"`
	IntakeConclusionNotes    *string                     `json:"intake_conclusion_notes"`
	EvaluationIntervalsWeeks int32                       `json:"evaluation_intervals_weeks"`
	Signature                *string                     `json:"signature"`
}

type listIntakeFormsRequest struct {
	httpapi.PageRequest
	Search    *string                  `form:"search" binding:"omitempty"`
	Status    *db.IntakeConclusionEnum `form:"status" binding:"omitempty"`
	SortOrder *string                  `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

type updateIntakeFormRequest struct {
	DateOfIntake             *time.Time                   `json:"date_of_intake"`
	CareType                 *db.IntakeCareTypeEnum       `json:"care_type"`
	IntakeParticipants       *[]db.IntakeParticipantsEnum `json:"intake_participants"`
	FamilySituation          *string                      `json:"family_situation"`
	PsychologicalState       *string                      `json:"psychological_state"`
	SelfSufficiency          *int32                       `json:"self_sufficiency"`
	SenderID                 *uuid.UUID                   `json:"sender_id"`
	AssignedLocationID       *uuid.UUID                   `json:"assigned_location_id"`
	RiskAssessment           *string                      `json:"risk_assessment"`
	EvaluationIntervalsWeeks *int32                       `json:"evaluation_intervals_weeks"`
	Signature                *string                      `json:"signature"`
	ClearFields              []string                     `json:"clear_fields"`
}

type intakeFormGoalItemRequest struct {
	TopicID       uuid.UUID                     `json:"topic_id" binding:"required,uuid"`
	CurrentLevel  int32                         `json:"current_level" binding:"required"`
	ProposedGoals []domain.IntakeAssessmentGoal `json:"proposed_goals"`
	Notes         *string                       `json:"notes"`
}

type createIntakeFormGoalsRequest struct {
	Assessments []intakeFormGoalItemRequest `json:"assessments" binding:"required,dive"`
}

type generateIntakeGoalsRequest struct {
	TopicID      uuid.UUID `json:"topic_id" binding:"required,uuid"`
	CurrentLevel int       `json:"current_level" binding:"required"`
	UserDesc     string    `json:"user_desc" binding:"required"`
}

type updateIntakeConclusionRequest struct {
	Decision              string  `json:"decision" binding:"required,oneof=accept refuse"`
	IntakeConclusionNotes *string `json:"intake_conclusion_notes"`
}


// ==================== Response DTOs ====================

type intakeAssessmentGoalResponse struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority"`
}

type assignedLocationAddressResponse struct {
	Street              *string `json:"street"`
	HouseNumber         *string `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          *string `json:"postal_code"`
	City                *string `json:"city"`
}

type intakeFormLocationDetailsResponse struct {
	Name                *string `json:"name"`
	Street              *string `json:"street"`
	HouseNumber         *string `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          *string `json:"postal_code"`
	City                *string `json:"city"`
}

type intakeGoalTopicResponse struct {
	AssessmentID  uuid.UUID                      `json:"assessment_id"`
	TopicID       uuid.UUID                      `json:"topic_id"`
	TopicName     string                         `json:"topic_name"`
	CurrentLevel  int32                          `json:"current_level"`
	ProposedGoals []intakeAssessmentGoalResponse `json:"proposed_goals"`
	Notes         *string                        `json:"notes"`
}

type intakeMaturityAssessmentResponse struct {
	ID            uuid.UUID                      `json:"id"`
	IntakeFormID  uuid.UUID                      `json:"intake_form_id"`
	TopicID       uuid.UUID                      `json:"topic_id"`
	TopicName     string                         `json:"topic_name"`
	CurrentLevel  int32                          `json:"current_level"`
	ProposedGoals []intakeAssessmentGoalResponse `json:"proposed_goals"`
	Notes         *string                        `json:"notes"`
	CreatedAt     time.Time                      `json:"created_at"`
}

type intakeFormResponse struct {
	ID                       uuid.UUID                   `json:"id"`
	RegistrationFormID       uuid.UUID                   `json:"registration_form_id"`
	DateOfIntake             time.Time                   `json:"date_of_intake"`
	CareType                 db.IntakeCareTypeEnum       `json:"care_type"`
	IntakeParticipants       []db.IntakeParticipantsEnum `json:"intake_participants"`
	FamilySituation          *string                     `json:"family_situation"`
	PsychologicalState       *string                     `json:"psychological_state"`
	SelfSufficiency          int32                       `json:"self_sufficiency"`
	SenderID                 *uuid.UUID                  `json:"sender_id"`
	AssignedLocationID       *uuid.UUID                  `json:"assigned_location_id"`
	RiskAssessment           *string                     `json:"risk_assessment"`
	IntakeConclusion         db.IntakeConclusionEnum     `json:"intake_conclusion"`
	IntakeConclusionNotes    *string                     `json:"intake_conclusion_notes"`
	EvaluationIntervalsWeeks int32                       `json:"evaluation_intervals_weeks"`
	Signature                *string                     `json:"signature"`
	UpdatedAt                time.Time                   `json:"updated_at"`
}

type intakeFormListItemResponse struct {
	ID                      uuid.UUID                        `json:"id"`
	RegistrationFormID      uuid.UUID                        `json:"registration_form_id"`
	DateOfIntake            time.Time                        `json:"date_of_intake"`
	ClientFirstName         string                           `json:"client_first_name"`
	ClientLastName          string                           `json:"client_last_name"`
	ClientBsnNumber         string                           `json:"client_bsn_number"`
	IntakeStatus            db.IntakeConclusionEnum          `json:"intake_status"`
	GoalAssessmentDone      bool                             `json:"goal_assessment_done"`
	CareType                db.IntakeCareTypeEnum            `json:"care_type"`
	AssignedLocationID      *uuid.UUID                       `json:"assigned_location_id"`
	AssignedLocationAddress *assignedLocationAddressResponse `json:"assigned_location_address"`
}

type intakeFormTotalsResponse struct {
	FurtherInvestigationTotal int64 `json:"further_investigation_total"`
	WithoutGoalsTotal         int64 `json:"without_goals_total"`
}

type intakeFormDetailResponse struct {
	ID                       uuid.UUID                          `json:"id"`
	RegistrationFormID       uuid.UUID                          `json:"registration_form_id"`
	DateOfIntake             time.Time                          `json:"date_of_intake"`
	CareType                 db.IntakeCareTypeEnum              `json:"care_type"`
	IntakeParticipants       []db.IntakeParticipantsEnum        `json:"intake_participants"`
	FamilySituation          *string                            `json:"family_situation"`
	PsychologicalState       *string                            `json:"psychological_state"`
	SelfSufficiency          int32                              `json:"self_sufficiency"`
	SenderID                 *uuid.UUID                         `json:"sender_id"`
	AssignedLocationID       *uuid.UUID                         `json:"assigned_location_id"`
	RiskAssessment           *string                            `json:"risk_assessment"`
	IntakeConclusion         db.IntakeConclusionEnum            `json:"intake_conclusion"`
	IntakeConclusionNotes    *string                            `json:"intake_conclusion_notes"`
	EvaluationIntervalsWeeks int32                              `json:"evaluation_intervals_weeks"`
	Signature                *string                            `json:"signature"`
	CreatedAt                time.Time                          `json:"created_at"`
	UpdatedAt                time.Time                          `json:"updated_at"`
	ClientFirstName          string                             `json:"client_first_name"`
	ClientLastName           string                             `json:"client_last_name"`
	ClientBsnNumber          string                             `json:"client_bsn_number"`
	DesiredGoals             []string                           `json:"desired_goals"`
	SenderName               *string                            `json:"sender_name"`
	Location                 *intakeFormLocationDetailsResponse `json:"location"`
	IntakeGoalsAssigned      []intakeGoalTopicResponse          `json:"intake_goals_assigned"`
	HasClient                bool                               `json:"has_client"`
}

type intakeFormConclusionResponse struct {
	ID                    uuid.UUID               `json:"id"`
	IntakeConclusion      db.IntakeConclusionEnum `json:"intake_conclusion"`
	IntakeConclusionNotes *string                 `json:"intake_conclusion_notes"`
	UpdatedAt             time.Time               `json:"updated_at"`
}

type generateIntakeGoalsResponse struct {
	Goals []intakeAssessmentGoalResponse `json:"goals"`
}

type promoteIntakeToClientResponse struct {
	ClientID                   uuid.UUID `json:"client_id"`
	IntakeFormID               uuid.UUID `json:"intake_form_id"`
	RegistrationFormID         uuid.UUID `json:"registration_form_id"`
	Message                    string    `json:"message"`
	MaturityAssessmentsCreated int       `json:"maturity_assessments_created"`
	EmergencyContactsCreated   int       `json:"emergency_contacts_created"`
}

type intakeFormGoalsResponse struct {
	Assessments []intakeMaturityAssessmentResponse `json:"assessments"`
}

// ==================== Mappers ====================

func toIntakeFormResponse(res *domain.IntakeForm) intakeFormResponse {
	return intakeFormResponse{
		ID:                       res.ID,
		RegistrationFormID:       res.RegistrationFormID,
		DateOfIntake:             res.DateOfIntake,
		CareType:                 res.CareType,
		IntakeParticipants:       res.IntakeParticipants,
		FamilySituation:          res.FamilySituation,
		PsychologicalState:       res.PsychologicalState,
		SelfSufficiency:          res.SelfSufficiency,
		SenderID:                 res.SenderID,
		AssignedLocationID:       res.AssignedLocationID,
		RiskAssessment:           res.RiskAssessment,
		IntakeConclusion:         res.IntakeConclusion,
		IntakeConclusionNotes:    res.IntakeConclusionNotes,
		EvaluationIntervalsWeeks: res.EvaluationIntervalsWeeks,
		Signature:                res.Signature,
		UpdatedAt:                res.UpdatedAt,
	}
}

func toAssignedLocationAddressResponse(addr *domain.AssignedLocationAddress) *assignedLocationAddressResponse {
	if addr == nil {
		return nil
	}
	return &assignedLocationAddressResponse{
		Street:              addr.Street,
		HouseNumber:         addr.HouseNumber,
		HouseNumberAddition: addr.HouseNumberAddition,
		PostalCode:          addr.PostalCode,
		City:                addr.City,
	}
}

func toIntakeFormListItemResponse(item domain.IntakeFormListItem) intakeFormListItemResponse {
	return intakeFormListItemResponse{
		ID:                      item.ID,
		RegistrationFormID:      item.RegistrationFormID,
		DateOfIntake:            item.DateOfIntake,
		ClientFirstName:         item.ClientFirstName,
		ClientLastName:          item.ClientLastName,
		ClientBsnNumber:         item.ClientBsnNumber,
		IntakeStatus:            item.IntakeStatus,
		GoalAssessmentDone:      item.GoalAssessmentDone,
		CareType:                item.CareType,
		AssignedLocationID:      item.AssignedLocationID,
		AssignedLocationAddress: toAssignedLocationAddressResponse(item.AssignedLocationAddress),
	}
}

func toIntakeAssessmentGoalResponse(goal domain.IntakeAssessmentGoal) intakeAssessmentGoalResponse {
	return intakeAssessmentGoalResponse{
		ID:          goal.ID,
		Title:       goal.Title,
		Description: goal.Description,
		Priority:    goal.Priority,
	}
}

func toIntakeGoalTopicResponse(topic domain.IntakeGoalTopic) intakeGoalTopicResponse {
	goals := make([]intakeAssessmentGoalResponse, len(topic.ProposedGoals))
	for i, g := range topic.ProposedGoals {
		goals[i] = toIntakeAssessmentGoalResponse(g)
	}
	return intakeGoalTopicResponse{
		AssessmentID:  topic.AssessmentID,
		TopicID:       topic.TopicID,
		TopicName:     topic.TopicName,
		CurrentLevel:  topic.CurrentLevel,
		ProposedGoals: goals,
		Notes:         topic.Notes,
	}
}

func toIntakeFormLocationDetailsResponse(loc *domain.IntakeFormLocationDetails) *intakeFormLocationDetailsResponse {
	if loc == nil {
		return nil
	}
	return &intakeFormLocationDetailsResponse{
		Name:                loc.Name,
		Street:              loc.Street,
		HouseNumber:         loc.HouseNumber,
		HouseNumberAddition: loc.HouseNumberAddition,
		PostalCode:          loc.PostalCode,
		City:                loc.City,
	}
}

func toIntakeFormDetailResponse(res *domain.IntakeFormDetail) intakeFormDetailResponse {
	goals := make([]intakeGoalTopicResponse, len(res.IntakeGoalsAssigned))
	for i, g := range res.IntakeGoalsAssigned {
		goals[i] = toIntakeGoalTopicResponse(g)
	}

	return intakeFormDetailResponse{
		ID:                       res.ID,
		RegistrationFormID:       res.RegistrationFormID,
		DateOfIntake:             res.DateOfIntake,
		CareType:                 res.CareType,
		IntakeParticipants:       res.IntakeParticipants,
		FamilySituation:          res.FamilySituation,
		PsychologicalState:       res.PsychologicalState,
		SelfSufficiency:          res.SelfSufficiency,
		SenderID:                 res.SenderID,
		AssignedLocationID:       res.AssignedLocationID,
		RiskAssessment:           res.RiskAssessment,
		IntakeConclusion:         res.IntakeConclusion,
		IntakeConclusionNotes:    res.IntakeConclusionNotes,
		EvaluationIntervalsWeeks: res.EvaluationIntervalsWeeks,
		Signature:                res.Signature,
		CreatedAt:                res.CreatedAt,
		UpdatedAt:                res.UpdatedAt,
		ClientFirstName:          res.ClientFirstName,
		ClientLastName:           res.ClientLastName,
		ClientBsnNumber:          res.ClientBsnNumber,
		DesiredGoals:             res.DesiredGoals,
		SenderName:               res.SenderName,
		Location:                 toIntakeFormLocationDetailsResponse(res.Location),
		IntakeGoalsAssigned:      goals,
		HasClient:                res.HasClient,
	}
}

func toIntakeFormConclusionResponse(res *domain.IntakeFormConclusion) intakeFormConclusionResponse {
	return intakeFormConclusionResponse{
		ID:                    res.ID,
		IntakeConclusion:      res.IntakeConclusion,
		IntakeConclusionNotes: res.IntakeConclusionNotes,
		UpdatedAt:             res.UpdatedAt,
	}
}

func toIntakeMaturityAssessmentResponse(asm domain.IntakeMaturityAssessment) intakeMaturityAssessmentResponse {
	goals := make([]intakeAssessmentGoalResponse, len(asm.ProposedGoals))
	for i, g := range asm.ProposedGoals {
		goals[i] = toIntakeAssessmentGoalResponse(g)
	}
	return intakeMaturityAssessmentResponse{
		ID:            asm.ID,
		IntakeFormID:  asm.IntakeFormID,
		TopicID:       asm.TopicID,
		TopicName:     asm.TopicName,
		CurrentLevel:  asm.CurrentLevel,
		ProposedGoals: goals,
		Notes:         asm.Notes,
		CreatedAt:     asm.CreatedAt,
	}
}

func toGenerateIntakeGoalsResponse(res *domain.GenerateIntakeGoalsResult) generateIntakeGoalsResponse {
	goals := make([]intakeAssessmentGoalResponse, len(res.Goals))
	for i, g := range res.Goals {
		goals[i] = toIntakeAssessmentGoalResponse(g)
	}
	return generateIntakeGoalsResponse{Goals: goals}
}

func toPromoteIntakeToClientResponse(res *domain.PromoteIntakeToClientResult) promoteIntakeToClientResponse {
	return promoteIntakeToClientResponse{
		ClientID:                   res.ClientID,
		IntakeFormID:               res.IntakeFormID,
		RegistrationFormID:         res.RegistrationFormID,
		Message:                    res.Message,
		MaturityAssessmentsCreated: res.MaturityAssessmentsCreated,
		EmergencyContactsCreated:   res.EmergencyContactsCreated,
	}
}

func toIntakeFormGoalsResponse(res *domain.IntakeFormGoalsResult) intakeFormGoalsResponse {
	assessments := make([]intakeMaturityAssessmentResponse, len(res.Assessments))
	for i, a := range res.Assessments {
		assessments[i] = toIntakeMaturityAssessmentResponse(a)
	}
	return intakeFormGoalsResponse{Assessments: assessments}
}

func toIntakeFormTotalsResponse(res *domain.IntakeFormTotals) intakeFormTotalsResponse {
	return intakeFormTotalsResponse{
		FurtherInvestigationTotal: res.FurtherInvestigationTotal,
		WithoutGoalsTotal:         res.WithoutGoalsTotal,
	}
}

// ==================== Request Mappers ====================

func toCreateIntakeFormParams(req createIntakeFormRequest) domain.CreateIntakeFormParams {
	return domain.CreateIntakeFormParams{
		RegistrationFormID:       req.RegistrationFormID,
		DateOfIntake:             req.DateOfIntake,
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
}

func toListIntakeFormsParams(req listIntakeFormsRequest) domain.ListIntakeFormsParams {
	return domain.ListIntakeFormsParams{
		Search:    req.Search,
		Status:    req.Status,
		SortOrder: req.SortOrder,
		Limit:     req.PageSize,
		Offset:    (req.Page - 1) * req.PageSize,
	}
}

func toUpdateIntakeFormParams(id uuid.UUID, req updateIntakeFormRequest) domain.UpdateIntakeFormParams {
	return domain.UpdateIntakeFormParams{
		ID:                       id,
		DateOfIntake:             req.DateOfIntake,
		CareType:                 req.CareType,
		IntakeParticipants:       req.IntakeParticipants,
		FamilySituation:          req.FamilySituation,
		PsychologicalState:       req.PsychologicalState,
		SelfSufficiency:          req.SelfSufficiency,
		SenderID:                 req.SenderID,
		AssignedLocationID:       req.AssignedLocationID,
		RiskAssessment:           req.RiskAssessment,
		EvaluationIntervalsWeeks: req.EvaluationIntervalsWeeks,
		Signature:                req.Signature,
		ClearFields:              req.ClearFields,
	}
}

func toReplaceIntakeFormGoalsParams(req createIntakeFormGoalsRequest) domain.ReplaceIntakeFormGoalsParams {
	assessments := make([]domain.IntakeFormGoalItem, len(req.Assessments))
	for i, a := range req.Assessments {
		goals := make([]domain.IntakeAssessmentGoal, len(a.ProposedGoals))
		for j, g := range a.ProposedGoals {
			goals[j] = domain.IntakeAssessmentGoal{
				ID:          g.ID,
				Title:       g.Title,
				Description: g.Description,
				Priority:    g.Priority,
			}
		}
		assessments[i] = domain.IntakeFormGoalItem{
			TopicID:       a.TopicID,
			CurrentLevel:  a.CurrentLevel,
			ProposedGoals: goals,
			Notes:         a.Notes,
		}
	}
	return domain.ReplaceIntakeFormGoalsParams{Assessments: assessments}
}

func toGenerateIntakeGoalsParams(intakeFormID uuid.UUID, req generateIntakeGoalsRequest) domain.GenerateIntakeGoalsParams {
	return domain.GenerateIntakeGoalsParams{
		IntakeAssessmentID: uuid.Nil,
		IntakeFormID:       intakeFormID,
		RegistrationFormID: uuid.Nil,
		TopicID:            req.TopicID,
		CurrentLevel:       req.CurrentLevel,
		UserDesc:           req.UserDesc,
	}
}

func toUpdateIntakeConclusionParams(req updateIntakeConclusionRequest) domain.UpdateIntakeConclusionParams {
	return domain.UpdateIntakeConclusionParams{
		Decision:              req.Decision,
		IntakeConclusionNotes: req.IntakeConclusionNotes,
	}
}
