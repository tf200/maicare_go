package invoice

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

const (
	PAYMENT_TOLERANCE float64 = 50
)

func (s *invoiceService) CreatePayment(ctx context.Context, invoiceID uuid.UUID, req CreatePaymentRequest, employeeID uuid.UUID) (*CreatePaymentResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to begin transaction", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to set current employee ID", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.Store.WithTx(tx)

	getInvoice, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to get invoice", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to get invoice: %v", err)
	}

	paymentParams := db.CreatePaymentParams{
		InvoiceID:        invoiceID,
		PaymentMethod:    db.PaymentMethodEnum(req.PaymentMethod),
		PaymentStatus:    db.PaymentStatusEnum(req.PaymentStatus),
		Amount:           req.Amount,
		PaymentDate:      pgtype.Date{Time: req.PaymentDate, Valid: true},
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
		RecordedBy:       &employeeID,
	}

	payment, err := qtx.CreatePayment(ctx, paymentParams)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to create payment", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	var newInvoiceStatus db.InvoiceStatusEnum
	invoiceStatusChanged := false

	if req.PaymentStatus == string(PaymentStatusCompleted) {
		totalPaid, err := qtx.GetCompletedPaymentSum(ctx, invoiceID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to get total completed payment", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to get total completed payment: %v", err)
		}

		newInvoiceStatus, err = DetermineInvoiceStatus(getInvoice.TotalAmount, totalPaid)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to determine invoice status", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to determine invoice status: %v", err)
		}

		if newInvoiceStatus != db.InvoiceStatusEnum(getInvoice.Status) {
			invoiceStatusChanged = true
			_, err = qtx.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
				ID:     invoiceID,
				Status: newInvoiceStatus,
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to update invoice status", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
				return nil, fmt.Errorf("failed to update invoice status: %v", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreatePayment", "Failed to commit transaction", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	response := &CreatePaymentResponse{
		PaymentID:            payment.ID,
		InvoiceID:            payment.InvoiceID,
		PaymentMethod:        string(payment.PaymentMethod),
		PaymentStatus:        string(payment.PaymentStatus),
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

func (s *invoiceService) ListPayments(ctx context.Context, invoiceID uuid.UUID) ([]ListPaymentsResponse, error) {
	payments, err := s.Store.ListPayments(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListPayments", "Failed to list payments", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
		return nil, fmt.Errorf("failed to list payments: %v", err)
	}

	response := []ListPaymentsResponse{}
	for _, payment := range payments {
		response = append(response, ListPaymentsResponse{
			PaymentID:           payment.ID,
			InvoiceID:           payment.InvoiceID,
			PaymentMethod:       string(payment.PaymentMethod),
			PaymentStatus:       string(payment.PaymentStatus),
			Amount:              payment.Amount,
			PaymentDate:         payment.PaymentDate.Time,
			PaymentReference:    payment.PaymentReference,
			Notes:               payment.Notes,
			RecordedBy:          payment.RecordedBy,
			CreatedAt:           payment.CreatedAt.Time,
			UpdatedAt:           payment.UpdatedAt.Time,
			RecordedByFirstName: payment.RecordedByFirstName,
			RecordedByLastName:  payment.RecordedByLastName,
		})
	}

	return response, nil
}

func (s *invoiceService) GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (*GetPaymentByIDResponse, error) {
	payment, err := s.Store.GetPayment(ctx, paymentID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetPaymentByID", "Failed to get payment by ID", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to get payment by ID: %v", err)
	}

	response := &GetPaymentByIDResponse{
		PaymentID:           payment.ID,
		InvoiceID:           payment.InvoiceID,
		PaymentMethod:       string(payment.PaymentMethod),
		PaymentStatus:       string(payment.PaymentStatus),
		Amount:              payment.Amount,
		PaymentDate:         payment.PaymentDate.Time,
		PaymentReference:    payment.PaymentReference,
		Notes:               payment.Notes,
		RecordedBy:          payment.RecordedBy,
		CreatedAt:           payment.CreatedAt.Time,
		UpdatedAt:           payment.UpdatedAt.Time,
		RecordedByFirstName: payment.RecordedByFirstName,
		RecordedByLastName:  payment.RecordedByLastName,
	}

	return response, nil
}

func (s *invoiceService) UpdatePayment(ctx context.Context, invoiceID uuid.UUID, employeeID uuid.UUID, paymentID uuid.UUID, req UpdatePaymentRequest) (*UpdatePaymentResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to begin transaction", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to set current employee ID", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.Store.WithTx(tx)
	currentPayment, err := qtx.GetPaymentWithInvoice(ctx, paymentID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to get current payment", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to get current payment: %v", err)
	}

	if currentPayment.InvoiceID != invoiceID {
		return nil, fmt.Errorf("payment does not belong to the specified invoice")
	}

	originalInvoiceStatus := currentPayment.InvoiceStatus
	updateParams := db.UpdatePaymentParams{
		ID:               paymentID,
		PaymentMethod:    db.NullPaymentMethodFromPtr(req.PaymentMethod),
		PaymentStatus:    db.NullPaymentStatusFromPtr(req.PaymentStatus),
		Amount:           req.Amount,
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
		RecordedBy:       &employeeID,
	}
	if req.PaymentDate != nil {
		updateParams.PaymentDate = pgtype.Date{Time: *req.PaymentDate, Valid: true}
	}

	updatedPayment, err := qtx.UpdatePayment(ctx, updateParams)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to update payment", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	var newInvoiceStatus db.InvoiceStatusEnum
	var statusChanged bool = false

	if updatedPayment.PaymentStatus == db.PaymentStatusEnumCompleted ||
		currentPayment.PaymentStatus == db.PaymentStatusEnumCompleted {
		totalPaid, err := qtx.GetTotalPaidAmountByInvoice(ctx, invoiceID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to get total paid amount", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to get total paid amount: %v", err)
		}

		newStatus, err := DetermineInvoiceStatus(currentPayment.InvoiceTotalAmount, totalPaid)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to determine invoice status", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to determine invoice status: %v", err)
		}

		if newStatus != originalInvoiceStatus {
			updatedInvoice, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
				ID:     invoiceID,
				Status: db.NullInvoiceStatusEnum{InvoiceStatusEnum: newStatus, Valid: true},
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to update invoice status", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
				return nil, fmt.Errorf("failed to update invoice status: %v", err)
			}
			newInvoiceStatus = updatedInvoice.Status
			statusChanged = true
		} else {
			newInvoiceStatus = originalInvoiceStatus
		}
	} else {
		newInvoiceStatus = originalInvoiceStatus
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdatePayment", "Failed to commit transaction", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	response := &UpdatePaymentResponse{
		PaymentID:             updatedPayment.ID,
		InvoiceID:             updatedPayment.InvoiceID,
		PaymentMethod:         string(updatedPayment.PaymentMethod),
		PaymentStatus:         string(updatedPayment.PaymentStatus),
		Amount:                updatedPayment.Amount,
		PaymentDate:           updatedPayment.PaymentDate.Time,
		PaymentReference:      updatedPayment.PaymentReference,
		Notes:                 updatedPayment.Notes,
		RecordedBy:            updatedPayment.RecordedBy,
		InvoiceStatusChanged:  statusChanged,
		CurrentInvoiceStatus:  string(newInvoiceStatus),
		PreviousInvoiceStatus: string(originalInvoiceStatus),
	}

	return response, nil
}

