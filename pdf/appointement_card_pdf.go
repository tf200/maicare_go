package pdf

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"maicare_go/bucket"
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

func GenerateAppointmentCardPDF(appointmentCardData AppointmentCard) ([]byte, error) {
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

	return pdfg.Bytes(), nil
}

// UploadIncidentPDF uploads a PDF to B2 with a generated filename
func UploadAppointmentCardPDF(ctx context.Context, pdfBytes []byte, appointmentCardID int64, b2Client *bucket.B2Client) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("appointment_cards/%s/appointment_card_%d.pdf", timestamp, appointmentCardID)

	// Create reader from PDF bytes
	pdfReader := bytes.NewReader(pdfBytes)

	// Upload to B2
	obj := b2Client.Bucket.Object(filename)
	writer := obj.NewWriter(ctx)
	writer.ConcurrentUploads = 4

	// Copy the PDF data to B2
	_, err := io.Copy(writer, pdfReader)
	if err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to copy PDF to B2: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close B2 writer: %w", err)
	}
	fileURL := fmt.Sprintf("%s/file/%s/%s",
		b2Client.Bucket.BaseURL(),
		b2Client.Bucket.Name(),
		filename)

	return fileURL, nil
}

// Helper function to do both operations if needed
func GenerateAndUploadAppointmentCardPDF(ctx context.Context, cardData AppointmentCard, b2Client *bucket.B2Client) (string, error) {
	// Generate PDF
	pdfBytes, err := GenerateAppointmentCardPDF(cardData)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	fileURL, err := UploadAppointmentCardPDF(ctx, pdfBytes, cardData.ID, b2Client)
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	return fileURL, nil
}
