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

//go:embed templates/contract.html
var contractTemplateFS embed.FS

// ContractData represents the data structure for the care contract template.
type ContractData struct {
	// Header Information
	ID     int64  `json:"ContractID"` // Corresponds to {.ContractID}
	Status string `json:"Status"`     // Corresponds to {.Status} and used for status-badge class

	// Contract Periods
	StartDate      string `json:"StartDate"`      // Corresponds to {.StartDate} and [START_DATE] in terms
	EndDate        string `json:"EndDate"`        // Corresponds to {.EndDate} and [END_DATE] in terms
	ReminderPeriod int    `json:"ReminderPeriod"` // Corresponds to {.ReminderPeriod}

	// Parties
	SenderName        string `json:"SenderName"`        // Corresponds to {.SenderName}
	SenderAddress     string `json:"SenderAddress"`     // Corresponds to {.SenderAddress}
	SenderContactInfo string `json:"SenderContactInfo"` // Corresponds to {.SenderContactInfo}

	ClientFirstName   string `json:"ClientFirstName"`   // Corresponds to {.ClientFirstName}
	ClientLastName    string `json:"ClientLastName"`    // Corresponds to {.ClientLastName}
	ClientAddress     string `json:"ClientAddress"`     // Corresponds to {.ClientAddress}
	ClientContactInfo string `json:"ClientContactInfo"` // Corresponds to {.ClientContactInfo}

	// Care Specifications
	CareType        string `json:"CareType"`        // Corresponds to {.CareType}
	CareName        string `json:"CareName"`        // Corresponds to {.CareName}
	FinancingAct    string `json:"FinancingAct"`    // Corresponds to {.FinancingAct}
	FinancingOption string `json:"FinancingOption"` // Corresponds to {.FinancingOption}

	// Ambulante Care Hours (conditional display)
	Hours            float64 `json:"Hours"`            // Corresponds to {.Hours}
	HoursType        string  `json:"HoursType"`        // Corresponds to {.HoursType}
	AmbulanteDisplay string  `json:"AmbulanteDisplay"` // "block" or "none" for [AMBULANTE_DISPLAY]

	// Financial Terms
	Price          float64 `json:"Price"`          // Corresponds to {.Price}
	PriceTimeUnit  string  `json:"PriceTimeUnit"`  // Corresponds to {.PriceTimeUnit}
	Vat            float64 `json:"Vat"`            // Corresponds to {.Vat}
	TypeName       string  `json:"TypeName"`       // Corresponds to {.TypeName}
	GenerationDate string  `json:"GenerationDate"` // Date when the contract was generated

}

// GenerateIncidentPDF generates a PDF from incident data and returns the PDF bytes
func GenerateContractPDF(contractData ContractData) (multipart.File, error) {

	// Parse and execute HTML template
	templ, err := template.ParseFS(contractTemplateFS, "templates/contract.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := templ.Execute(&body, contractData); err != nil {
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
func UploadContractPDF(ctx context.Context, pdfFile multipart.File, contractID int64, b2Client *bucket.ObjectStorageClient) (string, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("contract/%s/contract-%d.pdf", timestamp, contractID)

	// Upload to B2
	key, _, err := b2Client.Upload(ctx, pdfFile, filename, "application/pdf")
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF to B2: %w", err)
	}
	return key, nil
}

// Helper function to do both operations if needed
func GenerateAndUploadContractPDF(ctx context.Context, contractData ContractData, b2Client *bucket.ObjectStorageClient) (string, error) {
	// Generate PDF
	pdfFile, err := GenerateContractPDF(contractData)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	fileURL, err := UploadContractPDF(ctx, pdfFile, contractData.ID, b2Client)
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF: %w", err)
	}

	return fileURL, nil
}
