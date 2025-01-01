package api

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomSender(t *testing.T) db.Sender {
	arg := db.CreateSenderParams{
		Types:        "main_provider",
		Name:         util.RandomString(5),
		Address:      util.StringPtr("test"),
		PostalCode:   util.StringPtr("test"),
		Place:        util.StringPtr("test"),
		Land:         util.StringPtr("test"),
		Kvknumber:    util.StringPtr("test"),
		Btwnumber:    util.StringPtr("test"),
		PhoneNumber:  util.StringPtr("test"),
		ClientNumber: util.StringPtr("test"),
		EmailAdress:  util.StringPtr("test"),
		Contacts:     []byte(`[{"name": "Test Contact", "email": "test@example.com", "phone": "1234567890"}]`),
	}

	sender, err := testStore.CreateSender(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sender)
	require.Equal(t, arg.Types, sender.Types)
	require.Equal(t, arg.Name, sender.Name)
	require.Equal(t, arg.Address, sender.Address)
	require.Equal(t, arg.PostalCode, sender.PostalCode)
	require.Equal(t, arg.Place, sender.Place)
	require.Equal(t, arg.Land, sender.Land)
	require.Equal(t, arg.Kvknumber, sender.Kvknumber)
	require.Equal(t, arg.Btwnumber, sender.Btwnumber)
	require.Equal(t, arg.PhoneNumber, sender.PhoneNumber)
	require.Equal(t, arg.ClientNumber, sender.ClientNumber)
	require.Equal(t, arg.EmailAdress, sender.EmailAdress)
	require.Equal(t, arg.Contacts, sender.Contacts)
	return sender
}

func TestCreateSender(t *testing.T) {
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
				createSenderReq := CreateSenderRequest{
					Types:        "main_provider",
					Name:         "Test Company",
					Address:      nil,
					PostalCode:   util.StringPtr("1234 AB"),
					Place:        util.StringPtr("Amsterdam"),
					Land:         util.StringPtr("Netherlands"),
					KVKNumber:    util.StringPtr("12345678"),
					BTWNumber:    util.StringPtr("NL123456789B01"),
					PhoneNumber:  util.StringPtr("+31612345678"),
					ClientNumber: util.StringPtr("CLI123"),
					Contacts: []Contact{
						{
							Name:        util.StringPtr("John Doe"),
							Email:       util.StringPtr("john@example.com"),
							PhoneNumber: util.StringPtr("+31612345678"),
						},
					},
				}

				data, err := json.Marshal(createSenderReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/senders", bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response CreateSenderResponse
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.NotEmpty(t, response.ID)
				require.Equal(t, "Test Company", response.Name)
			},
		},
		{
			name: "InvalidType",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createSenderReq := CreateSenderRequest{
					Types: "invalid_type",
					Name:  "Test Company",
					Contacts: []Contact{
						{
							Name:        util.StringPtr("John Doe"),
							Email:       util.StringPtr("john@example.com"),
							PhoneNumber: util.StringPtr("+31612345678"),
						},
					},
				}

				data, err := json.Marshal(createSenderReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/senders", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createSenderReq := CreateSenderRequest{
					Types:        "main_provider",
					Name:         "Test Company",
					Address:      util.StringPtr("Test Street 123"),
					PostalCode:   util.StringPtr("1234 AB"),
					Place:        util.StringPtr("Amsterdam"),
					Land:         util.StringPtr("Netherlands"),
					KVKNumber:    util.StringPtr("12345678"),
					BTWNumber:    util.StringPtr("NL123456789B01"),
					PhoneNumber:  util.StringPtr("+31612345678"),
					ClientNumber: util.StringPtr("CLI123"),
					Contacts: []Contact{
						{
							Name:        util.StringPtr("John Doe"),
							Email:       util.StringPtr("invalid-email"),
							PhoneNumber: util.StringPtr("+31612345678"),
						},
					},
				}

				data, err := json.Marshal(createSenderReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/senders", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingRequiredField",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createSenderReq := CreateSenderRequest{
					Types: "main_provider",
					Name:  "", // Required field is empty
					Contacts: []Contact{
						{
							Name:        util.StringPtr("John Doe"),
							Email:       util.StringPtr("john@example.com"),
							PhoneNumber: util.StringPtr("+31612345678"),
						},
					},
				}

				data, err := json.Marshal(createSenderReq)
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/senders", bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

func TestListSendersAPI(t *testing.T) {
	user := createRandomUser(t)
	initialCount, err := testStore.CountSenders(context.Background(), util.BoolPtr(true))
	require.NoError(t, err)
	numSenders := 20
	var wg sync.WaitGroup
	for i := 0; i < numSenders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomSender(t)
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
				url := "/senders?page=1&page_size=10"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response pagination.Response[ListSendersResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotNil(t, response.Next)
				require.Nil(t, response.Previous)
				require.Equal(t, int64(numSenders)+initialCount, response.Count)
				require.Equal(t, int32(10), response.PageSize)
				require.Len(t, response.Results, 10)

				for _, sender := range response.Results {
					require.NotEmpty(t, sender.ID)
					require.NotEmpty(t, sender.Name)
				}
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Don't add authorization
			},
			buildRequest: func() (*http.Request, error) {
				url := "/senders?page=1&page_size=10"
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
				url := "/senders?page=1&page_size=10"
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
				url := "/senders?page=1&page_size=hh"
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "SecondPage",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/senders?page=2&page_size=10"
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
