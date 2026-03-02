package invoice

import (
	"time"

	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/google/uuid"
)

// InvoiceLineInput represents a single invoice line for manual creation.
type InvoiceLineInput struct {
	LineType    string     `json:"line_type" binding:"required,oneof=contract manual adjustment"`
	ContractID  *uuid.UUID `json:"contract_id"`
	ServiceType string     `json:"service_type" binding:"required"`
	Description string     `json:"description" binding:"required"`
	PeriodStart time.Time  `json:"period_start" binding:"required"`
	PeriodEnd   time.Time  `json:"period_end" binding:"required"`

	Quantity  float64 `json:"quantity" binding:"required,min=0"`
	Unit      string  `json:"unit" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`

	VatRate float64 `json:"vat_rate" binding:"min=0,max=100"`
}

// CreateInvoiceRequest represents the request body for creating an invoice.
type CreateInvoiceRequest struct {
	ClientID        uuid.UUID          `json:"client_id" binding:"required"`
	SenderID        uuid.UUID          `json:"sender_id" binding:"required"`
	InvoiceType     string             `json:"invoice_type" binding:"required,oneof=standard credit_note"`
	IssueDate       time.Time          `json:"issue_date" binding:"required"`
	DueDate         time.Time          `json:"due_date" binding:"required"`
	Currency        string             `json:"currency"`
	Lines           []InvoiceLineInput `json:"lines" binding:"required"`
	ExtraContent    util.JSONObject    `json:"extra_content"`
	Status          string             `json:"status" binding:"required,oneof=outstanding partially_paid paid expired overpaid imported concept canceled"`
	BillingTimezone string             `json:"billing_timezone"`
}

// CreateInvoiceResponse represents the response body for creating an invoice.
type CreateInvoiceResponse struct {
	ID              uuid.UUID       `json:"id"`
	InvoiceNumber   string          `json:"invoice_number"`
	IssueDate       time.Time       `json:"issue_date"`
	DueDate         time.Time       `json:"due_date"`
	Status          string          `json:"status"`
	Source          string          `json:"source"`
	InvoiceType     string          `json:"invoice_type"`
	Currency        string          `json:"currency"`
	NetTotal        float64         `json:"net_total_amount"`
	VatTotal        float64         `json:"vat_total_amount"`
	GrossTotal      float64         `json:"gross_total_amount"`
	PdfAttachmentID *uuid.UUID      `json:"pdf_attachment_id"`
	ExtraContent    util.JSONObject `json:"extra_content"`
	ClientID        uuid.UUID       `json:"client_id"`
	SenderID        uuid.UUID       `json:"sender_id"`
	Lines           []InvoiceLine   `json:"lines"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CreatedAt       time.Time       `json:"created_at"`
}

// GetInvoiceByIDResponse represents the response body for getting an invoice by ID.
type GetInvoiceByIDResponse struct {
	ID                   uuid.UUID       `json:"id"`
	InvoiceNumber        string          `json:"invoice_number"`
	IssueDate            time.Time       `json:"issue_date"`
	DueDate              time.Time       `json:"due_date"`
	Status               string          `json:"status"`
	Source               string          `json:"source"`
	InvoiceType          string          `json:"invoice_type"`
	OriginalInvoiceID    *uuid.UUID      `json:"original_invoice_id"`
	ReplacesInvoiceID    *uuid.UUID      `json:"replaces_invoice_id"`
	PeriodStart          *time.Time      `json:"period_start"`
	PeriodEnd            *time.Time      `json:"period_end"`
	BillingTimezone      string          `json:"billing_timezone"`
	Currency             string          `json:"currency"`
	NetTotal             float64         `json:"net_total_amount"`
	VatTotal             float64         `json:"vat_total_amount"`
	GrossTotal           float64         `json:"gross_total_amount"`
	ExtraContent         util.JSONObject `json:"extra_content"`
	ClientID             uuid.UUID       `json:"client_id"`
	SenderID             uuid.UUID       `json:"sender_id"`
	Lines                []InvoiceLine   `json:"lines"`
	UpdatedAt            time.Time       `json:"updated_at"`
	CreatedAt            time.Time       `json:"created_at"`
	SenderName           *string         `json:"sender_name"`
	SenderKvknumber      *string         `json:"sender_kvknumber"`
	SenderBtwnumber      *string         `json:"sender_btwnumber"`
	ClientFirstName      string          `json:"client_first_name"`
	ClientLastName       string          `json:"client_last_name"`
	PaymentCompletionPrc float64         `json:"payment_completion_prc"`
}

// ListInvoicesRequest represents the request parameters for listing invoices.
type ListInvoicesRequest struct {
	ClientID *uuid.UUID `form:"client_id"`
	SenderID *uuid.UUID `form:"sender_id"`

	// Either a single status (legacy) or multiple statuses.
	Status   *string  `form:"status" binding:"omitempty,oneof=outstanding partially_paid paid expired overpaid imported concept canceled"`
	Statuses []string `form:"statuses" binding:"omitempty,dive,oneof=outstanding partially_paid paid expired overpaid imported concept canceled"`

	Source      *string    `form:"source" binding:"omitempty,oneof=auto manual imported"`
	InvoiceType *string    `form:"invoice_type" binding:"omitempty,oneof=standard credit_note"`
	RunID       *uuid.UUID `form:"run_id"`

	// Filter by issue date range (YYYY-MM-DD). These are optional.
	StartDate string `form:"start_date" example:"2026-02-01"`
	EndDate   string `form:"end_date" example:"2026-02-29"`

	// Filter by billing period range (RFC3339 or YYYY-MM-DD). These are optional.
	PeriodStart string `form:"period_start" example:"2026-02-01T00:00:00Z"`
	PeriodEnd   string `form:"period_end" example:"2026-02-29T00:00:00Z"`

	Locked          *bool   `form:"locked"`
	MinWarningCount *int32  `form:"min_warning_count"`
	Q               *string `form:"q"`

	SortBy  string `form:"sort_by" binding:"omitempty,oneof=updated_at issue_date due_date invoice_number gross_total_amount balance_due_amount"`
	SortDir string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
	pagination.Request
}

// ListInvoicesResponse represents the response body for listing invoices.
type ListInvoicesResponse struct {
	ID            uuid.UUID `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	SenderName    string    `json:"sender_name"`
	IsOverdue     bool      `json:"is_overdue"`

	ClientFirstName  string `json:"client_first_name"`
	ClientLastName   string `json:"client_last_name"`
	ClientFilenumber string `json:"client_filenumber"`

	Currency   string  `json:"currency"`
	GrossTotal float64 `json:"gross_total_amount"`
	BalanceDue float64 `json:"balance_due_amount"`
	PaidTotal  float64 `json:"paid_total_amount"`

	Status    string    `json:"status"`
	IssueDate time.Time `json:"issue_date"`
	DueDate   time.Time `json:"due_date"`

	ClientID uuid.UUID `json:"client_id"`
	SenderID uuid.UUID `json:"sender_id"`
}

