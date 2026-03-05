package clientp

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

// CreateIncidentRequest represents a request to create an incident
type CreateIncidentRequest struct {
	ClientID                uuid.UUID `json:"client_id" binding:"required"`
	EmployeeID              uuid.UUID `json:"employee_id"`
	LocationID              uuid.UUID `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement" binding:"required" enums:"directly_involved,witness,found_afterwards,alarmed"`
	InformedParties         []string  `json:"informed_parties"`
	OccurredAt              time.Time `json:"occurred_at"`
	IncidentType            string    `json:"incident_type" binding:"required"`
	SeverityOfIncident      string    `json:"severity_of_incident" binding:"required" enums:"fatal,serious,less_serious,near_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk" binding:"required" enums:"high,very_high,means,very_low"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	CauseCategories         []string  `json:"cause_categories"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury" binding:"required" enums:"no_injuries,not_noticeable_yet,bruising_swelling,broken_bones,shortness_of_breath,death,other"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation" binding:"required" enums:"no,not_clear,hospitalization,consult_gp"`
	FollowUpActions         []string  `json:"follow_up_actions"`
	FollowUpNotes           *string   `json:"follow_up_notes"`
	IsEmployeeAbsent        bool      `json:"is_employee_absent"`
	AdditionalDetails       *string   `json:"additional_details"`
	Emails                  []string  `json:"emails"`
}

// CreateIncidentResponse represents a response for CreateIncidentApi
type CreateIncidentResponse struct {
	ID                      uuid.UUID `json:"id"`
	EmployeeID              uuid.UUID `json:"employee_id"`
	LocationID              uuid.UUID `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformedParties         []string  `json:"informed_parties"`
	OccurredAt              time.Time `json:"occurred_at"`
	IncidentType            string    `json:"incident_type"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	CauseCategories         []string  `json:"cause_categories"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	FollowUpActions         []string  `json:"follow_up_actions"`
	FollowUpNotes           *string   `json:"follow_up_notes"`
	IsEmployeeAbsent        bool      `json:"is_employee_absent"`
	AdditionalDetails       *string   `json:"additional_details"`
	ClientID                uuid.UUID `json:"client_id"`
	Emails                  []string  `json:"emails"`
	UpdatedAt               time.Time `json:"updated"`
	CreatedAt               time.Time `json:"created"`
}

// ListIncidentsRequest defines the request for listing incidents
type ListIncidentsRequest struct {
	pagination.Request
}

// ListIncidentsResponse defines the response for listing incidents
type ListIncidentsResponse struct {
	ID                     uuid.UUID `json:"id"`
	OccurredAt             time.Time `json:"occurred_at"`
	IncidentType           string    `json:"incident_type"`
	SeverityOfIncident     string    `json:"severity_of_incident"`
	IsConfirmed            bool      `json:"is_confirmed"`
	EmployeeFirstName      string    `json:"employee_first_name"`
	EmployeeLastName       string    `json:"employee_last_name"`
	EmployeeProfilePicture *string   `json:"employee_profile_picture"`
	LocationName           string    `json:"location_name"`
}

// GetIncidentResponse represents a response for GetIncidentApi
type GetIncidentResponse struct {
	ID                      uuid.UUID `json:"id"`
	EmployeeID              uuid.UUID `json:"employee_id"`
	EmployeeFirstName       string    `json:"employee_first_name"`
	EmployeeLastName        string    `json:"employee_last_name"`
	LocationID              uuid.UUID `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformedParties         []string  `json:"informed_parties"`
	OccurredAt              time.Time `json:"occurred_at"`
	IncidentType            string    `json:"incident_type"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	CauseCategories         []string  `json:"cause_categories"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	FollowUpActions         []string  `json:"follow_up_actions"`
	FollowUpNotes           *string   `json:"follow_up_notes"`
	IsEmployeeAbsent        bool      `json:"is_employee_absent"`
	AdditionalDetails       *string   `json:"additional_details"`
	ClientID                uuid.UUID `json:"client_id"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
	IsConfirmed             bool      `json:"is_confirmed"`
	LocationName            string    `json:"location_name"`
	Emails                  []string  `json:"emails"`
}

