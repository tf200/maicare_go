package db

import (
	"context"
	"sync"

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
		var wg sync.WaitGroup
		errChan := make(chan error, len(arg.IdentityAttachments))

		for _, attachmentID := range arg.IdentityAttachments {
			wg.Add(1)
			go func(id uuid.UUID) {
				defer wg.Done()
				_, err := q.SetAttachmentAsUsed(ctx, id)
				if err != nil {
					errChan <- err
				}
			}(attachmentID)
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			return err
		}

		var err error
		result.Client, err = q.CreateClientDetails(ctx, arg.CreateClientParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
