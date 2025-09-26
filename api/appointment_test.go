package api

import (
	"bytes"
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/service/appointment"
	"maicare_go/token"
	"maicare_go/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createRandomAppointment(t *testing.T, employeeID int64) db.ScheduledAppointment {
	arg := db.CreateAppointmentParams{
		CreatorEmployeeID: &employeeID,
		StartTime:         pgtype.Timestamp{Time: time.Date(time.Now().Year(), time.August, 25, 12, 0, 0, 0, time.UTC), Valid: true},
		EndTime:           pgtype.Timestamp{Time: time.Date(time.Now().Year(), time.August, 25, 18, 0, 0, 0, time.UTC), Valid: true},
		Location:          util.StringPtr("Test Location"),
		Description:       util.StringPtr("Test Description"),
	}

	appointment, err := testStore.CreateAppointment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointment)
	return appointment
}

func TestCreateAppointmentApi(t *testing.T) {
	testasynqClient.EXPECT().EnqueueNotificationTask(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	employee, user := createRandomEmployee(t)
	client := createRandomClientDetails(t)

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
				appointReq := appointment.CreateAppointmentRequest{
					StartTime:              time.Now(),
					EndTime:                time.Now().Add(1 * time.Hour),
					Location:               util.StringPtr("Test Location"),
					Description:            util.StringPtr("Test Description"),
					RecurrenceType:         "NONE",
					RecurrenceInterval:     util.Int32Ptr(0),
					RecurrenceEndDate:      time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
					ParticipantEmployeeIDs: []int64{employee.ID},
					ClientIDs:              []int64{client.ID},
				}
				reqBody, err := json.Marshal(appointReq)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[appointment.CreateAppointmentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response)

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

func TestAddParticipantToAppointmentApi(t *testing.T) {
	employee, user := createRandomEmployee(t)

	apntmt := createRandomAppointment(t, employee.ID)

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
				addParticipantReq := appointment.AddParticipantToAppointmentRequest{
					ParticipantEmployeeIDs: []int64{employee.ID},
				}
				reqBody, err := json.Marshal(addParticipantReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/appointments/%s/participants", apntmt.ID.String())
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
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

func TestGetAppointmentApi(t *testing.T) {
	employee, user := createRandomEmployee(t)
	appointment := createRandomAppointment(t, employee.ID)

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
				url := fmt.Sprintf("/appointments/%s", appointment.ID.String())
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetAppointmentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response)

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

func TestUpdateAppointmentApi(t *testing.T) {
	employee, user := createRandomEmployee(t)
	client := createRandomClientDetails(t)

	appointment := createRandomAppointment(t, employee.ID)

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
				updateReq := UpdateAppointmentRequest{
					StartTime:   time.Now(),
					EndTime:     time.Now().Add(1 * time.Hour),
					Location:    util.StringPtr("Updated Location"),
					Description: util.StringPtr("Updated Description"),

					ParticipantEmployeeIDs: &[]int64{employee.ID},
					ClientIDs:              &[]int64{client.ID},
				}
				reqBody, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/appointments/%s", appointment.ID.String())
				request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateAppointmentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response)

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

func TestDeleteAppointmentApi(t *testing.T) {
	employee, user := createRandomEmployee(t)

	appointment := createRandomAppointment(t, employee.ID)

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
				url := fmt.Sprintf("/appointments/%s", appointment.ID.String())
				request, err := http.NewRequest(http.MethodDelete, url, nil)
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

func TestListAppointmentsForEmployeeApi(t *testing.T) {
	employee, user := createRandomEmployee(t)

	for i := 0; i < 5; i++ {
		createRandomAppointment(t, employee.ID)
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				req := ListAppointmentsForEmployeeInRangeRequest{
					StartDate: time.Now().Add(-24 * time.Hour),
					EndDate:   time.Now().Add(24 * time.Hour),
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				url := fmt.Sprintf("/employees/%d/appointments", employee.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
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

func TestListAppointmentsForClientApi(t *testing.T) {
	employee, user := createRandomEmployee(t)
	client := createRandomClientDetails(t)

	for i := 0; i < 5; i++ {
		createRandomAppointment(t, employee.ID)
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				req := ListAppointmentsForClientRequest{
					StartDate: time.Now().Add(-24 * time.Hour),
					EndDate:   time.Now().Add(24 * time.Hour),
				}
				reqBody, err := json.Marshal(req)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/appointments", client.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
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
