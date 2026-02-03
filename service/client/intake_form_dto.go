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
	CreatedAt                pgtype.Timestamptz          `json:"created_at"`
	UpdatedAt                pgtype.Timestamptz          `json:"updated_at"`
}

type ListIntakeFormsRequest struct {
	pagination.Request
	Search    *string `param:"search" binding:"omitempty"`
	SortOrder *string `param:"sort_order" binding:"omitempty,oneof=asc desc"`
}

type ListIntakeFormsResponse struct {
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
	CreatedAt                pgtype.Timestamptz          `json:"created_at"`
	UpdatedAt                pgtype.Timestamptz          `json:"updated_at"`
	ClientFirstName          string                      `json:"client_first_name"`
	ClientLastName           string                      `json:"client_last_name"`
	ClientBsnNumber          string                      `json:"client_bsn_number"`
}
