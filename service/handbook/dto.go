package handbook

import (
	"time"

	"github.com/goccy/go-json"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"

	"github.com/google/uuid"
)

type CreateTemplateForDepartmentRequest struct {
	DepartmentID uuid.UUID `json:"department_id" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	Description  *string   `json:"description"`
}

type HandbookTemplateAPI struct {
	ID           uuid.UUID  `json:"id"`
	DepartmentID uuid.UUID  `json:"department_id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Version      int32      `json:"version"`
	Status       string     `json:"status"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`
	ArchivedAt   *time.Time `json:"archived_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CloneTemplateToDraftRequest struct {
	SourceTemplateID uuid.UUID `json:"source_template_id" binding:"required"`
}

type UpdateTemplateRequest struct {
	TemplateID     uuid.UUID `json:"template_id" binding:"required"`
	Title          *string   `json:"title"`
	SetTitle       bool      `json:"set_title"`
	Description    *string   `json:"description"`
	SetDescription bool      `json:"set_description"`
}

type PublishTemplateRequest struct {
	TemplateID uuid.UUID `json:"template_id" binding:"required"`
}

type CreateStepRequest struct {
	TemplateID uuid.UUID       `json:"template_id" binding:"required"`
	SortOrder  int32           `json:"sort_order" binding:"required"`
	Kind       string          `json:"kind" binding:"required,oneof=content ack link quiz rich_text"`
	Title      string          `json:"title" binding:"required"`
	Body       *string         `json:"body"`
	Content    json.RawMessage `json:"content"`
	IsRequired *bool           `json:"is_required"`
}

type LinkStepContent struct {
	URL string `json:"url"`
}

type QuizStepContent struct {
	Question           string   `json:"question"`
	Options            []string `json:"options"`
	CorrectOptionIndex int      `json:"correct_option_index"`
}

type CreateStepResponse struct {
	ID         uuid.UUID               `json:"id"`
	TemplateID uuid.UUID               `json:"template_id"`
	SortOrder  int32                   `json:"sort_order"`
	Kind       db.HandbookStepKindEnum `json:"kind"`
	Title      string                  `json:"title"`
	Body       *string                 `json:"body"`
	Content    any                     `json:"content"`
	IsRequired bool                    `json:"is_required"`
}

type ListStepResponse struct {
	ID         uuid.UUID               `json:"id"`
	SortOrder  int32                   `json:"sort_order"`
	Kind       db.HandbookStepKindEnum `json:"kind"`
	Title      string                  `json:"title"`
	Body       *string                 `json:"body"`
	Content    any                     `json:"content"`
	IsRequired bool                    `json:"is_required"`
}

type UpdateStepRequest struct {
	StepID          uuid.UUID       `json:"step_id" binding:"required"`
	Title           *string         `json:"title"`
	SetTitle        bool            `json:"set_title"`
	Body            *string         `json:"body"`
	SetBody         bool            `json:"set_body"`
	Content         json.RawMessage `json:"content"`
	ContentProvided bool            `json:"content_provided"`
	IsRequired      *bool           `json:"is_required"`
	SetIsRequired   bool            `json:"set_is_required"`
}

type UpdateStepResponse struct {
	ID         uuid.UUID               `json:"id"`
	TemplateID uuid.UUID               `json:"template_id"`
	SortOrder  int32                   `json:"sort_order"`
	Kind       db.HandbookStepKindEnum `json:"kind"`
	Title      string                  `json:"title"`
	Body       *string                 `json:"body"`
	Content    any                     `json:"content"`
	IsRequired bool                    `json:"is_required"`
	UpdatedAt  time.Time               `json:"updated_at"`
}

type DeleteStepRequest struct {
	StepID uuid.UUID `json:"step_id" binding:"required"`
}

type DeleteStepResponse struct {
	StepID  uuid.UUID `json:"step_id"`
	Deleted bool      `json:"deleted"`
}

type ReorderStepsRequest struct {
	TemplateID     uuid.UUID   `json:"template_id" binding:"required"`
	OrderedStepIDs []uuid.UUID `json:"ordered_step_ids" binding:"required,min=1"`
}

type ReorderStepsResponse struct {
	TemplateID uuid.UUID          `json:"template_id"`
	Steps      []ListStepResponse `json:"steps"`
}

type AssignTemplateToEmployeeRequest struct {
	EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
	TemplateID uuid.UUID `json:"template_id" binding:"required"`
}

type AssignTemplateToEmployeeResponse struct {
	EmployeeHandbookID uuid.UUID `json:"employee_handbook_id"`
	EmployeeID         uuid.UUID `json:"employee_id"`
	TemplateID         uuid.UUID `json:"template_id"`
	AssignedAt         time.Time `json:"assigned_at"`
	Status             string    `json:"status"`
}

type WaiveEmployeeHandbookRequest struct {
	EmployeeHandbookID uuid.UUID `json:"employee_handbook_id" binding:"required"`
	Reason             *string   `json:"reason"`
}

type WaiveEmployeeHandbookResponse struct {
	EmployeeHandbookID uuid.UUID  `json:"employee_handbook_id"`
	EmployeeID         uuid.UUID  `json:"employee_id"`
	Status             string     `json:"status"`
	CompletedAt        *time.Time `json:"completed_at"`
}

type HandbookAssignmentHistoryEntry struct {
	ID                 uuid.UUID  `json:"id"`
	EmployeeHandbookID *uuid.UUID `json:"employee_handbook_id"`
	EmployeeID         uuid.UUID  `json:"employee_id"`
	TemplateID         uuid.UUID  `json:"template_id"`
	TemplateVersion    int32      `json:"template_version"`
	Event              string     `json:"event"`
	ActorEmployeeID    *uuid.UUID `json:"actor_employee_id"`
	Metadata           any        `json:"metadata"`
	CreatedAt          time.Time  `json:"created_at"`
}

type ListEmployeeHandbookAssignmentsRequest struct {
	pagination.Request
	DepartmentID *uuid.UUID `form:"department_id"`
	Search       *string    `form:"search"`
	Status       *string    `form:"status"`
}

type EmployeeHandbookAssignmentSummary struct {
	EmployeeID             uuid.UUID  `json:"employee_id"`
	FirstName              string     `json:"first_name"`
	LastName               string     `json:"last_name"`
	DepartmentID           *uuid.UUID `json:"department_id"`
	DepartmentName         *string    `json:"department_name"`
	EmployeeHandbookID     *uuid.UUID `json:"employee_handbook_id"`
	TemplateID             *uuid.UUID `json:"template_id"`
	TemplateTitle          *string    `json:"template_title"`
	TemplateVersion        *int32     `json:"template_version"`
	HandbookStatus         string     `json:"handbook_status"`
	AssignedAt             *time.Time `json:"assigned_at"`
	StartedAt              *time.Time `json:"started_at"`
	CompletedAt            *time.Time `json:"completed_at"`
	DueAt                  *time.Time `json:"due_at"`
	RequiredStepsTotal     int32      `json:"required_steps_total"`
	RequiredStepsCompleted int32      `json:"required_steps_completed"`
}

type GetEmployeeHandbookDetailsResponse struct {
	EmployeeHandbookID uuid.UUID        `json:"employee_handbook_id"`
	EmployeeID         uuid.UUID        `json:"employee_id"`
	FirstName          string           `json:"first_name"`
	LastName           string           `json:"last_name"`
	Status             string           `json:"status"`
	AssignedAt         time.Time        `json:"assigned_at"`
	StartedAt          *time.Time       `json:"started_at"`
	CompletedAt        *time.Time       `json:"completed_at"`
	DueAt              *time.Time       `json:"due_at"`
	TemplateID         uuid.UUID        `json:"template_id"`
	TemplateTitle      string           `json:"template_title"`
	TemplateDesc       *string          `json:"template_description"`
	TemplateVersion    int32            `json:"template_version"`
	DepartmentID       uuid.UUID        `json:"department_id"`
	DepartmentName     string           `json:"department_name"`
	Steps              []MyHandbookStep `json:"steps"`
}

type GetMyActiveHandbookResponse struct {
	HandbookID  uuid.UUID  `json:"handbook_id"`
	Status      string     `json:"status"`
	AssignedAt  time.Time  `json:"assigned_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DueAt       *time.Time `json:"due_at"`

	TemplateID      uuid.UUID `json:"template_id"`
	TemplateTitle   string    `json:"template_title"`
	TemplateDesc    *string   `json:"template_description"`
	TemplateVersion int32     `json:"template_version"`

	DepartmentID   uuid.UUID `json:"department_id"`
	DepartmentName string    `json:"department_name"`

	Steps []MyHandbookStep `json:"steps"`
}

