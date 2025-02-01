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
				url := fmt.Sprintf("/clients/%d/client_allergies", client.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
				url := fmt.Sprintf("/clients/%d/client_allergies?page=1&page_size=5", client.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var allergyResp Response[pagination.Response[ListClientAllergiesResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &allergyResp)
				require.NoError(t, err)
				t.Log(allergyResp)
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
