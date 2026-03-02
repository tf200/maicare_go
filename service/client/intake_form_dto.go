package clientp

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateIntakeFormRequest represents the request payload for creating a new intake form.
type CreateIntakeFormRequest struct {
	RegistrationFormID       uuid.UUID                   `json:"registration_form_id"`
	DateOfIntake             time.Time                   `json:"date_of_intake"`
	CareType                 db.IntakeCareTypeEnum       `json:"care_type"`
	IntakeParticipants       []db.IntakeParticipantsEnum `json:"intake_participants"`
	FamilySituation          *string                     `json:"family_situation"`
	PsychologicalState       *string                     `json:"psychological_state"`
	SelfSufficiency          int32                       `json:"self_sufficiency"`
	SenderID                 *uuid.UUID                  `json:"sender_id"`            // Formal sender match from referrer dropdown
	AssignedLocationID       *uuid.UUID                  `json:"assigned_location_id"` // Where client will live
	RiskAssessment           *string                     `json:"risk_assessment"`
	IntakeConclusion         db.IntakeConclusionEnum     `json:"intake_conclusion"`
	IntakeConclusionNotes    *string                     `json:"intake_conclusion_notes"`
	EvaluationIntervalsWeeks int32                       `json:"evaluation_intervals_weeks"`
	Signature                *string                     `json:"signature"`
}

// CreateIntakeFormResponse represents the response payload after creating a new intake form.
type CreateIntakeFormResponse struct {
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
	UpdatedAt                pgtype.Timestamptz          `json:"updated_at"`
}

type ListIntakeFormsRequest struct {
	pagination.Request
	Search    *string                  `param:"search" binding:"omitempty"`
	Status    *db.IntakeConclusionEnum `param:"status" binding:"omitempty"`
	SortOrder *string                  `param:"sort_order" binding:"omitempty,oneof=asc desc"`
}

type AssignedLocationAddress struct {
	Street              *string `json:"street"`
	HouseNumber         *string `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          *string `json:"postal_code"`
	City                *string `json:"city"`
}

type ListIntakeFormsResponse struct {
	ID                      uuid.UUID                `json:"id"`
	RegistrationFormID      uuid.UUID                `json:"registration_form_id"`
	DateOfIntake            time.Time                `json:"date_of_intake"`
	ClientFirstName         string                   `json:"client_first_name"`
	ClientLastName          string                   `json:"client_last_name"`
	ClientBsnNumber         string                   `json:"client_bsn_number"`
	IntakeStatus            db.IntakeConclusionEnum  `json:"intake_status"`
	GoalAssessmentDone      bool                     `json:"goal_assessment_done"`
	CareType                db.IntakeCareTypeEnum    `json:"care_type"`
	AssignedLocationID      *uuid.UUID               `json:"assigned_location_id"`
	AssignedLocationAddress *AssignedLocationAddress `json:"assigned_location_address"`
}

type GetIntakeFormTotalsResponse struct {
	FurtherInvestigationTotal int64 `json:"further_investigation_total"`
	WithoutGoalsTotal         int64 `json:"without_goals_total"`
}

type IntakeFormLocationDetails struct {
	Name                *string `json:"name"`
	Street              *string `json:"street"`
	HouseNumber         *string `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          *string `json:"postal_code"`
	City                *string `json:"city"`
}

type IntakeGoalTopic struct {
	AssessmentID  uuid.UUID              `json:"assessment_id"`
	TopicID       uuid.UUID              `json:"topic_id"`
	TopicName     string                 `json:"topic_name"`
	CurrentLevel  int32                  `json:"current_level"`
	ProposedGoals []IntakeAssessmentGoal `json:"proposed_goals"`
	Notes         *string                `json:"notes"`
}

type GetIntakeFormResponse struct {
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
	CreatedAt                time.Time                   `json:"created_at"`
	UpdatedAt                time.Time                   `json:"updated_at"`
	ClientFirstName          string                      `json:"client_first_name"`
	ClientLastName           string                      `json:"client_last_name"`
	ClientBsnNumber          string                      `json:"client_bsn_number"`
	DesiredGoals             []string                    `json:"desired_goals"`
	SenderName               *string                     `json:"sender_name"`
	Location                 *IntakeFormLocationDetails  `json:"location"`
	IntakeGoalsAssigned      []IntakeGoalTopic           `json:"intake_goals_assigned"`
	HasClient                bool                        `json:"has_client"`
}

type UpdateIntakeConclusionRequest struct {
	Decision              string  `json:"decision" binding:"required,oneof=accept refuse"`
	IntakeConclusionNotes *string `json:"intake_conclusion_notes"`
}

type UpdateIntakeConclusionResponse struct {
	ID                    uuid.UUID               `json:"id"`
	IntakeConclusion      db.IntakeConclusionEnum `json:"intake_conclusion"`
	IntakeConclusionNotes *string                 `json:"intake_conclusion_notes"`
	UpdatedAt             time.Time               `json:"updated_at"`
}

type UpdateIntakeFormRequest struct {
	DateOfIntake             *time.Time                   `json:"date_of_intake"`
	CareType                 *db.IntakeCareTypeEnum       `json:"care_type"`
	IntakeParticipants       *[]db.IntakeParticipantsEnum `json:"intake_participants"`
	FamilySituation          *string                      `json:"family_situation"`
	PsychologicalState       *string                      `json:"psychological_state"`
	SelfSufficiency          *int32                       `json:"self_sufficiency"`
	SenderID                 *uuid.UUID                   `json:"sender_id"`
	AssignedLocationID       *uuid.UUID                   `json:"assigned_location_id"`
	RiskAssessment           *string                      `json:"risk_assessment"`
	IntakeConclusion         db.IntakeConclusionEnum      `json:"intake_conclusion"`
	IntakeConclusionNotes    *string                      `json:"intake_conclusion_notes"`
	EvaluationIntervalsWeeks *int32                       `json:"evaluation_intervals_weeks"`
	Signature                *string                      `json:"signature"`
	ClearFields              []string                     `json:"clear_fields"`
}

type UpdateIntakeFormResponse struct {
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