type MyHandbookStep struct {
	StepID     uuid.UUID               `json:"step_id"`
	SortOrder  int32                   `json:"sort_order"`
	Kind       db.HandbookStepKindEnum `json:"kind"`
	Title      string                  `json:"title"`
	Body       *string                 `json:"body"`
	Content    any                     `json:"content"`
	IsRequired bool                    `json:"is_required"`

	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Response    any        `json:"response"`
}

type StartMyHandbookResponse struct {
	HandbookID uuid.UUID  `json:"handbook_id"`
	Status     string     `json:"status"`
	StartedAt  *time.Time `json:"started_at"`
}

type CompleteMyHandbookStepResponse struct {
	HandbookID     uuid.UUID `json:"handbook_id"`
	StepID         uuid.UUID `json:"step_id"`
	StepStatus     string    `json:"step_status"`
	CompletedAt    time.Time `json:"completed_at"`
	HandbookStatus string    `json:"handbook_status"`
}

type ListTemplatesByDepartmentRequest struct {
	pagination.Request
	DepartmentID uuid.UUID `form:"department_id" binding:"required"`
}

type ListEligibleEmployeesRequest struct {
	pagination.Request
	DepartmentID *uuid.UUID `form:"department_id"`
	Search       *string    `form:"search"`
}

type ListEligibleEmployeesResponse struct {
	EmployeeID     uuid.UUID  `json:"employee_id"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	DepartmentID   *uuid.UUID `json:"department_id"`
	DepartmentName *string    `json:"department_name"`
}
