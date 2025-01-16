package db

import (
	"context"

	"github.com/google/uuid"
)

type CreateClientDetailsTxParams struct {
	CreateClientParams CreateClientDetailsParams
	IdentityAttachment uuid.UUID
}

type CreateClientDetailsTxResult struct {
	Client     ClientDetail
	Attachment AttachmentFile
}

func (store *Store) CreateClientDetailsTx(ctx context.Context, arg CreateClientDetailsTxParams) (CreateClientDetailsTxResult, error) {
	var result CreateClientDetailsTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {

		var err error

		result.Attachment, err = q.SetAttachmentAsUsed(ctx, arg.IdentityAttachment)
		if err != nil {
			return err
		}

		result.Client, err = q.CreateClientDetails(ctx, arg.CreateClientParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
