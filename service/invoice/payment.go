package invoice

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

const (
	PAYMENT_TOLERANCE float64 = 50
)

func (s *invoiceService) CreatePayment(ctx context.Context, invoiceID int64, req CreatePaymentRequest, employeeID int64) (*CreatePaymentResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to begin transaction", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to set current employee ID", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.Store.WithTx(tx)

	getInvoice, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to get invoice", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to get invoice: %v", err)
	}

	paymentParams := db.CreatePaymentParams{
		InvoiceID:        invoiceID,
		PaymentMethod:    req.PaymentMethod,
		PaymentStatus:    req.PaymentStatus,
		Amount:           req.Amount,
		PaymentDate:      pgtype.Date{Time: req.PaymentDate, Valid: true},
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
		RecordedBy:       &employeeID,
	}

	payment, err := qtx.CreatePayment(ctx, paymentParams)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to create payment", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	var newInvoiceStatus InvoiceStatus
	invoiceStatusChanged := false

	if req.PaymentStatus == string(PaymentStatusCompleted) {
		totalPaid, err := qtx.GetCompletedPaymentSum(ctx, invoiceID)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to get total completed payment", zap.Error(err), zap.Int64("invoice_id", invoiceID))
			return nil, fmt.Errorf("failed to get total completed payment: %v", err)
		}

		newInvoiceStatus, err = DetermineInvoiceStatus(getInvoice.TotalAmount, totalPaid)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to determine invoice status", zap.Error(err), zap.Int64("invoice_id", invoiceID))
			return nil, fmt.Errorf("failed to determine invoice status: %v", err)
		}

		if newInvoiceStatus != InvoiceStatus(getInvoice.Status) {
			invoiceStatusChanged = true
			_, err = qtx.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
				ID:     invoiceID,
				Status: string(newInvoiceStatus),
			})
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to update invoice status", zap.Error(err), zap.Int64("invoice_id", invoiceID))
				return nil, fmt.Errorf("failed to update invoice status: %v", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreatePayment", "Failed to commit transaction", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	response := &CreatePaymentResponse{
		PaymentID:            payment.ID,
		InvoiceID:            payment.InvoiceID,
		PaymentMethod:        payment.PaymentMethod,
		PaymentStatus:        payment.PaymentStatus,
		Amount:               payment.Amount,
		PaymentDate:          payment.PaymentDate.Time,
		PaymentReference:     payment.PaymentReference,
		Notes:                payment.Notes,
		RecordedBy:           payment.RecordedBy,
		InvoiceStatusChanged: invoiceStatusChanged,
		CurrentInvoiceStatus: string(newInvoiceStatus),
	}

	return response, nil
}

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
