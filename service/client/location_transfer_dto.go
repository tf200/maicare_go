package clientp

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"time"

	"github.com/google/uuid"
)

type LocationTransferRequest struct {
	FromLocationID uuid.UUID  `json:"from_location_id" binding:"required"`
	ToLocationID   uuid.UUID  `json:"to_location_id" binding:"required"`
	Reason         string     `json:"reason" binding:"required"`
	NewMentorID    *uuid.UUID `json:"new_mentor_id"`
}

type ApproveOrRejectLocationTransferRequest struct {
	TransferID uuid.UUID `json:"transfer_id" binding:"required"`
	Status     string    `json:"status" binding:"required,oneof=approved rejected"`
}

type ListLocationTransferRequestsRequest struct {
	pagination.Request
}

type ListLocationTransferRequestsResponse struct {
	ID                 uuid.UUID                           `json:"id"`
	ClientID           uuid.UUID                           `json:"client_id"`
	FromLocationID     *uuid.UUID                          `json:"from_location_id"`
	ToLocationID       *uuid.UUID                          `json:"to_location_id"`
	NewMentorID        *uuid.UUID                          `json:"new_mentor_id"`
	RequestDate        time.Time                           `json:"request_date"`
	Status             db.ClientLocationTransferStatusEnum `json:"status"`
	ApprovedRejectedBy *uuid.UUID                          `json:"approved_rejected_by"`
	ApprovedRejectedAt time.Time                           `json:"approved_rejected_at"`
	Reason             *string                             `json:"reason"`
	MentorFirstName    *string                             `json:"mentor_first_name"`
	MentorLastName     *string                             `json:"mentor_last_name"`
}
