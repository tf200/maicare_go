package invoice

import (
	"context"
	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *invoiceService) SendInvoiceReminder(ctx context.Context, invoiceID int64) error {
	senderID, err := s.Store.GetInvoiceSenderID(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SendInvoiceReminder", "Failed to get sender ID for invoice",
			zap.Error(err), zap.Int64("invoiceID", invoiceID))
		return err
	}
	if senderID != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "SendInvoiceReminder", "Sending reminder email",
			zap.Int64("invoiceID", invoiceID), zap.Int64("senderID", *senderID))
		// Logic to send reminder email using senderID
	}
	return nil
}
