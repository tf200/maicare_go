package api

import (
	"maicare_go/pagination"
	clientp "maicare_go/service/client"
	"maicare_go/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestListAllIncidentsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 5; i++ {
		createRandomClientIncident(t, client.ID)
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, client.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/incidents"
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[pagination.Response[clientp.ListIncidentsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response)
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
