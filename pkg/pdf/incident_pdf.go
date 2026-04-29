package pdf

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type IncidentReportData struct {
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
	ClientFirstName         string    `json:"client_firstname"`
	ClientLastName          string    `json:"client_lastname"`
	LocationName            string    `json:"location_name"`
}

func (s *pdfService) GenerateIncidentPDF(ctx context.Context, incidentData IncidentReportData) ([]byte, error) {
	_ = ctx

	pdfBytes, err := buildIncidentReportPDF(incidentData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate incident pdf: %w", err)
	}

	return pdfBytes, nil
}

func buildIncidentReportPDF(incidentData IncidentReportData) ([]byte, error) {
	primary := &props.Color{Red: 24, Green: 74, Blue: 96}
	muted := &props.Color{Red: 89, Green: 101, Blue: 113}
	lightBlue := &props.Color{Red: 232, Green: 242, Blue: 247}
	lightGray := &props.Color{Red: 246, Green: 248, Blue: 250}
	borderColor := &props.Color{Red: 214, Green: 222, Blue: 228}
	severityColor := incidentSeverityColor(incidentData.SeverityOfIncident)

	cfg := config.NewBuilder().
		WithLeftMargin(12).
		WithRightMargin(12).
		WithTopMargin(12).
		WithBottomMargin(12).
		Build()

	m := maroto.New(cfg)

	m.AddRow(20,
		text.NewCol(8, "Incident Report", props.Text{Style: fontstyle.Bold, Size: 20, Color: primary, Top: 4}),
		text.NewCol(4, fmt.Sprintf("Generated\n%s", time.Now().Format("2006-01-02 15:04")), props.Text{Size: 8, Align: align.Right, Color: muted, Top: 3}),
	).WithStyle(&props.Cell{BackgroundColor: lightBlue})

	m.AddRow(8,
		text.NewCol(8, fmt.Sprintf("Incident ID: %s", incidentData.ID.String()), props.Text{Size: 8, Color: muted, Top: 2}),
		text.NewCol(4, strings.ToUpper(displayValue(incidentData.SeverityOfIncident)), props.Text{Style: fontstyle.Bold, Size: 11, Align: align.Center, Color: &props.WhiteColor, Top: 2}),
	).WithStyle(&props.Cell{BackgroundColor: severityColor})

	addIncidentSpacer(m, 4)
	addIncidentInfoGrid(m, incidentData, primary, muted, lightGray, borderColor)
	addIncidentSpacer(m, 4)

	addIncidentSection(m, "Incident Details", primary, lightBlue, borderColor, []incidentField{
		{Label: "Occurred at", Value: formatTimeOrNA(incidentData.OccurredAt)},
		{Label: "Location", Value: fmt.Sprintf("%s (%s)", displayValue(incidentData.LocationName), incidentData.LocationID.String())},
		{Label: "Reporter involvement", Value: incidentData.ReporterInvolvement},
		{Label: "Incident type", Value: incidentData.IncidentType},
		{Label: "Informed parties", Value: joinOrNA(incidentData.InformedParties)},
	})

	addIncidentSection(m, "Impact And Risk", primary, lightBlue, borderColor, []incidentField{
		{Label: "Explanation", Value: stringOrNA(incidentData.IncidentExplanation)},
		{Label: "Recurrence risk", Value: incidentData.RecurrenceRisk},
		{Label: "Preventive steps", Value: stringOrNA(incidentData.IncidentPreventSteps)},
		{Label: "Taken measures", Value: stringOrNA(incidentData.IncidentTakenMeasures)},
	})

	addIncidentSection(m, "Cause Analysis", primary, lightBlue, borderColor, []incidentField{
		{Label: "Cause categories", Value: joinOrNA(incidentData.CauseCategories)},
		{Label: "Cause explanation", Value: stringOrNA(incidentData.CauseExplanation)},
	})

	addIncidentSection(m, "Damage And Follow-Up", primary, lightBlue, borderColor, []incidentField{
		{Label: "Physical injury", Value: incidentData.PhysicalInjury},
		{Label: "Physical injury details", Value: stringOrNA(incidentData.PhysicalInjuryDesc)},
		{Label: "Psychological damage", Value: incidentData.PsychologicalDamage},
		{Label: "Psychological damage details", Value: stringOrNA(incidentData.PsychologicalDamageDesc)},
		{Label: "Needed consultation", Value: incidentData.NeededConsultation},
		{Label: "Follow-up actions", Value: joinOrNA(incidentData.FollowUpActions)},
		{Label: "Follow-up notes", Value: stringOrNA(incidentData.FollowUpNotes)},
		{Label: "Employee absent", Value: yesNo(incidentData.IsEmployeeAbsent)},
		{Label: "Additional details", Value: stringOrNA(incidentData.AdditionalDetails)},
	})

	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate incident maroto document: %w", err)
	}

	return document.GetBytes(), nil
}

