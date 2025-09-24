package pdf

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"maicare_go/bucket"
	"mime/multipart"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

//go:embed templates/appointment_card.html
var appointmentCardTemplateFS embed.FS

type AppointmentCard struct {
	ID                     int64
	ClientName             string
	Date                   string
	Mentor                 string
	GeneralInformation     []string
	ImportantContacts      []string
	HouseholdInfo          []string
	OrganizationAgreements []string
	YouthOfficerAgreements []string
	TreatmentAgreements    []string
	SmokingRules           []string
	Work                   []string
	SchoolInternship       []string
	Travel                 []string
	Leave                  []string
}

func GenerateAppointmentCardPDF(appointmentCardData AppointmentCard) (multipart.File, error) {
	// Parse and execute HTML template
	templ, err := template.ParseFS(appointmentCardTemplateFS, "templates/appointment_card.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := templ.Execute(&body, appointmentCardData); err != nil {
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
func UploadAppointmentCardPDF(ctx context.Context, pdfFile multipart.File, appointmentCardID int64, b2Client bucket.ObjectStorageInterface) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("appointment_cards/%s/appointment_card_%d.pdf", timestamp, appointmentCardID)

	// Upload to B2
	key, _, err := b2Client.Upload(ctx, pdfFile, filename, "application/pdf")
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF to B2: %w", err)
	}

	return key, nil
}

// Helper function to do both operations if needed
func GenerateAndUploadAppointmentCardPDF(ctx context.Context, cardData AppointmentCard, b2Client bucket.ObjectStorageInterface) (string, error) {
	// Generate PDF
	pdfFile, err := GenerateAppointmentCardPDF(cardData)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	fileURL, err := UploadAppointmentCardPDF(ctx, pdfFile, cardData.ID, b2Client)
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	return fileURL, nil
}
