package clientp

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *clientService) RequestLocationTransfer(ctx context.Context, clientID uuid.UUID, req LocationTransferRequest) error {
	err := s.Store.CreateClientLocationTransfer(ctx, db.CreateClientLocationTransferParams{
		ClientID:       clientID,
		FromLocationID: &req.FromLocationID,
		ToLocationID:   &req.ToLocationID,
		Reason:         &req.Reason,
		NewMentorID:    req.NewMentorID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "RequestLocationTransferApi", "Failed to create location transfer request", zap.Error(err))
		return fmt.Errorf("failed to create location transfer request")
	}

	return nil
}

func (s *clientService) ApproveLocationTransfer(ctx context.Context, employeeID uuid.UUID, req ApproveOrRejectLocationTransferRequest) error {
	err := s.Store.ApproveOrRejectClientLocationTransfer(ctx, db.ApproveOrRejectClientLocationTransferParams{
		ID:                 req.TransferID,
		Status:             db.ClientLocationTransferStatusEnum(req.Status),
		ApprovedRejectedBy: &employeeID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ApproveLocationTransferApi", "Failed to approve location transfer request", zap.Error(err))
		return fmt.Errorf("failed to approve location transfer request")
	}

	return nil
}
