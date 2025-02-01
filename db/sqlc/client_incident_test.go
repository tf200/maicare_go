package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomClientIncident(t *testing.T, clientID int64) Incident {

	employee, _ := createRandomEmployee(t)
	location := CreateRandomLocation(t)

	arg := CreateIncidentParams{
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
		EmployeeAbsenteeism:     "client",
		ClientID:                clientID,
	}

	incident, err := testQueries.CreateIncident(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, incident)
	require.Equal(t, arg.EmployeeID, incident.EmployeeID)
	require.Equal(t, arg.ClientID, incident.ClientID)
	return incident
}
func TestCreateIncident(t *testing.T) {
	client := createRandomClientDetails(t)
	createRandomClientIncident(t, client.ID)
}

func TestListIncidents(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 4; i++ {
		createRandomClientIncident(t, client.ID)
	}

	incidents, err := testQueries.ListIncidents(context.Background(), ListIncidentsParams{
		Limit:    5,
		Offset:   0,
		ClientID: client.ID,
	})
	require.NoError(t, err)

	for _, incident := range incidents {
		require.NotEmpty(t, incident)
	}
	t.Log(incidents)
	require.Len(t, incidents, 4)
}

func TestGetIncident(t *testing.T) {
	client := createRandomClientDetails(t)
	incident1 := createRandomClientIncident(t, client.ID)
	incident2, err := testQueries.GetIncident(context.Background(), incident1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, incident2)
	require.Equal(t, incident1.ID, incident2.ID)
}

func TestUpdateIncident(t *testing.T) {
	client := createRandomClientDetails(t)
	incident1 := createRandomClientIncident(t, client.ID)
	arg := UpdateIncidentParams{
		ID:           incident1.ID,
		IncidentType: util.StringPtr("test incident"),
	}
	incident2, err := testQueries.UpdateIncident(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, incident2)
	require.Equal(t, incident1.ID, incident2.ID)
	require.NotEqual(t, incident1.IncidentType, incident2.IncidentType)
}

func TestDeleteIncident(t *testing.T) {
	client := createRandomClientDetails(t)
	incident1 := createRandomClientIncident(t, client.ID)
	_, err := testQueries.DeleteIncident(context.Background(), incident1.ID)
	require.NoError(t, err)
	incident2, err := testQueries.GetIncident(context.Background(), incident1.ID)
	require.Error(t, err)
	require.Empty(t, incident2)
}
