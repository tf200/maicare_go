package invoice

import (
	"context"

	"maicare_go/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *invoiceService) SendInvoiceReminder(ctx context.Context, invoiceID uuid.UUID) error {
	senderID, err := s.Store.GetInvoiceSenderID(ctx, invoiceID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "SendInvoiceReminder", "Failed to get sender ID for invoice",
			zap.Error(err), zap.String("invoiceID", invoiceID.String()))
		return err
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "SendInvoiceReminder", "Sending reminder email",
		zap.String("invoiceID", invoiceID.String()), zap.String("senderID", senderID.String()))
	// Logic to send reminder email using senderID
	return nil
}
