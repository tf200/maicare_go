package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
			_, err := q.SetAttachmentAsUsedorUnused(ctx, SetAttachmentAsUsedorUnusedParams{
				Uuid:   attachmentID,
				IsUsed: true,
			})
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

		attachement, err := q.SetAttachmentAsUsedorUnused(ctx, SetAttachmentAsUsedorUnusedParams{
			Uuid:   arg.AttachmentID,
			IsUsed: true,
		})

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

type AddClientDocumentTxParams struct {
	ClientID     int64
	AttachmentID uuid.UUID
	Label        string
}

type AddClientDocumentTxResults struct {
	ClientDocument ClientDocument
	Attachment     AttachmentFile
}

func (store *Store) AddClientDocumentTx(ctx context.Context, arg AddClientDocumentTxParams) (AddClientDocumentTxResults, error) {
	var result AddClientDocumentTxResults

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error
		result.Attachment, err = q.SetAttachmentAsUsedorUnused(ctx, SetAttachmentAsUsedorUnusedParams{
			Uuid:   arg.AttachmentID,
			IsUsed: true,
		})
		if err != nil {
			return fmt.Errorf("failed to set attachment %s as used: %w", arg.AttachmentID, err)
		}

		result.ClientDocument, err = q.CreateClientDocument(ctx, CreateClientDocumentParams{
			ClientID:       arg.ClientID,
			AttachmentUuid: pgtype.UUID{Bytes: result.Attachment.Uuid, Valid: true},
			Label:          arg.Label,
		})
		if err != nil {
			return fmt.Errorf("failed to create client details: %w", err)
		}

		return nil
	})

	return result, err
}

type DeleteClientDocumentParams struct {
	AttachmentID uuid.UUID
}

type DeleteClientDocumentResults struct {
	ClientDocument ClientDocument
	Attachment     AttachmentFile
}

func (store *Store) DeleteClientDocumentTx(ctx context.Context, arg DeleteClientDocumentParams) (DeleteClientDocumentResults, error) {
	var result DeleteClientDocumentResults

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error
		result.Attachment, err = q.SetAttachmentAsUsedorUnused(ctx, SetAttachmentAsUsedorUnusedParams{
			Uuid:   arg.AttachmentID,
			IsUsed: false,
		})
		if err != nil {
			return fmt.Errorf("failed to set attachment %s as used: %w", arg.AttachmentID, err)
		}

		result.ClientDocument, err = q.DeleteClientDocument(ctx, pgtype.UUID{Bytes: arg.AttachmentID, Valid: true})
		if err != nil {
			return fmt.Errorf("failed to create client details: %w", err)
		}

		return nil
	})

	return result, err
}
