package pdf

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IncidentReportData struct {
	ID                      uuid.UUID `json:"id"`
	EmployeeID              uuid.UUID `json:"employee_id"`
	EmployeeFirstName       string    `json:"employee_first_name"`
	EmployeeLastName        string    `json:"employee_last_name"`
	LocationID              uuid.UUID `json:"location_id"`
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
	ClientID                uuid.UUID `json:"client_id"`
	ClientFirstName         string    `json:"client_firstname"`
	ClientLastName          string    `json:"client_lastname"`
	LocationName            string    `json:"location_name"`
}

// GenerateIncidentPDF generates a PDF from incident data and returns the PDF bytes
func (s *pdfService) generateIncidentPDF(incidentData IncidentReportData) (multipart.File, error) {
	headerLines := []string{
		fmt.Sprintf("Incident ID: %s", incidentData.ID),
		fmt.Sprintf("Incident date: %s", incidentData.IncidentDate.Format(time.RFC3339)),
		fmt.Sprintf("Location: %s (%s)", incidentData.LocationName, incidentData.LocationID),
		fmt.Sprintf("Reporter: %s %s (%s)", incidentData.EmployeeFirstName, incidentData.EmployeeLastName, incidentData.EmployeeID),
		fmt.Sprintf("Client: %s %s (%s)", incidentData.ClientFirstName, incidentData.ClientLastName, incidentData.ClientID),
	}

	sections := []documentSection{
		{
			Title: "Incident details",
			Lines: []string{
				fmt.Sprintf("Reporter involvement: %s", incidentData.ReporterInvolvement),
				fmt.Sprintf("Runtime incident: %s", incidentData.RuntimeIncident),
				fmt.Sprintf("Incident type: %s", incidentData.IncidentType),
				fmt.Sprintf("Severity: %s", incidentData.SeverityOfIncident),
				fmt.Sprintf("Inform who: %s", joinOrNA(incidentData.InformWho)),
			},
		},
		{
			Title: "Incident categories",
			Lines: []string{
				fmt.Sprintf("Passing away: %s", yesNo(incidentData.PassingAway)),
				fmt.Sprintf("Self harm: %s", yesNo(incidentData.SelfHarm)),
				fmt.Sprintf("Violence: %s", yesNo(incidentData.Violence)),
				fmt.Sprintf("Fire/water damage: %s", yesNo(incidentData.FireWaterDamage)),
				fmt.Sprintf("Accident: %s", yesNo(incidentData.Accident)),
				fmt.Sprintf("Client absence: %s", yesNo(incidentData.ClientAbsence)),
				fmt.Sprintf("Medicines: %s", yesNo(incidentData.Medicines)),
				fmt.Sprintf("Organization: %s", yesNo(incidentData.Organization)),
				fmt.Sprintf("Use prohibited substances: %s", yesNo(incidentData.UseProhibitedSubstances)),
				fmt.Sprintf("Other notifications: %s", yesNo(incidentData.OtherNotifications)),
			},
		},
		{
			Title: "Impact and risk",
			Lines: []string{
				fmt.Sprintf("Incident explanation: %s", stringOrNA(incidentData.IncidentExplanation)),
				fmt.Sprintf("Recurrence risk: %s", incidentData.RecurrenceRisk),
				fmt.Sprintf("Preventive steps: %s", stringOrNA(incidentData.IncidentPreventSteps)),
				fmt.Sprintf("Taken measures: %s", stringOrNA(incidentData.IncidentTakenMeasures)),
			},
		},
		{
			Title: "Cause analysis",
			Lines: []string{
				fmt.Sprintf("Technical: %s", joinOrNA(incidentData.Technical)),
				fmt.Sprintf("Organizational: %s", joinOrNA(incidentData.Organizational)),
				fmt.Sprintf("Mese worker: %s", joinOrNA(incidentData.MeseWorker)),
				fmt.Sprintf("Client options: %s", joinOrNA(incidentData.ClientOptions)),
				fmt.Sprintf("Other cause: %s", stringOrNA(incidentData.OtherCause)),
				fmt.Sprintf("Cause explanation: %s", stringOrNA(incidentData.CauseExplanation)),
			},
		},
		{
			Title: "Damage and follow-up",
			Lines: []string{
				fmt.Sprintf("Physical injury: %s", incidentData.PhysicalInjury),
				fmt.Sprintf("Physical injury description: %s", stringOrNA(incidentData.PhysicalInjuryDesc)),
				fmt.Sprintf("Psychological damage: %s", incidentData.PsychologicalDamage),
				fmt.Sprintf("Psychological damage description: %s", stringOrNA(incidentData.PsychologicalDamageDesc)),
				fmt.Sprintf("Needed consultation: %s", incidentData.NeededConsultation),
				fmt.Sprintf("Succession: %s", joinOrNA(incidentData.Succession)),
				fmt.Sprintf("Succession description: %s", stringOrNA(incidentData.SuccessionDesc)),
				fmt.Sprintf("Other: %s", yesNo(incidentData.Other)),
				fmt.Sprintf("Other description: %s", stringOrNA(incidentData.OtherDesc)),
				fmt.Sprintf("Additional appointments: %s", stringOrNA(incidentData.AdditionalAppointments)),
				fmt.Sprintf("Employee absenteeism: %s", incidentData.EmployeeAbsenteeism),
			},
		},
	}

	pdfBytes, err := buildSectionsPDF("Incident report", headerLines, sections)
	if err != nil {
		return nil, fmt.Errorf("failed to generate incident pdf: %w", err)
	}

	return toMultipartFile(pdfBytes), nil
}

// UploadIncidentPDF uploads a PDF to B2 with a generated filename
func (s *pdfService) uploadIncidentPDF(ctx context.Context, pdfFile multipart.File, incidentID uuid.UUID) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("incident_reports/%s/incident_report_%s.pdf", timestamp, incidentID.String())

	// Upload to B2
	key, _, err := s.bucketClient.Upload(ctx, pdfFile, filename, "application/pdf")
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF to B2: %w", err)
	}
	return key, nil
}

// Helper function to do both operations if needed
func (s *pdfService) GenerateAndUploadIncidentPDF(ctx context.Context, incidentData IncidentReportData) (string, error) {
	// Generate PDF
	pdfFile, err := s.generateIncidentPDF(incidentData)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	filename, err := s.uploadIncidentPDF(ctx, pdfFile, incidentData.ID)
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	return filename, nil
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func stringOrNA(value *string) string {
	if value == nil || *value == "" {
		return "N/A"
	}
	return *value
}

func joinOrNA(values []string) string {
	if len(values) == 0 {
		return "N/A"
	}
	return strings.Join(values, "; ")
}
