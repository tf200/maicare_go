package invoice

import (
	"context"

	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InvoiceService Interface and implementation
//
//go:generate mockgen -source=service.go -destination=../mocks/mock_invoice_service.go -package=mocks
type InvoiceService interface {
	GenerateInvoice(req GenerateInvoiceRequest, ctx context.Context) (*GenerateInvoiceResponse, int64, error)
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest, employeeID uuid.UUID) (*CreateInvoiceResponse, error)
	CreditInvoice(ctx context.Context, invoiceID uuid.UUID, employeeID uuid.UUID) (*CreditInvoiceResponse, error)
	GetInvoiceByID(ctx context.Context, invoiceID uuid.UUID) (*GetInvoiceByIDResponse, error)
	ListInvoices(ctx *gin.Context, req ListInvoicesRequest) (*pagination.Response[ListInvoicesResponse], error)
	UpdateInvoice(ctx context.Context, invoiceID uuid.UUID, payload UpdateInvoiceRequest, employeeID uuid.UUID) (*UpdateInvoiceResponse, error)
	DeleteInvoice(ctx context.Context, invoiceID uuid.UUID) error

	GenerateInvoicePdf(ctx context.Context, invoiceID uuid.UUID) (*GenerateInvoicePDFResponse, error)

	GetInvoiceTemplateItemsApi(ctx context.Context) ([]GetInvoiceTemplateItemsResponse, error)

	// Invoice Audit Log methods
	GetInvoiceAuditLogs(ctx context.Context, invoiceID uuid.UUID) ([]GetInvoiceAuditLogResponse, error)

	// Invoice Reminder methods
	SendInvoiceReminder(ctx context.Context, invoiceID uuid.UUID) error

	// Payment methods
	CreatePayment(ctx context.Context, invoiceID uuid.UUID, req CreatePaymentRequest, employeeID uuid.UUID) (*CreatePaymentResponse, error)
	ListPayments(ctx context.Context, invoiceID uuid.UUID) ([]ListPaymentsResponse, error)
	GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (*GetPaymentByIDResponse, error)
	UpdatePayment(ctx context.Context, invoiceID uuid.UUID, employeeID uuid.UUID, paymentID uuid.UUID, req UpdatePaymentRequest) (*UpdatePaymentResponse, error)
	DeletePayment(ctx context.Context, invoiceID uuid.UUID, paymentID uuid.UUID, employeeID uuid.UUID) (*DeletePaymentResponse, error)
}

type invoiceService struct {
	*deps.ServiceDependencies
}

func NewInvoiceService(deps *deps.ServiceDependencies) InvoiceService {
	return &invoiceService{
		ServiceDependencies: deps,
	}
}
