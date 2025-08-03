package api

import (
	"bytes"
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/invoice"
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

func createRandomBillableHours(t *testing.T) int64 {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)

	// priceFrequency := []string{"minute", "hourly", "daily", "weekly", "monthly"}
	// careType := []string{"ambulante", "accommodation"}
	financingAct := []string{"WMO", "ZVW", "WLZ", "JW", "WPG"}
	financingOption := []string{"ZIN", "PGB"}

	contractType := createRandomContractType(t)

	arg1 := db.CreateContractParams{
		TypeID:          &contractType.ID,
		StartDate:       pgtype.Timestamptz{Time: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC), Valid: true},
		ReminderPeriod:  10,
		Vat:             util.Int32Ptr(20),
		Status:          "approved",
		Price:           0,
		PriceTimeUnit:   "daily", // util.RandomEnum(priceFrequency),
		Hours:           nil,
		HoursType:       nil,
		CareName:        "Test Care",
		CareType:        "accommodation",
		ClientID:        client.ID,
		SenderID:        client.SenderID,
		FinancingAct:    util.RandomEnum(financingAct),
		FinancingOption: util.RandomEnum(financingOption),
		AttachmentIds:   []uuid.UUID{},
	}

	contract, err := testStore.CreateContract(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, contract)

	arg2 := db.CreateContractParams{
		TypeID:          &contractType.ID,
		StartDate:       pgtype.Timestamptz{Time: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC), Valid: true},
		ReminderPeriod:  17,
		Vat:             util.Int32Ptr(20),
		Status:          "approved",
		Price:           558,
		PriceTimeUnit:   "hourly", // util.RandomEnum(priceFrequency),
		Hours:           util.Float64Ptr(40),
		HoursType:       util.StringPtr("weekly"),
		CareName:        "Test Care",
		CareType:        "ambulante", // util.RandomEnum(careType),
		ClientID:        client.ID,
		SenderID:        client.SenderID,
		FinancingAct:    util.RandomEnum(financingAct),
		FinancingOption: util.RandomEnum(financingOption),
		AttachmentIds:   []uuid.UUID{},
	}
	contract2, err := testStore.CreateContract(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, contract2)

	appointement := createRandomAppointment(t, employee.ID)
	err = testStore.BulkAddAppointmentClients(context.Background(), db.BulkAddAppointmentClientsParams{
		AppointmentID: appointement.ID,
		ClientIds:     []int64{client.ID},
	})

	require.NoError(t, err)
	return client.ID
}

func TestCreateInvoiceApi(t *testing.T) {
	sender := createRandomSender(t)
	clientID := createRandomBillableHours(t)
	contract := createRandomContract(t, clientID, &sender.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute, 1)
			},
			buildRequest: func() (*http.Request, error) {
				amblanteTiotalMinutes := 100.0
				preVatTotal := contract.Price * amblanteTiotalMinutes
				Total := preVatTotal * (1 + float64(*contract.Vat)/100)
				req := CreateInvoiceRequest{
					ClientID:  clientID,
					IssueDate: time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
					DueDate:   time.Date(2025, time.August, 31, 23, 59, 59, 0, time.UTC),
					InvoiceDetails: []invoice.InvoiceDetails{
						{
							ContractID:   contract.ID,
							ContractType: contract.CareType,
							Periods: []invoice.InvoicePeriod{
								{
									StartDate:             time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
									EndDate:               time.Date(2025, time.August, 31, 23, 59, 59, 0, time.UTC),
									AmbulanteTotalMinutes: &amblanteTiotalMinutes,
								},
							},
							PreVatTotal:   preVatTotal,
							Vat:           float64(*contract.Vat),
							Total:         Total,
							Price:         contract.Price,
							PriceTimeUnit: contract.PriceTimeUnit,
							Warnings:      []string{},
						},
					},
					ExtraContent: util.JSONObject{},
					Status:       "concept",
					TotalAmount:  Total,
					InvoiceType:  "standard",
				}

				data, err := json.Marshal(req)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPost, "/invoices", bytes.NewBuffer(data))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[CreateInvoiceResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
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

func TestGenerateInvoiceApi(t *testing.T) {
	clientID := createRandomBillableHours(t)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute, 1)
			},
			buildRequest: func() (*http.Request, error) {
				req := GenerateInvoiceRequest{
					ClientID:  clientID,
					StartDate: time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, time.August, 31, 23, 59, 59, 0, time.UTC),
				}
				data, err := json.Marshal(req)
				require.NoError(t, err)
				request, err := http.NewRequest(http.MethodPost, "/invoices/generate", bytes.NewBuffer(data))
				require.NoError(t, err)
				request.Header.Set("Content-Type", "application/json")
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GenerateInvoiceResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.NotEmpty(t, response.Data.InvoiceNumber)
				require.NotEmpty(t, response.Data.InvoiceDetails)

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
