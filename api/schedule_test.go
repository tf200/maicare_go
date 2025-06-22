package api

import (
	"bytes"
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomSchedule(t *testing.T, employeeID int64) db.Schedule {
	location := createRandomLocation(t)

	arg := db.CreateScheduleParams{
		EmployeeID:    employeeID,
		LocationID:    location.ID,
		StartDatetime: pgtype.Timestamp{Time: time.Now(), Valid: true},
		EndDatetime:   pgtype.Timestamp{Time: time.Now().Add(24 * time.Hour), Valid: true},
	}
	schedule, err := testStore.CreateSchedule(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, schedule)
	require.Equal(t, arg.EmployeeID, schedule.EmployeeID)
	return schedule
}

func TestCreateScheduleApi(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	location := createRandomLocation(t)
	// shift := createRandomShift(t, location.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK IS CUSTOM",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				now := time.Now()
				endTime := time.Now().Add(1 * time.Hour)
				createScheduleReq := CreateScheduleRequest{
					EmployeeID:    employee.ID,
					LocationID:    location.ID,
					IsCustom:      true,
					StartDatetime: &now,
					EndDatetime:   &endTime,
				}
				data, err := json.Marshal(createScheduleReq)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPost, "/schedules", bytes.NewBuffer(data))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[CreateScheduleResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, "Schedule created successfully", response.Message)
			},
		},
		{
			name: "OK IS NOT CUSTOM",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				now := time.Now().Format("2006-01-02")
				createScheduleReq := CreateScheduleRequest{
					EmployeeID:      2,
					LocationID:      3,
					IsCustom:        false,
					StartDatetime:   nil,
					EndDatetime:     nil,
					LocationShiftID: util.IntPtr(1),
					ShiftDate:       &now,
				}
				data, err := json.Marshal(createScheduleReq)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPost, "/schedules", bytes.NewBuffer(data))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[CreateScheduleResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, "Schedule created successfully", response.Message)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}

}

func TestGetMonthlySchedulesByLocationApi(t *testing.T) {

	employee, _ := createRandomEmployee(t)
	schedule := createRandomSchedule(t, employee.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/locations/%d/monthly_schedules?year=%d&month=%d", schedule.LocationID, schedule.StartDatetime.Time.Year(), schedule.StartDatetime.Time.Month())
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[[]GetMonthlySchedulesByLocationResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetDailySchedulesByLocationApi(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	schedule := createRandomSchedule(t, employee.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/locations/%d/daily_schedules?year=%d&month=%d&day=%d", schedule.LocationID,
					schedule.StartDatetime.Time.Year(),
					schedule.StartDatetime.Time.Month(),
					schedule.StartDatetime.Time.Day())
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetDailySchedulesByLocationResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
