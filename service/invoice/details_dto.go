package invoice

import (
	"maicare_go/pagination"
	"maicare_go/util"
	"time"

	"github.com/google/uuid"
)

// CreateInvoiceRequest represents the request body for creating an invoice.
type CreateInvoiceRequest struct {
	ClientID       int64            `json:"client_id" binding:"required"`
	InvoiceType    string           `json:"invoice_type" binding:"required,oneof=standard credit_note"`
	IssueDate      time.Time        `json:"issue_date" binding:"required"`
	DueDate        time.Time        `json:"due_date" binding:"required"`
	InvoiceDetails []InvoiceDetails `json:"invoice_details" binding:"required"`
	TotalAmount    float64          `json:"total_amount" binding:"required"`
	ExtraContent   util.JSONObject  `json:"extra_content" binding:"required"`
	Status         string           `json:"status" binding:"required,oneof=outstanding partially_paid paid expired overpaid imported concept"`
}

// CreateInvoiceResponse represents the response body for creating an invoice.
type CreateInvoiceResponse struct {
	ID              int64            `json:"id"`
	InvoiceNumber   string           `json:"invoice_number"`
	IssueDate       time.Time        `json:"issue_date"`
	DueDate         time.Time        `json:"due_date"`
	Status          string           `json:"status"`
	InvoiceDetails  []InvoiceDetails `json:"invoice_details"`
	TotalAmount     float64          `json:"total_amount"`
	PdfAttachmentID *uuid.UUID       `json:"pdf_attachment_id"`
	ExtraContent    util.JSONObject  `json:"extra_content"`
	ClientID        int64            `json:"client_id"`
	SenderID        *int64           `json:"sender_id"`
	InvoiceType     string           `json:"invoice_type"`
	UpdatedAt       time.Time        `json:"updated_at"`
	CreatedAt       time.Time        `json:"created_at"`
}

// GetInvoiceByIDResponse represents the response body for getting an invoice by ID.
type GetInvoiceByIDResponse struct {
	ID                   int64            `json:"id"`
	InvoiceNumber        string           `json:"invoice_number"`
	IssueDate            time.Time        `json:"issue_date"`
	DueDate              time.Time        `json:"due_date"`
	Status               string           `json:"status"`
	InvoiceDetails       []InvoiceDetails `json:"invoice_details"`
	TotalAmount          float64          `json:"total_amount"`
	PdfAttachmentID      *uuid.UUID       `json:"pdf_attachment_id"`
	ExtraContent         util.JSONObject  `json:"extra_content"`
	ClientID             int64            `json:"client_id"`
	SenderID             *int64           `json:"sender_id"`
	InvoiceType          string           `json:"invoice_type"`
	OriginalInvoiceID    *int64           `json:"original_invoice_id"`
	UpdatedAt            time.Time        `json:"updated_at"`
	CreatedAt            time.Time        `json:"created_at"`
	SenderName           *string          `json:"sender_name"`
	SenderKvknumber      *string          `json:"sender_kvknumber"`
	SenderBtwnumber      *string          `json:"sender_btwnumber"`
	ClientFirstName      string           `json:"client_first_name"`
	ClientLastName       string           `json:"client_last_name"`
	PaymentCompletionPrc float64          `json:"payment_completion_prc"`
}

// ListInvoicesRequest represents the request parameters for listing invoices.
type ListInvoicesRequest struct {
	ClientID  *int64    `form:"client_id"`
	SenderID  *int64    `form:"sender_id"`
	Status    *string   `form:"status" binding:"omitempty,oneof=outstanding partially_paid paid expired overpaid imported concept"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	pagination.Request
}

// ListInvoicesResponse represents the response body for listing invoices.
type ListInvoicesResponse struct {
	ID                int64            `json:"id"`
	InvoiceNumber     string           `json:"invoice_number"`
	IssueDate         time.Time        `json:"issue_date"`
	DueDate           time.Time        `json:"due_date"`
	Status            string           `json:"status"`
	InvoiceDetails    []InvoiceDetails `json:"invoice_details"`
	TotalAmount       float64          `json:"total_amount"`
	PdfAttachmentID   *uuid.UUID       `json:"pdf_attachment_id"`
	ExtraContent      util.JSONObject  `json:"extra_content"`
	ClientID          int64            `json:"client_id"`
	SenderID          *int64           `json:"sender_id"`
	InvoiceType       string           `json:"invoice_type"`
	OriginalInvoiceID *int64           `json:"original_invoice_id"`
	UpdatedAt         time.Time        `json:"updated_at"`
	CreatedAt         time.Time        `json:"created_at"`
	SenderName        *string          `json:"sender_name"`
	ClientFirstName   string           `json:"client_first_name"`
	ClientLastName    string           `json:"client_last_name"`
	WarningCount      int32            `json:"warning_count"`
}

// UpdateInvoiceRequest represents the request body for updating an invoice.
type UpdateInvoiceRequest struct {
	IssueDate      time.Time        `json:"issue_date"`
	DueDate        time.Time        `json:"due_date"`
	InvoiceDetails []InvoiceDetails `json:"invoice_details"`
	TotalAmount    float64          `json:"total_amount"`
	ExtraContent   util.JSONObject  `json:"extra_content"`
	Status         string           `json:"status"`
	WarningCount   int32            `json:"warning_count"`
}

// UpdateInvoiceResponse represents the response body for updating an invoice.
type UpdateInvoiceResponse struct {
	ID              int64            `json:"id"`
	InvoiceNumber   string           `json:"invoice_number"`
	IssueDate       time.Time        `json:"issue_date"`
	DueDate         time.Time        `json:"due_date"`
	Status          string           `json:"status"`
	InvoiceDetails  []InvoiceDetails `json:"invoice_details"`
	TotalAmount     float64          `json:"total_amount"`
	PdfAttachmentID *uuid.UUID       `json:"pdf_attachment_id"`
	ExtraContent    util.JSONObject  `json:"extra_content"`
	ClientID        int64            `json:"client_id"`
	SenderID        *int64           `json:"sender_id"`
	UpdatedAt       time.Time        `json:"updated_at"`
	CreatedAt       time.Time        `json:"created_at"`
}

// GetInvoiceAuditLogResponse represents the response body for getting invoice audit logs.
type GetInvoiceAuditLogResponse struct {
	AuditID            int64           `json:"audit_id"`
	InvoiceID          int64           `json:"invoice_id"`
	Operation          string          `json:"operation"`
	ChangedBy          *int64          `json:"changed_by"`
	ChangedAt          time.Time       `json:"changed_at"`
	OldValues          util.JSONObject `json:"old_values"`
	NewValues          util.JSONObject `json:"new_values"`
	ChangedFields      []string        `json:"changed_fields"`
	ChangedByFirstName *string         `json:"changed_by_first_name"`
	ChangedByLastName  *string         `json:"changed_by_last_name"`
}

// GetInvoiceTemplateItemsResponse represents the response body for getting invoice template items.
type GetInvoiceTemplateItemsResponse struct {
	ID           int64  `json:"id"`
	ItemTag      string `json:"item_tag"`
	Description  string `json:"description"`
	SourceTable  string `json:"source_table"`
	SourceColumn string `json:"source_column"`
}
