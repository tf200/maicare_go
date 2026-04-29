package handler

import (
	"time"

	"maicare_go/internal/httpapi"

	"github.com/google/uuid"
)

// ==================== Invoice DTOs ====================

type InvoiceLineDTO struct {
	ID          uuid.UUID  `json:"id"`
	LineNo      int32      `json:"line_no"`
	LineType    string     `json:"line_type"`
	ContractID  *uuid.UUID `json:"contract_id"`
	ServiceType string     `json:"service_type"`
	Description string     `json:"description"`
	PeriodStart time.Time  `json:"period_start"`
	PeriodEnd   time.Time  `json:"period_end"`
	Quantity    float64    `json:"quantity"`
	Unit        string     `json:"unit"`
	UnitPrice   float64    `json:"unit_price"`
	NetAmount   float64    `json:"net_amount"`
	VatRate     float64    `json:"vat_rate"`
	VatAmount   float64    `json:"vat_amount"`
	GrossAmount float64    `json:"gross_amount"`
}

type CreateInvoiceRequest struct {
	ClientID    uuid.UUID                `json:"client_id" binding:"required"`
	InvoiceType string                   `json:"invoice_type" binding:"required,oneof=standard credit_note"`
	IssueDate   time.Time                `json:"issue_date" binding:"required"`
	DueDate     time.Time                `json:"due_date" binding:"required"`
	Lines       []InvoiceLineInputDTO    `json:"lines" binding:"required,min=1,dive"`
}

type InvoiceLineInputDTO struct {
	LineType    string     `json:"line_type" binding:"required,oneof=contract manual adjustment"`
	ContractID  *uuid.UUID `json:"contract_id"`
	ServiceType string     `json:"service_type" binding:"required"`
	Description string     `json:"description" binding:"required"`
	PeriodStart time.Time  `json:"period_start" binding:"required"`
	PeriodEnd   time.Time  `json:"period_end" binding:"required"`
	Quantity    float64    `json:"quantity" binding:"required,min=0"`
	Unit        string     `json:"unit" binding:"required"`
	UnitPrice   float64    `json:"unit_price" binding:"required,min=0"`
	VatRate     float64    `json:"vat_rate" binding:"min=0,max=100"`
}

type CreateInvoiceResponseDTO struct {
	ID              uuid.UUID          `json:"id"`
	InvoiceNumber   string             `json:"invoice_number"`
	IssueDate       time.Time          `json:"issue_date"`
	DueDate         time.Time          `json:"due_date"`
	Status          string             `json:"status"`
	Source          string             `json:"source"`
	InvoiceType     string             `json:"invoice_type"`
	Currency        string             `json:"currency"`
	NetTotal        float64            `json:"net_total_amount"`
	VatTotal        float64            `json:"vat_total_amount"`
	GrossTotal      float64            `json:"gross_total_amount"`
	PdfAttachmentID *uuid.UUID         `json:"pdf_attachment_id"`
	ClientID        uuid.UUID          `json:"client_id"`
	SenderID        uuid.UUID          `json:"sender_id"`
	Lines           []InvoiceLineDTO   `json:"lines"`
	UpdatedAt       time.Time          `json:"updated_at"`
	CreatedAt       time.Time          `json:"created_at"`
}

type GetInvoiceByIDResponseDTO struct {
	ID                   uuid.UUID        `json:"id"`
	InvoiceNumber        string           `json:"invoice_number"`
	IssueDate            time.Time        `json:"issue_date"`
	DueDate              time.Time        `json:"due_date"`
	Status               string           `json:"status"`
	Source               string           `json:"source"`
	InvoiceType          string           `json:"invoice_type"`
	OriginalInvoiceID    *uuid.UUID       `json:"original_invoice_id"`
	ReplacesInvoiceID    *uuid.UUID       `json:"replaces_invoice_id"`
	PeriodStart          *time.Time       `json:"period_start"`
	PeriodEnd            *time.Time       `json:"period_end"`
	BillingTimezone      string           `json:"billing_timezone"`
	Currency             string           `json:"currency"`
	NetTotal             float64          `json:"net_total_amount"`
	VatTotal             float64          `json:"vat_total_amount"`
	GrossTotal           float64          `json:"gross_total_amount"`
	ClientID             uuid.UUID        `json:"client_id"`
	SenderID             uuid.UUID        `json:"sender_id"`
	Lines                []InvoiceLineDTO `json:"lines"`
	UpdatedAt            time.Time        `json:"updated_at"`
	CreatedAt            time.Time        `json:"created_at"`
	SenderName           *string          `json:"sender_name"`
	SenderKvknumber      *string          `json:"sender_kvknumber"`
	SenderBtwnumber      *string          `json:"sender_btwnumber"`
	ClientFirstName      string           `json:"client_first_name"`
	ClientLastName       string           `json:"client_last_name"`
	PaymentCompletionPrc float64          `json:"payment_completion_prc"`
}

