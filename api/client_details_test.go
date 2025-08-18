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

func createRandomAttachmentFile(t *testing.T) db.AttachmentFile {

	tagvalue := "test"
	arg := db.CreateAttachmentParams{
		Name: util.RandomString(5),
		File: util.GetRandomImageURL(),
		Size: 23,
		Tag:  &tagvalue,
	}
	attachment, err := testStore.CreateAttachment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, attachment)

	require.Equal(t, arg.Name, attachment.Name)
	require.Equal(t, arg.File, attachment.File)
	require.Equal(t, arg.Size, attachment.Size)
	require.Equal(t, arg.Tag, attachment.Tag)

	require.NotZero(t, attachment.Uuid)
	require.NotZero(t, attachment.Created)
	return attachment
}

func createRandomClientDetails(t *testing.T) db.ClientDetail {
	location := createRandomLocation(t)
	employee, _ := createRandomEmployee(t)
	sender := createRandomSender(t)

	arg := db.CreateClientDetailsParams{
		FirstName:                  faker.FirstName(),
		LastName:                   faker.LastName(),
		Email:                      faker.Email(),
		PhoneNumber:                util.StringPtr(faker.Phonenumber()),
		DateOfBirth:                pgtype.Date{Time: time.Now().AddDate(-20, 0, 0), Valid: true},
		Identity:                   false,
		Bsn:                        util.StringPtr(util.RandomString(9)),
		BsnVerifiedBy:              &employee.ID,
		Source:                     util.StringPtr("Test Source"),
		Birthplace:                 util.StringPtr("test city"),
		Organisation:               util.StringPtr("test org"),
		Departement:                util.StringPtr("test dep"),
		Gender:                     "male",
		Filenumber:                 "testfile",
		ProfilePicture:             util.StringPtr(util.GetRandomImageURL()),
		Infix:                      util.StringPtr("van"),
		SenderID:                   &sender.ID,
		LocationID:                 util.IntPtr(location.ID),
		DepartureReason:            util.StringPtr("test Reason"),
		DepartureReport:            util.StringPtr("test report"),
		Addresses:                  []byte("[]"),
		LegalMeasure:               util.StringPtr("test measure"),
		EducationCurrentlyEnrolled: false,
		EducationInstitution:       util.StringPtr("test institution"),
		EducationMentorName:        util.StringPtr("test mentor"),
		EducationMentorEmail:       util.StringPtr("test mentor email"),
		EducationMentorPhone:       util.StringPtr("test mentor phone"),
		EducationAdditionalNotes:   util.StringPtr("test additional notes"),
		EducationLevel:             util.StringPtr("test education level"),
		WorkCurrentlyEmployed:      false,
		WorkCurrentEmployer:        util.StringPtr("test employer"),
		WorkCurrentEmployerPhone:   util.StringPtr("test employer phone"),
		WorkCurrentEmployerEmail:   util.StringPtr("test employer email"),
		WorkCurrentPosition:        util.StringPtr("test position"),
		WorkStartDate:              pgtype.Date{Time: time.Now(), Valid: true},
		WorkAdditionalNotes:        util.StringPtr("test work additional notes"),
		LivingSituation:            util.StringPtr("test living situation"),
		LivingSituationNotes:       util.StringPtr("test living situation notes"),
	}

	client, err := testStore.CreateClientDetails(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, client)
	require.Equal(t, arg.FirstName, client.FirstName)
	require.Equal(t, arg.LastName, client.LastName)
	require.Equal(t, arg.Email, client.Email)
	require.Equal(t, arg.PhoneNumber, client.PhoneNumber)
	require.Equal(t, arg.Filenumber, client.Filenumber)
	require.Equal(t, arg.ProfilePicture, client.ProfilePicture)
	require.Equal(t, arg.Infix, client.Infix)
	require.Equal(t, arg.SenderID, client.SenderID)
	require.Equal(t, arg.LocationID, client.LocationID)
	require.Equal(t, arg.DepartureReason, client.DepartureReason)
	require.Equal(t, arg.DepartureReport, client.DepartureReport)
	require.Equal(t, arg.Addresses, client.Addresses)
	require.Equal(t, arg.LegalMeasure, client.LegalMeasure)
	return client
}

