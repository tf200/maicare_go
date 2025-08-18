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

	"github.com/go-faker/faker/v4"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomRegistrationForm(t *testing.T) db.RegistrationForm {
	arg := db.CreateRegistrationFormParams{
		ClientFirstName:               faker.FirstName(),
		ClientLastName:                faker.LastName(),
		ClientBsnNumber:               "123456789",
		ClientGender:                  "male",
		ClientNationality:             "Dutch",
		ClientPhoneNumber:             faker.PhoneNumber,
		ClientEmail:                   faker.Email(),
		ClientStreet:                  "123 Main St",
		ClientHouseNumber:             "1A",
		ClientPostalCode:              "1234AB",
		ClientCity:                    "Amsterdam",
		ReferrerFirstName:             faker.FirstName(),
		ReferrerLastName:              faker.LastName(),
		ReferrerOrganization:          "Referrer Org",
		ReferrerJobTitle:              "Referrer Job",
		ReferrerPhoneNumber:           faker.PhoneNumber,
		ReferrerEmail:                 faker.Email(),
		Guardian1FirstName:            faker.FirstName(),
		Guardian1LastName:             faker.LastName(),
		Guardian1Relationship:         "Parent",
		Guardian1PhoneNumber:          faker.PhoneNumber,
		Guardian1Email:                faker.Email(),
		Guardian2FirstName:            faker.FirstName(),
		Guardian2LastName:             faker.LastName(),
		Guardian2Relationship:         "Parent",
		Guardian2PhoneNumber:          faker.PhoneNumber,
		Guardian2Email:                faker.Email(),
		EducationInstitution:          util.StringPtr("Education Institution"),
		EducationMentorName:           util.StringPtr(faker.Name()),
		EducationMentorPhone:          util.StringPtr(faker.PhoneNumber),
		EducationMentorEmail:          util.StringPtr(faker.Email()),
		EducationCurrentlyEnrolled:    true,
		EducationAdditionalNotes:      util.StringPtr("Additional notes"),
		EducationLevel:                util.StringPtr("higher"),
		CareProtectedLiving:           util.BoolPtr(true),
		CareAssistedIndependentLiving: util.BoolPtr(false),
		CareRoomTrainingCenter:        util.BoolPtr(true),
		CareAmbulatoryGuidance:        util.BoolPtr(false),
		RiskAggressiveBehavior:        util.BoolPtr(true),
		RiskSuicidalSelfharm:          util.BoolPtr(false),
		RiskSubstanceAbuse:            util.BoolPtr(true),
		RiskPsychiatricIssues:         util.BoolPtr(false),
		RiskCriminalHistory:           util.BoolPtr(true),
		RiskFlightBehavior:            util.BoolPtr(false),
		RiskWeaponPossession:          util.BoolPtr(true),
		RiskSexualBehavior:            util.BoolPtr(false),
		RiskDayNightRhythm:            util.BoolPtr(true),
		RiskOther:                     util.BoolPtr(false),
		RiskOtherDescription:          nil,
		RiskAdditionalNotes:           nil,
		DocumentReferral:              nil,
		DocumentEducationReport:       nil,
		DocumentPsychiatricReport:     nil,
		DocumentDiagnosis:             nil,
		DocumentSafetyPlan:            nil,
		DocumentIDCopy:                nil,
		ApplicationDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		ReferrerSignature: util.BoolPtr(true),
	}

	registrationForm, err := testStore.CreateRegistrationForm(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, registrationForm)
	return registrationForm
}

