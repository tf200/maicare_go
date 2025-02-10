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

func createRandomClientDiagnosis(t *testing.T, clientID int64) db.ClientDiagnosis {

	arg := db.CreateClientDiagnosisParams{
		ClientID:            clientID,
		Title:               util.StringPtr("test title"),
		DiagnosisCode:       "test code",
		Description:         "test description",
		Severity:            util.StringPtr("Mild"),
		Status:              "Active",
		DiagnosingClinician: util.StringPtr("test clinician"),
		Notes:               util.StringPtr("test note"),
	}

	diagnosis, err := testStore.CreateClientDiagnosis(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, diagnosis)
	require.Equal(t, arg.ClientID, diagnosis.ClientID)
	return diagnosis
}

func TestCreateClientDiagnosisApi(t *testing.T) {
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
				diagnosisReq := CreateClientDiagnosisRequest{
					Title:               util.StringPtr("test title"),
					DiagnosisCode:       "test code",
					Description:         "test description",
					Severity:            util.StringPtr("Mild"),
					Status:              "Active",
					DiagnosingClinician: util.StringPtr("test clinician"),
					Notes:               util.StringPtr("test note"),
				}
				reqBody, err := json.Marshal(diagnosisReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/diagnosis", client.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var diagnosisResp Response[CreateClientDiagnosisResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &diagnosisResp)
				require.NoError(t, err)
				require.NotEmpty(t, diagnosisResp.Data.ID)
				require.Equal(t, diagnosisResp.Data.ClientID, client.ID)
				require.NotNil(t, diagnosisResp.Data.Title)
				require.Equal(t, *diagnosisResp.Data.Title, "test title")
				require.Equal(t, diagnosisResp.Data.DiagnosisCode, "test code")
				require.Equal(t, diagnosisResp.Data.Description, "test description")
				require.NotNil(t, diagnosisResp.Data.Severity)
				require.Equal(t, *diagnosisResp.Data.Severity, "Mild")
				require.Equal(t, diagnosisResp.Data.Status, "Active")
				require.NotNil(t, diagnosisResp.Data.DiagnosingClinician)
				require.Equal(t, *diagnosisResp.Data.DiagnosingClinician, "test clinician")
				require.NotNil(t, diagnosisResp.Data.Notes)
				require.Equal(t, *diagnosisResp.Data.Notes, "test note")
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

func TestListClientDiagnoses(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 20; i++ {
		createRandomClientDiagnosis(t, client.ID)
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
				url := fmt.Sprintf("/clients/%d/diagnosis?page=1&page_size=5", client.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var diagnosisResp Response[pagination.Response[ListClientDiagnosesResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &diagnosisResp)
				require.NoError(t, err)
				require.Len(t, diagnosisResp.Data.Results, 5)
				require.Equal(t, int32(5), diagnosisResp.Data.PageSize)
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

func TestGetClientDiagnosisApi(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis := createRandomClientDiagnosis(t, client.ID)

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
				url := fmt.Sprintf("/clients/%d/diagnosis/%d", client.ID, diagnosis.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var diagnosisResp Response[GetClientDiagnosisResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &diagnosisResp)
				require.NoError(t, err)
				require.Equal(t, diagnosisResp.Data.ID, diagnosis.ID)
				require.Equal(t, diagnosisResp.Data.ClientID, client.ID)
				require.NotNil(t, diagnosisResp.Data.Title)
				require.Equal(t, *diagnosisResp.Data.Title, "test title")
				require.Equal(t, diagnosisResp.Data.DiagnosisCode, "test code")
				require.Equal(t, diagnosisResp.Data.Description, "test description")
				require.NotNil(t, diagnosisResp.Data.Severity)
				require.Equal(t, *diagnosisResp.Data.Severity, "Mild")
				require.Equal(t, diagnosisResp.Data.Status, "Active")
				require.NotNil(t, diagnosisResp.Data.DiagnosingClinician)
				require.Equal(t, *diagnosisResp.Data.DiagnosingClinician, "test clinician")
				require.NotNil(t, diagnosisResp.Data.Notes)
				require.Equal(t, *diagnosisResp.Data.Notes, "test note")
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

func TestUpdateClientDiagnosisApi(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis := createRandomClientDiagnosis(t, client.ID)

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
				diagnosisReq := UpdateClientDiagnosisRequest{
					Title:               util.StringPtr("test title updated"),
					DiagnosisCode:       util.StringPtr("15235"),
					Description:         util.StringPtr("test description updated"),
					Severity:            util.StringPtr("Severe"),
					Status:              util.StringPtr("Inactive"),
					DiagnosingClinician: util.StringPtr("test clinician updated"),
					Notes:               util.StringPtr("test note updated"),
				}
				reqBody, err := json.Marshal(diagnosisReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/diagnosis/%d", client.ID, diagnosis.ID)
				request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var diagnosisResp Response[UpdateClientDiagnosisResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &diagnosisResp)
				require.NoError(t, err)
				require.Equal(t, diagnosisResp.Data.ID, diagnosis.ID)
				require.Equal(t, diagnosisResp.Data.ClientID, client.ID)
				require.NotNil(t, diagnosisResp.Data.Title)
				require.Equal(t, *diagnosisResp.Data.Title, "test title updated")
				require.Equal(t, diagnosisResp.Data.DiagnosisCode, "15235")
				require.Equal(t, diagnosisResp.Data.Description, "test description updated")
				require.NotNil(t, diagnosisResp.Data.Severity)
				require.Equal(t, *diagnosisResp.Data.Severity, "Severe")
				require.Equal(t, diagnosisResp.Data.Status, "Inactive")
				require.NotNil(t, diagnosisResp.Data.Notes)
				require.Equal(t, *diagnosisResp.Data.Notes, "test note updated")
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
