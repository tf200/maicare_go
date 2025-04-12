package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"

	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *db.CustomUser {
	hashedPassword, err := util.HashPassword("t2aha000")
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Password:       hashedPassword,
		Email:          util.RandomEmail(),
		IsActive:       true,
		ProfilePicture: util.StringPtr(util.GetRandomImageURL()),
		RoleID:         1,
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)

	require.Equal(t, arg.IsActive, user.IsActive)
	require.Equal(t, arg.ProfilePicture, user.ProfilePicture)

	require.NotZero(t, user.ID)

	return &user
}

func TestLogin(t *testing.T) {
	user := createRandomUser(t)

	testCases := []struct {
		name          string
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildRequest: func() (*http.Request, error) {
				loginReq := LoginUserRequest{
					Email:    user.Email,
					Password: "t2aha000",
				}
				data, err := json.Marshal(loginReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response LoginUserResponse
				res := SuccessResponse(response, "login successful")
				err := json.NewDecoder(recorder.Body).Decode(&res)

				require.NoError(t, err)
				require.NotEmpty(t, res.Data.AccessToken)
				require.NotEmpty(t, res.Data.RefreshToken)
			},
		},
		{
			name: "UserNotFound",
			buildRequest: func() (*http.Request, error) {
				loginReq := LoginUserRequest{
					Email:    "nonexistent@email.com",
					Password: "password123",
				}
				data, err := json.Marshal(loginReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "WrongPassword",
			buildRequest: func() (*http.Request, error) {
				loginReq := LoginUserRequest{
					Email:    user.Email,
					Password: "wrongpassword",
				}
				data, err := json.Marshal(loginReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomSession(t *testing.T, token string, payload *token.Payload) db.Session {
	user := createRandomUser(t)

	// Get current time for timestamps
	now := time.Now()
	expireTime := now.Add(24 * time.Hour) // Session expires in 24 hours

	arg := db.CreateSessionParams{
		ID:           payload.ID,
		RefreshToken: token,
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
	session, err := testStore.CreateSession(context.Background(), arg)
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
	return session
}

func TestRefreshTokenHandler(t *testing.T) {
	token, payload, err := testServer.tokenMaker.CreateToken(1, 1, testServer.config.RefreshTokenDuration, token.RefreshToken)
	createRandomSession(t, token, payload)
	require.NoError(t, err)
	testCases := []struct {
		name          string
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildRequest: func() (*http.Request, error) {
				RefreshReq := RefreshTokenRequest{
					Token: token,
				}
				data, err := json.Marshal(RefreshReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body)
				require.Equal(t, http.StatusOK, recorder.Code)

				var response RefreshTokenResponse
				res := SuccessResponse(response, "Refresh token")
				err := json.NewDecoder(recorder.Body).Decode(&res)
				require.NoError(t, err)
				require.NotEmpty(t, res.Data.AccessToken)
			},
		},
		{
			name: "InvalidToken",
			buildRequest: func() (*http.Request, error) {
				RefreshReq := RefreshTokenRequest{
					Token: "invalid",
				}
				data, err := json.Marshal(RefreshReq)
				if err != nil {
					return nil, err
				}
				req, err := http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestChangePasswordApi(t *testing.T) {
	user := createRandomUser(t)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/auth/change_password"
				changePasswordReq := ChangePasswordRequest{
					OldPassword: "t2aha000",
					NewPassword: "newpassword123",
				}
				data, err := json.Marshal(changePasswordReq)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request, err := tc.buildRequest()
			require.NoError(t, err)

			tc.setupAuth(t, request, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}
