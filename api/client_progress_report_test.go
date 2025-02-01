package api

import (
	"bytes"
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
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

func createRandomProgressReport(t *testing.T, clientID int64, employeeID int64) db.ProgressReport {
	arg := db.CreateProgressReportParams{
		ClientID:       clientID,
		EmployeeID:     &employeeID,
		Title:          util.StringPtr("Test Progress Report"),
		Date:           pgtype.Timestamptz{Time: util.RandomTIme(), Valid: true},
		ReportText:     "Test Progress Report",
		Type:           "morning_report",
		EmotionalState: "normal",
	}
	progressReport, err := testStore.CreateProgressReport(context.Background(), arg)
	require.NoError(t, err)
	return progressReport
}

func TestCreateProgressReportApi(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)

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
				createReq := CreateProgressReportRequest{
					EmployeeID:     &employee.ID,
					Title:          util.StringPtr("Test Progress Report"),
					Date:           util.RandomTIme(),
					ReportText:     "Test Progress Report",
					Type:           "morning_report",
					EmotionalState: "normal",
				}
				reqBody, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/progress_reports", client.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateProgressReportResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, client.ID, response.Data.ClientID)
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

func TestListProgressReportApi(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	for i := 0; i < 10; i++ {
		createRandomProgressReport(t, client.ID, employee.ID)
	}

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
				url := fmt.Sprintf("/clients/%d/progress_reports?page=1&page_size=10", client.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[pagination.Response[ListProgressReportsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.Results)
				require.Len(t, response.Data.Results, 10)
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

func TestGetProgressReportApi(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	progressReport1 := createRandomProgressReport(t, client.ID, employee.ID)

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
				url := fmt.Sprintf("/clients/%d/progress_reports/%d", client.ID, progressReport1.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetProgressReportResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, progressReport1.ID, response.Data.ID)
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

func TestUpdateProgressReportApi(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	progressReport1 := createRandomProgressReport(t, client.ID, employee.ID)

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
				updateReq := UpdateProgressReportRequest{
					ReportText:     util.StringPtr("Updated Progress Report"),
					EmotionalState: util.StringPtr("happy"),
				}
				reqBody, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/progress_reports/%d", client.ID, progressReport1.ID)
				request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateProgressReportResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, progressReport1.ID, response.Data.ID)
				require.Equal(t, "Updated Progress Report", response.Data.ReportText)
				require.Equal(t, "happy", response.Data.EmotionalState)
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
