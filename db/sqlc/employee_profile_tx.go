package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type SetEmployeeProfilePictureTxParams struct {
	EmployeeID    uuid.UUID
	AttachementID uuid.UUID
}

type SetEmployeeProfilePictureTxResult struct {
	User CustomUser
}

func (store *Store) SetEmployeeProfilePictureTx(ctx context.Context, arg SetEmployeeProfilePictureTxParams) (SetEmployeeProfilePictureTxResult, error) {
	var result SetEmployeeProfilePictureTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {
		attachement, err := q.SetAttachmentAsUsedorUnused(ctx, SetAttachmentAsUsedorUnusedParams{
			Uuid:   arg.AttachementID,
			IsUsed: true,
		})
		if err != nil {
			return fmt.Errorf("failed to set attachment %s as used: %w", arg.AttachementID, err)
		}

		result.User, err = q.SetEmployeeProfilePicture(ctx, SetEmployeeProfilePictureParams{
			ID:             arg.EmployeeID,
			ProfilePicture: stringPtr(attachement.File),
		})
		if err != nil {
			return fmt.Errorf("failed to create client details: %w", err)
		}

		return nil
	})

	return result, err
}

func stringPtr(s string) *string {
	return &s
}
