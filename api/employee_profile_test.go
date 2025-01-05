package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomEmployee(t *testing.T) (db.EmployeeProfile, *db.CustomUser) {
	// Create prerequisite records
	location := createRandomLocation(t)
	user := createRandomUser(t)

	arg := db.CreateEmployeeProfileParams{
		UserID:                    user.ID,
		FirstName:                 util.RandomString(5),
		LastName:                  util.RandomString(5),
		Position:                  util.StringPtr(util.RandomString(5)),
		Department:                util.StringPtr("IT"),
		EmployeeNumber:            util.StringPtr(util.RandomString(5)),
		EmploymentNumber:          nil,
		PrivateEmailAddress:       util.StringPtr(util.RandomString(5)),
		Email:                     util.RandomEmail(),
		AuthenticationPhoneNumber: util.StringPtr(util.RandomString(5)),
		PrivatePhoneNumber:        util.StringPtr(util.RandomString(5)),
		WorkPhoneNumber:           util.StringPtr(util.RandomString(5)),
		DateOfBirth:               pgtype.Date{Time: time.Now(), Valid: true},
		HomeTelephoneNumber:       util.StringPtr(util.RandomString(5)),
		IsSubcontractor:           util.BoolPtr(util.RandomBool()),
		Gender:                    util.StringPtr(util.RandomString(5)),
		LocationID:                util.IntPtr(location.ID),
		HasBorrowed:               false,
		OutOfService:              util.BoolPtr(util.RandomBool()),
		IsArchived:                util.RandomBool(),
	}

	employee, err := testStore.CreateEmployeeProfile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employee)

	// Verify all fields match
	require.Equal(t, arg.UserID, employee.UserID)
	require.Equal(t, arg.FirstName, employee.FirstName)
	require.Equal(t, arg.LastName, employee.LastName)
	require.Equal(t, arg.Position, employee.Position)
	require.Equal(t, arg.Department, employee.Department)
	require.Equal(t, arg.EmployeeNumber, employee.EmployeeNumber)
	require.Equal(t, arg.EmploymentNumber, employee.EmploymentNumber)
	require.Equal(t, arg.PrivateEmailAddress, employee.PrivateEmailAddress)
	require.Equal(t, arg.Email, employee.Email)
	require.Equal(t, arg.AuthenticationPhoneNumber, employee.AuthenticationPhoneNumber)
	require.Equal(t, arg.PrivatePhoneNumber, employee.PrivatePhoneNumber)
	require.Equal(t, arg.WorkPhoneNumber, employee.WorkPhoneNumber)
	require.Equal(t, arg.DateOfBirth.Time.Format("2006-01-02"), employee.DateOfBirth.Time.Format("2006-01-02"))
	require.Equal(t, arg.HomeTelephoneNumber, employee.HomeTelephoneNumber)
	require.Equal(t, arg.IsSubcontractor, employee.IsSubcontractor)
	require.Equal(t, arg.Gender, employee.Gender)
	require.Equal(t, arg.LocationID, employee.LocationID)
	require.Equal(t, arg.HasBorrowed, employee.HasBorrowed)
	require.Equal(t, arg.OutOfService, employee.OutOfService)
	require.Equal(t, arg.IsArchived, employee.IsArchived)

	// Verify auto-generated fields
	require.NotZero(t, employee.ID)
	require.NotZero(t, employee.CreatedAt)

	// Verify foreign key constraints
	require.Equal(t, util.IntPtr(location.ID), employee.LocationID)
	return employee, user
}

