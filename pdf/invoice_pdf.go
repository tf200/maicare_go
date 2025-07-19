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

type InvoicePDFData struct {
	ID                   int64
	SenderName           string
	SenderContactPerson  string
	SenderAddressLine1   string
	SenderPostalCodeCity string
	InvoiceNumber        string
	InvoiceDate          time.Time
	DueDate              time.Time
	InvoiceDetails       []InvoiceDetail
	TotalAmount          float64
	ExtraItems           []struct {
		Label string
		Value string
	}
}

type InvoiceDetail struct {
	CareType      string
	Periods       []InvoicePeriod
	Price         float64
	PriceTimeUnit string // e.g. "hour", "day", etc.
	PreVatTotal   float64
	Total         float64
}

type InvoicePeriod struct {
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	AcommodationTimeFrame string    `json:"accommodation_time_frame,omitempty"`
	AmbulanteTotalMinutes float64   `json:"ambulante_total_minutes,omitempty"`
}

func sumPreVat(details []InvoiceDetail) float64 {
	var total float64
	for _, d := range details {
		total += d.PreVatTotal
	}
	return total
}

func sumVat(details []InvoiceDetail) float64 {
	var total float64
	for _, d := range details {
		total += d.Total - d.PreVatTotal
	}
	return total
}

//go:embed templates/invoice.html
var invoiceTemplateFS embed.FS

func GenerateInvoicePDF(invoiceData InvoicePDFData) ([]byte, error) {
	funcMap := template.FuncMap{
		"sumPreVat": sumPreVat,
		"sumVat":    sumVat,
	}

	// Parse and execute HTML template
	templ, err := template.New("invoice.html").Funcs(funcMap).ParseFS(invoiceTemplateFS, "templates/invoice.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := templ.Execute(&body, invoiceData); err != nil {
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

// UploadInvoicePDF uploads a PDF to B2 with a generated filename
func UploadInvoicePDF(ctx context.Context, pdfBytes []byte, invoiceID int64, b2Client *bucket.B2Client) (string, string, int64, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("invoice_reports/%s/invoice_report_%d.pdf", timestamp, invoiceID)

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
		return "", "", 0, fmt.Errorf("failed to copy PDF to B2: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", "", 0, fmt.Errorf("failed to close B2 writer: %w", err)
	}

	// Get file size from pdfBytes
	fileSize := int64(len(pdfBytes))

	fileURL := fmt.Sprintf("%s/file/%s/%s",
		b2Client.Bucket.BaseURL(),
		b2Client.Bucket.Name(),
		filename)

	return fileURL, filename, fileSize, nil
}

// Helper function to do both operations if needed
func GenerateAndUploadInvoicePDF(ctx context.Context, invoiceData InvoicePDFData, b2Client *bucket.B2Client) (string, string, int64, error) {
	// Generate PDF
	pdfBytes, err := GenerateInvoicePDF(invoiceData)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	fileURL, filename, fileSize, err := UploadInvoicePDF(ctx, pdfBytes, invoiceData.ID, b2Client)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to upload PDF: %w", err)
	}

	return fileURL, filename, fileSize, nil
}