// UpdateInvoiceRequest represents the request body for updating an invoice.
type UpdateInvoiceRequest struct {
	IssueDate    time.Time       `json:"issue_date"`
	DueDate      time.Time       `json:"due_date"`
	ExtraContent util.JSONObject `json:"extra_content"`
	Status       string          `json:"status"`
	WarningCount int32           `json:"warning_count"`
	LockedAt     *time.Time      `json:"locked_at"`
}

// UpdateInvoiceResponse represents the response body for updating an invoice.
type UpdateInvoiceResponse struct {
	ID              uuid.UUID       `json:"id"`
	InvoiceNumber   string          `json:"invoice_number"`
	IssueDate       time.Time       `json:"issue_date"`
	DueDate         time.Time       `json:"due_date"`
	Status          string          `json:"status"`
	Currency        string          `json:"currency"`
	NetTotal        float64         `json:"net_total_amount"`
	VatTotal        float64         `json:"vat_total_amount"`
	GrossTotal      float64         `json:"gross_total_amount"`
	PdfAttachmentID *uuid.UUID      `json:"pdf_attachment_id"`
	ExtraContent    util.JSONObject `json:"extra_content"`
	ClientID        uuid.UUID       `json:"client_id"`
	SenderID        uuid.UUID       `json:"sender_id"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CreatedAt       time.Time       `json:"created_at"`
}

// GetInvoiceAuditLogResponse represents the response body for getting invoice audit logs.
type GetInvoiceAuditLogResponse struct {
	AuditID            uuid.UUID       `json:"audit_id"`
	InvoiceID          uuid.UUID       `json:"invoice_id"`
	Operation          string          `json:"operation"`
	ChangedBy          *uuid.UUID      `json:"changed_by"`
	ChangedAt          time.Time       `json:"changed_at"`
	OldValues          util.JSONObject `json:"old_values"`
	NewValues          util.JSONObject `json:"new_values"`
	ChangedFields      []string        `json:"changed_fields"`
	ChangedByFirstName *string         `json:"changed_by_first_name"`
	ChangedByLastName  *string         `json:"changed_by_last_name"`
}

// GetInvoiceTemplateItemsResponse represents the response body for getting invoice template items.
type GetInvoiceTemplateItemsResponse struct {
	ID           uuid.UUID `json:"id"`
	ItemTag      string    `json:"item_tag"`
	Description  string    `json:"description"`
	SourceTable  string    `json:"source_table"`
	SourceColumn string    `json:"source_column"`
}

// GenerateInvoicePDFResponse represents the response body for generating an invoice PDF.
type GenerateInvoicePDFResponse struct {
	FileUrl string `json:"file_url"`
}

// Contact represents a contact information.
type SenderContact struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"email"`
	PhoneNumber *string `json:"phone_number"`
}
