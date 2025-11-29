package invoice

import (
	"time"

	"github.com/google/uuid"
)

// CreatePaymentRequest represents the request body for creating a payment.
type CreatePaymentRequest struct {
	PaymentMethod    string    `json:"payment_method" binding:"oneof=credit_card bank_transfer cash check other"`
	PaymentStatus    string    `json:"payment_status" binding:"required,oneof=pending completed failed refunded reversed"`
	Amount           float64   `json:"amount" binding:"required,min=0"`
	PaymentDate      time.Time `json:"payment_date" binding:"required" example:"2023-10-01T00:00:00Z"`
	PaymentReference *string   `json:"payment_reference"`
	Notes            *string   `json:"notes"`
}

// CreatePaymentResponse represents the response body for creating a payment.
type CreatePaymentResponse struct {
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

// ListPaymentsResponse represents the response body for listing payments.
type ListPaymentsResponse struct {
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

// GetPaymentByIDResponse represents the response body for getting a payment by ID.
type GetPaymentByIDResponse struct {
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

// UpdatePaymentRequest represents the request body for updating a payment.
type UpdatePaymentRequest struct {
	PaymentMethod    *string    `json:"payment_method"`
	PaymentStatus    *string    `json:"payment_status"`
	Amount           *float64   `json:"amount"`
	PaymentDate      *time.Time `json:"payment_date"`
	PaymentReference *string    `json:"payment_reference"`
	Notes            *string    `json:"notes"`
}

// UpdatePaymentResponse represents the response body for updating a payment.
type UpdatePaymentResponse struct {
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

// DeletePaymentResponse represents the response body for deleting a payment.
type DeletePaymentResponse struct {
	DeletedPaymentID      uuid.UUID `json:"deleted_payment_id"`
	InvoiceID             uuid.UUID `json:"invoice_id"`
	DeletedAmount         float64   `json:"deleted_amount"`
	DeletedPaymentStatus  string    `json:"deleted_payment_status"`
	InvoiceStatusChanged  bool      `json:"invoice_status_changed"`
	CurrentInvoiceStatus  string    `json:"current_invoice_status"`
	PreviousInvoiceStatus string    `json:"previous_invoice_status"`
}
