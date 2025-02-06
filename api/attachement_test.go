package api

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"maicare_go/token"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/stretchr/testify/require"
)

func createRandomFile(t *testing.T) (string, []byte) {
	filename := fmt.Sprintf("test-%d.pdf", time.Now().UnixNano())

	// Create a minimal valid PDF file
	content := []byte("%PDF-1.5\n%¥±ë\n\n1 0 obj\n<</Type/Catalog/Pages 2 0 R>>\nendobj")

	// Add some random content after the PDF header
	randomContent := make([]byte, 100)
	_, err := rand.Read(randomContent)
	require.NoError(t, err)

	// Combine the PDF header with random content
	content = append(content, randomContent...)

	return filename, content
}

func TestUploadHandler(t *testing.T) {
	filename, fileContent := createRandomFile(t)

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

				req, err := http.NewRequest(http.MethodPost, "/attachments/upload", body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response Response[UploadHandlerResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data.FileURL)
				require.NotEmpty(t, response.Data.FileID)
				require.NotEmpty(t, response.Data.CreatedAt)
				require.Equal(t, "File uploaded successfully", response.Message)

			},
		},
		{
			name: "NoFile",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				err := writer.Close()
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/attachments/upload", body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "EmptyFile",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				_, err := writer.CreateFormFile("file", "empty.pdf")
				if err != nil {
					return nil, err
				}

				err = writer.Close()
				if err != nil {
					return nil, err
				}

				req, err := http.NewRequest(http.MethodPost, "/attachments/upload", body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
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
			recorder := httptest.NewRecorder()
			request, err := tc.buildRequest()
			require.NoError(t, err)

			tc.setupAuth(t, request, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetAttachmentById(t *testing.T) {
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
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/attachments/%s", file.Uuid), nil)
				t.Log(req)
				if err != nil {
					return nil, err
				}
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response Response[GetAttachmentByIdResponse]
				err := json.NewDecoder(recorder.Body).Decode(&response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, file.Uuid, response.Data.FileID)
			},
		},
		{
			name: "NotFound",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				req, err := http.NewRequest(http.MethodGet, "/attachment/00000000-0000-0000-0000-000000000000", nil)
				if err != nil {
					return nil, err
				}
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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