type ListInvoicesRequest struct {
	httpapi.PageRequest
	ClientID        *uuid.UUID `form:"client_id"`
	SenderID        *uuid.UUID `form:"sender_id"`
	Status          *string    `form:"status" binding:"omitempty,oneof=outstanding partially_paid paid expired overpaid imported concept canceled"`
	Statuses        []string   `form:"statuses" binding:"omitempty,dive,oneof=outstanding partially_paid paid expired overpaid imported concept canceled"`
	Source          *string    `form:"source" binding:"omitempty,oneof=auto manual imported"`
	InvoiceType     *string    `form:"invoice_type" binding:"omitempty,oneof=standard credit_note"`
	RunID           *uuid.UUID `form:"run_id"`
	StartDate       string     `form:"start_date" example:"2026-02-01"`
	EndDate         string     `form:"end_date" example:"2026-02-29"`
	PeriodStart     string     `form:"period_start" example:"2026-02-01T00:00:00Z"`
	PeriodEnd       string     `form:"period_end" example:"2026-02-29T00:00:00Z"`
	Locked          *bool      `form:"locked"`
	MinWarningCount *int32     `form:"min_warning_count"`
	Q               *string    `form:"q"`
	SortBy          string     `form:"sort_by" binding:"omitempty,oneof=updated_at issue_date due_date invoice_number gross_total_amount balance_due_amount"`
	SortDir         string     `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

type ListInvoicesItemDTO struct {
	ID              uuid.UUID `json:"id"`
	InvoiceNumber   string    `json:"invoice_number"`
	SenderName      string    `json:"sender_name"`
	IsOverdue       bool      `json:"is_overdue"`
	ClientFirstName string    `json:"client_first_name"`
	ClientLastName  string    `json:"client_last_name"`
	ClientFilenumber string   `json:"client_filenumber"`
	Currency        string    `json:"currency"`
	GrossTotal      float64   `json:"gross_total_amount"`
	BalanceDue      float64   `json:"balance_due_amount"`
	PaidTotal       float64   `json:"paid_total_amount"`
	Status          string    `json:"status"`
	IssueDate       time.Time `json:"issue_date"`
	DueDate         time.Time `json:"due_date"`
	ClientID        uuid.UUID `json:"client_id"`
	SenderID        uuid.UUID `json:"sender_id"`
}

type UpdateInvoiceRequest struct {
	IssueDate time.Time             `json:"issue_date"`
	DueDate   time.Time             `json:"due_date"`
	Status    string                `json:"status"`
	Lines     []InvoiceLineInputDTO `json:"lines" binding:"omitempty,dive"`
}

type UpdateInvoiceResponseDTO struct {
	ID              uuid.UUID  `json:"id"`
	InvoiceNumber   string     `json:"invoice_number"`
	IssueDate       time.Time  `json:"issue_date"`
	DueDate         time.Time  `json:"due_date"`
	Status          string     `json:"status"`
	Currency        string     `json:"currency"`
	NetTotal        float64    `json:"net_total_amount"`
	VatTotal        float64    `json:"vat_total_amount"`
	GrossTotal      float64    `json:"gross_total_amount"`
	PdfAttachmentID *uuid.UUID `json:"pdf_attachment_id"`
	ClientID        uuid.UUID  `json:"client_id"`
	SenderID        uuid.UUID  `json:"sender_id"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type GenerateInvoiceRequest struct {
	ClientID    uuid.UUID `json:"client_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	BillingTimezone string `json:"billing_timezone"`
	BillingCycle    string `json:"billing_cycle"`
}

type GenerateInvoiceResponseDTO struct {
	ID              uuid.UUID          `json:"id"`
	InvoiceNumber   string             `json:"invoice_number"`
	IssueDate       time.Time          `json:"issue_date"`
	DueDate         time.Time          `json:"due_date"`
	Status          string             `json:"status"`
	Source          string             `json:"source"`
	InvoiceType     string             `json:"invoice_type"`
	PeriodStart     *time.Time         `json:"period_start"`
	PeriodEnd       *time.Time         `json:"period_end"`
	Currency        string             `json:"currency"`
	NetTotal        float64            `json:"net_total_amount"`
	VatTotal        float64            `json:"vat_total_amount"`
	GrossTotal      float64            `json:"gross_total_amount"`
	PdfAttachmentID *uuid.UUID         `json:"pdf_attachment_id"`
	ClientID        uuid.UUID          `json:"client_id"`
	SenderID        uuid.UUID          `json:"sender_id"`
	Lines           []InvoiceLineDTO   `json:"lines"`
	Warnings        []string           `json:"warnings"`
	UpdatedAt       time.Time          `json:"updated_at"`
	CreatedAt       time.Time          `json:"created_at"`
}

type CreditInvoiceResponseDTO struct {
	ID uuid.UUID `json:"id"`
}

type InvoiceAuditLogDTO struct {
	AuditID            uuid.UUID `json:"audit_id"`
	InvoiceID          uuid.UUID `json:"invoice_id"`
	Operation          string    `json:"operation"`
	ChangedBy          *uuid.UUID `json:"changed_by"`
	ChangedAt          time.Time  `json:"changed_at"`
	ChangedFields      []string   `json:"changed_fields"`
	ChangedByFirstName *string    `json:"changed_by_first_name"`
	ChangedByLastName  *string    `json:"changed_by_last_name"`
}

type InvoiceTemplateItemDTO struct {
	ID           uuid.UUID `json:"id"`
	ItemTag      string    `json:"item_tag"`
	Description  string    `json:"description"`
	SourceTable  string    `json:"source_table"`
	SourceColumn string    `json:"source_column"`
}

type GeneratePDFResponseDTO struct {
	FileURL string `json:"file_url"`
}

// ==================== Payment DTOs ====================

type CreatePaymentRequest struct {
	PaymentMethod    string    `json:"payment_method" binding:"oneof=credit_card bank_transfer cash check other"`
	PaymentStatus    string    `json:"payment_status" binding:"required,oneof=pending completed failed refunded reversed"`
	Amount           float64   `json:"amount" binding:"required,min=0"`
	PaymentDate      time.Time `json:"payment_date" binding:"required" example:"2023-10-01T00:00:00Z"`
	PaymentReference *string   `json:"payment_reference"`
	Notes            *string   `json:"notes"`
}

type CreatePaymentResponseDTO struct {
	PaymentID            uuid.UUID  `json:"payment_id"`
	InvoiceID            uuid.UUID  `json:"invoice_id"`
	PaymentMethod        string     `json:"payment_method"`
	PaymentStatus        string     `json:"payment_status"`
	Amount               float64    `json:"amount"`
	PaymentDate          time.Time  `json:"payment_date"`
	PaymentReference     *string    `json:"payment_reference"`
	Notes                *string    `json:"notes"`
	InvoiceStatusChanged bool       `json:"invoice_status_changed"`
	CurrentInvoiceStatus string     `json:"current_invoice_status"`
	RecordedBy           *uuid.UUID `json:"recorded_by"`
}

type ListPaymentsItemDTO struct {
	PaymentID           uuid.UUID  `json:"payment_id"`
	InvoiceID           uuid.UUID  `json:"invoice_id"`
	PaymentMethod       string     `json:"payment_method"`
	PaymentStatus       string     `json:"payment_status"`
	Amount              float64    `json:"amount"`
	PaymentDate         time.Time  `json:"payment_date"`
	PaymentReference    *string    `json:"payment_reference"`
	Notes               *string    `json:"notes"`
	RecordedBy          *uuid.UUID `json:"recorded_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	RecordedByFirstName *string    `json:"recorded_by_first_name"`
	RecordedByLastName  *string    `json:"recorded_by_last_name"`
}

type GetPaymentByIDResponseDTO struct {
	PaymentID           uuid.UUID  `json:"payment_id"`
	InvoiceID           uuid.UUID  `json:"invoice_id"`
	PaymentMethod       string     `json:"payment_method"`
	PaymentStatus       string     `json:"payment_status"`
	Amount              float64    `json:"amount"`
	PaymentDate         time.Time  `json:"payment_date"`
	PaymentReference    *string    `json:"payment_reference"`
	Notes               *string    `json:"notes"`
	RecordedBy          *uuid.UUID `json:"recorded_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	RecordedByFirstName *string    `json:"recorded_by_first_name"`
	RecordedByLastName  *string    `json:"recorded_by_last_name"`
}

type UpdatePaymentRequest struct {
	PaymentMethod    *string    `json:"payment_method"`
	PaymentStatus    *string    `json:"payment_status"`
	Amount           *float64   `json:"amount"`
	PaymentDate      *time.Time `json:"payment_date"`
	PaymentReference *string    `json:"payment_reference"`
	Notes            *string    `json:"notes"`
}

type UpdatePaymentResponseDTO struct {
	PaymentID             uuid.UUID  `json:"payment_id"`
	InvoiceID             uuid.UUID  `json:"invoice_id"`
	PaymentMethod         string     `json:"payment_method"`
	PaymentStatus         string     `json:"payment_status"`
	Amount                float64    `json:"amount"`
	PaymentDate           time.Time  `json:"payment_date"`
	PaymentReference      *string    `json:"payment_reference"`
	Notes                 *string    `json:"notes"`
	RecordedBy            *uuid.UUID `json:"recorded_by"`
	InvoiceStatusChanged  bool       `json:"invoice_status_changed"`
	CurrentInvoiceStatus  string     `json:"current_invoice_status"`
	PreviousInvoiceStatus string     `json:"previous_invoice_status"`
}

type DeletePaymentResponseDTO struct {
	DeletedPaymentID      uuid.UUID `json:"deleted_payment_id"`
	InvoiceID             uuid.UUID `json:"invoice_id"`
	DeletedAmount         float64   `json:"deleted_amount"`
	DeletedPaymentStatus  string    `json:"deleted_payment_status"`
	InvoiceStatusChanged  bool      `json:"invoice_status_changed"`
	CurrentInvoiceStatus  string    `json:"current_invoice_status"`
	PreviousInvoiceStatus string    `json:"previous_invoice_status"`
}
