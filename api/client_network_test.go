package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/token"
	"maicare_go/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEmergencyContact(t *testing.T, clientID int64) db.ClientEmergencyContact {

	arg := db.CreateEmemrgencyContactParams{
		ClientID:         clientID,
		FirstName:        util.StringPtr(util.RandomString(5)),
		LastName:         util.StringPtr(util.RandomString(5)),
		Email:            util.StringPtr(util.RandomEmail()),
		PhoneNumber:      util.StringPtr(util.RandomString(4)),
		Address:          util.StringPtr(util.RandomString(5)),
		Relationship:     util.StringPtr(util.RandomString(5)),
		RelationStatus:   util.StringPtr("Primary Relationship"),
		MedicalReports:   util.RandomBool(),
		IncidentsReports: util.RandomBool(),
		GoalsReports:     util.RandomBool(),
	}

	contact, err := testStore.CreateEmemrgencyContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact)
	return contact
}

func TestCreateEmemrgencyContactApi(t *testing.T) {
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateClientEmergencyContactParams{
					FirstName:        util.StringPtr(util.RandomString(5)),
					LastName:         util.StringPtr(util.RandomString(5)),
					Email:            util.StringPtr(util.RandomEmail()),
					PhoneNumber:      util.StringPtr(util.RandomString(4)),
					Address:          util.StringPtr(util.RandomString(5)),
					Relationship:     util.StringPtr(util.RandomString(5)),
					RelationStatus:   util.StringPtr("Primary Relationship"),
					MedicalReports:   util.RandomBool(),
					IncidentsReports: util.RandomBool(),
					GoalsReports:     util.RandomBool(),
				}
				reqBody, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/emergency_contacts", client.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				var res Response[CreateClientEmergencyContactResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.NotEmpty(t, res.Data)
				require.Equal(t, res.Data.ClientID, client.ID)
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

func TestListClientEmergencyContactsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 20; i++ {
		createRandomEmergencyContact(t, client.ID)
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
				url := fmt.Sprintf("/clients/%d/emergency_contacts?page=1&page_size=5", client.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[pagination.Response[ListClientEmergencyContactsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				t.Log(res)
				require.NotEmpty(t, res.Data)
				require.Len(t, res.Data.Results, 5)
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

func TestGetEmergencyContactApi(t *testing.T) {
	client := createRandomClientDetails(t)
	contact := createRandomEmergencyContact(t, client.ID)

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
				url := fmt.Sprintf("/clients/%d/emergency_contacts/%d", client.ID, contact.ID)
				t.Log(url)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[GetClientEmergencyContactResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.NotEmpty(t, res.Data)
				require.Equal(t, res.Data.ID, contact.ID)
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

func TestUpdateEmergencyContactApi(t *testing.T) {
	client := createRandomClientDetails(t)
	contact := createRandomEmergencyContact(t, client.ID)

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
				updateReq := UpdateClientEmergencyContactParams{
					FirstName: util.StringPtr(util.RandomString(5)),
					LastName:  util.StringPtr(util.RandomString(5)),
					Email:     util.StringPtr(util.RandomEmail()),
				}
				reqBody, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/emergency_contacts/%d", client.ID, contact.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[UpdateClientEmergencyContactResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.NotEmpty(t, res.Data)
				require.Equal(t, res.Data.ID, contact.ID)
				require.NotEqual(t, res.Data.FirstName, contact.FirstName)
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

func TestDeleteEmergencyContactApi(t *testing.T) {
	client := createRandomClientDetails(t)
	contact := createRandomEmergencyContact(t, client.ID)

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
				url := fmt.Sprintf("/clients/%d/emergency_contacts/%d", client.ID, contact.ID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
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
