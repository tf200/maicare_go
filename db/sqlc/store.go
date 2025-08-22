package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	ConnPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		ConnPool: connPool,
		Queries:  New(connPool),
	}
}

type TxFn func(queries *Queries) error

// ExecTx executes a function within a database transaction
func (store *Store) ExecTx(ctx context.Context, fn TxFn) error {
	tx, err := store.ConnPool.Begin(ctx)
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
	RoleID               int32
}

type CreateEmployeeWithAccountTxResult struct {
	User     CustomUser
	Employee EmployeeProfile
}

func (store *Store) CreateEmployeeWithAccountTx(ctx context.Context, arg CreateEmployeeWithAccountTxParams) (CreateEmployeeWithAccountTxResult, error) {
	var result CreateEmployeeWithAccountTxResult

	err := store.ExecTx(ctx, func(q *Queries) error {
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

		err = q.GrantRoleToUser(ctx, GrantRoleToUserParams{
			UserID: result.User.ID,
			RoleID: arg.RoleID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
