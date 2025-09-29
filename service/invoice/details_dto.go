package invoice

import (
	"maicare_go/util"
	"time"

	"github.com/google/uuid"
)

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
