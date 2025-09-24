package invoice

import (
	"context"
	"maicare_go/service/deps"
)

// InvoiceService Interface and implementation
//
//go:generate mockgen -source=service.go -destination=../mocks/mock_invoice_service.go -package=mocks
type InvoiceService interface {
	GenerateInvoice(req GenerateInvoiceRequest, ctx context.Context) (*GenerateInvoiceResult, int64, error)
	GetInvoiceByID(ctx context.Context, invoiceID int64) (*GetInvoiceByIDResponse, error)
	SendInvoiceReminder(ctx context.Context, invoiceID int64) error
}

type invoiceService struct {
	*deps.ServiceDependencies
}

func NewInvoiceService(deps *deps.ServiceDependencies) InvoiceService {
	return &invoiceService{
		ServiceDependencies: deps,
	}
}
