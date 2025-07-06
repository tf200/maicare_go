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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomContractType(t *testing.T) db.ContractType {

	contractType, err := testStore.CreateContractType(context.Background(), "Test Contract Type")
	require.NoError(t, err)
	require.NotEmpty(t, contractType)
	require.NotEmpty(t, contractType.ID)
	require.Equal(t, "Test Contract Type", contractType.Name)
	return contractType
}

func TestCreateContractTypeApi(t *testing.T) {
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
				requestBody := CreateContractTypeRequest{
					Name: "Test Contract Type",
				}
				requestBodyBytes, err := json.Marshal(requestBody)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPost, "/contract_types", bytes.NewReader(requestBodyBytes))
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[CreateContractTypeResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.ID)
				require.Equal(t, "Test Contract Type", response.Data.Name)
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

func TestListContractTypeApi(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomContractType(t)
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
				request, err := http.NewRequest(http.MethodGet, "/contract_types", nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[[]ListContractTypesResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
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

var (
	// PriceFrequency represents available pricing frequency options
	PriceFrequency = []string{"minute", "hourly", "daily", "weekly", "monthly"}

	// HoursType represents types of hour calculations
	HoursType = []string{"weekly", "all_period"}

	// CareType represents types of care provided
	CareType = []string{"ambulante", "accommodation"}

	// FinancingAct represents different financing acts/laws
	FinancingAct = []string{"WMO", "ZVW", "WLZ", "JW", "WPG"}

	// FinancingOption represents financing options
	FinancingOption = []string{"ZIN", "PGB"}
)

func createRandomContract(t *testing.T, clientID int64, senderID *int64) db.Contract {

	contractType := createRandomContractType(t)
	attachment := createRandomAttachmentFile(t)

	arg := db.CreateContractParams{
		TypeID:          &contractType.ID,
		StartDate:       pgtype.Timestamptz{Time: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Date(2021, 9, 30, 0, 0, 0, 0, time.UTC), Valid: true},
		ReminderPeriod:  10,
		Vat:             util.Int32Ptr(15),
		Price:           5.58,
		PriceTimeUnit:   util.RandomEnum(PriceFrequency),
		Hours:           util.Float64Ptr(100),
		HoursType:       util.StringPtr(util.RandomEnum(HoursType)),
		CareName:        "Test Care",
		CareType:        "ambulante",
		ClientID:        clientID,
		SenderID:        senderID,
		FinancingAct:    util.RandomEnum(FinancingAct),
		FinancingOption: util.RandomEnum(FinancingOption),
		AttachmentIds:   []uuid.UUID{attachment.Uuid},
	}

	contract, err := testStore.CreateContract(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contract)
	require.Equal(t, arg.TypeID, contract.TypeID)
	require.Equal(t, arg.ReminderPeriod, contract.ReminderPeriod)
	require.Equal(t, arg.Vat, contract.Vat)
	require.Equal(t, arg.Price, contract.Price)
	require.Equal(t, arg.PriceTimeUnit, contract.PriceTimeUnit)
	require.Equal(t, arg.Hours, contract.Hours)
	require.Equal(t, arg.HoursType, contract.HoursType)
	require.Equal(t, arg.CareName, contract.CareName)
	require.Equal(t, arg.CareType, contract.CareType)
	require.Equal(t, arg.ClientID, contract.ClientID)
	require.Equal(t, arg.SenderID, contract.SenderID)
	require.Equal(t, arg.FinancingAct, contract.FinancingAct)
	require.Equal(t, arg.FinancingOption, contract.FinancingOption)
	require.NotZero(t, contract.ID)
	return contract
}

func TestCreateClientContractApi(t *testing.T) {
	client := createRandomClientDetails(t)
	contractType := createRandomContractType(t)
	attachment := createRandomAttachmentFile(t)

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
				requestBody := CreateContractRequest{
					TypeID:          &contractType.ID,
					StartDate:       time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
					EndDate:         time.Date(2021, 9, 30, 0, 0, 0, 0, time.UTC),
					ReminderPeriod:  10,
					Vat:             util.Int32Ptr(15),
					Price:           5.58,
					PriceTimeUnit:   util.RandomEnum(PriceFrequency),
					Hours:           util.Float64Ptr(100),
					HoursType:       util.StringPtr(util.RandomEnum(HoursType)),
					CareName:        "Test Care",
					CareType:        "ambulante",
					SenderID:        client.SenderID,
					FinancingAct:    util.RandomEnum(FinancingAct),
					FinancingOption: util.RandomEnum(FinancingOption),
					AttachmentIds:   []uuid.UUID{attachment.Uuid},
				}
				requestBodyBytes, err := json.Marshal(requestBody)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d/contracts", client.ID)
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBodyBytes))
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[CreateContractResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.ID)
				require.Equal(t, &contractType.ID, response.Data.TypeID)
				require.Equal(t, int32(10), response.Data.ReminderPeriod)
				require.Equal(t, util.Int32Ptr(15), response.Data.Vat)
				require.Equal(t, 5.58, response.Data.Price)
				require.Contains(t, PriceFrequency, response.Data.PriceTimeUnit)
				require.Equal(t, util.Float64Ptr(100), response.Data.Hours)
				require.NotNil(t, response.Data.HoursType)
				require.Contains(t, HoursType, *response.Data.HoursType)
				require.Equal(t, "Test Care", response.Data.CareName)
				require.Contains(t, CareType, response.Data.CareType)
				require.Equal(t, client.ID, response.Data.ClientID)
				require.Equal(t, client.SenderID, response.Data.SenderID)
				require.Contains(t, FinancingAct, response.Data.FinancingAct)
				require.Contains(t, FinancingOption, response.Data.FinancingOption)
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

func TestListClientContractsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 10; i++ {
		createRandomContract(t, client.ID, client.SenderID)
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

				url := fmt.Sprintf("/clients/%d/contracts?page=1&page_size=5", client.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[pagination.Response[ListClientContractsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.Results)
				require.Len(t, response.Data.Results, 5)
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

func TestGetClientContract(t *testing.T) {
	client := createRandomClientDetails(t)
	contract := createRandomContract(t, client.ID, client.SenderID)
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
				url := fmt.Sprintf("/clients/%d/contracts/%d", contract.ClientID, contract.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetClientContractResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.ID)
				require.Equal(t, contract.ID, response.Data.ID)
				require.Equal(t, contract.TypeID, response.Data.TypeID)
				require.Equal(t, int32(10), response.Data.ReminderPeriod)
				require.Equal(t, util.Int32Ptr(15), response.Data.Vat)
				require.Equal(t, 5.58, response.Data.Price)
				require.Contains(t, PriceFrequency, response.Data.PriceTimeUnit)
				require.Equal(t, util.Float64Ptr(100), response.Data.Hours)
				require.NotNil(t, response.Data.HoursType)
				require.Contains(t, HoursType, *response.Data.HoursType)
				require.Equal(t, "Test Care", response.Data.CareName)
				require.Contains(t, CareType, response.Data.CareType)
				require.Equal(t, contract.ClientID, response.Data.ClientID)
				require.Equal(t, contract.SenderID, response.Data.SenderID)
				require.Contains(t, FinancingAct, response.Data.FinancingAct)
				require.Contains(t, FinancingOption, response.Data.FinancingOption)
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

func TestListContractsApi(t *testing.T) {
	for i := 0; i < 10; i++ {
		client := createRandomClientDetails(t)
		createRandomContract(t, client.ID, client.SenderID)
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

				request, err := http.NewRequest(http.MethodGet, "/contracts?page=1&page_size=5", nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[pagination.Response[ListContractsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.Results)
				require.Len(t, response.Data.Results, 5)
			},
		},
		{
			name: "Filter By Status",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {

				request, err := http.NewRequest(http.MethodGet, "/contracts?page=1&page_size=5&status=draft", nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[pagination.Response[ListContractsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Len(t, response.Data.Results, 5)
				require.Equal(t, "draft", response.Data.Results[0].Status)
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