func TestCreateRegistrationFormApi(t *testing.T) {

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
				rfRequest := CreateRegistrationFormRequest{
					ClientFirstName:               faker.FirstName(),
					ClientLastName:                faker.LastName(),
					ClientBsnNumber:               "123456789",
					ClientGender:                  "male",
					ClientNationality:             "Dutch",
					ClientPhoneNumber:             faker.PhoneNumber,
					ClientEmail:                   faker.Email(),
					ClientStreet:                  "123 Main St",
					ClientHouseNumber:             "1A",
					ClientPostalCode:              "1234AB",
					ClientCity:                    "Amsterdam",
					ReferrerFirstName:             faker.FirstName(),
					ReferrerLastName:              faker.LastName(),
					ReferrerOrganization:          "Referrer Org",
					ReferrerJobTitle:              "Referrer Job",
					ReferrerPhoneNumber:           faker.PhoneNumber,
					ReferrerEmail:                 faker.Email(),
					Guardian1FirstName:            faker.FirstName(),
					Guardian1LastName:             faker.LastName(),
					Guardian1Relationship:         "Parent",
					Guardian1PhoneNumber:          faker.PhoneNumber,
					Guardian1Email:                faker.Email(),
					Guardian2FirstName:            faker.FirstName(),
					Guardian2LastName:             faker.LastName(),
					Guardian2Relationship:         "Parent",
					Guardian2PhoneNumber:          faker.PhoneNumber,
					Guardian2Email:                faker.Email(),
					EducationInstitution:          util.StringPtr("Education Institution"),
					EducationMentorName:           util.StringPtr(faker.Name()),
					EducationMentorPhone:          util.StringPtr(faker.PhoneNumber),
					EducationMentorEmail:          util.StringPtr(faker.Email()),
					EducationCurrentlyEnrolled:    true,
					EducationAdditionalNotes:      util.StringPtr("Additional notes"),
					EducationLevel:                util.StringPtr("higher"),
					WorkCurrentEmployer:           util.StringPtr("Current Employer"),
					WorkEmployerPhone:             util.StringPtr(faker.PhoneNumber),
					WorkEmployerEmail:             util.StringPtr(faker.Email()),
					WorkCurrentPosition:           util.StringPtr("Current Position"),
					WorkCurrentlyEmployed:         true,
					WorkStartDate:                 util.TimePtr(time.Now().AddDate(2, 0, 0)),
					WorkAdditionalNotes:           util.StringPtr("Work additional notes"),
					CareProtectedLiving:           util.BoolPtr(true),
					CareAssistedIndependentLiving: util.BoolPtr(false),
					CareRoomTrainingCenter:        util.BoolPtr(true),
					CareAmbulatoryGuidance:        util.BoolPtr(false),
					RiskAggressiveBehavior:        util.BoolPtr(true),
					RiskSuicidalSelfharm:          util.BoolPtr(false),
					RiskSubstanceAbuse:            util.BoolPtr(true),
					RiskPsychiatricIssues:         util.BoolPtr(false),
					RiskCriminalHistory:           util.BoolPtr(true),
					RiskFlightBehavior:            util.BoolPtr(false),
					RiskWeaponPossession:          util.BoolPtr(true),
					RiskSexualBehavior:            util.BoolPtr(false),
					RiskDayNightRhythm:            util.BoolPtr(true),
					RiskOther:                     util.BoolPtr(false),
					RiskOtherDescription:          nil,
					RiskAdditionalNotes:           nil,
					DocumentReferral:              nil,
					DocumentEducationReport:       nil,
					DocumentPsychiatricReport:     nil,
					DocumentDiagnosis:             nil,
					DocumentSafetyPlan:            nil,
					DocumentIDCopy:                nil,
					ApplicationDate:               time.Now(),
					ReferrerSignature:             util.BoolPtr(true),
				}
				data, err := json.Marshal(rfRequest)
				require.NoError(t, err)
				url := "/registration_form"
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var registrationFormCard Response[CreateRegistrationFormResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &registrationFormCard)
				require.NoError(t, err)
				require.NotEmpty(t, registrationFormCard.Data)
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

func TestListRegistrationFormsApi(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomRegistrationForm(t)
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
				url := "/registration_form?page=1&page_size=5"
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var registrationFormsCard Response[pagination.Response[ListRegistrationFormsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &registrationFormsCard)
				require.NoError(t, err)
				require.NotEmpty(t, registrationFormsCard.Data)
				require.Len(t, registrationFormsCard.Data.Results, 5)
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

func TestGetRegistrationFormApi(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)

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
				url := fmt.Sprintf("/registration_form/%d", registrationForm.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var registrationFormCard Response[GetRegistrationFormResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &registrationFormCard)
				require.NoError(t, err)
				require.NotEmpty(t, registrationFormCard.Data)
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

func TestUpdateRegistrationFormApi(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)

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
				rfRequest := UpdateRegistrationFormRequest{
					ClientFirstName: util.StringPtr(faker.FirstName()),
					ClientLastName:  util.StringPtr(faker.LastName()),
					ClientBsnNumber: util.StringPtr("123456789"),
				}
				data, err := json.Marshal(rfRequest)
				require.NoError(t, err)
				url := fmt.Sprintf("/registration_form/%d", registrationForm.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var registrationFormCard Response[UpdateRegistrationFormResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &registrationFormCard)
				require.NoError(t, err)
				require.NotEmpty(t, registrationFormCard.Data)
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

func TestDeleteRegistrationFormApi(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)

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
				url := fmt.Sprintf("/registration_form/%d", registrationForm.ID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)

				return req, nil
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

func TestUpdateRegistrationFormStatusApi(t *testing.T) {
	registrationForm := createRandomRegistrationForm(t)

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
				reqBody := UpdateRegistrationFormStatusRequest{
					Status:                    "approved",
					IntakeAppointmentDate:     time.Now().AddDate(0, 0, 7),
					IntakeAppointmentLocation: util.StringPtr("Intake Location"),
					AddmissionType:            util.StringPtr("regular_placement"),
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)
				url := fmt.Sprintf("/registration_form/%d/status", registrationForm.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
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
