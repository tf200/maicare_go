package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrIncidentNotFound = errors.New("incident not found")

// Incident represents a full incident domain model.
type Incident struct {
	ID                      uuid.UUID
	EmployeeID              uuid.UUID
	EmployeeFirstName       string
	EmployeeLastName        string
	LocationID              uuid.UUID
	LocationName            string
	ClientID                uuid.UUID
	ClientFirstName         string
	ClientLastName          string
	ReporterInvolvement     string
	InformedParties         []string
	OccurredAt              time.Time
	IncidentType            string
	SeverityOfIncident      string
	IncidentExplanation     *string
	RecurrenceRisk          string
	IncidentPreventSteps    *string
	IncidentTakenMeasures   *string
	CauseCategories         []string
	CauseExplanation        *string
	PhysicalInjury          string
	PhysicalInjuryDesc      *string
	PsychologicalDamage     string
	PsychologicalDamageDesc *string
	NeededConsultation      string
	FollowUpActions         []string
	FollowUpNotes           *string
	IsEmployeeAbsent        bool
	AdditionalDetails       *string
	Emails                  []string
	UpdatedAt               time.Time
	CreatedAt               time.Time
	IsConfirmed             bool
	FileUrl                 *string
	ConfirmedAt             *time.Time
	ConfirmedBy             *uuid.UUID
	ConfirmationEmailSentAt *time.Time
}

// IncidentListItem represents an incident in a client-specific list.
type IncidentListItem struct {
	ID                     uuid.UUID
	OccurredAt             time.Time
	IncidentType           string
	SeverityOfIncident     string
	IsConfirmed            bool
	EmployeeFirstName      string
	EmployeeLastName       string
	EmployeeProfilePicture *string
	LocationName           string
}

// IncidentSummary represents an incident in the global list.
type IncidentSummary struct {
	ID                 uuid.UUID
	ClientID           uuid.UUID
	OccurredAt         time.Time
	IncidentType       string
	SeverityOfIncident string
	IsConfirmed        bool
	ClientFirstName    string
	ClientLastName     string
	ClientBSN          *string
	EmployeeFirstName  string
	EmployeeLastName   string
	LocationName       string
}

// IncidentCounts holds aggregate incident statistics.
type IncidentCounts struct {
	SeriousFatalCount        int64
	PendingConfirmationCount int64
	Past24hCount             int64
}

// CreateIncidentParams holds parameters for creating an incident.
type CreateIncidentParams struct {
	ClientID                uuid.UUID
	EmployeeID              uuid.UUID
	LocationID              uuid.UUID
	ReporterInvolvement     string
	InformedParties         []string
	OccurredAt              time.Time
	IncidentType            string
	SeverityOfIncident      string
	IncidentExplanation     *string
	RecurrenceRisk          string
	IncidentPreventSteps    *string
	IncidentTakenMeasures   *string
	CauseCategories         []string
	CauseExplanation        *string
	PhysicalInjury          string
	PhysicalInjuryDesc      *string
	PsychologicalDamage     string
	PsychologicalDamageDesc *string
	NeededConsultation      string
	FollowUpActions         []string
	FollowUpNotes           *string
	IsEmployeeAbsent        bool
	AdditionalDetails       *string
	Emails                  []string
}

// UpdateIncidentParams holds parameters for updating an incident.
type UpdateIncidentParams struct {
	ID                      uuid.UUID
	EmployeeID              *uuid.UUID
	LocationID              *uuid.UUID
	ReporterInvolvement     *string
	InformedParties         []string
	OccurredAt              time.Time
	IncidentType            *string
	SeverityOfIncident      *string
	IncidentExplanation     *string
	RecurrenceRisk          *string
	IncidentPreventSteps    *string
	IncidentTakenMeasures   *string
	CauseCategories         []string
	CauseExplanation        *string
	PhysicalInjury          *string
	PhysicalInjuryDesc      *string
	PsychologicalDamage     *string
	PsychologicalDamageDesc *string
	NeededConsultation      *string
	FollowUpActions         []string
	FollowUpNotes           *string
	IsEmployeeAbsent        *bool
	AdditionalDetails       *string
	Emails                  []string
}

