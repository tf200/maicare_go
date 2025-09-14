package invoice

import (
	"context"
	"maicare_go/service/deps"
)

type InvoiceService interface {
	GenerateInvoice(req GenerateInvoiceRequest, ctx context.Context) (*GenerateInvoiceResult, int64, error)
}

type invoiceService struct {
	*deps.ServiceDependencies
}

func NewInvoiceService(deps *deps.ServiceDependencies) InvoiceService {
	return &invoiceService{
		ServiceDependencies: deps,
	}
}
