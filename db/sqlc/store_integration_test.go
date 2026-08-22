package db

import (
	"context"
	"errors"
	"os"
	"testing"

	"maicare_go/internal/ctxkeys"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestExecActorTxIdentityIsTransactionLocal(t *testing.T) {
	store := openIntegrationStore(t)
	actorA := ctxkeys.ActorIdentity{UserID: uuid.New(), EmployeeID: uuid.New()}
	actorB := ctxkeys.ActorIdentity{UserID: uuid.New(), EmployeeID: uuid.New()}

	assertActorSettings(t, store, actorA)
	assertNoActorSettings(t, store)
	assertActorSettings(t, store, actorB)
	assertNoActorSettings(t, store)
}

func TestExecActorTxRejectsMissingIdentity(t *testing.T) {
	store := openIntegrationStore(t)
	err := store.ExecActorTx(context.Background(), func(*Queries) error {
		t.Fatal("transaction callback must not run without actor identity")
		return nil
	})
	if !errors.Is(err, ErrMissingActorIdentity) {
		t.Fatalf("error = %v, want ErrMissingActorIdentity", err)
	}
}

func TestExecActorTxIdentityDoesNotLeakAfterRollback(t *testing.T) {
	store := openIntegrationStore(t)
	actor := ctxkeys.ActorIdentity{UserID: uuid.New(), EmployeeID: uuid.New()}
	ctx := ctxkeys.WithActorIdentity(context.Background(), actor)
	wantErr := errors.New("force rollback")

	err := store.ExecActorTx(ctx, func(*Queries) error { return wantErr })
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	assertNoActorSettings(t, store)
}

func assertActorSettings(t *testing.T, store *Store, actor ctxkeys.ActorIdentity) {
	t.Helper()
	ctx := ctxkeys.WithActorIdentity(context.Background(), actor)
	err := store.ExecActorTx(ctx, func(q *Queries) error {
		var userID, employeeID string
		if err := q.db.QueryRow(ctx, `SELECT
			current_setting('myapp.current_user_id', true),
			current_setting('myapp.current_employee_id', true)`).Scan(&userID, &employeeID); err != nil {
			return err
		}
		if userID != actor.UserID.String() || employeeID != actor.EmployeeID.String() {
			t.Fatalf("settings = %q, %q; want %q, %q", userID, employeeID, actor.UserID, actor.EmployeeID)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ExecActorTx() error = %v", err)
	}
}

func assertNoActorSettings(t *testing.T, store *Store) {
	t.Helper()
	var userID, employeeID string
	err := store.ConnPool.QueryRow(context.Background(), `SELECT
		current_setting('myapp.current_user_id', true),
		current_setting('myapp.current_employee_id', true)`).Scan(&userID, &employeeID)
	if err != nil {
		t.Fatalf("read settings: %v", err)
	}
	if userID != "" || employeeID != "" {
		t.Fatalf("actor identity leaked: user_id=%q employee_id=%q", userID, employeeID)
	}
}

func openIntegrationStore(t *testing.T) *Store {
	t.Helper()
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL: %v", err)
	}
	config.MaxConns = 1
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return RegisterEnumTypes(ctx, conn)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatalf("connect TEST_DATABASE_URL: %v", err)
	}
	t.Cleanup(pool.Close)
	return NewStore(pool)
}