func TestCreateClientApi(t *testing.T) {

	sender := createRandomSender(t)
	location := createRandomLocation(t)
	employee, _ := createRandomEmployee(t)
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
				clientReq := CreateClientDetailsRequest{
					FirstName:     faker.FirstName(),
					LastName:      faker.LastName(),
					Email:         faker.Email(),
					Organisation:  util.StringPtr("Test Organisation"),
					LocationID:    &location.ID,
					LegalMeasure:  util.StringPtr("Test Legal Measure"),
					Birthplace:    util.StringPtr("Test Birthplace"),
					Departement:   util.StringPtr("Test Departement"),
					Gender:        "male",
					Filenumber:    util.RandomString(5),
					DateOfBirth:   "2006-01-02",
					PhoneNumber:   util.StringPtr("1234567890"),
					Infix:         util.StringPtr("Test Infix"),
					Source:        util.StringPtr("Test Sources"),
					Bsn:           util.StringPtr("Test Bsn"),
					BsnVerifiedBy: &employee.ID,
					Addresses: []Address{
						{
							BelongsTo:   util.StringPtr("Test Belongs To"),
							Address:     util.StringPtr("Test Address"),
							City:        util.StringPtr("Test City"),
							ZipCode:     util.StringPtr("12345"),
							PhoneNumber: util.StringPtr("1234567890"),
						},
					},
					SenderID:                   &sender.ID,
					EducationInstitution:       util.StringPtr("Test Institution"),
					EducationMentorName:        util.StringPtr("Test Mentor"),
					EducationMentorPhone:       util.StringPtr("1234567890"),
					EducationMentorEmail:       util.StringPtr("test@example.com"),
					EducationCurrentlyEnrolled: true,
					EducationAdditionalNotes:   util.StringPtr("Test Additional Notes"),
					EducationLevel:             util.StringPtr("primary"),
					WorkCurrentlyEmployed:      true,
					WorkCurrentEmployer:        util.StringPtr("Test Employer"),
					WorkCurrentEmployerPhone:   util.StringPtr(faker.Phonenumber()),
					WorkCurrentEmployerEmail:   util.StringPtr(faker.Email()),
					WorkCurrentPosition:        util.StringPtr("Test Position"),
					WorkStartDate:              time.Now(),
					WorkAdditionalNotes:        util.StringPtr("Test Work Additional Notes"),
					LivingSituation:            util.StringPtr("home"),
					LivingSituationNotes:       util.StringPtr("home"),
				}
				reqBody, err := json.Marshal(clientReq)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPost, "/clients", bytes.NewReader(reqBody))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Logf("Response Status Code: %d", recorder.Code)
				t.Logf("Raw Response Body: %s", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var clientRes Response[CreateClientDetailsResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
				require.NotEmpty(t, clientRes.Data.ID)
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

func TestListClient(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomClientDetails(t)
	}
	testCaes := []struct {
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
				url := "/clients?page=1&page_size=5"
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var clients Response[pagination.Response[ListClientsApiResponse]]
				err := json.NewDecoder(recorder.Body).Decode(&clients)
				require.NoError(t, err)
				require.NotEmpty(t, clients.Data)
				require.Len(t, clients.Data.Results, 5)
			},
		},
		{
			name: "Invalid Page Size",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/clients?page=1&page_size=101"
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Page Number",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/clients?page=0&page_size=5"
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCaes {
		tc := testCaes[i]
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

func TestGetClientDetails(t *testing.T) {
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
				url := fmt.Sprintf("/clients/%d", client.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[GetClientApiResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
				require.Equal(t, client.ID, clientRes.Data.ID)
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

func TestUpdateClientDetailsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	location := createRandomLocation(t)
	_ = createRandomSender(t)
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
				clientReq := UpdateClientDetailsRequest{
					FirstName:    util.StringPtr(faker.FirstName()),
					LastName:     util.StringPtr(faker.LastName()),
					Email:        util.StringPtr(faker.Email()),
					Organisation: util.StringPtr("Test Organisation"),
					LocationID:   &location.ID,
					LegalMeasure: util.StringPtr("Test Legal Measure"),
					Birthplace:   util.StringPtr("Test Birthplace"),
					Departement:  util.StringPtr("Test Departement"),
				}
				reqBody, err := json.Marshal(clientReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/clients/%d", client.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Logf("Response Status Code: %d", recorder.Code)
				t.Logf("Raw Response Body: %s", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[UpdateClientDetailsResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
				require.NotEmpty(t, clientRes.Data.ID)
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

func TestSetClientProfilePictureApi(t *testing.T) {
	client := createRandomClientDetails(t)
	file := createRandomAttachmentFile(t)
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
				reqBody := SetClientProfilePictureRequest{
					AttachmentID: file.Uuid,
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)

				url := fmt.Sprintf("/clients/%d/profile_picture", client.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[SetClientProfilePictureResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
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

func TestAddClientDocumentApi(t *testing.T) {
	client := createRandomClientDetails(t)
	file := createRandomAttachmentFile(t)
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
				reqBody := AddClientDocumentApiRequest{
					AttachmentID: file.Uuid,
					Label:        "other",
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)

				url := fmt.Sprintf("/clients/%d/documents", client.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				var clientRes Response[AddClientDocumentApiResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
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

func addRandomClientDocument(t *testing.T, ClientID int64) db.ClientDocument {

	attachment := createRandomAttachmentFile(t)

	arg := db.AddClientDocumentTxParams{
		ClientID:     ClientID,
		AttachmentID: attachment.Uuid,
		Label:        "registration_form",
	}

	clientDoc, err := testStore.AddClientDocumentTx(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, clientDoc)
	require.Equal(t, arg.ClientID, clientDoc.ClientDocument.ClientID)
	return clientDoc.ClientDocument
}

func TestListClientDocumentsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 10; i++ {
		addRandomClientDocument(t, client.ID)
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
				url := fmt.Sprintf("/clients/%d/documents?page=1&page_size=5", client.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[pagination.Response[ListClientDocumentsApiResponse]]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
				require.Len(t, clientRes.Data.Results, 5)
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

func TestDeleteClientDocumentApi(t *testing.T) {
	cleint := createRandomClientDetails(t)
	clientDoc := addRandomClientDocument(t, cleint.ID)

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
				url := fmt.Sprintf("/clients/%d/documents/%d", cleint.ID, clientDoc.ID)
				data := DeleteClientDocumentApiRequest{
					AttachmentID: *clientDoc.AttachmentUuid,
				}
				reqBody, err := json.Marshal(data)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(reqBody))
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

func TestGetMissingClientDocumentsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	clientDoc := addRandomClientDocument(t, client.ID)
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
				url := fmt.Sprintf("/clients/%d/missing_documents", client.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[GetMissingClientDocumentsApiResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
				require.NotEmpty(t, clientRes.Data.MissingDocs)
				require.NotContains(t, clientRes.Data.MissingDocs, clientDoc.Label)

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

func TestUpdateClientStatusApi(t *testing.T) {
	client := createRandomClientDetails(t)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK No Scheduling",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				reqBody := UpdateClientStatusRequest{
					Status:       "In Care",
					Reason:       "Test Reason",
					IsSchedueled: false,
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)

				url := fmt.Sprintf("/clients/%d/status", client.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[UpdateClientStatusResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
			},
		},
		{
			name: "OK With Scheduling",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				scheduledTime := time.Date(2028, 1, 2, 0, 0, 0, 0, time.UTC)
				reqBody := UpdateClientStatusRequest{
					Status:        "In Care",
					Reason:        "Test Reason",
					IsSchedueled:  true,
					SchedueledFor: scheduledTime,
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)

				url := fmt.Sprintf("/clients/%d/status", client.ID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[UpdateClientStatusResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
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
