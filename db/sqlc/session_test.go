package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	user := CreateRandomUser(t)

	uuid, err := uuid.NewRandom()
	require.NoError(t, err)

	// Get current time for timestamps
	now := time.Now()
	expireTime := now.Add(24 * time.Hour) // Session expires in 24 hours

	arg := CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: uuid,
			Valid: true,
		},
		RefreshToken: util.RandomString(16),
		UserAgent:    util.RandomString(5),
		ClientIp:     util.RandomString(5),
		IsBlocked:    false,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expireTime,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
		UserID: user.ID,
	}

	// Create the session
	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	// Verify all fields match
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.Equal(t, arg.UserID, session.UserID)

	// Verify timestamps
	require.WithinDuration(t, arg.ExpiresAt.Time, session.ExpiresAt.Time, time.Second)
	require.WithinDuration(t, arg.CreatedAt.Time, session.CreatedAt.Time, time.Second)

	// Verify session was created with correct user
	require.Equal(t, user.ID, session.UserID)
}
