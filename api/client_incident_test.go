package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/token"
	"maicare_go/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomClientIncident(t *testing.T, clientID int64) db.Incident {

	employee, _ := createRandomEmployee(t)
	location := createRandomLocation(t)

	arg := db.CreateIncidentParams{
		EmployeeID:              employee.ID,
		LocationID:              &location.ID,
		ReporterInvolvement:     "directly_involved",
		InformWho:               []byte("[\"client\"]"),
		IncidentDate:            pgtype.Date{Time: time.Now(), Valid: true},
		RuntimeIncident:         "no",
		IncidentType:            "accident",
		PassingAway:             false,
		SelfHarm:                false,
		Violence:                false,
		FireWaterDamage:         false,
		Accident:                false,
		ClientAbsence:           false,
		Medicines:               false,
		Organization:            false,
		UseProhibitedSubstances: false,
		OtherNotifications:      false,
		SeverityOfIncident:      "fatal",
		IncidentExplanation:     util.StringPtr("test explanation"),
		RecurrenceRisk:          "high",
		IncidentPreventSteps:    util.StringPtr("test steps"),
		IncidentTakenMeasures:   util.StringPtr("test measures"),
		Technical:               []byte("[\"client\"]"),
		Organizational:          []byte("[\"client\"]"),
		MeseWorker:              []byte("[\"client\"]"),
		ClientOptions:           []byte("[\"client\"]"),
		OtherCause:              util.StringPtr("test cause"),
		CauseExplanation:        util.StringPtr("test cause explanation"),
		PhysicalInjury:          "no_injuries",
		PhysicalInjuryDesc:      util.StringPtr("test injuries"),
		PsychologicalDamage:     "other",
		PsychologicalDamageDesc: util.StringPtr("test damage"),
		NeededConsultation:      "no",
		Succession:              []byte("[\"client\"]"),
		SuccessionDesc:          util.StringPtr("test succession"),
		Other:                   false,
		OtherDesc:               util.StringPtr("test other"),
		AdditionalAppointments:  util.StringPtr("test appointments"),
		EmployeeAbsenteeism:     []byte("[\"client\"]"),
		ClientID:                clientID,
	}

	incident, err := testStore.CreateIncident(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, incident)
	require.Equal(t, arg.EmployeeID, incident.EmployeeID)
	require.Equal(t, arg.ClientID, incident.ClientID)
	return incident
}

func TestCreateIncident(t *testing.T) {
	client := createRandomClientDetails(t)
	location := createRandomLocation(t)

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
				incidentReq := CreateIncidentRequest{
					EmployeeID:              1,
					LocationID:              &location.ID,
					ReporterInvolvement:     "directly_involved",
					InformWho:               []string{"client"},
					IncidentDate:            time.Now(),
					RuntimeIncident:         "no",
					IncidentType:            "accident",
					PassingAway:             false,
					SelfHarm:                false,
					Violence:                false,
					FireWaterDamage:         false,
					Accident:                false,
					ClientAbsence:           false,
					Medicines:               false,
					Organization:            false,
					UseProhibitedSubstances: false,
					OtherNotifications:      false,
					SeverityOfIncident:      "fatal",
					IncidentExplanation:     util.StringPtr("test explanation"),
					RecurrenceRisk:          "high",
					IncidentPreventSteps:    util.StringPtr("test steps"),
					IncidentTakenMeasures:   util.StringPtr("test measures"),
					Technical:               []string{"client"},
					Organizational:          []string{"client"},
					MeseWorker:              []string{"client"},
					ClientOptions:           []string{"client"},
					OtherCause:              util.StringPtr("test cause"),
					CauseExplanation:        util.StringPtr("test cause explanation"),
					PhysicalInjury:          "no_injuries",
					PhysicalInjuryDesc:      util.StringPtr("test injuries"),
					PsychologicalDamage:     "other",
					PsychologicalDamageDesc: util.StringPtr("test damage"),
					NeededConsultation:      "no",
					Succession:              []string{"client"},
					SuccessionDesc:          util.StringPtr("test succession"),
					Other:                   false,
					OtherDesc:               util.StringPtr("test other"),
					AdditionalAppointments:  util.StringPtr("test appointments"),
					EmployeeAbsenteeism:     []string{"client"},
					ClientID:                client.ID,
				}
				url := fmt.Sprintf("/clients/%d/incidents", client.ID)
				reqBody, err := json.Marshal(incidentReq)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				return req, nil

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				var res Response[CreateIncidentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.NotEmpty(t, res.Data.ID)
				require.Equal(t, res.Data.EmployeeID, int64(1))
				require.Equal(t, *res.Data.LocationID, location.ID)
				require.Equal(t, res.Data.ReporterInvolvement, "directly_involved")
				require.Equal(t, res.Data.InformWho, []string{"client"})
				require.Equal(t, res.Data.RuntimeIncident, "no")
				require.Equal(t, res.Data.IncidentType, "accident")
				require.Equal(t, res.Data.PassingAway, false)
				require.Equal(t, res.Data.SelfHarm, false)
				require.Equal(t, res.Data.Violence, false)
				require.Equal(t, res.Data.FireWaterDamage, false)
				require.Equal(t, res.Data.Accident, false)
				require.Equal(t, res.Data.ClientAbsence, false)
				require.Equal(t, res.Data.Medicines, false)
				require.Equal(t, res.Data.Organization, false)
				require.Equal(t, res.Data.UseProhibitedSubstances, false)
				require.Equal(t, res.Data.OtherNotifications, false)
				require.Equal(t, res.Data.SeverityOfIncident, "fatal")
				require.NotNil(t, res.Data.IncidentExplanation)
				require.Equal(t, *res.Data.IncidentExplanation, "test explanation")
				require.Equal(t, res.Data.RecurrenceRisk, "high")
				require.NotNil(t, res.Data.IncidentPreventSteps)
				require.Equal(t, *res.Data.IncidentPreventSteps, "test steps")
				require.NotNil(t, res.Data.IncidentTakenMeasures)
				require.Equal(t, *res.Data.IncidentTakenMeasures, "test measures")
				require.Equal(t, res.Data.Technical, []string{"client"})
				require.Equal(t, res.Data.Organizational, []string{"client"})
				require.Equal(t, res.Data.MeseWorker, []string{"client"})
				require.Equal(t, res.Data.ClientOptions, []string{"client"})
				require.NotNil(t, res.Data.OtherCause)
				require.Equal(t, *res.Data.OtherCause, "test cause")
				require.NotNil(t, res.Data.CauseExplanation)
				require.Equal(t, *res.Data.CauseExplanation, "test cause explanation")
				require.Equal(t, res.Data.PhysicalInjury, "no_injuries")
				require.NotNil(t, res.Data.PhysicalInjuryDesc)
				require.Equal(t, *res.Data.PhysicalInjuryDesc, "test injuries")
				require.Equal(t, res.Data.PsychologicalDamage, "other")
				require.NotNil(t, res.Data.PsychologicalDamageDesc)
				require.Equal(t, *res.Data.PsychologicalDamageDesc, "test damage")
				require.Equal(t, res.Data.NeededConsultation, "no")
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

func TestListIncidentsApi(t *testing.T) {
	client := createRandomClientDetails(t)

	for i := 0; i < 20; i++ {
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/clients/%d/incidents?page=1&page_size=10", client.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var incidents Response[pagination.Response[ListIncidentsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &incidents)
				require.NoError(t, err)
				require.NotEmpty(t, incidents)
				require.NotEmpty(t, incidents.Data.Results)
				require.Len(t, incidents.Data.Results, 10)

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
