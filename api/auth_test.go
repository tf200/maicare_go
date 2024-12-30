package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *db.CustomUser {
	hashedPassword, err := util.HashPassword("t2aha000")
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Password: hashedPassword,
		Username: util.StringPtr(util.RandomString(9)),
		// Username:    "taha",
		Email:          util.RandomEmail(),
		FirstName:      util.RandomString(5),
		LastName:       util.RandomString(5),
		IsSuperuser:    true,
		IsStaff:        true,
		IsActive:       true,
		ProfilePicture: util.StringPtr(util.GetRandomImageURL()),
		PhoneNumber:    util.IntPtr(5862),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.IsSuperuser, user.IsSuperuser)
	require.Equal(t, arg.IsStaff, user.IsStaff)
	require.Equal(t, arg.IsActive, user.IsActive)
	require.Equal(t, arg.ProfilePicture, user.ProfilePicture)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)

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

func TestRefreshTokenHandler(t *testing.T) {
	token, _, err := testServer.tokenMaker.CreateToken(1, testServer.config.RefreshTokenDuration, token.RefreshToken)
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