type incidentField struct {
	Label string
	Value string
}

func addIncidentInfoGrid(m core.Maroto, incidentData IncidentReportData, primary, muted, background, borderColor *props.Color) {
	m.AddRow(9,
		text.NewCol(6, "Client", props.Text{Style: fontstyle.Bold, Size: 9, Color: muted, Top: 2, Left: 2}),
		text.NewCol(6, "Reporter", props.Text{Style: fontstyle.Bold, Size: 9, Color: muted, Top: 2, Left: 2}),
	).WithStyle(&props.Cell{BackgroundColor: background, BorderType: border.Full, BorderColor: borderColor, BorderThickness: 0.1})
	m.AddRow(13,
		text.NewCol(6, fmt.Sprintf("%s %s\n%s", displayValue(incidentData.ClientFirstName), displayValue(incidentData.ClientLastName), incidentData.ClientID.String()), props.Text{Size: 10, Color: primary, Top: 2, Left: 2}),
		text.NewCol(6, fmt.Sprintf("%s %s\n%s", displayValue(incidentData.EmployeeFirstName), displayValue(incidentData.EmployeeLastName), incidentData.EmployeeID.String()), props.Text{Size: 10, Color: primary, Top: 2, Left: 2}),
	).WithStyle(&props.Cell{BorderType: border.Full, BorderColor: borderColor, BorderThickness: 0.1})
}

func addIncidentSection(m core.Maroto, title string, primary, headerBackground, borderColor *props.Color, fields []incidentField) {
	addIncidentSpacer(m, 3)
	m.AddRow(9, text.NewCol(12, title, props.Text{Style: fontstyle.Bold, Size: 11, Color: primary, Top: 2, Left: 2})).WithStyle(&props.Cell{BackgroundColor: headerBackground, BorderType: border.Full, BorderColor: borderColor, BorderThickness: 0.1})

	for _, field := range fields {
		m.AddAutoRow(
			text.NewCol(3, field.Label, props.Text{Style: fontstyle.Bold, Size: 9, Color: primary, Top: 2, Left: 2, Bottom: 2}),
			text.NewCol(9, displayValue(field.Value), props.Text{Size: 9, Top: 2, Left: 2, Bottom: 2}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderColor: borderColor, BorderThickness: 0.1})
	}
}

func addIncidentSpacer(m core.Maroto, height float64) {
	m.AddRow(height, text.NewCol(12, "", props.Text{Size: 1}))
}

func incidentSeverityColor(severity string) *props.Color {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "fatal", "serious", "critical", "high":
		return &props.Color{Red: 172, Green: 44, Blue: 45}
	case "moderate", "medium":
		return &props.Color{Red: 189, Green: 117, Blue: 37}
	default:
		return &props.Color{Red: 24, Green: 74, Blue: 96}
	}
}

func (s *pdfService) generateIncidentPDF(incidentData IncidentReportData) (multipart.File, error) {
	pdfBytes, err := s.GenerateIncidentPDF(context.Background(), incidentData)
	if err != nil {
		return nil, err
	}
	return toMultipartFile(pdfBytes), nil
}

func (s *pdfService) uploadIncidentPDF(ctx context.Context, pdfFile multipart.File, incidentID uuid.UUID) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("incident_reports/%s/incident_report_%s.pdf", timestamp, incidentID.String())

	key, _, err := s.bucketClient.Upload(ctx, pdfFile, filename, "application/pdf")
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF to B2: %w", err)
	}
	return key, nil
}

func (s *pdfService) GenerateAndUploadIncidentPDF(ctx context.Context, incidentData IncidentReportData) (string, error) {
	pdfFile, err := s.generateIncidentPDF(incidentData)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

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

func displayValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "N/A"
	}
	return value
}

func formatTimeOrNA(value time.Time) string {
	if value.IsZero() {
		return "N/A"
	}
	return value.Format(time.RFC3339)
}
