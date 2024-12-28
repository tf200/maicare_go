package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

type TxFn func(queries *Queries) error

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn TxFn) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}

type CreateEmployeeWithAccountTxParams struct {
	CreateUserParams     CreateUserParams
	CreateEmployeeParams CreateEmployeeProfileParams
}

type CreateEmployeeWithAccountTxResult struct {
	User     CustomUser
	Employee EmployeeProfile
}

func (store *Store) CreateEmployeeWithAccountTx(ctx context.Context, arg CreateEmployeeWithAccountTxParams) (CreateEmployeeWithAccountTxResult, error) {
	var result CreateEmployeeWithAccountTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		arg.CreateEmployeeParams.UserID = result.User.ID
		result.Employee, err = q.CreateEmployeeProfile(ctx, arg.CreateEmployeeParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

type CreateClientWithAccountTxParams struct {
	CreateUserParams   CreateUserParams
	CreateClientParams CreateClientDetailsParams
}

type CreateClientWithAccountTxResult struct {
	User   CustomUser
	Client ClientDetail
}

func (store *Store) CreateClientWithAccountTx(ctx context.Context, arg CreateClientWithAccountTxParams) (CreateClientWithAccountTxResult, error) {
	var result CreateClientWithAccountTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		arg.CreateClientParams.UserID = result.User.ID
		result.Client, err = q.CreateClientDetails(ctx, arg.CreateClientParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
