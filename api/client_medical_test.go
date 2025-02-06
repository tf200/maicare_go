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

	"github.com/stretchr/testify/require"
)

// func TestListAllergyTypesApi (t *testing.T){
// 	test
// }

func createRandomClientAllergy(t *testing.T, clientID int64) db.ClientAllergy {

	arg := db.CreateClientAllergyParams{
		ClientID:      clientID,
		AllergyTypeID: 1,
		Severity:      "Mild",
		Reaction:      "test reaction",
		Notes:         util.StringPtr("test note"),
	}

	allergy, err := testStore.CreateClientAllergy(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, allergy)
	require.Equal(t, arg.ClientID, allergy.ClientID)
	return allergy
}

func TestCreateClientAllergyApi(t *testing.T) {
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
				allergyReq := CreateClientAllergyRequest{
					AllergyTypeID: 1,
					Severity:      "Mild",
					Reaction:      "test reaction",
					Notes:         util.StringPtr("test note"),
				}
				reqBody, err := json.Marshal(allergyReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/allergies", client.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var allergyResp Response[CreateClientAllergyResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &allergyResp)
				require.NoError(t, err)
				require.NotEmpty(t, allergyResp.Data.ID)
				require.Equal(t, allergyResp.Data.ClientID, client.ID)
				require.Equal(t, allergyResp.Data.AllergyTypeID, int64(1))
				require.Equal(t, allergyResp.Data.Severity, "Mild")
				require.Equal(t, allergyResp.Data.Reaction, "test reaction")
				require.Equal(t, *allergyResp.Data.Notes, "test note")
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

func TestListClientAllergies(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 20; i++ {
		createRandomClientAllergy(t, client.ID)
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
				url := fmt.Sprintf("/clients/%d/allergies?page=1&page_size=5", client.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var allergyResp Response[pagination.Response[ListClientAllergiesResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &allergyResp)
				require.NoError(t, err)
				require.Len(t, allergyResp.Data.Results, 5)
				require.Equal(t, int32(5), allergyResp.Data.PageSize)
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

func TestGetClientAllergyApi(t *testing.T) {
	client := createRandomClientDetails(t)
	allergy := createRandomClientAllergy(t, client.ID)
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
				url := fmt.Sprintf("/clients/%d/allergies/%d", client.ID, allergy.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var allergyResp Response[GetClientAllergyResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &allergyResp)
				require.NoError(t, err)
				require.Equal(t, allergyResp.Data.ID, allergy.ID)
				require.Equal(t, allergyResp.Data.ClientID, client.ID)
				require.Equal(t, allergyResp.Data.AllergyTypeID, int64(1))
				require.Equal(t, allergyResp.Data.Severity, "Mild")
				require.Equal(t, allergyResp.Data.Reaction, "test reaction")
				require.NotNil(t, allergyResp.Data.Notes)
				require.Equal(t, *allergyResp.Data.Notes, "test note")
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

func TestUpdateClientAllergyyApi(t *testing.T) {
	client := createRandomClientDetails(t)
	allergy := createRandomClientAllergy(t, client.ID)

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
				allergyReq := UpdateClientAllergyRequest{
					AllergyTypeID: util.IntPtr(2),
					Severity:      util.StringPtr("Severe"),
					Reaction:      util.StringPtr("test reaction updated"),
					Notes:         util.StringPtr("test note updated"),
				}
				reqBody, err := json.Marshal(allergyReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/allergies/%d", client.ID, allergy.ID)
				request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var allergyResp Response[UpdateClientAllergyResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &allergyResp)
				require.NoError(t, err)
				require.Equal(t, allergyResp.Data.ID, allergy.ID)
				require.Equal(t, allergyResp.Data.ClientID, client.ID)
				require.Equal(t, allergyResp.Data.AllergyTypeID, int64(2))
				require.Equal(t, allergyResp.Data.Severity, "Severe")
				require.Equal(t, allergyResp.Data.Reaction, "test reaction updated")
				require.NotNil(t, allergyResp.Data.Notes)
				require.Equal(t, *allergyResp.Data.Notes, "test note updated")
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

func TestDeleteClientAllergyApi(t *testing.T) {
	client := createRandomClientDetails(t)
	allergy := createRandomClientAllergy(t, client.ID)

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
				url := fmt.Sprintf("/clients/%d/allergies/%d", client.ID, allergy.ID)
				request, err := http.NewRequest(http.MethodDelete, url, nil)
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

// TO DO CLIENT DIAGNOSIS TEST
