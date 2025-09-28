package clientp

import (
	"maicare_go/pagination"
	"time"
)

// CreateIncidentRequest represents a request to create an incident
type CreateIncidentRequest struct {
	EmployeeID              int64     `json:"employee_id"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement" binding:"required" enums:"directly_involved,witness,found_afterwards,alarmed"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         string    `json:"runtime_incident" binding:"required"`
	IncidentType            string    `json:"incident_type" binding:"required"`
	PassingAway             bool      `json:"passing_away"`
	SelfHarm                bool      `json:"self_harm"`
	Violence                bool      `json:"violence"`
	FireWaterDamage         bool      `json:"fire_water_damage"`
	Accident                bool      `json:"accident"`
	ClientAbsence           bool      `json:"client_absence"`
	Medicines               bool      `json:"medicines"`
	Organization            bool      `json:"organization"`
	UseProhibitedSubstances bool      `json:"use_prohibited_substances"`
	OtherNotifications      bool      `json:"other_notifications"`
	SeverityOfIncident      string    `json:"severity_of_incident" binding:"required" enums:"fatal,serious,less_serious,near_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk" binding:"required" enums:"high,very_high,means,very_low"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury" binding:"required" enums:"no_injuries,not_noticeable_yet,bruising_swelling,broken_bones,shortness_of_breath,death,other"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation" binding:"required" enums:"no,not_clear,hospitalization,consult_gp"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	Emails                  []string  `json:"emails"`
}

// CreateIncidentResponse represents a response for CreateIncidentApi
type CreateIncidentResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         string    `json:"runtime_incident"`
	IncidentType            string    `json:"incident_type"`
	PassingAway             bool      `json:"passing_away"`
	SelfHarm                bool      `json:"self_harm"`
	Violence                bool      `json:"violence"`
	FireWaterDamage         bool      `json:"fire_water_damage"`
	Accident                bool      `json:"accident"`
	ClientAbsence           bool      `json:"client_absence"`
	Medicines               bool      `json:"medicines"`
	Organization            bool      `json:"organization"`
	UseProhibitedSubstances bool      `json:"use_prohibited_substances"`
	OtherNotifications      bool      `json:"other_notifications"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	Emails                  []string  `json:"emails"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated"`
	CreatedAt               time.Time `json:"created"`
}

// ListIncidentsRequest defines the request for listing incidents
type ListIncidentsRequest struct {
	pagination.Request
}

// ListIncidentsResponse defines the response for listing incidents
type ListIncidentsResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	EmployeeFirstName       string    `json:"employee_first_name"`
	EmployeeLastName        string    `json:"employee_last_name"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         string    `json:"runtime_incident"`
	IncidentType            string    `json:"incident_type"`
	PassingAway             bool      `json:"passing_away"`
	SelfHarm                bool      `json:"self_harm"`
	Violence                bool      `json:"violence"`
	FireWaterDamage         bool      `json:"fire_water_damage"`
	Accident                bool      `json:"accident"`
	ClientAbsence           bool      `json:"client_absence"`
	Medicines               bool      `json:"medicines"`
	Organization            bool      `json:"organization"`
	UseProhibitedSubstances bool      `json:"use_prohibited_substances"`
	OtherNotifications      bool      `json:"other_notifications"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	Emails                  []string  `json:"emails"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
	IsConfirmed             bool      `json:"is_confirmed"`
	EmployeeProfilePicture  *string   `json:"employee_profile_picture"`
	LocationName            string    `json:"location_name"`
}

// GetIncidentResponse represents a response for GetIncidentApi
type GetIncidentResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	EmployeeFirstName       string    `json:"employee_first_name"`
	EmployeeLastName        string    `json:"employee_last_name"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         string    `json:"runtime_incident"`
	IncidentType            string    `json:"incident_type"`
	PassingAway             bool      `json:"passing_away"`
	SelfHarm                bool      `json:"self_harm"`
	Violence                bool      `json:"violence"`
	FireWaterDamage         bool      `json:"fire_water_damage"`
	Accident                bool      `json:"accident"`
	ClientAbsence           bool      `json:"client_absence"`
	Medicines               bool      `json:"medicines"`
	Organization            bool      `json:"organization"`
	UseProhibitedSubstances bool      `json:"use_prohibited_substances"`
	OtherNotifications      bool      `json:"other_notifications"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
	IsConfirmed             bool      `json:"is_confirmed"`
	LocationName            string    `json:"location_name"`
	Emails                  []string  `json:"emails"`
}

// UpdateIncidentRequest represents a request to update an incident
type UpdateIncidentRequest struct {
	ID                      int64     `json:"id"`
	EmployeeID              *int64    `json:"employee_id"`
	LocationID              *int64    `json:"location_id"`
	ReporterInvolvement     *string   `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         *string   `json:"runtime_incident"`
	IncidentType            *string   `json:"incident_type"`
	PassingAway             *bool     `json:"passing_away"`
	SelfHarm                *bool     `json:"self_harm"`
	Violence                *bool     `json:"violence"`
	FireWaterDamage         *bool     `json:"fire_water_damage"`
	Accident                *bool     `json:"accident"`
	ClientAbsence           *bool     `json:"client_absence"`
	Medicines               *bool     `json:"medicines"`
	Organization            *bool     `json:"organization"`
	UseProhibitedSubstances *bool     `json:"use_prohibited_substances"`
	OtherNotifications      *bool     `json:"other_notifications"`
	SeverityOfIncident      *string   `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          *string   `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          *string   `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     *string   `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      *string   `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   *bool     `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     *string   `json:"employee_absenteeism"`
	Emails                  []string  `json:"emails"`
}

// UpdateIncidentResponse represents a response for UpdateIncidentApi
type UpdateIncidentResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         string    `json:"runtime_incident"`
	IncidentType            string    `json:"incident_type"`
	PassingAway             bool      `json:"passing_away"`
	SelfHarm                bool      `json:"self_harm"`
	Violence                bool      `json:"violence"`
	FireWaterDamage         bool      `json:"fire_water_damage"`
	Accident                bool      `json:"accident"`
	ClientAbsence           bool      `json:"client_absence"`
	Medicines               bool      `json:"medicines"`
	Organization            bool      `json:"organization"`
	UseProhibitedSubstances bool      `json:"use_prohibited_substances"`
	OtherNotifications      bool      `json:"other_notifications"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated"`
	CreatedAt               time.Time `json:"created"`
	IsConfirmed             bool      `json:"is_confirmed"`
	Emails                  []string  `json:"emails"`
}

// GenerateIncidentFileResponse represents a response for GenerateIncidentFileApi
type GenerateIncidentFileResponse struct {
	FileUrl *string `json:"file_url"`
	ID      int64   `json:"incident_id"`
}

// ConfirmIncidentResponse represents a response for ConfirmIncidentApi
type ConfirmIncidentResponse struct {
	ID      int64   `json:"id"`
	FileUrl *string `json:"file_url"`
}
