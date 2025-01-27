package db

import (
	"context"
	"fmt"
	"maicare_go/util"

	"github.com/google/uuid"
)

type SetEmployeeProfilePictureTxParams struct {
	EmployeeID    int64
	AttachementID uuid.UUID
}

type SetEmployeeProfilePictureTxResult struct {
	User CustomUser
}

func (store *Store) SetEmployeeProfilePictureTx(ctx context.Context, arg SetEmployeeProfilePictureTxParams) (SetEmployeeProfilePictureTxResult, error) {
	var result SetEmployeeProfilePictureTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {

		attachement, err := q.SetAttachmentAsUsed(ctx, arg.AttachementID)
		if err != nil {
			return fmt.Errorf("failed to set attachment %s as used: %w", arg.AttachementID, err)
		}

		result.User, err = q.SetEmployeeProfilePicture(ctx, SetEmployeeProfilePictureParams{
			ID:             arg.EmployeeID,
			ProfilePicture: util.StringPtr(attachement.File),
		})
		if err != nil {
			return fmt.Errorf("failed to create client details: %w", err)
		}

		return nil
	})

	return result, err
}
