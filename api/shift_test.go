package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/service/organization"
	"maicare_go/token"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomShift(t *testing.T, locationID uuid.UUID) db.LocationShift {
	arg := db.CreateShiftParams{
		LocationID: locationID,
		ShiftName:  "Slaapdienst of Waakdienst",
		StartTime:  pgtype.Time{Microseconds: 23 * 3600 * 1000000, Valid: true}, // 23:00:00
		EndTime:    pgtype.Time{Microseconds: 7 * 3600 * 1000000, Valid: true},  // 07:00:00
	}
	shift, err := testStore.CreateShift(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, shift)
	require.Equal(t, arg.LocationID, shift.LocationID)
	return shift
}

func TestCreateShiftsApi(t *testing.T) {
	_, user := createRandomEmployee(t)
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
				createShiftReq := organization.CreateShiftApiRequest{
					ShiftName: "Morning Shift",
					StartTime: "08:00",
					EndTime:   "16:00",
				}
				data, err := json.Marshal(createShiftReq)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/locations/%d/shifts", location.ID), bytes.NewBuffer(data))
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
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

func TestGetShiftsByLocationApi(t *testing.T) {
	_, user := createRandomEmployee(t)
	location := createRandomLocation(t)
	createRandomShift(t, location.ID)

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
				request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/locations/%d/shifts", location.ID), nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
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

func TestDeleteShiftApi(t *testing.T) {
	_, user := createRandomEmployee(t)
	location := createRandomLocation(t)
	shift := createRandomShift(t, location.ID)

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
				request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/locations/%d/shifts/%d", location.ID, shift.ID), nil)
				require.NoError(t, err)
				return request, nil
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

func TestUpdateShiftApi(t *testing.T) {
	_, user := createRandomEmployee(t)
	location := createRandomLocation(t)
	shift := createRandomShift(t, location.ID)

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
				updateShiftReq := organization.UpdateShiftApiRequest{
					ShiftName: "Updated Shift",
					StartTime: "09:00",
					EndTime:   "17:00",
				}
				data, err := json.Marshal(updateShiftReq)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/locations/%d/shifts/%d", location.ID, shift.ID), bytes.NewBuffer(data))
				require.NoError(t, err)
				return request, nil
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
