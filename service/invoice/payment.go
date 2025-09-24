package invoice

import (
	"context"
	"fmt"
	"maicare_go/logger"

	"go.uber.org/zap"
)

const (
	PAYMENT_TOLERANCE float64 = 50
)

func DetermineInvoiceStatus(invoiceTotal, totalPaid float64) (InvoiceStatus, error) {
	diffrence := totalPaid - invoiceTotal

	if totalPaid <= PAYMENT_TOLERANCE {
		return InvoiceStatusOutstanding, nil
	}

	if diffrence < -PAYMENT_TOLERANCE {
		return InvoiceStatusPartiallyPaid, nil
	}

	if diffrence >= -PAYMENT_TOLERANCE && diffrence <= PAYMENT_TOLERANCE {
		return InvoiceStatusPaid, nil
	}
	if diffrence > PAYMENT_TOLERANCE {
		return InvoiceStatusOverpaid, nil
	}

	return "", fmt.Errorf("could not determine invoice status for totalPaid: %f, invoiceTotal: %f", totalPaid, invoiceTotal)

}

func (s *invoiceService) calculatePaymentCompletionPercentage(ctx context.Context, totalAmount float64, invoiceID int64) float64 {
	if totalAmount == 0 {
		return 0
	}

	totalPaid, err := s.Store.GetCompletedPaymentSum(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "calculatePaymentCompletionPercentage", "Failed to get total completed payment", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return 0
	}
	return (totalPaid / totalAmount) * 100
}
