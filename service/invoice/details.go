package invoice

import (
	"context"
	"encoding/json"
	"fmt"
	"maicare_go/logger"
	"maicare_go/util"

	"go.uber.org/zap"
)

func (s *invoiceService) GetInvoiceByID(ctx context.Context, invoiceID int64) (*GetInvoiceByIDResponse, error) {
	inv, err := s.Store.GetInvoice(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetInvoiceByID", "Failed to get invoice", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, err
	}

	var invoiceDetails []InvoiceDetails
	if err := json.Unmarshal(inv.InvoiceDetails, &invoiceDetails); err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetInvoiceByID", "Failed to unmarshal invoice details", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, fmt.Errorf("failed to unmarshal invoice details")
	}

	paymentCompletionPrc := s.calculatePaymentCompletionPercentage(ctx, inv.TotalAmount, invoiceID)

	return &GetInvoiceByIDResponse{
		ID:                   inv.ID,
		InvoiceNumber:        inv.InvoiceNumber,
		IssueDate:            inv.IssueDate.Time,
		DueDate:              inv.DueDate.Time,
		Status:               inv.Status,
		InvoiceDetails:       invoiceDetails,
		TotalAmount:          inv.TotalAmount,
		PdfAttachmentID:      inv.PdfAttachmentID,
		ExtraContent:         util.ParseJSONToObject(inv.ExtraContent),
		ClientID:             inv.ClientID,
		SenderID:             inv.SenderID,
		InvoiceType:          inv.InvoiceType,
		OriginalInvoiceID:    inv.OriginalInvoiceID,
		UpdatedAt:            inv.UpdatedAt.Time,
		CreatedAt:            inv.CreatedAt.Time,
		SenderName:           inv.SenderName,
		SenderKvknumber:      inv.SenderKvknumber,
		SenderBtwnumber:      inv.SenderBtwnumber,
		ClientFirstName:      inv.ClientFirstName,
		ClientLastName:       inv.ClientLastName,
		PaymentCompletionPrc: paymentCompletionPrc,
	}, nil
}
