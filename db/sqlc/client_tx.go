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

type AddClientDocumentTxParams struct {
	ClientID     uuid.UUID
	AttachmentID uuid.UUID
	Label        string
}

type AddClientDocumentTxResults struct {
	ClientDocument ClientDocument
	Attachment     AttachmentFile
}

type AddClientDocumentsTxParams struct {
	ClientID  uuid.UUID
	Documents []AddClientDocumentTxParams
}

type AddClientDocumentsTxResults struct {
	Documents []AddClientDocumentTxResults
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
			AttachmentUuid: &result.Attachment.Uuid,
			Label:          ClientDocumentLabelEnum(arg.Label),
		})
		if err != nil {
			return fmt.Errorf("failed to create client details: %w", err)
		}

		return nil
	})

	return result, err
}

func (store *Store) AddClientDocumentsTx(ctx context.Context, arg AddClientDocumentsTxParams) (AddClientDocumentsTxResults, error) {
	var result AddClientDocumentsTxResults

	err := store.ExecTx(ctx, func(q *Queries) error {
		attachmentIDs := make([]uuid.UUID, 0, len(arg.Documents))
		for _, doc := range arg.Documents {
			attachmentIDs = append(attachmentIDs, doc.AttachmentID)
		}

		attachments, err := q.SetAttachmentsAsUsedorUnusedByUUIDs(ctx, SetAttachmentsAsUsedorUnusedByUUIDsParams{
			Column1: attachmentIDs,
			IsUsed:  true,
		})
		if err != nil {
			return fmt.Errorf("failed to set attachments as used: %w", err)
		}

		attachmentsByID := make(map[uuid.UUID]AttachmentFile, len(attachments))
		for _, attachment := range attachments {
			attachmentsByID[attachment.Uuid] = attachment
		}

		result.Documents = make([]AddClientDocumentTxResults, 0, len(arg.Documents))
		for _, doc := range arg.Documents {
			item := AddClientDocumentTxResults{}
			attachmentRow, ok := attachmentsByID[doc.AttachmentID]
			if !ok {
				return fmt.Errorf("attachment %s was not updated", doc.AttachmentID)
			}

			clientDoc, err := q.CreateClientDocument(ctx, CreateClientDocumentParams{
				ClientID:       arg.ClientID,
				AttachmentUuid: &attachmentRow.Uuid,
				Label:          ClientDocumentLabelEnum(doc.Label),
			})
			if err != nil {
				return fmt.Errorf("failed to create client document for attachment %s: %w", doc.AttachmentID, err)
			}

			item.Attachment = attachmentRow
			item.ClientDocument = clientDoc
			result.Documents = append(result.Documents, item)
		}

		return nil
	})

	return result, err
}

type DeleteClientDocumentTxParams struct {
	ClientID   uuid.UUID
	DocumentID uuid.UUID
}

type DeleteClientDocumentResults struct {
	ClientDocument ClientDocument
	Attachment     AttachmentFile
}

func (store *Store) DeleteClientDocumentTx(ctx context.Context, arg DeleteClientDocumentTxParams) (DeleteClientDocumentResults, error) {
	var result DeleteClientDocumentResults

	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error
		result.ClientDocument, err = q.DeleteClientDocument(ctx, DeleteClientDocumentParams{
			ID:       arg.DocumentID,
			ClientID: arg.ClientID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete client document %s: %w", arg.DocumentID, err)
		}
		if result.ClientDocument.AttachmentUuid == nil {
			return fmt.Errorf("client document %s has no linked attachment", arg.DocumentID)
		}

		result.Attachment, err = q.SetAttachmentAsUsedorUnused(ctx, SetAttachmentAsUsedorUnusedParams{
			Uuid:   *result.ClientDocument.AttachmentUuid,
			IsUsed: false,
		})
		if err != nil {
			return fmt.Errorf("failed to set attachment %s as unused: %w", result.ClientDocument.AttachmentUuid.String(), err)
		}

		return nil
	})

	return result, err
}
