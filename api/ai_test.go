package api

import (
	"bytes"
	"maicare_go/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestSpellCheckApi(t *testing.T) {
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
				reqBody := CorrectSpellingRequest{
					InitialText: "The qwick brown fox jumpped over the laazy dog. It was a beutiful day, but sudenly, the wheather changed dramaticaly. Peopel ran for sheltr as the rain began pooring down. In the distanse, a lound thunder clap made everyonne jump. This storm came out of no where! exclaimed a passerby.",
				}
				jsonBody, err := json.Marshal(reqBody)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPost, "/ai/spelling_check", bytes.NewBuffer(jsonBody))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[CorrectSpellingResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.Equal(t, "Spelling check successful", res.Message)
				require.NotEmpty(t, res.Data.CorrectedText)
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
