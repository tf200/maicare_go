package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/goccy/go-json"

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
		Gender:                    util.StringPtr("male"),
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
	_, user := createRandomEmployee(t)

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
				Empreq := CreateEmployeeProfileRequest{
					EmployeeNumber:            nil, // util.StringPtr(fmt.Sprintf("EMP%d", util.RandomInt(1000, 9999))),
					EmploymentNumber:          util.StringPtr(fmt.Sprintf("EN%d", util.RandomInt(10000, 99999))),
					LocationID:                util.IntPtr(locationID),
					IsSubcontractor:           util.BoolPtr(util.RandomBool()),
					FirstName:                 util.RandomString(6),
					LastName:                  util.RandomString(8),
					DateOfBirth:               util.StringPtr("2000-01-05"),
					Gender:                    util.StringPtr("male"),
					Email:                     "farsjiataha@gmail.com",
					PrivateEmailAddress:       util.StringPtr(util.RandomEmail()),
					AuthenticationPhoneNumber: util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					WorkPhoneNumber:           util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					PrivatePhoneNumber:        util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					HomeTelephoneNumber:       util.StringPtr(fmt.Sprintf("+%d%d", util.RandomInt(1, 99), util.RandomInt(1000000000, 9999999999))),
					OutOfService:              util.BoolPtr(util.RandomBool()),
					RoleID:                    1, // Assign a default role, e.g., RoleID 2
				}
				data, err := json.Marshal(Empreq)
				require.NoError(t, err)

				req, err := http.NewRequest(http.MethodPost, "/employees", bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body)
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response Response[CreateEmployeeProfileResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.Empty(t, response.Data.EmployeeNumber)
				require.NotEmpty(t, response.Data.EmploymentNumber)
				require.NotEmpty(t, response.Data.LocationID)
				require.NotEmpty(t, response.Data.FirstName)
				require.NotEmpty(t, response.Data.LastName)
				require.NotEmpty(t, response.Data.DateOfBirth)
				require.NotEmpty(t, response.Data.Gender)
				require.NotEmpty(t, response.Data.Email)
				require.NotEmpty(t, response.Data.PrivateEmailAddress)
				require.NotEmpty(t, response.Data.AuthenticationPhoneNumber)
				require.NotEmpty(t, response.Data.WorkPhoneNumber)
				require.NotEmpty(t, response.Data.PrivatePhoneNumber)
				require.NotEmpty(t, response.Data.HomeTelephoneNumber)
				require.NotEmpty(t, response.Data.UserID)
				require.NotEmpty(t, response.Data.CreatedAt)
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
				url := "/employees?page=1&page_size=10"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[pagination.Response[ListEmployeeResponse]]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotNil(t, response.Data.Next)
				require.Nil(t, response.Data.Previous)
				require.Equal(t, int64(numEmployees)+initialCount, response.Data.Count)
				require.Equal(t, int32(10), response.Data.PageSize)
				require.Len(t, response.Data.Results, 10)

				// Check results are ordered by created DESC
				for i := 1; i < len(response.Data.Results); i++ {
					require.True(t, response.Data.Results[i-1].CreatedAt.After(response.Data.Results[i].CreatedAt) ||
						response.Data.Results[i-1].CreatedAt.Equal(response.Data.Results[i].CreatedAt))
				}
			},
		},
		{
			name: "InvalidPageSize",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/employees?page=1&page_size=1000"
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
				url := "/employees?page=1&page_size=10&department=IT"
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
				url := "/employees?page=2&page_size=10"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[pagination.Response[ListEmployeeResponse]]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotNil(t, response.Data.Previous, "Previous should not be nil")
				require.Contains(t, *response.Data.Previous, "page=1", "Previous should contain 'page=1'")
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

func TestGetEmployeeProfileByIDApi(t *testing.T) {
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
				url := fmt.Sprintf("/employees/%d", employee.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[GetEmployeeProfileByIDApiResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.Equal(t, employee.ID, response.Data.ID)
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

func TestUpdateEmployeeProfileApi(t *testing.T) {
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
				updatereq := UpdateEmployeeProfileRequest{
					EmployeeNumber:   nil, // util.StringPtr(fmt.Sprintf("EMP%d", util.RandomInt(1000, 9999))),
					EmploymentNumber: util.StringPtr(fmt.Sprintf("EN%d", util.RandomInt(10000, 99999)))}
				data, err := json.Marshal(updatereq)
				require.NoError(t, err)
				url := fmt.Sprintf("/employees/%d", employee.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[UpdateEmployeeProfileResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.Equal(t, employee.ID, response.Data.ID)
				require.Equal(t, employee.UserID, response.Data.UserID)
				require.Equal(t, employee.FirstName, response.Data.FirstName)
				require.Equal(t, employee.LastName, response.Data.LastName)
				require.NotEqual(t, employee.EmploymentNumber, response.Data.EmploymentNumber)

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
				require.NotEmpty(t, response.Data.Permissions)

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

func TestSetEmployeeProfilePictureApi(t *testing.T) {
	employee, user := createRandomEmployee(t)
	attachement := createRandomAttachmentFile(t)
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
				url := fmt.Sprintf("/employees/%d/profile_picture", employee.ID)
				setReq := SetEmployeeProfilePictureRequest{
					AttachmentID: attachement.Uuid,
				}
				data, err := json.Marshal(setReq)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[SetEmployeeProfilePictureResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
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

func createRandomEducation(t *testing.T) (int64, int64) {
	employee, user := createRandomEmployee(t)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder) int64
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				addEducationReq := AddEducationToEmployeeProfileRequest{
					Degree:          "BsC",
					FieldOfStudy:    "Computer Science",
					InstitutionName: "University of Ghana",
					StartDate:       "2018-01-01",
					EndDate:         "2022-01-01",
				}
				data, err := json.Marshal(addEducationReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/employees/%d/education", employee.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) int64 {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response Response[AddEducationToEmployeeProfileResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.Equal(t, "BsC", response.Data.Degree)
				require.Equal(t, "Computer Science", response.Data.FieldOfStudy)
				require.Equal(t, "University of Ghana", response.Data.InstitutionName)
				return response.Data.ID

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
	return employee.ID, user.ID
}

func TestAddEducationToEmployeeProfileApi(t *testing.T) {
	createRandomEducation(t)
}

func TestListEmployeeEducationApi(t *testing.T) {
	employeeID, userID := createRandomEducation(t)
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
				url := fmt.Sprintf("/employees/%d/education", employeeID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[[]ListEmployeeEducationResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotEmpty(t, response.Data)

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

// to do test UpdateEmployeeEducationApi

func TestAddEmployeeExperienceApi(t *testing.T) {
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
				addExperienceReq := AddEmployeeExperienceRequest{
					CompanyName: "Google",
					JobTitle:    "Software Engineer",
					StartDate:   "2018-01-01",
					EndDate:     "2022-01-01",
					Description: util.StringPtr("Worked on the search engine"),
				}
				data, err := json.Marshal(addExperienceReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/employees/%d/experience", employee.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response Response[AddEmployeeExperienceResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.Equal(t, "Google", response.Data.CompanyName)
				require.Equal(t, "Software Engineer", response.Data.JobTitle)
				require.Equal(t, "2018-01-01T00:00:00Z", response.Data.StartDate)
				require.Equal(t, "2022-01-01T00:00:00Z", response.Data.EndDate)
				require.Equal(t, util.StringPtr("Worked on the search engine"), response.Data.Description)

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

func TestAddEmployeeCertificationApi(t *testing.T) {
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
				addCertificationReq := AddEmployeeCertificationRequest{
					Name:       "AWS Certified Developer",
					IssuedBy:   "AWS",
					DateIssued: "2022-01-01",
				}
				data, err := json.Marshal(addCertificationReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/employees/%d/certification", employee.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response Response[AddEmployeeCertificationResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.Equal(t, "AWS Certified Developer", response.Data.Name)
				require.Equal(t, "AWS", response.Data.IssuedBy)
				require.Equal(t, "2022-01-01", response.Data.DateIssued.Format("2006-01-02"))
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

func TestSearchEmployeesByNameOrEmailApi(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEmployee(t)
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
				url := "/employees/emails?search=John"
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body)
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[[]SearchEmployeesByNameOrEmailResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
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
