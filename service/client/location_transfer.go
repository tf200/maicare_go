package clientp

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
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

func (s *clientService) ListLocationTransferRequests(ctx *gin.Context, req ListLocationTransferRequestsRequest) (*pagination.Response[ListLocationTransferRequestsResponse], error) {
	params := req.GetParams()

	transferRequests, err := s.Store.ListClientLocationTransfer(ctx, db.ListClientLocationTransferParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListLocationTransferRequestsApi", "Failed to list location transfer requests", zap.Error(err))
		return nil, fmt.Errorf("failed to list location transfer requests")
	}

	var responseList []ListLocationTransferRequestsResponse
	for _, transferRequest := range transferRequests {
		responseList = append(responseList, ListLocationTransferRequestsResponse{
			ID:                 transferRequest.ID,
			ClientID:           transferRequest.ClientID,
			FromLocationID:     transferRequest.FromLocationID,
			ToLocationID:       transferRequest.ToLocationID,
			NewMentorID:        transferRequest.NewMentorID,
			RequestDate:        transferRequest.RequestDate.Time,
			Status:             transferRequest.Status,
			ApprovedRejectedBy: transferRequest.ApprovedRejectedBy,
			ApprovedRejectedAt: transferRequest.ApprovedRejectedAt.Time,
			Reason:             transferRequest.Reason,
			MentorFirstName:    transferRequest.MentorFirstName,
			MentorLastName:     transferRequest.MentorLastName,
		})
	}
	if len(responseList) == 0 {
		return &pagination.Response[ListLocationTransferRequestsResponse]{}, nil
	}
	totalCount := transferRequests[0].TotalCount

	res := pagination.NewResponse(ctx, req.Request, responseList, totalCount)
	return &res, nil
}
