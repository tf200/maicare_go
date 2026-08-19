package repository

import (
	"context"

	db "maicare_go/db/sqlc"
)

func actorQuery[T any](ctx context.Context, store *db.Store, query func(*db.Queries) (T, error)) (T, error) {
	var result T
	err := store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		result, err = query(q)
		return err
	})
	return result, err
}
