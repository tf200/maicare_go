package api

import (
	"bytes"
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	clientp "maicare_go/service/client"
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

// func TestListAllergyTypesApi (t *testing.T){
// 	test
// }

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
				diagnosisReq := clientp.CreateClientDiagnosisRequest{
					Title:               util.StringPtr("test title"),
					DiagnosisCode:       "test code",
					Description:         "test description",
					Severity:            util.StringPtr("Mild"),
					Status:              "Active",
					DiagnosingClinician: util.StringPtr("test clinician"),
					Notes:               util.StringPtr("test note"),
					Medications: []clientp.DiagnosisMedicationCreate{
						{
							Name:             "test medication",
							Dosage:           "test dosage",
							StartDate:        time.Now(),
							EndDate:          time.Now().Add(24 * time.Hour),
							Notes:            util.StringPtr("test note"),
							SelfAdministered: false,
							AdministeredByID: &employee.ID,
							IsCritical:       false,
						},
					},
				}
				reqBody, err := json.Marshal(diagnosisReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/diagnosis", client.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/jsson")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var diagnosisResp Response[clientp.CreateClientDiagnosisResponse]
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
	employee, user := createRandomEmployee(t)
	for i := 0; i < 20; i++ {
		diag := createRandomClientDiagnosis(t, client.ID)
		createRandomClientMedication(t, diag.ID, employee.ID)
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
				url := fmt.Sprintf("/clients/%d/diagnosis?page=1&page_size=5", client.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var diagnosisResp Response[pagination.Response[clientp.ListClientDiagnosesResponse]]
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
				var diagnosisResp Response[clientp.GetClientDiagnosisResponse]
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
				diagnosisReq := clientp.UpdateClientDiagnosisRequest{
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
				var diagnosisResp Response[clientp.UpdateClientDiagnosisResponse]
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

func TestDeleteClientDiagnosisApi(t *testing.T) {
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

func createRandomClientMedication(t *testing.T, diagnosisID int64, employeeID int64) db.ClientMedication {

	arg := db.CreateClientMedicationParams{
		DiagnosisID:      &diagnosisID,
		Name:             "test name",
		Dosage:           "test dosage",
		StartDate:        pgtype.Date{Time: util.RandomTIme(), Valid: true},
		EndDate:          pgtype.Date{Time: util.RandomTIme(), Valid: true},
		Notes:            util.StringPtr("test note"),
		SelfAdministered: true,
		AdministeredByID: util.IntPtr(employeeID),
		IsCritical:       true,
	}

	medication, err := testStore.CreateClientMedication(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, medication)
	require.Equal(t, arg.DiagnosisID, medication.DiagnosisID)
	return medication
}

func TestCreateClientMedicationApi(t *testing.T) {
	employee, user := createRandomEmployee(t)
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				medicationReq := clientp.CreateClientMedicationRequest{
					Name:             "test medication",
					Dosage:           "test dosage",
					StartDate:        time.Now(),
					EndDate:          time.Now().Add(24 * time.Hour),
					Notes:            util.StringPtr("test note"),
					SelfAdministered: false,
					AdministeredByID: &employee.ID,
					IsCritical:       false,
				}
				reqBody, err := json.Marshal(medicationReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/diagnosis/%d/medications", client.ID, diagnosis.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var medicationResp Response[clientp.CreateClientMedicationResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &medicationResp)
				require.NoError(t, err)
				require.NotEmpty(t, medicationResp.Data.ID)
				require.NotNil(t, medicationResp.Data.Name)
				require.Equal(t, medicationResp.Data.Name, "test medication")
				require.Equal(t, medicationResp.Data.Dosage, "test dosage")
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

func TestListClientMedications(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis := createRandomClientDiagnosis(t, client.ID)
	employee, user := createRandomEmployee(t)
	for i := 0; i < 20; i++ {
		createRandomClientMedication(t, diagnosis.ID, employee.ID)
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
				url := fmt.Sprintf("/clients/%d/diagnosis/%d/medications?page=1&page_size=5", client.ID, diagnosis.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var medicationResp Response[pagination.Response[ListClientMedicationsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &medicationResp)
				require.NoError(t, err)
				require.Len(t, medicationResp.Data.Results, 5)
				require.Equal(t, int32(5), medicationResp.Data.PageSize)
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

func TestGetClientMedicationApi(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis := createRandomClientDiagnosis(t, client.ID)
	employee, user := createRandomEmployee(t)
	medication := createRandomClientMedication(t, diagnosis.ID, employee.ID)

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
				url := fmt.Sprintf("/clients/%d/diagnosis/%d/medications/%d", client.ID, diagnosis.ID, medication.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var medicationResp Response[GetClientMedicationResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &medicationResp)
				require.NoError(t, err)
				require.Equal(t, medicationResp.Data.ID, medication.ID)
				require.Equal(t, medicationResp.Data.Name, medication.Name)
				require.Equal(t, medicationResp.Data.Dosage, medication.Dosage)
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

func TestUpdateClientMedicationApi(t *testing.T) {
	client := createRandomClientDetails(t)
	diagnosis := createRandomClientDiagnosis(t, client.ID)
	employee, user := createRandomEmployee(t)
	medication := createRandomClientMedication(t, diagnosis.ID, employee.ID)

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
				medicationReq := UpdateClientMedicationRequest{
					Name:             util.StringPtr("test medication updated"),
					Dosage:           util.StringPtr("test dosage updated"),
					StartDate:        time.Now(),
					EndDate:          time.Now().Add(24 * time.Hour),
					Notes:            util.StringPtr("test note updated"),
					SelfAdministered: util.BoolPtr(false),
					AdministeredByID: &employee.ID,
					IsCritical:       util.BoolPtr(false),
				}
				reqBody, err := json.Marshal(medicationReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/diagnosis/%d/medications/%d", client.ID, diagnosis.ID, medication.ID)
				request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var medicationResp Response[UpdateClientMedicationResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &medicationResp)
				require.NoError(t, err)
				require.Equal(t, medicationResp.Data.ID, medication.ID)
				require.Equal(t, medicationResp.Data.Name, "test medication updated")
				require.Equal(t, medicationResp.Data.Dosage, "test dosage updated")
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
		},
		)
	}
}
