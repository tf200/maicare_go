package api

import (
	"bytes"
	"context"
	"encoding/json"
	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomAttachmentFile(t *testing.T) db.AttachmentFile {

	tagvalue := "test"
	arg := db.CreateAttachmentParams{
		Name: util.RandomString(5),
		File: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
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

func TestCreateClientApi(t *testing.T) {
	var filesUuids [10]uuid.UUID
	for i := 0; i < 10; i++ {
		filesUuids[i] = createRandomAttachmentFile(t).Uuid

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
				clientReq := CreateClientDetailsRequest{
					FirstName:    util.RandomString(5),
					LastName:     util.RandomString(8),
					Email:        util.RandomEmail(),
					Organisation: util.StringPtr("Test Organisation"),
					LocationID:   util.IntPtr(1),
					LegalMeasure: util.StringPtr("Test Legal Measure"),
					Birthplace:   util.StringPtr("Test Birthplace"),
					Departement:  util.StringPtr("Test Departement"),
					Gender:       "male",
					Filenumber:   util.RandomString(5),
					DateOfBirth:  "2006-01-02",
					PhoneNumber:  util.StringPtr("1234567890"),
					Infix:        util.StringPtr("Test Infix"),
					Source:       util.StringPtr("Test Source"),
					Bsn:          util.StringPtr("Test Bsn"),
					Addresses: []Address{
						{
							BelongsTo:   util.StringPtr("Test Belongs To"),
							Address:     util.StringPtr("Test Address"),
							City:        util.StringPtr("Test City"),
							ZipCode:     util.StringPtr("12345"),
							PhoneNumber: util.StringPtr("1234567890"),
						},
					},
					IdentityAttachmentIds: filesUuids[:],
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
				require.Equal(t, http.StatusOK, recorder.Code)
				var clientRes Response[CreateClientDetailsResponse]
				err := json.NewDecoder(recorder.Body).Decode(&clientRes)
				require.NoError(t, err)
				require.NotEmpty(t, clientRes.Data)
				require.NotEmpty(t, clientRes.Data.ID)
				require.Equal(t, clientRes.Data.IdentityAttachmentIds, filesUuids[:])
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
