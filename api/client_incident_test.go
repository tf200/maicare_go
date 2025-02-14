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

	"github.com/goccy/go-json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomClientIncident(t *testing.T, clientID int64) db.Incident {

	employee, _ := createRandomEmployee(t)
	location := createRandomLocation(t)

	arg := db.CreateIncidentParams{
		EmployeeID:              employee.ID,
		LocationID:              location.ID,
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
		EmployeeAbsenteeism:     "client",
		ClientID:                clientID,
		Emails:                  []string{"testemail@gg.com", "gaga@gog.com"},
	}

	incident, err := testStore.CreateIncident(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, incident)
	require.Equal(t, arg.EmployeeID, incident.EmployeeID)
	require.Equal(t, arg.ClientID, incident.ClientID)
	require.Equal(t, arg.Emails, incident.Emails)
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
					LocationID:              location.ID,
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
					EmployeeAbsenteeism:     "client",
					ClientID:                client.ID,
					Emails:                  []string{"testemail@gg.com", "gaga@gog.com"},
				}
				url := fmt.Sprintf("/clients/%d/incidents", client.ID)
				reqBody, err := json.Marshal(incidentReq)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				return req, nil

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var res Response[CreateIncidentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.NotEmpty(t, res.Data.ID)
				require.Equal(t, res.Data.EmployeeID, int64(1))
				require.Equal(t, res.Data.LocationID, location.ID)
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
				require.Equal(t, res.Data.Emails, []string{"testemail@gg.com", "gaga@gog.com"})
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
				require.NotEmpty(t, incidents.Data.Results[0].Emails)

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

func TestGetIncidentApi(t *testing.T) {
	client := createRandomClientDetails(t)
	incident := createRandomClientIncident(t, client.ID)
	var InformWho []string
	err := json.Unmarshal(incident.InformWho, &InformWho)
	require.NoError(t, err)
	var Technical []string
	err = json.Unmarshal(incident.Technical, &Technical)
	require.NoError(t, err)
	var Organizational []string
	err = json.Unmarshal(incident.Organizational, &Organizational)
	require.NoError(t, err)
	var MeseWorker []string
	err = json.Unmarshal(incident.MeseWorker, &MeseWorker)
	require.NoError(t, err)
	var ClientOptions []string
	err = json.Unmarshal(incident.ClientOptions, &ClientOptions)
	require.NoError(t, err)
	var Succession []string
	err = json.Unmarshal(incident.Succession, &Succession)
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/clients/%d/incidents/%d", client.ID, incident.ID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[GetIncidentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.NotEmpty(t, res.Data.ID)
				require.Equal(t, res.Data.ID, incident.ID)
				require.Equal(t, res.Data.EmployeeID, incident.EmployeeID)
				require.Equal(t, res.Data.LocationID, incident.LocationID)
				require.Equal(t, res.Data.ReporterInvolvement, incident.ReporterInvolvement)
				require.Equal(t, res.Data.InformWho, InformWho)
				require.Equal(t, res.Data.RuntimeIncident, incident.RuntimeIncident)
				require.Equal(t, res.Data.IncidentType, incident.IncidentType)
				require.Equal(t, res.Data.PassingAway, incident.PassingAway)
				require.Equal(t, res.Data.SelfHarm, incident.SelfHarm)
				require.Equal(t, res.Data.Violence, incident.Violence)
				require.Equal(t, res.Data.FireWaterDamage, incident.FireWaterDamage)
				require.Equal(t, res.Data.Accident, incident.Accident)
				require.Equal(t, res.Data.ClientAbsence, incident.ClientAbsence)
				require.Equal(t, res.Data.Medicines, incident.Medicines)
				require.Equal(t, res.Data.Organization, incident.Organization)
				require.Equal(t, res.Data.UseProhibitedSubstances, incident.UseProhibitedSubstances)
				require.Equal(t, res.Data.OtherNotifications, incident.OtherNotifications)
				require.Equal(t, res.Data.SeverityOfIncident, incident.SeverityOfIncident)
				require.NotNil(t, res.Data.IncidentExplanation)
				require.Equal(t, *res.Data.IncidentExplanation, *incident.IncidentExplanation)
				require.Equal(t, res.Data.RecurrenceRisk, incident.RecurrenceRisk)
				require.NotNil(t, res.Data.IncidentPreventSteps)
				require.Equal(t, *res.Data.IncidentPreventSteps, *incident.IncidentPreventSteps)
				require.NotNil(t, res.Data.IncidentTakenMeasures)
				require.Equal(t, *res.Data.IncidentTakenMeasures, *incident.IncidentTakenMeasures)
				require.Equal(t, res.Data.Technical, Technical)
				require.Equal(t, res.Data.Organizational, Organizational)
				require.Equal(t, res.Data.MeseWorker, MeseWorker)
				require.Equal(t, res.Data.ClientOptions, ClientOptions)
				require.NotNil(t, res.Data.OtherCause)
				require.Equal(t, *res.Data.OtherCause, *incident.OtherCause)
				require.Equal(t, res.Data.Emails, incident.Emails)
			},
		},
		{
			name: "Not Found",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/clients/%d/incidents/%d", client.ID, 0)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return request, nil
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

func TestUpdateIncidentApi(t *testing.T) {
	client := createRandomClientDetails(t)
	incident := createRandomClientIncident(t, client.ID)
	var InformWho []string
	err := json.Unmarshal(incident.InformWho, &InformWho)
	require.NoError(t, err)
	var Technical []string
	err = json.Unmarshal(incident.Technical, &Technical)
	require.NoError(t, err)
	var Organizational []string
	err = json.Unmarshal(incident.Organizational, &Organizational)
	require.NoError(t, err)
	var MeseWorker []string
	err = json.Unmarshal(incident.MeseWorker, &MeseWorker)
	require.NoError(t, err)
	var ClientOptions []string
	err = json.Unmarshal(incident.ClientOptions, &ClientOptions)
	require.NoError(t, err)
	var Succession []string
	err = json.Unmarshal(incident.Succession, &Succession)
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				incidentReq := UpdateIncidentRequest{
					ReporterInvolvement: util.StringPtr("directly_involved"),
					InformWho:           []string{"updated"},
					Emails:              []string{"taha@gmail.com"},
				}
				url := fmt.Sprintf("/clients/%d/incidents/%d", client.ID, incident.ID)
				reqBody, err := json.Marshal(incidentReq)
				require.NoError(t, err)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(reqBody))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[UpdateIncidentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.NotEmpty(t, res.Data.ID)
				require.Equal(t, res.Data.ID, incident.ID)
				require.Equal(t, res.Data.ReporterInvolvement, "directly_involved")
				require.Equal(t, res.Data.InformWho, []string{"updated"})
				require.Equal(t, res.Data.Technical, Technical)
				require.Equal(t, res.Data.Emails, []string{"taha@gmail.com"})
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
