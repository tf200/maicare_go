package invoice

import (
	"maicare_go/util"
	"time"

	"github.com/google/uuid"
)

// InvoiceDetails contains details for each contract in the invoice
type InvoiceDetails struct {
	ContractID    int64           `json:"contract_id"`
	ContractType  string          `json:"contract_name"`
	Periods       []InvoicePeriod `json:"periods"`
	PreVatTotal   float64         `json:"pre_vat_total_price"`
	Total         float64         `json:"total_price"`
	Vat           float64         `json:"vat"`
	Price         float64         `json:"price"`
	PriceTimeUnit string          `json:"price_time_unit"`
	Warnings      []string        `json:"warnings"`
}

type GenerateInvoiceRequest struct {
	ClientID  int64     `json:"client_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

// GenerateInvoiceResponse represents the response body for generating an invoice.
type GenerateInvoiceResponse struct {
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
