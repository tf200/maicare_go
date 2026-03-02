package invoice

import (
	"time"

	"maicare_go/util"

	"github.com/google/uuid"
)

type GenerateInvoiceRequest struct {
	ClientID         uuid.UUID `json:"client_id" binding:"required"`
	StartDate        time.Time `json:"start_date" binding:"required"`
	EndDate          time.Time `json:"end_date" binding:"required"`
	BillingTimezone  string    `json:"billing_timezone"`
	BillingCycle     string    `json:"billing_cycle"` // e.g. "iso_4_week"
	ForceNewDocument bool      `json:"force_new_document"`
}

// InvoiceLine is the canonical invoice breakdown used for calculation and storage.
type InvoiceLine struct {
	ID          uuid.UUID  `json:"id"`
	LineNo      int32      `json:"line_no"`
	LineType    string     `json:"line_type"`
	ContractID  *uuid.UUID `json:"contract_id"`
	ServiceType string     `json:"service_type"`
	Description string     `json:"description"`
	PeriodStart time.Time  `json:"period_start"`
	PeriodEnd   time.Time  `json:"period_end"`

	Quantity  float64 `json:"quantity"`
	Unit      string  `json:"unit"`
	UnitPrice float64 `json:"unit_price"`

	NetAmount   float64 `json:"net_amount"`
	VatRate     float64 `json:"vat_rate"`
	VatAmount   float64 `json:"vat_amount"`
	GrossAmount float64 `json:"gross_amount"`
}

// GenerateInvoiceResponse represents the response body for generating an invoice.
type GenerateInvoiceResponse struct {
	ID            uuid.UUID       `json:"id"`
	InvoiceNumber string          `json:"invoice_number"`
	IssueDate     time.Time       `json:"issue_date"`
	DueDate       time.Time       `json:"due_date"`
	Status        string          `json:"status"`
	Source        string          `json:"source"`
	InvoiceType   string          `json:"invoice_type"`
	PeriodStart   *time.Time      `json:"period_start"`
	PeriodEnd     *time.Time      `json:"period_end"`
	Currency      string          `json:"currency"`
	NetTotal      float64         `json:"net_total_amount"`
	VatTotal      float64         `json:"vat_total_amount"`
	GrossTotal    float64         `json:"gross_total_amount"`
	PdfAttachmentID *uuid.UUID    `json:"pdf_attachment_id"`
	ExtraContent  util.JSONObject `json:"extra_content"`
	ClientID      uuid.UUID       `json:"client_id"`
	SenderID      uuid.UUID       `json:"sender_id"`
	Lines         []InvoiceLine   `json:"lines"`
	Warnings      []string        `json:"warnings"`
	UpdatedAt     time.Time       `json:"updated_at"`
	CreatedAt     time.Time       `json:"created_at"`
}
