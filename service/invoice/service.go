package invoice

import (
	"context"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
)

// InvoiceService Interface and implementation
//
//go:generate mockgen -source=service.go -destination=../mocks/mock_invoice_service.go -package=mocks
type InvoiceService interface {
	GenerateInvoice(req GenerateInvoiceRequest, ctx context.Context) (*GenerateInvoiceResponse, int64, error)
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest, employeeID int64) (*CreateInvoiceResponse, error)
	CreditInvoice(ctx context.Context, invoiceID int64, employeeID int64) (*CreditInvoiceResponse, error)
	GetInvoiceByID(ctx context.Context, invoiceID int64) (*GetInvoiceByIDResponse, error)
	ListInvoices(ctx *gin.Context, req ListInvoicesRequest) (*pagination.Response[ListInvoicesResponse], error)
	UpdateInvoice(ctx context.Context, invoiceID int64, payload UpdateInvoiceRequest, employeeID int64) (*UpdateInvoiceResponse, error)
	DeleteInvoice(ctx context.Context, invoiceID int64) error

	GetInvoiceTemplateItemsApi(ctx context.Context) ([]GetInvoiceTemplateItemsResponse, error)

	// Invoice Audit Log methods
	GetInvoiceAuditLogs(ctx context.Context, invoiceID int64) ([]GetInvoiceAuditLogResponse, error)

	// Invoice Reminder methods
	SendInvoiceReminder(ctx context.Context, invoiceID int64) error

	// Payment methods
	CreatePayment(ctx context.Context, invoiceID int64, req CreatePaymentRequest, employeeID int64) (*CreatePaymentResponse, error)
}

type invoiceService struct {
	*deps.ServiceDependencies
}

func NewInvoiceService(deps *deps.ServiceDependencies) InvoiceService {
	return &invoiceService{
		ServiceDependencies: deps,
	}
}
