package clientp

import (
	"maicare_go/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// IntakeAssessmentGoal represents a single goal within a maturity topic assessment
type IntakeAssessmentGoal struct {
	ID          string `json:"id,omitempty"`          // Optional unique identifier
	Title       string `json:"title"`                 // Goal title/description
	Description string `json:"description,omitempty"` // Detailed description
	Priority    string `json:"priority"`              // "high", "medium", "low"
}

// GenerateIntakeGoalsRequest represents a request to generate goals for an intake assessment
type GenerateIntakeGoalsRequest struct {
	IntakeAssessmentID uuid.UUID `json:"assessment_id"`
	RegistrationFormID uuid.UUID `json:"registration_form_id"`
	TopicID            uuid.UUID `json:"topic_id"`
	CurrentLevel       int       `json:"current_level"`
	UserDesc           string    `json:"user_desc"`
}

// GenerateIntakeGoalsResponse represents the response with generated goals
type GenerateIntakeGoalsResponse struct {
	Goals []IntakeAssessmentGoal `json:"goals"`
}

type InitializeIntakeAssessmentsRequest struct {
	IntakeFormID uuid.UUID `json:"intake_form_id" binding:"required"`
}

type InitializeIntakeAssessmentsResponse struct {
	Count int `json:"count"`
}

type CreateIntakeFormGoalItem struct {
	TopicID       uuid.UUID              `json:"topic_id"`
	CurrentLevel  int32                  `json:"current_level"`
	ProposedGoals []IntakeAssessmentGoal `json:"proposed_goals"`
	Notes         *string                `json:"notes,omitempty"`
}

type CreateIntakeFormGoalsRequest struct {
	Assessments []CreateIntakeFormGoalItem `json:"assessments"`
}

type CreateIntakeFormGoalsResponse struct {
	Assessments []ListIntakeMaturityAssessmentsResponse `json:"assessments"`
}

// CreateIntakeMaturityAssessmentRequest represents a request to create a maturity assessment for an intake
type CreateIntakeMaturityAssessmentRequest struct {
	IntakeFormID     uuid.UUID              `json:"intake_form_id"`     // The intake this assessment belongs to
	MaturityMatrixID uuid.UUID              `json:"maturity_matrix_id"` // The maturity topic (Finance, Housing, etc.)
	CurrentLevel     int32                  `json:"current_level"`      // Current maturity level (1-5)
	ProposedGoals    []IntakeAssessmentGoal `json:"proposed_goals"`     // Array of goals for this topic
	Notes            *string                `json:"notes,omitempty"`    // Assessment notes
}

// CreateIntakeMaturityAssessmentResponse represents the response after creating a maturity assessment
type CreateIntakeMaturityAssessmentResponse struct {
	ID               uuid.UUID              `json:"id"`
	IntakeFormID     uuid.UUID              `json:"intake_form_id"`
	MaturityMatrixID uuid.UUID              `json:"maturity_matrix_id"`
	TopicName        string                 `json:"topic_name"`
	CurrentLevel     int32                  `json:"current_level"`
	ProposedGoals    []IntakeAssessmentGoal `json:"proposed_goals"`
	Notes            *string                `json:"notes"`
	CreatedAt        pgtype.Timestamptz     `json:"created_at"`
}

// ListIntakeMaturityAssessmentsRequest represents a request to list maturity assessments for an intake
type ListIntakeMaturityAssessmentsRequest struct {
	pagination.Request
	IntakeFormID uuid.UUID `json:"intake_form_id" param:"intake_id" binding:"required,uuid"`
}

// ListIntakeMaturityAssessmentsResponse represents a maturity assessment in the list response
type ListIntakeMaturityAssessmentsResponse struct {
	ID            uuid.UUID              `json:"id"`
	IntakeFormID  uuid.UUID              `json:"intake_form_id"`
	TopicID       uuid.UUID              `json:"topic_id"`
	TopicName     string                 `json:"topic_name"`
	CurrentLevel  int32                  `json:"current_level"`
	ProposedGoals []IntakeAssessmentGoal `json:"proposed_goals"`
	Notes         *string                `json:"notes"`
	CreatedAt     pgtype.Timestamptz     `json:"created_at"`
}

// GetIntakeMaturityAssessmentRequest represents a request to get a specific maturity assessment
type GetIntakeMaturityAssessmentRequest struct {
	AssessmentID uuid.UUID `json:"assessment_id" param:"assessment_id" binding:"required,uuid"`
}

// GetIntakeMaturityAssessmentResponse represents the response with a single maturity assessment
type GetIntakeMaturityAssessmentResponse struct {
	ID               uuid.UUID              `json:"id"`
	IntakeFormID     uuid.UUID              `json:"intake_form_id"`
	MaturityMatrixID uuid.UUID              `json:"maturity_matrix_id"`
	TopicName        string                 `json:"topic_name"`
	CurrentLevel     int32                  `json:"current_level"`
	ProposedGoals    []IntakeAssessmentGoal `json:"proposed_goals"`
	Notes            *string                `json:"notes"`
	CreatedAt        pgtype.Timestamptz     `json:"created_at"`
}

// UpdateIntakeMaturityAssessmentRequest represents a request to update a maturity assessment
type UpdateIntakeMaturityAssessmentRequest struct {
	AssessmentID  uuid.UUID              `json:"assessment_id" param:"assessment_id" binding:"required,uuid"`
	CurrentLevel  *int32                 `json:"current_level,omitempty"`  // Optional, only update if provided
	ProposedGoals []IntakeAssessmentGoal `json:"proposed_goals,omitempty"` // Optional, only update if provided
	Notes         *string                `json:"notes,omitempty"`          // Optional, only update if provided
}

// UpdateIntakeMaturityAssessmentResponse represents the response after updating a maturity assessment
type UpdateIntakeMaturityAssessmentResponse struct {
	ID               uuid.UUID              `json:"id"`
	IntakeFormID     uuid.UUID              `json:"intake_form_id"`
	MaturityMatrixID uuid.UUID              `json:"maturity_matrix_id"`
	CurrentLevel     int32                  `json:"current_level"`
	ProposedGoals    []IntakeAssessmentGoal `json:"proposed_goals"`
	Notes            *string                `json:"notes"`
	CreatedAt        pgtype.Timestamptz     `json:"created_at"`
}

// DeleteIntakeMaturityAssessmentRequest represents a request to delete a maturity assessment
type DeleteIntakeMaturityAssessmentRequest struct {
	AssessmentID uuid.UUID `json:"assessment_id" param:"assessment_id" binding:"required,uuid"`
}
