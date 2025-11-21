package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateSessionParams
		checks func(t *testing.T, session Session, err error)
	}{
		{
			name: "successful session creation",
			setup: func(ctx context.Context, qtx *Queries) CreateSessionParams {
				user := createRandomUser(ctx, qtx)
				sessionID, _ := uuid.NewRandom()
				return CreateSessionParams{
					ID:           sessionID,
					RefreshToken: util.RandomString(16),
					UserAgent:    util.RandomString(10),
					ClientIp:     util.RandomString(10),
					IsBlocked:    false,
					ExpiresAt: pgtype.Timestamptz{
						Time:  time.Now().Add(24 * time.Hour),
						Valid: true,
					},
					CreatedAt: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
					UserID: user.ID,
				}
			},
			checks: func(t *testing.T, session Session, err error) {
				require.NoError(t, err, "CreateSession should not return an error")
				require.NotZero(t, session.ID)
				require.False(t, session.IsBlocked)
			},
		},
		{
			name: "session creation with blocked",
			setup: func(ctx context.Context, qtx *Queries) CreateSessionParams {
				user := createRandomUser(ctx, qtx)
				sessionID, _ := uuid.NewRandom()
				return CreateSessionParams{
					ID:           sessionID,
					RefreshToken: util.RandomString(16),
					UserAgent:    util.RandomString(10),
					ClientIp:     util.RandomString(10),
					IsBlocked:    true,
					ExpiresAt: pgtype.Timestamptz{
						Time:  time.Now().Add(24 * time.Hour),
						Valid: true,
					},
					CreatedAt: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
					UserID: user.ID,
				}
			},
			checks: func(t *testing.T, session Session, err error) {
				require.NoError(t, err, "CreateSession should not return an error")
				require.True(t, session.IsBlocked)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			session, err := qtx.CreateSession(ctx, params)
			tt.checks(t, session, err)
		})
	}
}

func TestGetSessionByID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, session Session, err error)
	}{
		{
			name: "get existing session",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				session := createRandomSession(ctx, qtx)
				return session.ID
			},
			checks: func(t *testing.T, session Session, err error) {
				require.NoError(t, err, "GetSessionByID should not error")
				require.NotZero(t, session.ID)
			},
		},
		{
			name: "get non-existent session",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, session Session, err error) {
				require.Error(t, err, "GetSessionByID should error for non-existent session")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			sessionID := tt.setup(ctx, qtx)
			session, err := qtx.GetSessionByID(ctx, sessionID)
			tt.checks(t, session, err)
		})
	}
}

func TestDeleteSession(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete existing session",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				session := createRandomSession(ctx, qtx)
				return session.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteSession should not error")
			},
		},
		{
			name: "delete non-existent session",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteSession should not error for non-existent session")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			sessionID := tt.setup(ctx, qtx)
			err = qtx.DeleteSession(ctx, sessionID)
			tt.checks(t, err)
		})
	}
}

// Helpers
func createRandomSession(ctx context.Context, qtx *Queries) Session {
	user := createRandomUser(ctx, qtx)
	sessionID, _ := uuid.NewRandom()
	params := CreateSessionParams{
		ID:           sessionID,
		RefreshToken: util.RandomString(16),
		UserAgent:    util.RandomString(10),
		ClientIp:     util.RandomString(10),
		IsBlocked:    false,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		UserID: user.ID,
	}
	session, err := qtx.CreateSession(ctx, params)
	if err != nil {
		panic(err)
	}
	return session
}
