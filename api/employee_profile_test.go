package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/token"
	"maicare_go/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomLocation(t *testing.T) db.Location {
	arg := db.CreateLocationParams{
		Name:    util.RandomString(5),
		Address: util.RandomString(8),
		Capacity: pgtype.Int4{
			Int32: 25,
			Valid: true,
		},
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
	return location
}

func createRandomEmployee(t *testing.T) db.EmployeeProfile {
	// Create prerequisite records
	location := createRandomLocation(t)
	user := createRandomUser(t)

	arg := db.CreateEmployeeProfileParams{
		UserID:                    user.ID,
		FirstName:                 util.RandomString(5),
		LastName:                  util.RandomString(5),
		Position:                  util.RandomPgText(),
		Department:                pgtype.Text{String: "IT", Valid: true},
		EmployeeNumber:            util.RandomPgText(),
		EmploymentNumber:          util.RandomPgText(),
		PrivateEmailAddress:       util.RandomPgText(),
		EmailAddress:              util.RandomPgText(),
		AuthenticationPhoneNumber: util.RandomPgText(),
		PrivatePhoneNumber:        util.RandomPgText(),
		WorkPhoneNumber:           util.RandomPgText(),
		DateOfBirth: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		HomeTelephoneNumber: util.RandomPgText(),
		IsSubcontractor:     util.RandomPgBool(),
		Gender:              util.RandomPgText(),
		LocationID: pgtype.Int8{
			Int64: location.ID,
			Valid: true,
		},
		HasBorrowed:  false,
		OutOfService: util.RandomPgBool(),
		IsArchived:   util.RandomPgBool(),
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
	require.Equal(t, arg.EmailAddress, employee.EmailAddress)
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
	require.NotZero(t, employee.Created)

	// Verify foreign key constraints
	require.Equal(t, location.ID, employee.LocationID.Int64)
	require.True(t, employee.LocationID.Valid)
	return employee

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
					EmployeeNumber:            fmt.Sprintf("EMP%d", util.RandomInt(1000, 9999)),
					EmploymentNumber:          fmt.Sprintf("EN%d", util.RandomInt(10000, 99999)),
					Location:                  locationID,
					IsSubcontractor:           bool(util.RandomPgBool().Bool),
					FirstName:                 util.RandomString(6),
					LastName:                  util.RandomString(8),
					DateOfBirth:               "2000-01-05",
					Gender:                    "male",
					EmailAddress:              util.RandomEmail(),
					PrivateEmailAddress:       util.RandomEmail(),
					AuthenticationPhoneNumber: fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999)),
					WorkPhoneNumber:           fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999)),
					PrivatePhoneNumber:        fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999)),
					HomeTelephoneNumber:       fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999)),
					OutOfService:              bool(util.RandomPgBool().Bool),
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
				require.NotEmpty(t, response.EmployeeNumber)
				require.NotEmpty(t, response.EmploymentNumber)
				require.NotEmpty(t, response.Location)
				require.NotEmpty(t, response.FirstName)
				require.NotEmpty(t, response.LastName)
				require.NotEmpty(t, response.DateOfBirth)
				require.NotEmpty(t, response.Gender)
				require.NotEmpty(t, response.EmailAddress)
				require.NotEmpty(t, response.PrivateEmailAddress)
				require.NotEmpty(t, response.AuthenticationPhoneNumber)
				require.NotEmpty(t, response.WorkPhoneNumber)
				require.NotEmpty(t, response.PrivatePhoneNumber)
				require.NotEmpty(t, response.HomeTelephoneNumber)
				require.NotEmpty(t, response.User)
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
		IncludeArchived:     pgtype.Bool{Valid: true, Bool: true},
		IncludeOutOfService: pgtype.Bool{Valid: true, Bool: true},
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
					require.True(t, response.Results[i-1].Created.Time.After(response.Results[i].Created.Time) ||
						response.Results[i-1].Created.Time.Equal(response.Results[i].Created.Time))
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
					require.Equal(t, "IT", emp.Department.String)
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
