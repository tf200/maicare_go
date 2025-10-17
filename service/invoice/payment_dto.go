package invoice

import "time"

// CreatePaymentRequest represents the request body for creating a payment.
type CreatePaymentRequest struct {
	PaymentMethod    *string   `json:"payment_method" binding:"oneof=credit_card bank_transfer cash check other"`
	PaymentStatus    string    `json:"payment_status" binding:"required,oneof=pending completed failed refunded reversed"`
	Amount           float64   `json:"amount" binding:"required,min=0"`
	PaymentDate      time.Time `json:"payment_date" binding:"required" example:"2023-10-01T00:00:00Z"`
	PaymentReference *string   `json:"payment_reference"`
	Notes            *string   `json:"notes"`
}

// CreatePaymentResponse represents the response body for creating a payment.
type CreatePaymentResponse struct {
	PaymentID            int64     `json:"payment_id"`
	InvoiceID            int64     `json:"invoice_id"`
	PaymentMethod        *string   `json:"payment_method"`
	PaymentStatus        string    `json:"payment_status"`
	Amount               float64   `json:"amount"`
	PaymentDate          time.Time `json:"payment_date"`
	PaymentReference     *string   `json:"payment_reference"`
	Notes                *string   `json:"notes"`
	InvoiceStatusChanged bool      `json:"invoice_status_changed"`
	CurrentInvoiceStatus string    `json:"current_invoice_status"`
	RecordedBy           *int64    `json:"recorded_by"`
}
