package api

import (
	"bytes"
	"fmt"
	"maicare_go/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestCreateClientMaturityMatrixAssessmentApi(t *testing.T) {
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
				assessmentReq := CreateClientMaturityMatrixAssessmentRequest{
					Assessments: []MatrixAssessment{
						{
							MaturityMatrixID: 1,
							InitialLevel:     2,
							StartDate:        time.Now(),
							EndDate:          time.Now().Add(time.Hour * 24 * 7),
						},
						{
							MaturityMatrixID: 2,
							InitialLevel:     3,
							StartDate:        time.Now(),
							EndDate:          time.Now().Add(time.Hour * 24 * 7),
						},
					},
				}
				data, err := json.Marshal(assessmentReq)
				require.NoError(t, err)

				url := fmt.Sprintf("/clients/%d/maturity_matrix_assessment", client.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var assessmentCard Response[CreateClientMaturityMatrixAssessmentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &assessmentCard)
				require.NoError(t, err)
				require.NotEmpty(t, assessmentCard.Data)
				require.Equal(t, client.ID, assessmentCard.Data.ClientID)
				require.Len(t, assessmentCard.Data.Assessments, 2)
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

// TODO WRITE THE TESTS DONT FORGET TO ADD THE TESTS TO THE TEST SUITE
