package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateClientDetailsTxParams struct {
	CreateClientParams  CreateClientDetailsParams
	IdentityAttachments []uuid.UUID
}

type CreateClientDetailsTxResult struct {
	Client ClientDetail
}

func (store *Store) CreateClientDetailsTx(ctx context.Context, arg CreateClientDetailsTxParams) (CreateClientDetailsTxResult, error) {
	var result CreateClientDetailsTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {
		// First check and update all attachments sequentially
		for _, attachmentID := range arg.IdentityAttachments {
			_, err := q.SetAttachmentAsUsed(ctx, attachmentID)
			if err != nil {
				return fmt.Errorf("failed to set attachment %s as used: %w", attachmentID, err)
			}
		}

		// Then create the client
		var err error
		result.Client, err = q.CreateClientDetails(ctx, arg.CreateClientParams)
		if err != nil {
			return fmt.Errorf("failed to create client details: %w", err)
		}

		return nil
	})

	return result, err
}

type SetClientProfilePictureTxParams struct {
	ClientID     int64
	AttachmentID uuid.UUID
}

type SetClientProfilePictureTxResult struct {
	User ClientDetail
}

func (store *Store) SetClientProfilePictureTx(ctx context.Context, arg SetClientProfilePictureTxParams) (SetClientProfilePictureTxResult, error) {
	var result SetClientProfilePictureTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {

		attachement, err := q.SetAttachmentAsUsed(ctx, arg.AttachmentID)
		if err != nil {
			return fmt.Errorf("failed to set attachment %s as used: %w", arg.AttachmentID, err)
		}

		result.User, err = q.SetClientProfilePicture(ctx, SetClientProfilePictureParams{
			ID:             arg.ClientID,
			ProfilePicture: &attachement.File,
		})
		if err != nil {
			return fmt.Errorf("failed to create client details: %w", err)
		}

		return nil
	})

	return result, err
}
