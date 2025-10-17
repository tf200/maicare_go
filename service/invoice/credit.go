package invoice

import (
	"context"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// CreditInvoiceResponse represents the response body for crediting an invoice.
type CreditInvoiceResponse struct {
	ID int64 `json:"id"`
}

func (s *invoiceService) CreditInvoice(ctx context.Context, invoiceID int64, employeeID int64) (*CreditInvoiceResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to begin transaction", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL myapp.current_employee_id = %d", employeeID))
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to set current employee ID", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to set current employee ID: %v", err)
	}
	qtx := s.Store.WithTx(tx)
	originalInvoice, err := qtx.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to get original invoice", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to get original invoice: %v", err)
	}

	if originalInvoice.Status == "credit_note" {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Cannot credit a credit note", zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("cannot credit a credit note")
	}

	var invoiceDetails []InvoiceDetails
	err = json.Unmarshal(originalInvoice.InvoiceDetails, &invoiceDetails)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to unmarshal invoice details", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to unmarshal invoice details: %v", err)
	}

	creditNoteInvoicedetails := make([]InvoiceDetails, len(invoiceDetails))
	for i, detail := range invoiceDetails {
		creditNoteInvoicedetails[i] = InvoiceDetails{
			ContractID:    detail.ContractID,
			ContractType:  detail.ContractType,
			Price:         -detail.Price, // Negate the price for credit note
			PriceTimeUnit: detail.PriceTimeUnit,
			PreVatTotal:   -detail.PreVatTotal, // Negate the pre-VAT total for credit note
			Total:         -detail.Total,       // Negate the total for credit note
			Periods:       detail.Periods,
		}
	}
	creditNoteDetailsBytes, err := json.Marshal(creditNoteInvoicedetails)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to marshal credit note invoice details", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to marshal credit note invoice details: %v", err)
	}

	invoiceNumber, invoiceSequence, err := s.GenerateInvoiceNumber(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to generate invoice number", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to generate invoice number: %v", err)
	}

	arg := db.CreateInvoiceParams{
		ClientID:        originalInvoice.ClientID,
		SenderID:        &originalInvoice.ID,
		DueDate:         pgtype.Date{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
		TotalAmount:     -originalInvoice.TotalAmount,
		InvoiceDetails:  creditNoteDetailsBytes,
		InvoiceNumber:   invoiceNumber,
		InvoiceSequence: invoiceSequence,
		InvoiceType:     "credit_note",
		ExtraContent:    []byte("{}"), // Assuming no extra content for credit notes
		WarningCount:    0,
		IssueDate:       pgtype.Date{Time: time.Now(), Valid: true},
	}
	creditInvoice, err := qtx.CreateInvoice(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to create credit invoice", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to create credit invoice: %v", err)
	}

	_, err = qtx.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
		ID:     originalInvoice.ID,
		Status: "canceled",
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to update original invoice status", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to update original invoice status: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreditInvoice", "Failed to commit transaction", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &CreditInvoiceResponse{
		ID: creditInvoice.ID,
	}, nil

}
