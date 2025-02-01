package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"

	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomLocation(t *testing.T) *db.Location {
	arg := db.CreateLocationParams{
		Name:     util.RandomString(5),
		Address:  util.RandomString(8),
		Capacity: util.Int32Ptr(52),
	}

	location, err := testStore.CreateLocation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, location)

	// Check if the returned location matches the input
	require.Equal(t, arg.Name, location.Name)
	require.Equal(t, arg.Address, location.Address)
	require.Equal(t, arg.Capacity, location.Capacity)

	// Verify ID is generated
	require.NotZero(t, location.ID)
	return &location
}

func TestCreateLocationApi(t *testing.T) {
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
				locationReq := CreateLocationRequest{
					Name:     "Test Location",
					Address:  "Test Address",
					Capacity: util.Int32Ptr(52),
				}
				reqBody, err := json.Marshal(locationReq)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPost, "/locations", bytes.NewReader(reqBody))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var locationRes Response[CreateLocationResponse]
				err := json.NewDecoder(recorder.Body).Decode(&locationRes)
				require.NoError(t, err)
				require.NotEmpty(t, locationRes.Data)
				require.NotEmpty(t, locationRes.Data.ID)
				require.NotEmpty(t, locationRes.Data.Name)
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

func TestUpdateLocationApi(t *testing.T) {
	user := createRandomUser(t)
	location := createRandomLocation(t)
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
				locationReq := UpdateLocationRequest{
					Name: util.StringPtr("Updated Name"),
				}
				reqBody, err := json.Marshal(locationReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/locations/%d", location.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var locationRes Response[UpdateLocationResponse]
				err := json.NewDecoder(recorder.Body).Decode(&locationRes)
				require.NoError(t, err)
				require.NotEmpty(t, locationRes.Data)
				require.NotEmpty(t, locationRes.Data.ID)
				require.NotEmpty(t, locationRes.Data.Name)
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

func TestDeleteLocationApi(t *testing.T) {
	location := createRandomLocation(t)
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
				url := fmt.Sprintf("/locations/%d", location.ID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var locationRes Response[DeleteLocationResponse]
				err := json.NewDecoder(recorder.Body).Decode(&locationRes)
				require.NoError(t, err)
				require.NotEmpty(t, locationRes.Data)
				require.NotEmpty(t, locationRes.Data.ID)
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

func TestGetLocationApi(t *testing.T) {
	location := createRandomLocation(t)
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
				url := fmt.Sprintf("/locations/%d", location.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var locationRes Response[GetLocationResponse]
				err := json.NewDecoder(recorder.Body).Decode(&locationRes)
				require.NoError(t, err)
				require.NotEmpty(t, locationRes.Data)
				require.NotEmpty(t, locationRes.Data.ID)
				require.NotEmpty(t, locationRes.Data.Name)
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
