package clientp

import "github.com/google/uuid"

type LocationTransferRequest struct {
	FromLocationID int64      `json:"from_location_id" binding:"required"`
	ToLocationID   int64      `json:"to_location_id" binding:"required"`
	Reason         string     `json:"reason" binding:"required"`
	NewMentorID    *uuid.UUID `json:"new_mentor_id"`
}

type ApproveOrRejectLocationTransferRequest struct {
	TransferID int64  `json:"transfer_id" binding:"required"`
	Status     string `json:"status" binding:"required,oneof=approved rejected"`
}
