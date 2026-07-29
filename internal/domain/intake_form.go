package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	db "maicare_go/db/sqlc"
)

var (
	ErrIntakeFormNotFound                     = errors.New("intake form not found")
	ErrIntakeGoalsUpdateBlockedByActiveClient = errors.New("intake goals cannot be updated when an active client exists for this intake form")
	ErrIntakeFormUpdateBlockedByActiveClient  = errors.New("intake form cannot be updated when an active client exists for this intake form")
	ErrIntakeFormDeleteBlockedByActiveClient  = errors.New("intake form cannot be deleted when an active client exists for this intake form")
	ErrIntakeFormUpdateConflict               = errors.New("intake form was updated by another request")
	ErrNoIntakeFormFieldsToUpdate             = errors.New("no intake form fields provided for update")
	ErrInvalidIntakeFormClearField            = errors.New("invalid intake form clear field")
	ErrInvalidIntakeFormUpdate                = errors.New("invalid intake form update")
	ErrIntakeNotSuitable                      = errors.New("intake conclusion is not suitable")
)

type IntakeAssessmentGoal struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority"`
}

type IntakeForm struct {
	ID                       uuid.UUID
	RegistrationFormID       uuid.UUID
	DateOfIntake             time.Time
	CareType                 db.IntakeCareTypeEnum
	IntakeParticipants       []db.IntakeParticipantsEnum
	FamilySituation          *string
	PsychologicalState       *string
	SelfSufficiency          int32
	SenderID                 *uuid.UUID
	AssignedLocationID       *uuid.UUID
	RiskAssessment           *string
	IntakeConclusion         db.IntakeConclusionEnum
	IntakeConclusionNotes    *string
	EvaluationIntervalsWeeks int32
	Signature                *string
	UpdatedAt                time.Time
}

type IntakeFormListItem struct {
	ID                      uuid.UUID
	RegistrationFormID      uuid.UUID
	DateOfIntake            time.Time
	ClientFirstName         string
	ClientLastName          string
	ClientBsnNumber         string
	IntakeStatus            db.IntakeConclusionEnum
	GoalAssessmentDone      bool
	HasClient               bool
	CareType                db.IntakeCareTypeEnum
	AssignedLocationID      *uuid.UUID
	AssignedLocationAddress *AssignedLocationAddress
}

type AssignedLocationAddress struct {
	Street              *string
	HouseNumber         *string
	HouseNumberAddition *string
	PostalCode          *string
	City                *string
}

type IntakeFormTotals struct {
	FurtherInvestigationTotal int64
	WithoutGoalsTotal         int64
}

type IntakeFormLocationDetails struct {
	Name                *string
	Street              *string
	HouseNumber         *string
	HouseNumberAddition *string
	PostalCode          *string
	City                *string
}

type IntakeGoalTopic struct {
	AssessmentID  uuid.UUID
	TopicID       uuid.UUID
	TopicName     string
	CurrentLevel  int32
	ProposedGoals []IntakeAssessmentGoal
	Notes         *string
}

type IntakeFormDetail struct {
	ID                       uuid.UUID
	RegistrationFormID       uuid.UUID
	DateOfIntake             time.Time
	CareType                 db.IntakeCareTypeEnum
	IntakeParticipants       []db.IntakeParticipantsEnum
	FamilySituation          *string
	PsychologicalState       *string
	SelfSufficiency          int32
	SenderID                 *uuid.UUID
	AssignedLocationID       *uuid.UUID
	RiskAssessment           *string
	IntakeConclusion         db.IntakeConclusionEnum
	IntakeConclusionNotes    *string
	EvaluationIntervalsWeeks int32
	Signature                *string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	ClientFirstName          string
	ClientLastName           string
	ClientBsnNumber          string
	DesiredGoals             []string
	SenderName               *string
	Location                 *IntakeFormLocationDetails
	IntakeGoalsAssigned      []IntakeGoalTopic
	HasClient                bool
}

type IntakeFormConclusion struct {
	ID                    uuid.UUID
	IntakeConclusion      db.IntakeConclusionEnum
	IntakeConclusionNotes *string
	UpdatedAt             time.Time
}

type IntakeMaturityAssessment struct {
	ID            uuid.UUID
	IntakeFormID  uuid.UUID
	TopicID       uuid.UUID
	TopicName     string
	CurrentLevel  int32
	ProposedGoals []IntakeAssessmentGoal
	Notes         *string
	CreatedAt     time.Time
}

type CreateIntakeFormParams struct {
	RegistrationFormID       uuid.UUID
	DateOfIntake             time.Time
	CareType                 db.IntakeCareTypeEnum
	IntakeParticipants       []db.IntakeParticipantsEnum
	FamilySituation          *string
	PsychologicalState       *string
	SelfSufficiency          int32
	SenderID                 *uuid.UUID
	AssignedLocationID       *uuid.UUID
	RiskAssessment           *string
	IntakeConclusion         db.IntakeConclusionEnum
	IntakeConclusionNotes    *string
	EvaluationIntervalsWeeks int32
	Signature                *string
}

type ListIntakeFormsParams struct {
	Search    *string
	Status    *db.IntakeConclusionEnum
	SortOrder *string
	Limit     int32
	Offset    int32
}

type UpdateIntakeFormParams struct {
	ID                       uuid.UUID
	DateOfIntake             *time.Time
	CareType                 *db.IntakeCareTypeEnum
	IntakeParticipants       *[]db.IntakeParticipantsEnum
	FamilySituation          *string
	PsychologicalState       *string
	SelfSufficiency          *int32
	SenderID                 *uuid.UUID
	AssignedLocationID       *uuid.UUID
	RiskAssessment           *string
	EvaluationIntervalsWeeks *int32
	Signature                *string
	ClearFields              []string
}

