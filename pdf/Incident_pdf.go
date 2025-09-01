package pdf

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"maicare_go/bucket"
	"mime/multipart"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

//go:embed templates/incident.html
var incidentTemplateFS embed.FS

type IncidentReportData struct {
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
	ClientFirstName         string    `json:"client_firstname"`
	ClientLastName          string    `json:"client_lastname"`
	LocationName            string    `json:"location_name"`
}

// GenerateIncidentPDF generates a PDF from incident data and returns the PDF bytes
func GenerateIncidentPDF(incidentData IncidentReportData) (multipart.File, error) {

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"now": func() time.Time { // You might also need a 'now' function if it's not predefined
			return time.Now().In(time.FixedZone("GMT+1", 1*60*60)) // Set to Tangier's timezone
		},
	}

	// Parse and execute HTML template
	templ, err := template.New("incident.html").Funcs(funcMap).ParseFS(incidentTemplateFS, "templates/incident.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := templ.Execute(&body, incidentData); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Create PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF generator: %w", err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	// Create a new input page from our HTML
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// Set page options
	page.EnableLocalFileAccess.Set(true)
	page.LoadErrorHandling.Set("ignore")

	// Add page to generator
	pdfg.AddPage(page)

	// Generate PDF
	err = pdfg.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF: %w", err)
	}

	// Get the generated PDF as a byte slice
	pdfBytes := pdfg.Bytes()
	// Wrap the byte slice in our InMemoryFile to satisfy the interface
	file := &bucket.InMemoryFile{
		Reader: bytes.NewReader(pdfBytes),
	}
	return file, nil
}

// UploadIncidentPDF uploads a PDF to B2 with a generated filename
func UploadIncidentPDF(ctx context.Context, pdfFile multipart.File, incidentID int64, b2Client *bucket.ObjectStorageClient) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("incident_reports/%s/incident_report_%d.pdf", timestamp, incidentID)

	// Upload to B2
	key, _, err := b2Client.Upload(ctx, pdfFile, filename, "application/pdf")
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF to B2: %w", err)
	}
	return key, nil
}

// Helper function to do both operations if needed
func GenerateAndUploadIncidentPDF(ctx context.Context, incidentData IncidentReportData, b2Client *bucket.ObjectStorageClient) (string, error) {
	// Generate PDF
	pdfFile, err := GenerateIncidentPDF(incidentData)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	filename, err := UploadIncidentPDF(ctx, pdfFile, incidentData.ID, b2Client)
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	return filename, nil
}
