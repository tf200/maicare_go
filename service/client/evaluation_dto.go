package clientp

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

type CreateGoalEvaluationRequest struct {
	OverallNotes *string                       `json:"overall_notes"`
	Submit       bool                          `json:"submit"`
	Items        []SaveGoalEvaluationDraftItem `json:"items"`
}

type SaveGoalEvaluationDraftItem struct {
	GoalID   uuid.UUID `json:"goal_id"`
	Progress string    `json:"progress"`
	Notes    *string   `json:"notes"`
}

type GoalEvaluationResponse struct {
	ID                      uuid.UUID                    `json:"id"`
	ClientID                uuid.UUID                    `json:"client_id"`
	EvaluationDate          time.Time                    `json:"evaluation_date"`
	PeriodStart             *time.Time                   `json:"period_start"`
	PeriodEnd               *time.Time                   `json:"period_end"`
	EvaluationIntervalWeeks int32                        `json:"evaluation_interval_weeks"`
	Status                  string                       `json:"status"`
	OverallNotes            *string                      `json:"overall_notes"`
	CreatedByEmployeeID     *uuid.UUID                   `json:"created_by_employee_id"`
	CreatedAt               time.Time                    `json:"created_at"`
	UpdatedAt               time.Time                    `json:"updated_at"`
	SubmitError             *string                      `json:"submit_error,omitempty"`
	Items                   []GoalEvaluationItemResponse `json:"items"`
}

type GoalEvaluationItemResponse struct {
	ID                uuid.UUID `json:"id"`
	EvaluationID      uuid.UUID `json:"evaluation_id"`
	GoalID            uuid.UUID `json:"goal_id"`
	GoalTitle         string    `json:"goal_title"`
	GoalDescription   *string   `json:"goal_description"`
	TopicNameSnapshot *string   `json:"topic_name_snapshot"`
	Progress          string    `json:"progress"`
	Notes             *string   `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ListUpcomingEvaluationsRequest struct {
	pagination.Request
}

type ListUpcomingEvaluationsResponse struct {
	ClientID         uuid.UUID `json:"client_id"`
	ClientFirstName  string    `json:"client_first_name"`
	ClientLastName   string    `json:"client_last_name"`
	DueDate          time.Time `json:"due_date"`
	DaysLeft         int32     `json:"days_left"`
	Priority         string    `json:"priority"`
	HasDraft         bool      `json:"has_draft"`
	FilledGoalsCount int32     `json:"filled_goals_count"`
	TotalGoalsCount  int32     `json:"total_goals_count"`
}

type ListRecentSubmittedEvaluationsRequest struct {
	pagination.Request
}

type ListRecentSubmittedEvaluationsResponse struct {
	EvaluationID       uuid.UUID  `json:"evaluation_id"`
	ClientID           uuid.UUID  `json:"client_id"`
	ClientFirstName    string     `json:"client_first_name"`
	ClientLastName     string     `json:"client_last_name"`
	EvaluationDate     time.Time  `json:"evaluation_date"`
	SubmittedAt        time.Time  `json:"submitted_at"`
	NextEvaluationDate *time.Time `json:"next_evaluation_date"`
	FilledGoalsCount   int32      `json:"filled_goals_count"`
	TotalGoalsCount    int32      `json:"total_goals_count"`
}

type ListRecentDraftEvaluationsRequest struct {
	pagination.Request
}

type ListRecentDraftEvaluationsResponse struct {
	EvaluationID     uuid.UUID `json:"evaluation_id"`
	ClientID         uuid.UUID `json:"client_id"`
	ClientFirstName  string    `json:"client_first_name"`
	ClientLastName   string    `json:"client_last_name"`
	DueDate          time.Time `json:"due_date"`
	UpdatedAt        time.Time `json:"updated_at"`
	DaysLeft         int32     `json:"days_left"`
	Priority         string    `json:"priority"`
	FilledGoalsCount int32     `json:"filled_goals_count"`
	TotalGoalsCount  int32     `json:"total_goals_count"`
}

type GoalEvaluationBootstrapResponse struct {
	ClientID                uuid.UUID                           `json:"client_id"`
	ClientFirstName         string                              `json:"client_first_name"`
	ClientLastName          string                              `json:"client_last_name"`
	NextEvaluationDate      *time.Time                          `json:"next_evaluation_date"`
	DaysLeft                *int32                              `json:"days_left"`
	Priority                *string                             `json:"priority"`
	ExistingDraft           *GoalEvaluationBootstrapDraft       `json:"existing_draft"`
	LastCompletedEvaluation *GoalEvaluationBootstrapCompleted   `json:"last_completed_evaluation"`
	ActiveGoals             []GoalEvaluationBootstrapActiveGoal `json:"active_goals"`
}

type GoalEvaluationBootstrapDraft struct {
	ID             uuid.UUID `json:"id"`
	EvaluationDate time.Time `json:"evaluation_date"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type GoalEvaluationBootstrapCompleted struct {
	ID                  uuid.UUID  `json:"id"`
	EvaluationDate      time.Time  `json:"evaluation_date"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	OverallNotes        *string    `json:"overall_notes"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id"`
	CreatorName         *string    `json:"creator_name"`
}

type GoalEvaluationBootstrapActiveGoal struct {
	GoalID            uuid.UUID `json:"goal_id"`
	Title             string    `json:"title"`
	TopicNameSnapshot *string   `json:"topic_name_snapshot"`
	Priority          string    `json:"priority"`
	SortOrder         int32     `json:"sort_order"`
	LastProgress      *string   `json:"last_progress"`
	LastNotes         *string   `json:"last_notes"`
}