// ListIncidentsParams holds parameters for listing incidents by client.
type ListIncidentsParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

// ListIncidentsResult holds the result of listing incidents by client.
type ListIncidentsResult struct {
	Items      []IncidentListItem
	TotalCount int64
}

// ListAllIncidentsParams holds parameters for listing all incidents.
type ListAllIncidentsParams struct {
	Limit       int32
	Offset      int32
	IsConfirmed *bool
	Search      *string
}

// ListAllIncidentsResult holds the result of listing all incidents.
type ListAllIncidentsResult struct {
	Items      []IncidentSummary
	TotalCount int64
}

// ConfirmIncidentResult holds the result of confirming an incident.
type ConfirmIncidentResult struct {
	ID      uuid.UUID
	FileUrl *string
}

// IncidentPDFData holds data needed to generate an incident PDF.
type IncidentPDFData struct {
	ID                      uuid.UUID
	EmployeeID              uuid.UUID
	EmployeeFirstName       string
	EmployeeLastName        string
	LocationID              uuid.UUID
	ReporterInvolvement     string
	InformedParties         []string
	OccurredAt              time.Time
	IncidentType            string
	SeverityOfIncident      string
	IncidentExplanation     *string
	RecurrenceRisk          string
	IncidentPreventSteps    *string
	IncidentTakenMeasures   *string
	CauseCategories         []string
	CauseExplanation        *string
	PhysicalInjury          string
	PhysicalInjuryDesc      *string
	PsychologicalDamage     string
	PsychologicalDamageDesc *string
	NeededConsultation      string
	FollowUpActions         []string
	FollowUpNotes           *string
	IsEmployeeAbsent        bool
	AdditionalDetails       *string
	ClientID                uuid.UUID
	ClientFirstName         string
	ClientLastName          string
	LocationName            string
}

// IncidentPDFGenerator generates incident PDFs.
type IncidentPDFGenerator interface {
	GenerateIncidentPDF(ctx context.Context, data IncidentPDFData) ([]byte, error)
}

// IncidentRepository defines incident persistence operations.
type IncidentRepository interface {
	CreateIncident(ctx context.Context, params CreateIncidentParams) (*Incident, error)
	ListIncidents(ctx context.Context, params ListIncidentsParams) (*ListIncidentsResult, error)
	GetIncident(ctx context.Context, id uuid.UUID) (*Incident, error)
	UpdateIncident(ctx context.Context, params UpdateIncidentParams) (*Incident, error)
	DeleteIncident(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	ConfirmIncident(ctx context.Context, id uuid.UUID) (int64, error)
	ListAllIncidents(ctx context.Context, params ListAllIncidentsParams) (*ListAllIncidentsResult, error)
	GetIncidentCounts(ctx context.Context) (*IncidentCounts, error)
	GetAllAdminUsers(ctx context.Context) ([]uuid.UUID, error)
}

// IncidentService defines incident business operations.
type IncidentService interface {
	CreateIncident(ctx context.Context, params CreateIncidentParams) (*Incident, error)
	ListIncidents(ctx context.Context, params ListIncidentsParams) (*ListIncidentsResult, error)
	GetIncident(ctx context.Context, id uuid.UUID) (*Incident, error)
	UpdateIncident(ctx context.Context, params UpdateIncidentParams) (*Incident, error)
	DeleteIncident(ctx context.Context, id uuid.UUID) error
	GenerateIncidentFile(ctx context.Context, id uuid.UUID) ([]byte, string, error)
	ConfirmIncident(ctx context.Context, id uuid.UUID, confirmedByUserID uuid.UUID) (*ConfirmIncidentResult, error)
	ListAllIncidents(ctx context.Context, params ListAllIncidentsParams) (*ListAllIncidentsResult, error)
	GetIncidentCounts(ctx context.Context) (*IncidentCounts, error)
}
