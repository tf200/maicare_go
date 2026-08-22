package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrPaymentNotFound = errors.New("payment not found")

// ==================== Payment Types ====================

type PaymentMethod string

const (
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodCheck        PaymentMethod = "check"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodOther        PaymentMethod = "other"
)

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusReversed  PaymentStatus = "reversed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// ==================== Domain Types ====================

type Invoice struct {
	ID                uuid.UUID
	InvoiceNumber     string
	IssueDate         time.Time
	DueDate           time.Time
	Status            string
	Source            string
	InvoiceType       string
	OriginalInvoiceID *uuid.UUID
	ReplacesInvoiceID *uuid.UUID
	PeriodStart       *time.Time
	PeriodEnd         *time.Time
	BillingTimezone   string
	Currency          string
	NetTotal          float64
	VatTotal          float64
	GrossTotal        float64
	PdfAttachmentID   *uuid.UUID
	ExtraContent      []byte
	ClientID          uuid.UUID
	SenderID          uuid.UUID
	WarningCount      int32
	LockedAt          *time.Time
	RunID             *uuid.UUID
	SenderName        *string
	SenderStreet      *string
	SenderHouseNumber *string
	SenderPostalCode  *string
	SenderCity        *string
	ClientFirstName   string
	ClientLastName    string
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

type InvoiceLine struct {
	ID          uuid.UUID
	LineNo      int32
	LineType    string
	ContractID  *uuid.UUID
	ServiceType string
	Description string
	PeriodStart time.Time
	PeriodEnd   time.Time
	Quantity    float64
	Unit        string
	UnitPrice   float64
	NetAmount   float64
	VatRate     float64
	VatAmount   float64
	GrossAmount float64
}

type InvoiceListItem struct {
	ID              uuid.UUID
	InvoiceNumber   string
	SenderName      string
	IsOverdue       bool
	ClientFirstName string
	ClientLastName  string
	ClientFilenumber string
	Currency        string
	GrossTotal      float64
	BalanceDue      float64
	PaidTotal       float64
	Status          string
	IssueDate       time.Time
	DueDate         time.Time
	ClientID        uuid.UUID
	SenderID        uuid.UUID
}

type Payment struct {
	ID                  uuid.UUID
	InvoiceID           uuid.UUID
	PaymentMethod       string
	PaymentStatus       string
	Amount              float64
	PaymentDate         time.Time
	PaymentReference    *string
	Notes               *string
	RecordedBy          *uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	RecordedByFirstName *string
	RecordedByLastName  *string
}

type InvoiceAuditLog struct {
	AuditID            uuid.UUID
	InvoiceID          uuid.UUID
	Operation          string
	ChangedBy          *uuid.UUID
	ChangedAt          time.Time
	OldValues          []byte
	NewValues          []byte
	ChangedFields      []string
	ChangedByFirstName *string
	ChangedByLastName  *string
}

// InvoiceTemplateItemRedeclared to avoid conflict with sender.go InvoiceTemplateItem
type InvoiceTemplateItemData struct {
	ID           uuid.UUID
	ItemTag      string
	Description  string
	SourceTable  string
	SourceColumn string
}

// ==================== Parameter Structs ====================

type CreateInvoiceLineInput struct {
	LineType    string
	ContractID  *uuid.UUID
	ServiceType string
	Description string
	PeriodStart time.Time
	PeriodEnd   time.Time
	Quantity    float64
	Unit        string
	UnitPrice   float64
	VatRate     float64
}

type ListInvoicesParams struct {
	ClientID        *uuid.UUID
	SenderID        *uuid.UUID
	Status          *string
	Statuses        []string
	Source          *string
	InvoiceType     *string
	RunID           *uuid.UUID
	StartDate       string
	EndDate         string
	PeriodStart     string
	PeriodEnd       string
	Locked          *bool
	MinWarningCount *int32
	Q               *string
	SortBy          string
	SortDir         string
	Limit           int32
	Offset          int32
}

type GenerateInvoiceParams struct {
	ClientID        uuid.UUID
	StartDate       time.Time
	EndDate         time.Time
	BillingTimezone string
	BillingCycle    string
}

type CreateInvoiceParams struct {
	ClientID    uuid.UUID
	InvoiceType string
	IssueDate   time.Time
	DueDate     time.Time
	Lines       []CreateInvoiceLineInput
	ExtraContent []byte
}

type CreatePaymentParams struct {
	PaymentMethod    string
	PaymentStatus    string
	Amount           float64
	PaymentDate      time.Time
	PaymentReference *string
	Notes            *string
}

type UpdatePaymentParams struct {
	PaymentMethod    *string
	PaymentStatus    *string
	Amount           *float64
	PaymentDate      *time.Time
	PaymentReference *string
	Notes            *string
}

// ==================== Service Interface ====================

type InvoiceService interface {
	CreateInvoice(ctx context.Context, params CreateInvoiceParams) (*Invoice, []InvoiceLine, error)
	GetInvoiceByID(ctx context.Context, invoiceID uuid.UUID) (*Invoice, []InvoiceLine, float64, error)
	ListInvoices(ctx context.Context, params ListInvoicesParams) (*ListResult[InvoiceListItem], error)
	UpdateInvoice(ctx context.Context, invoiceID uuid.UUID, params CreateInvoiceParams) (*Invoice, error)
	DeleteInvoice(ctx context.Context, invoiceID uuid.UUID) error

	GenerateInvoice(ctx context.Context, params GenerateInvoiceParams) (*GenerateInvoiceResult, int64, error)
	CreditInvoice(ctx context.Context, invoiceID uuid.UUID) (*CreditInvoiceResult, error)

	GetInvoiceAuditLogs(ctx context.Context, invoiceID uuid.UUID) ([]InvoiceAuditLog, error)
	GenerateInvoicePDF(ctx context.Context, invoiceID uuid.UUID) (*GeneratePDFResult, error)
	GetInvoiceTemplateItems(ctx context.Context) ([]InvoiceTemplateItemData, error)
	SendInvoiceReminder(ctx context.Context, invoiceID uuid.UUID) error

	CreatePayment(ctx context.Context, invoiceID uuid.UUID, params CreatePaymentParams) (*CreatePaymentResult, error)
	ListPayments(ctx context.Context, invoiceID uuid.UUID) ([]Payment, error)
	GetPaymentByID(ctx context.Context, invoiceID uuid.UUID, paymentID uuid.UUID) (*Payment, error)
	UpdatePayment(ctx context.Context, invoiceID uuid.UUID, paymentID uuid.UUID, params UpdatePaymentParams) (*UpdatePaymentResult, error)
	DeletePayment(ctx context.Context, invoiceID uuid.UUID, paymentID uuid.UUID) (*DeletePaymentResult, error)
}

// ==================== Result Structs ====================

type CreatePaymentResult struct {
	PaymentID            uuid.UUID
	InvoiceID            uuid.UUID
	PaymentMethod        string
	PaymentStatus        string
	Amount               float64
	PaymentDate          time.Time
	PaymentReference     *string
	Notes                *string
	RecordedBy           *uuid.UUID
	InvoiceStatusChanged bool
	CurrentInvoiceStatus string
}

type UpdatePaymentResult struct {
	PaymentID             uuid.UUID
	InvoiceID             uuid.UUID
	PaymentMethod         string
	PaymentStatus         string
	Amount                float64
	PaymentDate           time.Time
	PaymentReference      *string
	Notes                 *string
	RecordedBy            *uuid.UUID
	InvoiceStatusChanged  bool
	CurrentInvoiceStatus  string
	PreviousInvoiceStatus string
}

type DeletePaymentResult struct {
	DeletedPaymentID      uuid.UUID
	InvoiceID             uuid.UUID
	DeletedAmount         float64
	DeletedPaymentStatus  string
	InvoiceStatusChanged  bool
	CurrentInvoiceStatus  string
	PreviousInvoiceStatus string
}

type GeneratePDFResult struct {
	FileURL string
}

type CreditInvoiceResult struct {
	ID uuid.UUID
}

type GenerateInvoiceResult struct {
	Invoice
	Lines    []InvoiceLine
	Warnings []string
}
