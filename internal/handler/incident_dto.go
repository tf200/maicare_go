package handler

import (
	"time"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/google/uuid"
)

// --- Incident Request / Response DTOs ---

type createIncidentRequest struct {
	ClientID                uuid.UUID `json:"client_id" binding:"required"`
	EmployeeID              uuid.UUID `json:"employee_id"`
	LocationID              uuid.UUID `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement" binding:"required"`
	InformedParties         []string  `json:"informed_parties"`
	OccurredAt              time.Time `json:"occurred_at"`
	IncidentType            string    `json:"incident_type" binding:"required"`
	SeverityOfIncident      string    `json:"severity_of_incident" binding:"required"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk" binding:"required"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	CauseCategories         []string  `json:"cause_categories"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury" binding:"required"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation" binding:"required"`
	FollowUpActions         []string  `json:"follow_up_actions"`
	FollowUpNotes           *string   `json:"follow_up_notes"`
	IsEmployeeAbsent        bool      `json:"is_employee_absent"`
	AdditionalDetails       *string   `json:"additional_details"`
	Emails                  []string  `json:"emails"`
}

type createIncidentResponse struct {
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

type listIncidentsRequest struct {
	httpapi.PageRequest
}

type listIncidentsResponse struct {
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

type getIncidentResponse struct {
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

type updateIncidentRequest struct {
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

type updateIncidentResponse struct {
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

type confirmIncidentResponse struct {
	ID      uuid.UUID `json:"id"`
	FileUrl *string   `json:"file_url"`
}

type listAllIncidentsRequest struct {
	httpapi.PageRequest
	IsConfirmed bool    `form:"is_confirmed" json:"is_confirmed"`
	Search      *string `form:"search" json:"search" binding:"omitempty,max=120"`
}

type listAllIncidentsResponse struct {
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

type getIncidentCountsResponse struct {
	SeriousFatalCount        int64 `json:"serious_fatal_count"`
	PendingConfirmationCount int64 `json:"pending_confirmation_count"`
	Past24hCount             int64 `json:"past_24h_count"`
}

// --- Mappers ---

func toCreateIncidentParams(req createIncidentRequest) domain.CreateIncidentParams {
	return domain.CreateIncidentParams{
		ClientID:                req.ClientID,
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     req.ReporterInvolvement,
		InformedParties:         req.InformedParties,
		OccurredAt:              req.OccurredAt,
		IncidentType:            req.IncidentType,
		SeverityOfIncident:      req.SeverityOfIncident,
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          req.RecurrenceRisk,
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		CauseCategories:         req.CauseCategories,
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          req.PhysicalInjury,
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     req.PsychologicalDamage,
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      req.NeededConsultation,
		FollowUpActions:         req.FollowUpActions,
		FollowUpNotes:           req.FollowUpNotes,
		IsEmployeeAbsent:        req.IsEmployeeAbsent,
		AdditionalDetails:       req.AdditionalDetails,
		Emails:                  req.Emails,
	}
}

func toCreateIncidentResponse(i domain.Incident) createIncidentResponse {
	return createIncidentResponse{
		ID:                      i.ID,
		EmployeeID:              i.EmployeeID,
		LocationID:              i.LocationID,
		ReporterInvolvement:     i.ReporterInvolvement,
		InformedParties:         i.InformedParties,
		OccurredAt:              i.OccurredAt,
		IncidentType:            i.IncidentType,
		SeverityOfIncident:      i.SeverityOfIncident,
		IncidentExplanation:     i.IncidentExplanation,
		RecurrenceRisk:          i.RecurrenceRisk,
		IncidentPreventSteps:    i.IncidentPreventSteps,
		IncidentTakenMeasures:   i.IncidentTakenMeasures,
		CauseCategories:         i.CauseCategories,
		CauseExplanation:        i.CauseExplanation,
		PhysicalInjury:          i.PhysicalInjury,
		PhysicalInjuryDesc:      i.PhysicalInjuryDesc,
		PsychologicalDamage:     i.PsychologicalDamage,
		PsychologicalDamageDesc: i.PsychologicalDamageDesc,
		NeededConsultation:      i.NeededConsultation,
		FollowUpActions:         i.FollowUpActions,
		FollowUpNotes:           i.FollowUpNotes,
		IsEmployeeAbsent:        i.IsEmployeeAbsent,
		AdditionalDetails:       i.AdditionalDetails,
		ClientID:                i.ClientID,
		Emails:                  i.Emails,
		UpdatedAt:               i.UpdatedAt,
		CreatedAt:               i.CreatedAt,
	}
}

func toListIncidentsResponse(i domain.IncidentListItem) listIncidentsResponse {
	return listIncidentsResponse{
		ID:                     i.ID,
		OccurredAt:             i.OccurredAt,
		IncidentType:           i.IncidentType,
		SeverityOfIncident:     i.SeverityOfIncident,
		IsConfirmed:            i.IsConfirmed,
		EmployeeFirstName:      i.EmployeeFirstName,
		EmployeeLastName:       i.EmployeeLastName,
		EmployeeProfilePicture: i.EmployeeProfilePicture,
		LocationName:           i.LocationName,
	}
}

func toGetIncidentResponse(i domain.Incident) getIncidentResponse {
	return getIncidentResponse{
		ID:                      i.ID,
		EmployeeID:              i.EmployeeID,
		EmployeeFirstName:       i.EmployeeFirstName,
		EmployeeLastName:        i.EmployeeLastName,
		LocationID:              i.LocationID,
		ReporterInvolvement:     i.ReporterInvolvement,
		InformedParties:         i.InformedParties,
		OccurredAt:              i.OccurredAt,
		IncidentType:            i.IncidentType,
		SeverityOfIncident:      i.SeverityOfIncident,
		IncidentExplanation:     i.IncidentExplanation,
		RecurrenceRisk:          i.RecurrenceRisk,
		IncidentPreventSteps:    i.IncidentPreventSteps,
		IncidentTakenMeasures:   i.IncidentTakenMeasures,
		CauseCategories:         i.CauseCategories,
		CauseExplanation:        i.CauseExplanation,
		PhysicalInjury:          i.PhysicalInjury,
		PhysicalInjuryDesc:      i.PhysicalInjuryDesc,
		PsychologicalDamage:     i.PsychologicalDamage,
		PsychologicalDamageDesc: i.PsychologicalDamageDesc,
		NeededConsultation:      i.NeededConsultation,
		FollowUpActions:         i.FollowUpActions,
		FollowUpNotes:           i.FollowUpNotes,
		IsEmployeeAbsent:        i.IsEmployeeAbsent,
		AdditionalDetails:       i.AdditionalDetails,
		ClientID:                i.ClientID,
		UpdatedAt:               i.UpdatedAt,
		CreatedAt:               i.CreatedAt,
		IsConfirmed:             i.IsConfirmed,
		LocationName:            i.LocationName,
		Emails:                  i.Emails,
	}
}

func toUpdateIncidentParams(req updateIncidentRequest, incidentID uuid.UUID) domain.UpdateIncidentParams {
	return domain.UpdateIncidentParams{
		ID:                      incidentID,
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     req.ReporterInvolvement,
		InformedParties:         req.InformedParties,
		OccurredAt:              req.OccurredAt,
		IncidentType:            req.IncidentType,
		SeverityOfIncident:      req.SeverityOfIncident,
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          req.RecurrenceRisk,
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		CauseCategories:         req.CauseCategories,
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          req.PhysicalInjury,
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     req.PsychologicalDamage,
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      req.NeededConsultation,
		FollowUpActions:         req.FollowUpActions,
		FollowUpNotes:           req.FollowUpNotes,
		IsEmployeeAbsent:        req.IsEmployeeAbsent,
		AdditionalDetails:       req.AdditionalDetails,
		Emails:                  req.Emails,
	}
}

func toUpdateIncidentResponse(i domain.Incident) updateIncidentResponse {
	return updateIncidentResponse{
		ID:                      i.ID,
		EmployeeID:              i.EmployeeID,
		LocationID:              i.LocationID,
		ReporterInvolvement:     i.ReporterInvolvement,
		InformedParties:         i.InformedParties,
		OccurredAt:              i.OccurredAt,
		IncidentType:            i.IncidentType,
		SeverityOfIncident:      i.SeverityOfIncident,
		IncidentExplanation:     i.IncidentExplanation,
		RecurrenceRisk:          i.RecurrenceRisk,
		IncidentPreventSteps:    i.IncidentPreventSteps,
		IncidentTakenMeasures:   i.IncidentTakenMeasures,
		CauseCategories:         i.CauseCategories,
		CauseExplanation:        i.CauseExplanation,
		PhysicalInjury:          i.PhysicalInjury,
		PhysicalInjuryDesc:      i.PhysicalInjuryDesc,
		PsychologicalDamage:     i.PsychologicalDamage,
		PsychologicalDamageDesc: i.PsychologicalDamageDesc,
		NeededConsultation:      i.NeededConsultation,
		FollowUpActions:         i.FollowUpActions,
		FollowUpNotes:           i.FollowUpNotes,
		IsEmployeeAbsent:        i.IsEmployeeAbsent,
		AdditionalDetails:       i.AdditionalDetails,
		ClientID:                i.ClientID,
		UpdatedAt:               i.UpdatedAt,
		CreatedAt:               i.CreatedAt,
		IsConfirmed:             i.IsConfirmed,
		Emails:                  i.Emails,
	}
}

func toConfirmIncidentResponse(r domain.ConfirmIncidentResult) confirmIncidentResponse {
	return confirmIncidentResponse{
		ID:      r.ID,
		FileUrl: r.FileUrl,
	}
}

func toListAllIncidentsResponse(i domain.IncidentSummary) listAllIncidentsResponse {
	return listAllIncidentsResponse{
		ID:                 i.ID,
		OccurredAt:         i.OccurredAt,
		IncidentType:       i.IncidentType,
		SeverityOfIncident: i.SeverityOfIncident,
		IsConfirmed:        i.IsConfirmed,
		ClientFirstName:    i.ClientFirstName,
		ClientLastName:     i.ClientLastName,
		ClientBSN:          i.ClientBSN,
		EmployeeFirstName:  i.EmployeeFirstName,
		EmployeeLastName:   i.EmployeeLastName,
		LocationName:       i.LocationName,
	}
}

func toGetIncidentCountsResponse(c domain.IncidentCounts) getIncidentCountsResponse {
	return getIncidentCountsResponse{
		SeriousFatalCount:        c.SeriousFatalCount,
		PendingConfirmationCount: c.PendingConfirmationCount,
		Past24hCount:             c.Past24hCount,
	}
}
