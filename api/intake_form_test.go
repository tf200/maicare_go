package api

import (
	"bytes"
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomIntakeFormToken(t *testing.T) db.IntakeFormToken {
	arg := db.CreateIntakeFormTokenParams{
		Token: uuid.New().String(),
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

	token, err := testStore.CreateIntakeFormToken(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	return token
}

func TestGenerateIntakeFormToken(t *testing.T) {
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
				req, err := http.NewRequest(http.MethodPost, "/intake_form/token", nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
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

func TestIntakeFormUploadHandlerApi(t *testing.T) {
	token := createRandomIntakeFormToken(t)
	filename, fileContent := createRandomFile(t)

	testCases := []struct {
		name          string
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("file", filename)
				if err != nil {
					return nil, err
				}

				_, err = part.Write(fileContent)
				if err != nil {
					return nil, err
				}

				err = writer.Close()
				if err != nil {
					return nil, err
				}
				url := "/intake_form/upload?token=" + token.Token
				req, err := http.NewRequest(http.MethodPost, url, body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}

func TestVerifyIntakeFormToken(t *testing.T) {
	token := createRandomIntakeFormToken(t)

	testCases := []struct {
		name          string
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildRequest: func() (*http.Request, error) {
				url := "/intake_form/verify?token=" + token.Token
				req, err := http.NewRequest(http.MethodGet, url, nil)
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
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}

func TestCreateIntakeFormApi(t *testing.T) {
	token := createRandomIntakeFormToken(t)

	testCases := []struct {
		name          string
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildRequest: func() (*http.Request, error) {
				reqBody := CreateIntakeFormRequest{
					FirstName:                  faker.FirstName(),
					LastName:                   faker.LastName(),
					DateOfBirth:                time.Now(),
					PhoneNumber:                faker.Phonenumber(),
					Gender:                     faker.Gender(),
					PlaceOfBirth:               faker.GetRealAddress().Address,
					RepresentativeFirstName:    faker.FirstName(),
					RepresentativeLastName:     faker.LastName(),
					RepresentativePhoneNumber:  faker.Phonenumber(),
					RepresentativeEmail:        faker.Email(),
					RepresentativeRelationship: faker.Word(),
					RepresentativeAddress:      faker.GetRealAddress().Address,
					AttachementIds:             []uuid.UUID{},
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)
				url := "/intake_form?token=" + token.Token
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}
