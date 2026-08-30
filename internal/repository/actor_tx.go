package repository

import (
	"context"

	db "maicare_go/db/sqlc"

	"github.com/jackc/pgx/v5"
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

func actorQueryRepeatableRead[T any](ctx context.Context, store *db.Store, query func(*db.Queries) (T, error)) (T, error) {
	var result T
	tx, err := store.BeginActorTxWithOptions(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	result, err = query(store.WithTx(tx))
	if err != nil {
		return result, err
	}
	if err := tx.Commit(ctx); err != nil {
		return result, err
	}
	return result, nil
}