func (s *invoiceService) DeletePayment(ctx context.Context, invoiceID uuid.UUID, paymentID uuid.UUID, employeeID uuid.UUID) (*DeletePaymentResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to begin transaction", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to set current employee ID", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.Store.WithTx(tx)

	paymentToDelete, err := qtx.GetPaymentWithInvoice(ctx, paymentID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to get payment to delete", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to get payment to delete: %v", err)
	}

	if paymentToDelete.InvoiceID != invoiceID {
		return nil, fmt.Errorf("payment does not belong to the specified invoice")
	}

	// Store original invoice status for comparison
	originalInvoiceStatus := paymentToDelete.InvoiceStatus

	deletedPayment, err := qtx.DeletePayment(ctx, paymentID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to delete payment", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to delete payment: %v", err)
	}

	var newInvoiceStatus db.InvoiceStatusEnum
	var statusChanged bool = false

	if deletedPayment.PaymentStatus == db.PaymentStatusEnumCompleted {
		totalPaid, err := qtx.GetTotalPaidAmountByInvoice(ctx, invoiceID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to get total paid amount", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to get total paid amount: %v", err)
		}

		newStatus, err := DetermineInvoiceStatus(paymentToDelete.InvoiceTotalAmount, totalPaid)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to determine invoice status", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
			return nil, fmt.Errorf("failed to determine invoice status: %v", err)
		}

		if newStatus != originalInvoiceStatus {
			updatedInvoice, err := qtx.UpdateInvoice(ctx, db.UpdateInvoiceParams{
				ID:     invoiceID,
				Status: db.NullInvoiceStatusEnum{InvoiceStatusEnum: newStatus, Valid: true},
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to update invoice status", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
				return nil, fmt.Errorf("failed to update invoice status: %v", err)
			}
			newInvoiceStatus = updatedInvoice.Status
			statusChanged = true
		} else {
			newInvoiceStatus = originalInvoiceStatus
		}
	} else {
		newInvoiceStatus = originalInvoiceStatus
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeletePayment", "Failed to commit transaction", zap.Error(err), zap.String("payment_id", paymentID.String()))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &DeletePaymentResponse{
		DeletedPaymentID:      deletedPayment.ID,
		InvoiceID:             deletedPayment.InvoiceID,
		DeletedAmount:         deletedPayment.Amount,
		DeletedPaymentStatus:  string(deletedPayment.PaymentStatus),
		InvoiceStatusChanged:  statusChanged,
		CurrentInvoiceStatus:  string(newInvoiceStatus),
		PreviousInvoiceStatus: string(originalInvoiceStatus),
	}, nil
}

func DetermineInvoiceStatus(invoiceTotal, totalPaid float64) (db.InvoiceStatusEnum, error) {
	diffrence := totalPaid - invoiceTotal

	if totalPaid <= PAYMENT_TOLERANCE {
		return db.InvoiceStatusEnumOutstanding, nil
	}

	if diffrence < -PAYMENT_TOLERANCE {
		return db.InvoiceStatusEnumPartiallyPaid, nil
	}

	if diffrence >= -PAYMENT_TOLERANCE && diffrence <= PAYMENT_TOLERANCE {
		return db.InvoiceStatusEnumPaid, nil
	}
	if diffrence > PAYMENT_TOLERANCE {
		return db.InvoiceStatusEnumOverpaid, nil
	}

	return "", fmt.Errorf("could not determine invoice status for totalPaid: %f, invoiceTotal: %f", totalPaid, invoiceTotal)
}

func (s *invoiceService) calculatePaymentCompletionPercentage(ctx context.Context, totalAmount float64, invoiceID uuid.UUID) float64 {
	if totalAmount == 0 {
		return 0
	}

	totalPaid, err := s.Store.GetCompletedPaymentSum(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "calculatePaymentCompletionPercentage", "Failed to get total completed payment", zap.Error(err), zap.String("invoice_id", invoiceID.String()))
		return 0
	}
	return (totalPaid / totalAmount) * 100
}
