
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

func TestCreateAppointmentCardApi(t *testing.T) {
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, client.ID, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				appointmentReq := CreateAppointmentCardRequest{
					GeneralInformation:     []string{"Client is doing well", "No concerns raised"},
					ImportantContacts:      []string{"Mother - 555-123-4567", "Case Worker - 555-987-6543"},
					HouseholdInfo:          []string{"Lives with mother and younger sibling", "Stable home environment"},
					OrganizationAgreements: []string{"Signed agreement with YMCA", "Participating in job training program"},
					YouthOfficerAgreements: []string{"Compliant with curfew", "Attending mandatory check-ins"},
					TreatmentAgreements:    []string{"Attending therapy sessions", "Taking prescribed medication"},
					SmokingRules:           []string{"Agreed to smoke outside only", "Limited to 5 cigarettes per day"},
					Work:                   []string{"Employed at local grocery store", "Works 20 hours per week"},
					SchoolInternship:       []string{"Interning at a law firm", "Gaining valuable experience"},
					Travel:                 []string{"Planning a trip to visit family", "Needs assistance with travel arrangements"},
					Leave:                  []string{"Took a week off for vacation", "Returned to work on August 1st"},
				}
				data, err := json.Marshal(appointmentReq)
				require.NoError(t, err)

				url := fmt.Sprintf("/clients/%d/appointment_cards", client.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				var appointmentCard Response[CreateAppointmentCardResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &appointmentCard)
				require.NoError(t, err)
				require.NotEmpty(t, appointmentCard.Data)
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

// TODO: Add tests for other API endpoints
