package api

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/goccy/go-json"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func createRandomSender(t *testing.T) db.Sender {
	// Define a slice of Contact structs
	contacts := []SenderContact{
		{
			Name:        util.StringPtr(faker.Name()),
			Email:       util.StringPtr(faker.Email()),
			PhoneNumber: util.StringPtr(faker.Phonenumber()),
		},
	}

	// Marshal the contacts slice into JSON
	contactsJSON, err := json.Marshal(contacts)
	require.NoError(t, err)

	// Create the CreateSenderParams
	arg := db.CreateSenderParams{
		Types:        "main_provider",
		Name:         faker.FirstName(),
		Address:      util.StringPtr(faker.GetRealAddress().City),
		PostalCode:   util.StringPtr(faker.GetRealAddress().PostalCode),
		Place:        util.StringPtr("test"),
		Land:         util.StringPtr("test"),
		Kvknumber:    util.StringPtr("test"),
		Btwnumber:    util.StringPtr("test"),
		PhoneNumber:  util.StringPtr("test"),
		ClientNumber: util.StringPtr("test"),
		EmailAddress: util.StringPtr("test"),
		Contacts:     contactsJSON, // Use the marshaled JSON
	}

	// Create the sender in the database
	sender, err := testStore.CreateSender(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sender)

	// Verify the fields
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
	require.Equal(t, arg.EmailAddress, sender.EmailAddress)

	// Unmarshal the expected and actual Contacts fields for comparison
	var expectedContacts []SenderContact
	err = json.Unmarshal(arg.Contacts, &expectedContacts)
	require.NoError(t, err)

	var actualContacts []SenderContact
	err = json.Unmarshal(sender.Contacts, &actualContacts)
	require.NoError(t, err)

	// Compare the unmarshaled Contacts
	require.Equal(t, expectedContacts, actualContacts)

	return sender
}
func TestCreateSenderApi(t *testing.T) {
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
					Contacts: []SenderContact{
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

				var response Response[CreateSenderResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.ID)
				require.Equal(t, "Test Company", response.Data.Name)
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
					Contacts: []SenderContact{
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
					Contacts: []SenderContact{
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
					Contacts: []SenderContact{
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

				var response Response[pagination.Response[ListSendersResponse]]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotNil(t, response.Data.Next)
				require.Nil(t, response.Data.Previous)
				require.Equal(t, int64(numSenders)+initialCount, response.Data.Count)
				require.Equal(t, int32(10), response.Data.PageSize)
				require.Len(t, response.Data.Results, 10)

				for _, sender := range response.Data.Results {
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

				var response Response[pagination.Response[db.ListEmployeeProfileRow]]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)

				require.NotNil(t, response.Data.Previous)
				require.Contains(t, *response.Data.Previous, "page=1")
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

func TestUpdateSenderApi(t *testing.T) {
	sender := createRandomSender(t)
	contacts := make([]SenderContact, 0)
	err := json.Unmarshal(sender.Contacts, &contacts)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, sender.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateSenderReq := UpdateSenderRequest{
					Name: util.StringPtr("Updated Company2"),
				}
				data, err := json.Marshal(updateSenderReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/senders/%d", sender.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateSenderResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.Equal(t, sender.ID, response.Data.ID)
				require.Equal(t, "Updated Company2", response.Data.Name)
				require.Equal(t, sender.Types, response.Data.Types)

				require.Equal(t, contacts, response.Data.Contacts)
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

func TestGetSenderByIdAPI(t *testing.T) {
	sender := createRandomSender(t)
	contacts := make([]SenderContact, 0)
	err := json.Unmarshal(sender.Contacts, &contacts)
	require.NoError(t, err)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, sender.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/senders/%d", sender.ID)
				return http.NewRequest(http.MethodGet, url, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetSenderByIdResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, sender.ID, response.Data.ID)
				require.Equal(t, sender.Name, response.Data.Name)
				require.Equal(t, sender.Types, response.Data.Types)
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
