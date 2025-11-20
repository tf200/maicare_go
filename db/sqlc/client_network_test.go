package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestAssignEmployee(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AssignEmployeeParams
		checks func(t *testing.T, row AssignEmployeeRow, err error)
	}{
		{
			name: "successful creation",
			setup: func(ctx context.Context, qtx *Queries) AssignEmployeeParams {
				employee := createRandomEmployee(ctx, qtx)
				client := createRandomClientDetails(ctx, qtx)
				return AssignEmployeeParams{
					ClientID:   client.ID,
					EmployeeID: employee.ID,
					StartDate:  pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					Role:       "Manger",
				}
			},
			checks: func(t *testing.T, row AssignEmployeeRow, err error) {
				require.NoError(t, err, "AssignEmployee should not return an error")
				require.NotZero(t, row.ID, "Assigned employee ID should not be zero")
				require.Equal(t, row.ClientID, row.ClientID, "ClientID should match")
				require.Equal(t, row.EmployeeID, row.EmployeeID, "EmployeeID should match")
				require.Equal(t, row.Role, "Manger", "Role should match")
				require.True(t, row.StartDate.Valid, "StartDate should be valid")
				require.Equal(t, row.StartDate.Time, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "StartDate should match")
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
			row, err := qtx.AssignEmployee(ctx, params)
			tt.checks(t, row, err)
		})
	}
}

func TestAssignSender(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AssignSenderParams
		checks func(t *testing.T, row ClientDetail, err error)
	}{
		{
			name: "successful sender assignment",
			setup: func(ctx context.Context, qtx *Queries) AssignSenderParams {
				client := createRandomClientDetails(ctx, qtx)
				sender := createRandomSenders(ctx,)
				return AssignSenderParams{
					ClientID: client.ID,
					SenderID: &sender.ID,
				}
			},
			checks: func(t *testing.T, row ClientDetail, err error) {
				require.NoError(t, err, "AssignSender should not return an error")
				require.Equal(t, row.SenderID, row.SenderID, "SenderID should match")
			},
		},
		{
			name: "successful unassignment of sender",
			setup: func(ctx context.Context, qtx *Queries) AssignSenderParams {
				client := createRandomClientDetails(ctx, qtx)
				return AssignSenderParams{
					ID:       client.ID,
					SenderID: nil,
				}
			},
			checks: func(t *testing.T, row ClientDetail, err error) {
				require.NoError(t, err, "AssignSender should not return an error")
				require.False(t, row.SenderID.Valid, "SenderID should be null")
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
			row, err := qtx.AssignSender(ctx, params)
			tt.checks(t, row, err)
		})
	}
}
