package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name   string
		params CreateUserParams
		checks func(t *testing.T, user CustomUser)
	}{
		{
			name: "successful creation",
			params: CreateUserParams{
				Password:       "hashedpassword123",
				Email:          "test@example.com",
				IsActive:       true,
				ProfilePicture: nil,
			},
			checks: func(t *testing.T, user CustomUser) {
				require.Equal(t, "test@example.com", user.Email)
				require.True(t, user.IsActive)
				require.Nil(t, user.ProfilePicture)
			},
		},
		{
			name: "with profile picture",
			params: CreateUserParams{
				Password:       "hashedpassword123",
				Email:          "pic@example.com",
				IsActive:       true,
				ProfilePicture: util.StringPtr("https://example.com/pic.jpg"),
			},
			checks: func(t *testing.T, user CustomUser) {
				require.NotNil(t, user.ProfilePicture)
				require.Equal(t, "https://example.com/pic.jpg", *user.ProfilePicture)
			},
		},
		{
			name: "user inactive by default",
			params: CreateUserParams{
				Password:       "hashedpassword123",
				Email:          "inactive@example.com",
				IsActive:       false,
				ProfilePicture: nil,
			},
			checks: func(t *testing.T, user CustomUser) {
				require.False(t, user.IsActive)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			user, err := qtx.CreateUser(ctx, tt.params)
			require.NoError(t, err, "CreateUser() should not error")

			tt.checks(t, user)
		})
	}
}

func TestCreateTemp2FaSecret(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateTemp2FaSecretParams
		checks func(t *testing.T, rowsAffected int64, err error)
	}{
		{
			name: "successful creation",
			setup: func(ctx context.Context, qtx *Queries) CreateTemp2FaSecretParams {
				user := createRandomUser(ctx, qtx)
				return CreateTemp2FaSecretParams{
					ID:                  user.ID,
					TwoFactorSecretTemp: util.StringPtr("tempsecret123"),
				}
			},
			checks: func(t *testing.T, rowsAffected int64, err error) {
				require.NoError(t, err, "CreateTemp2FaSecret() should not error")
				require.Equal(t, int64(1), rowsAffected, "should affect exactly 1 row")
			},
		},
		{
			name: "invalid user ID",
			setup: func(ctx context.Context, qtx *Queries) CreateTemp2FaSecretParams {
				return CreateTemp2FaSecretParams{
					ID:                  uuid.New(),
					TwoFactorSecretTemp: util.StringPtr("tempsecret123"),
				}
			},
			checks: func(t *testing.T, rowsAffected int64, err error) {
				require.NoError(t, err, "CreateTemp2FaSecret() should not error")
				require.Equal(t, int64(0), rowsAffected, "should affect 0 rows for non-existent user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			rowsAffected, err := qtx.CreateTemp2FaSecret(ctx, params)
			tt.checks(t, rowsAffected, err)
		})
	}
}

func TestEnable2Fa(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) Enable2FaParams
		checks func(t *testing.T, rowsAffected int64, err error)
	}{
		{
			name: "successful 2FA enablement",
			setup: func(ctx context.Context, qtx *Queries) Enable2FaParams {
				user := createRandomUser(ctx, qtx)
				return Enable2FaParams{
					ID:              user.ID,
					TwoFactorSecret: util.StringPtr("secret123"),
					RecoveryCodes:   []string{"code1", "code2", "code3"},
				}
			},
			checks: func(t *testing.T, rowsAffected int64, err error) {
				require.NoError(t, err, "Enable2Fa() should not error")
				require.Equal(t, int64(1), rowsAffected, "should affect exactly 1 row")
			},
		},
		{
			name: "enable 2FA with empty recovery codes",
			setup: func(ctx context.Context, qtx *Queries) Enable2FaParams {
				user := createRandomUser(ctx, qtx)
				return Enable2FaParams{
					ID:              user.ID,
					TwoFactorSecret: util.StringPtr("secret456"),
					RecoveryCodes:   []string{},
				}
			},
			checks: func(t *testing.T, rowsAffected int64, err error) {
				require.NoError(t, err, "Enable2Fa() should not error with empty recovery codes")
				require.Equal(t, int64(1), rowsAffected, "should affect exactly 1 row")
			},
		},
		{
			name: "invalid user ID",
			setup: func(ctx context.Context, qtx *Queries) Enable2FaParams {
				return Enable2FaParams{
					ID:              uuid.Nil,
					TwoFactorSecret: util.StringPtr("secret789"),
					RecoveryCodes:   []string{"code1"},
				}
			},
			checks: func(t *testing.T, rowsAffected int64, err error) {
				require.Equal(t, int64(0), rowsAffected, "should affect 0 rows for invalid user ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			rowsAffected, err := qtx.Enable2Fa(ctx, params)
			tt.checks(t, rowsAffected, err)
		})
	}
}