type UpdateIntakeConclusionParams struct {
	Decision              string
	IntakeConclusionNotes *string
}

type ReplaceIntakeFormGoalsParams struct {
	Assessments []IntakeFormGoalItem
}

type IntakeFormGoalItem struct {
	TopicID       uuid.UUID              `json:"topic_id"`
	CurrentLevel  int32                  `json:"current_level"`
	ProposedGoals []IntakeAssessmentGoal `json:"proposed_goals"`
	Notes         *string                `json:"notes"`
}

type GenerateIntakeGoalsParams struct {
	IntakeAssessmentID uuid.UUID
	IntakeFormID       uuid.UUID
	RegistrationFormID uuid.UUID
	TopicID            uuid.UUID
	CurrentLevel       int
	UserDesc           string
}

type GenerateIntakeGoalsResult struct {
	Goals []IntakeAssessmentGoal
}

type PromoteIntakeToClientParams struct {
	IntakeFormID uuid.UUID
	EmployeeID   uuid.UUID
}

type PromoteIntakeToClientResult struct {
	ClientID                   uuid.UUID
	IntakeFormID               uuid.UUID
	RegistrationFormID         uuid.UUID
	Message                    string
	MaturityAssessmentsCreated int
	EmergencyContactsCreated   int
}

type IntakeFormGoalsResult struct {
	Assessments []IntakeMaturityAssessment
}

// Repository interface
type IntakeFormRepository interface {
	CreateIntakeForm(ctx context.Context, params CreateIntakeFormParams) (*IntakeForm, error)
	ListIntakeForms(ctx context.Context, params ListIntakeFormsParams) ([]IntakeFormListItem, int64, error)
	GetIntakeFormTotals(ctx context.Context) (*IntakeFormTotals, error)
	GetIntakeFormDetail(ctx context.Context, id uuid.UUID) (*IntakeFormDetail, error)
	UpdateIntakeForm(ctx context.Context, params UpdateIntakeFormParams) (*IntakeForm, error)
	UpdateIntakeConclusion(ctx context.Context, id uuid.UUID, params UpdateIntakeConclusionParams) (*IntakeFormConclusion, error)
	ReplaceIntakeFormGoals(ctx context.Context, intakeFormID uuid.UUID, params ReplaceIntakeFormGoalsParams) (*IntakeFormGoalsResult, error)
	PromoteIntakeToClient(ctx context.Context, params PromoteIntakeToClientParams) (*PromoteIntakeToClientResult, error)
	GetIntakeForm(ctx context.Context, id uuid.UUID) (db.IntakeForm, error)
	GetIntakeFormByRegistrationFormID(ctx context.Context, registrationFormID uuid.UUID) (db.IntakeForm, error)
	GetRegistrationForm(ctx context.Context, id uuid.UUID) (db.GetRegistrationFormRow, error)
	GetIntakeMaturityAssessment(ctx context.Context, id uuid.UUID) (db.GetIntakeMaturityAssessmentRow, error)
	GetTopicLevel(ctx context.Context, params db.GetTopicLevelParams) (db.GetTopicLevelRow, error)
	HasActiveClientByIntakeFormID(ctx context.Context, intakeFormID *uuid.UUID) (bool, error)
	LockIntakeFormByID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetClientByIntakeFormID(ctx context.Context, intakeFormID *uuid.UUID) (db.ClientDetail, error)
	CreateClientDetails(ctx context.Context, params db.CreateClientDetailsParams) (db.ClientDetail, error)
	CreateClientGoalsFromIntakeAssessments(ctx context.Context, params db.CreateClientGoalsFromIntakeAssessmentsParams) ([]db.CreateClientGoalsFromIntakeAssessmentsRow, error)
	CreateEmergencyContact(ctx context.Context, params db.CreateEmemrgencyContactParams) (db.ClientEmergencyContact, error)
	CreateIntakeTopicAssessmentsBatch(ctx context.Context, params db.CreateIntakeTopicAssessmentsBatchParams) ([]db.CreateIntakeTopicAssessmentsBatchRow, error)
	DeleteIntakeTopicAssessmentsByIntakeForm(ctx context.Context, intakeFormID uuid.UUID) error
	DeleteIntakeForm(ctx context.Context, id uuid.UUID) error
	ExecTx(ctx context.Context, fn func(*db.Queries) error) error
}

// Service interface
type IntakeFormService interface {
	CreateIntakeForm(ctx context.Context, params CreateIntakeFormParams) (*IntakeForm, error)
	ListIntakeForms(ctx context.Context, params ListIntakeFormsParams) (*ListResult[IntakeFormListItem], error)
	GetIntakeFormTotals(ctx context.Context) (*IntakeFormTotals, error)
	GetIntakeForm(ctx context.Context, id uuid.UUID) (*IntakeFormDetail, error)
	UpdateIntakeForm(ctx context.Context, params UpdateIntakeFormParams) (*IntakeForm, error)
	ReplaceIntakeFormGoals(ctx context.Context, intakeFormID uuid.UUID, params ReplaceIntakeFormGoalsParams) (*IntakeFormGoalsResult, error)
	GenerateIntakeGoals(ctx context.Context, params GenerateIntakeGoalsParams) (*GenerateIntakeGoalsResult, error)
	UpdateIntakeConclusion(ctx context.Context, id uuid.UUID, params UpdateIntakeConclusionParams) (*IntakeFormConclusion, error)
	PromoteIntakeToClient(ctx context.Context, params PromoteIntakeToClientParams) (*PromoteIntakeToClientResult, error)
	DeleteIntakeForm(ctx context.Context, id uuid.UUID) error
}