// UpdateIncidentRequest represents a request to update an incident
type UpdateIncidentRequest struct {
	ID                      uuid.UUID  `json:"id"`
	EmployeeID              *uuid.UUID `json:"employee_id"`
	LocationID              *uuid.UUID `json:"location_id"`
	ReporterInvolvement     *string    `json:"reporter_involvement"`
	InformedParties         []string   `json:"informed_parties"`
	OccurredAt              time.Time  `json:"occurred_at"`
	IncidentType            *string    `json:"incident_type"`
	SeverityOfIncident      *string    `json:"severity_of_incident"`
	IncidentExplanation     *string    `json:"incident_explanation"`
	RecurrenceRisk          *string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string    `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string    `json:"incident_taken_measures"`
	CauseCategories         []string   `json:"cause_categories"`
	CauseExplanation        *string    `json:"cause_explanation"`
	PhysicalInjury          *string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string    `json:"physical_injury_desc"`
	PsychologicalDamage     *string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string    `json:"psychological_damage_desc"`
	NeededConsultation      *string    `json:"needed_consultation"`
	FollowUpActions         []string   `json:"follow_up_actions"`
	FollowUpNotes           *string    `json:"follow_up_notes"`
	IsEmployeeAbsent        *bool      `json:"is_employee_absent"`
	AdditionalDetails       *string    `json:"additional_details"`
	Emails                  []string   `json:"emails"`
}

// UpdateIncidentResponse represents a response for UpdateIncidentApi
type UpdateIncidentResponse struct {
	ID                      uuid.UUID `json:"id"`
	EmployeeID              uuid.UUID `json:"employee_id"`
	LocationID              uuid.UUID `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformedParties         []string  `json:"informed_parties"`
	OccurredAt              time.Time `json:"occurred_at"`
	IncidentType            string    `json:"incident_type"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	CauseCategories         []string  `json:"cause_categories"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	FollowUpActions         []string  `json:"follow_up_actions"`
	FollowUpNotes           *string   `json:"follow_up_notes"`
	IsEmployeeAbsent        bool      `json:"is_employee_absent"`
	AdditionalDetails       *string   `json:"additional_details"`
	ClientID                uuid.UUID `json:"client_id"`
	UpdatedAt               time.Time `json:"updated"`
	CreatedAt               time.Time `json:"created"`
	IsConfirmed             bool      `json:"is_confirmed"`
	Emails                  []string  `json:"emails"`
}

// GenerateIncidentFileResponse represents a response for GenerateIncidentFileApi
type GenerateIncidentFileResponse struct {
	FileUrl *string   `json:"file_url"`
	ID      uuid.UUID `json:"incident_id"`
}

// ConfirmIncidentResponse represents a response for ConfirmIncidentApi
type ConfirmIncidentResponse struct {
	ID      uuid.UUID `json:"id"`
	FileUrl *string   `json:"file_url"`
}

// ListAllIncidentsRequest represents the request body for listing all incidents
type ListAllIncidentsRequest struct {
	pagination.Request
	IsConfirmed bool    `form:"is_confirmed" json:"is_confirmed"`
	Search      *string `form:"search" json:"search" binding:"omitempty,max=120"`
}

// ListAllIncidentsResponse represents the response body for listing all incidents
type ListAllIncidentsResponse struct {
	ID                 uuid.UUID `json:"id"`
	OccurredAt         time.Time `json:"occurred_at"`
	IncidentType       string    `json:"incident_type"`
	SeverityOfIncident string    `json:"severity_of_incident"`
	IsConfirmed        bool      `json:"is_confirmed"`
	ClientFirstName    string    `json:"client_first_name"`
	ClientLastName     string    `json:"client_last_name"`
	ClientBSN          *string   `json:"client_bsn"`
	EmployeeFirstName  string    `json:"employee_first_name"`
	EmployeeLastName   string    `json:"employee_last_name"`
	LocationName       string    `json:"location_name"`
}

type GetIncidentCountsResponse struct {
	SeriousFatalCount        int64 `json:"serious_fatal_count"`
	PendingConfirmationCount int64 `json:"pending_confirmation_count"`
	Past24hCount             int64 `json:"past_24h_count"`
}