func TestGetTemp2FaSecret(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, secret *string, err error)
	}{
		{
			name: "get existing temp 2FA secret",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				user := createRandomUser(ctx, qtx)
				_, err := qtx.CreateTemp2FaSecret(ctx, CreateTemp2FaSecretParams{
					ID:                  user.ID,
					TwoFactorSecretTemp: util.StringPtr("tempsecret123"),
				})
				require.NoError(t, err)
				return user.ID
			},
			checks: func(t *testing.T, secret *string, err error) {
				require.NoError(t, err, "GetTemp2FaSecret() should not error")
				require.NotNil(t, secret)
				require.Equal(t, "tempsecret123", *secret)
			},
		},
		{
			name: "get temp 2FA secret when not set",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				user := createRandomUser(ctx, qtx)
				return user.ID
			},
			checks: func(t *testing.T, secret *string, err error) {
				require.NoError(t, err, "GetTemp2FaSecret() should not error")
				require.Nil(t, secret)
			},
		},
		{
			name: "get temp 2FA secret for non-existent user",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, secret *string, err error) {
				require.Error(t, err, "GetTemp2FaSecret() should error for non-existent user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			userID := tt.setup(ctx, qtx)
			secret, err := qtx.GetTemp2FaSecret(ctx, userID)
			tt.checks(t, secret, err)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) string
		checks func(t *testing.T, user GetUserByEmailRow, err error)
	}{
		{
			name: "get existing user by email",
			setup: func(ctx context.Context, qtx *Queries) string {
				email := util.RandomEmail()
				_, err := qtx.CreateUser(ctx, CreateUserParams{
					Password:       "hashedpassword",
					Email:          email,
					IsActive:       true,
					ProfilePicture: nil,
				})
				require.NoError(t, err)
				return email
			},
			checks: func(t *testing.T, user GetUserByEmailRow, err error) {
				require.NoError(t, err, "GetUserByEmail() should not error")
				require.NotEmpty(t, user.ID)
			},
		},
		{
			name: "get non-existent user by email",
			setup: func(ctx context.Context, qtx *Queries) string {
				return util.RandomEmail()
			},
			checks: func(t *testing.T, user GetUserByEmailRow, err error) {
				require.Error(t, err, "GetUserByEmail() should error for non-existent email")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			email := tt.setup(ctx, qtx)
			user, err := qtx.GetUserByEmail(ctx, email)
			tt.checks(t, user, err)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, user GetUserByIDRow, err error)
	}{
		{
			name: "get existing user by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				user := createRandomUser(ctx, qtx)
				return user.ID
			},
			checks: func(t *testing.T, user GetUserByIDRow, err error) {
				require.NoError(t, err, "GetUserByID() should not error")
				require.NotEmpty(t, user.ID)
				require.True(t, user.IsActive)
			},
		},
		{
			name: "get non-existent user by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, user GetUserByIDRow, err error) {
				require.Error(t, err, "GetUserByID() should error for non-existent ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			userID := tt.setup(ctx, qtx)
			user, err := qtx.GetUserByID(ctx, userID)
			tt.checks(t, user, err)
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdatePasswordParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "successful password update",
			setup: func(ctx context.Context, qtx *Queries) UpdatePasswordParams {
				user := createRandomUser(ctx, qtx)
				return UpdatePasswordParams{
					ID:       user.ID,
					Password: "newhashed password123",
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdatePassword() should not error")
			},
		},
		{
			name: "update password for non-existent user",
			setup: func(ctx context.Context, qtx *Queries) UpdatePasswordParams {
				return UpdatePasswordParams{
					ID:       uuid.New(),
					Password: "newhashed password456",
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdatePassword() should not error even for non-existent user")
			},
		},
		{
			name: "update password with empty string",
			setup: func(ctx context.Context, qtx *Queries) UpdatePasswordParams {
				user := createRandomUser(ctx, qtx)
				return UpdatePasswordParams{
					ID:       user.ID,
					Password: "",
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdatePassword() should not error with empty password")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			err = qtx.UpdatePassword(ctx, params)
			tt.checks(t, err)
		})
	}
}

// Random data generators for tests

func createRandomUser(ctx context.Context, qtx *Queries) CustomUser {
	user, err := qtx.CreateUser(ctx, CreateUserParams{
		Password:       "randomhashedpassword",
		Email:          util.RandomEmail(),
		IsActive:       true,
		ProfilePicture: nil,
	})
	if err != nil {
		panic("failed to create random user: " + err.Error())
	}
	return user
}
