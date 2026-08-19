package db

import (
	"context"
	"errors"

	"maicare_go/internal/ctxkeys"

	"github.com/goccy/go-json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

var ErrMissingActorIdentity = errors.New("authenticated database transaction requires actor identity")

func (store *Store) ValidateActorIdentity(ctx context.Context, actor ctxkeys.ActorIdentity) error {
	if !actor.IsValid() {
		return ErrMissingActorIdentity
	}
	var valid bool
	err := store.ConnPool.QueryRow(ctx, `SELECT EXISTS (
		SELECT 1
		FROM custom_user cu
		JOIN employee_profile ep ON ep.user_id = cu.id
		WHERE cu.id = $1
		  AND ep.id = $2
		  AND cu.is_active
		  AND ep.is_active
	)`, actor.UserID, actor.EmployeeID).Scan(&valid)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("system actor must reference the same active user and employee")
	}
	return nil
}

// ExecTx executes a transaction and installs actor identity when one is present.
// Bootstrap and maintenance callers without request identity remain supported.
func (store *Store) ExecTx(ctx context.Context, fn TxFn) error {
	return store.execTx(ctx, false, fn)
}

// ExecActorTx executes a transaction that requires authenticated actor identity.
func (store *Store) ExecActorTx(ctx context.Context, fn TxFn) error {
	return store.execTx(ctx, true, fn)
}

// ExecAsActor executes work as an explicitly selected service actor.
func (store *Store) ExecAsActor(ctx context.Context, actor ctxkeys.ActorIdentity, fn TxFn) error {
	if !actor.IsValid() {
		return ErrMissingActorIdentity
	}
	return store.ExecActorTx(ctxkeys.WithActorIdentity(ctx, actor), fn)
}

// BeginActorTx starts a transaction initialized with authenticated actor identity.
// Callers are responsible for commit or rollback.
func (store *Store) BeginActorTx(ctx context.Context) (pgx.Tx, error) {
	actor, ok := ctxkeys.ActorIdentityFromContext(ctx)
	if !ok {
		return nil, ErrMissingActorIdentity
	}
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	if err := initializeActorTx(ctx, tx, actor); err != nil {
		_ = tx.Rollback(ctx)
		return nil, err
	}
	return tx, nil
}

func (store *Store) execTx(ctx context.Context, requireActor bool, fn TxFn) error {
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		return err
	}

	actor, hasActor := ctxkeys.ActorIdentityFromContext(ctx)
	if requireActor && !hasActor {
		_ = tx.Rollback(ctx)
		return ErrMissingActorIdentity
	}
	if hasActor {
		if err := initializeActorTx(ctx, tx, actor); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	q := New(tx)
	if err := fn(q); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func initializeActorTx(ctx context.Context, tx pgx.Tx, actor ctxkeys.ActorIdentity) error {
	if !actor.IsValid() {
		return ErrMissingActorIdentity
	}
	_, err := tx.Exec(ctx, `SELECT
		set_config('myapp.current_user_id', $1, true),
		set_config('myapp.current_employee_id', $2, true)`,
		actor.UserID.String(), actor.EmployeeID.String())
	return err
}

type CreateEmployeeWithAccountTxParams struct {
	CreateUserParams     CreateUserParams
	CreateEmployeeParams CreateEmployeeProfileParams
	RoleID               uuid.UUID
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

		// Auto-assign the active employee handbook for the employee's department (if configured).
		if arg.CreateEmployeeParams.DepartmentID != nil {
			template, tmplErr := q.GetActiveHandbookTemplateByDepartment(ctx, *arg.CreateEmployeeParams.DepartmentID)
			if tmplErr != nil {
				if !errors.Is(tmplErr, pgx.ErrNoRows) {
					return tmplErr
				}
			} else {
				assignedBy := ctxkeys.EmployeeIDFromContext(ctx)
				var assignedByPtr *uuid.UUID
				if assignedBy != uuid.Nil {
					assignedByPtr = &assignedBy
				}
				assigned, err := q.CreateEmployeeHandbookFromTemplate(ctx, CreateEmployeeHandbookFromTemplateParams{
					EmployeeID:           result.Employee.ID,
					TemplateID:           template.ID,
					AssignedByEmployeeID: assignedByPtr,
				})
				if err != nil {
					return err
				}

				_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, CreateEmployeeHandbookAssignmentHistoryParams{
					EmployeeHandbookID: &assigned.ID,
					EmployeeID:         result.Employee.ID,
					TemplateID:         template.ID,
					TemplateVersion:    assigned.TemplateVersion,
					Event:              HandbookAssignmentEventEnumAssigned,
					ActorEmployeeID:    assignedByPtr,
					Metadata: mustMarshalAssignmentMetadata(map[string]any{
						"source": "employee_creation",
					}),
				})
				if err != nil {
					return err
				}
			}
		}

		err = q.AssignRoleToUser(ctx, AssignRoleToUserParams{
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

func mustMarshalAssignmentMetadata(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte(`{}`)
	}
	return b
}