func TestCreateEmployeeProfileApi(t *testing.T) {
	locationID := createRandomLocation(t).ID
	userID := rand.Int63()
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				Empreq := CreateEmployeeProfileRequest{
					EmployeeNumber:            nil, // util.StringPtr(fmt.Sprintf("EMP%d", util.RandomInt(1000, 9999))),
					EmploymentNumber:          util.StringPtr(fmt.Sprintf("EN%d", util.RandomInt(10000, 99999))),
					Location:                  util.IntPtr(locationID),
					IsSubcontractor:           util.BoolPtr(util.RandomBool()),
					FirstName:                 util.RandomString(6),
					LastName:                  util.RandomString(8),
					DateOfBirth:               util.StringPtr("2000-01-05"),
					Gender:                    util.StringPtr("male"),
					Email:                     util.RandomEmail(),
					PrivateEmailAddress:       util.StringPtr(util.RandomEmail()),
					AuthenticationPhoneNumber: util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					WorkPhoneNumber:           util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					PrivatePhoneNumber:        util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					HomeTelephoneNumber:       util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					OutOfService:              util.BoolPtr(util.RandomBool()),
				}
				data, err := json.Marshal(Empreq)
				require.NoError(t, err)

				req, err := http.NewRequest(http.MethodPost, "/employee/employees_create/", bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response CreateEmployeeProfileResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.Empty(t, response.EmployeeNumber)
				require.NotEmpty(t, response.EmploymentNumber)
				require.NotEmpty(t, response.LocationID)
				require.NotEmpty(t, response.FirstName)
				require.NotEmpty(t, response.LastName)
				require.NotEmpty(t, response.DateOfBirth)
				require.NotEmpty(t, response.Gender)
				require.NotEmpty(t, response.Email)
				require.NotEmpty(t, response.PrivateEmailAddress)
				require.NotEmpty(t, response.AuthenticationPhoneNumber)
				require.NotEmpty(t, response.WorkPhoneNumber)
				require.NotEmpty(t, response.PrivatePhoneNumber)
				require.NotEmpty(t, response.HomeTelephoneNumber)
				require.NotEmpty(t, response.UserID)
				require.NotEmpty(t, response.Created)
			},
		},
		// Add more test cases as needed
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

func TestListEmployeeProfileApi(t *testing.T) {
	user := createRandomUser(t)
	initialCount, err := testStore.CountEmployeeProfile(context.Background(), db.CountEmployeeProfileParams{
		IncludeArchived:     util.BoolPtr(true),
		IncludeOutOfService: util.BoolPtr(true),
	})
	require.NoError(t, err)
	numEmployees := 20
	var wg sync.WaitGroup
	for i := 0; i < numEmployees; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomEmployee(t)
		}()
	}
	wg.Wait()

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
				url := "/employee/employees_list/?page=1&page_size=10"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response pagination.Response[db.ListEmployeeProfileRow]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotNil(t, response.Next)
				require.Nil(t, response.Previous)
				require.Equal(t, int64(numEmployees)+initialCount, response.Count)
				require.Equal(t, int32(10), response.PageSize)
				require.Len(t, response.Results, 10)

				// Check results are ordered by created DESC
				for i := 1; i < len(response.Results); i++ {
					require.True(t, response.Results[i-1].CreatedAt.Time.After(response.Results[i].CreatedAt.Time) ||
						response.Results[i-1].CreatedAt.Time.Equal(response.Results[i].CreatedAt.Time))
				}
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Don't add authorization
			},
			buildRequest: func() (*http.Request, error) {
				url := "/employee/employees_list/?page=1&page_size=10"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, -time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/employee/employees_list/?page=1&page_size=10"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/employee/employees_list/?page=1&page_size=1000"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FilterByDepartment",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/employee/employees_list/?page=1&page_size=10&department=IT"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response pagination.Response[db.ListEmployeeProfileRow]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				for _, emp := range response.Results {
					require.Equal(t, util.StringPtr("IT"), emp.Department)
				}
			},
		},
		{
			name: "SecondPage",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/employee/employees_list/?page=2&page_size=10"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response pagination.Response[db.ListEmployeeProfileRow]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotNil(t, response.Previous)
				require.Contains(t, *response.Previous, "page=1")
			},
		},
	}

	for _, tc := range testCases {
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

func TestGetEmployeeProfileApi(t *testing.T) {
	employee, user := createRandomEmployee(t)
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
				req, err := http.NewRequest(http.MethodGet, "/employees/profile", nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[GetEmployeeProfileResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.Equal(t, employee.ID, response.Data.EmployeeID)
				require.Equal(t, employee.UserID, response.Data.UserID)
				require.Equal(t, employee.FirstName, response.Data.FirstName)
				require.Equal(t, employee.LastName, response.Data.LastName)

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
