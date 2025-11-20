package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateIncident(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateIncidentParams
		checks func(t *testing.T, incident CreateIncidentRow, err error)
	}{
		{
			name: "successful creation with minimal required fields",
			setup: func(ctx context.Context, qtx *Queries) CreateIncidentParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				client := createRandomClientDetails(ctx, qtx)

				return CreateIncidentParams{
					EmployeeID:              employee.ID,
					LocationID:              location.ID,
					ReporterInvolvement:     IncidentReporterInvolvementEnumAlarmed,
					InformWho:               []string{"Manager", "HR"},
					IncidentDate:            pgtype.Date{Time: time.Now(), Valid: true},
					RuntimeIncident:         "Incident occurred during work hours",
					IncidentType:            "Workplace Incident",
					PassingAway:             false,
					SelfHarm:                false,
					Violence:                false,
					FireWaterDamage:         false,
					Accident:                true,
					ClientAbsence:           false,
					Medicines:               false,
					Organization:            false,
					UseProhibitedSubstances: false,
					OtherNotifications:      false,
					SeverityOfIncident:      SeverityOfIncidentEnumFatal,
					RecurrenceRisk:          RecurrenceRiskEnumHigh,
					PhysicalInjury:          PhysicalInjuryEnumBrokenBones,
					ClientID:                client.ID,
					EmployeeAbsenteeism:     "No",
					Emails:                  []string{"report@example.com"},
				}
			},
			checks: func(t *testing.T, incident CreateIncidentRow, err error) {
				require.NoError(t, err, "CreateIncident() should not error")
				require.NotZero(t, incident.ID, "incident ID should be set")
				require.True(t, incident.Accident, "accident field should be true")
				require.False(t, incident.IsConfirmed, "incident should not be confirmed initially")
				require.Equal(t, SeverityOfIncidentEnumFatal, incident.SeverityOfIncident)
			},
		},
		{
			name: "successful creation with all optional fields populated",
			setup: func(ctx context.Context, qtx *Queries) CreateIncidentParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				client := createRandomClientDetails(ctx, qtx)

				return CreateIncidentParams{
					EmployeeID:              employee.ID,
					LocationID:              location.ID,
					ReporterInvolvement:     IncidentReporterInvolvementEnumAlarmed,
					InformWho:               []string{"Manager", "HR", "Safety Officer"},
					IncidentDate:            pgtype.Date{Time: time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					RuntimeIncident:         "Detailed incident description",
					IncidentType:            "Physical Injury",
					PassingAway:             false,
					SelfHarm:                false,
					Violence:                true,
					FireWaterDamage:         false,
					Accident:                false,
					ClientAbsence:           false,
					Medicines:               true,
					Organization:            true,
					UseProhibitedSubstances: false,
					OtherNotifications:      true,
					SeverityOfIncident:      SeverityOfIncidentEnumFatal,
					IncidentExplanation:     util.StringPtr("Client had an altercation with another resident"),
					RecurrenceRisk:          RecurrenceRiskEnumHigh,
					IncidentPreventSteps:    util.StringPtr("Increase supervision during group activities"),
					IncidentTakenMeasures:   util.StringPtr("Client moved to different room"),
					Technical:               []string{"CCTV recorded"},
					Organizational:          []string{"Policy reviewed"},
					MeseWorker:              []string{"John Doe"},
					ClientOptions:           []string{"Counseling offered"},
					OtherCause:              util.StringPtr("Stress related"),
					CauseExplanation:        util.StringPtr("Client was stressed about upcoming family visit"),
					PhysicalInjury:          PhysicalInjuryEnumBrokenBones,
					PhysicalInjuryDesc:      util.StringPtr("Minor bruising on arm"),
					PsychologicalDamage:     PsychologicalDamageEnumDrowsiness,
					PsychologicalDamageDesc: util.StringPtr("Client upset and withdrawn"),
					NeededConsultation:      NeededConsultationEnumConsultGp,
					Succession:              []string{"Follow-up meeting", "Parent notification"},
					SuccessionDesc:          util.StringPtr("Parents informed of incident"),
					Other:                   true,
					OtherDesc:               util.StringPtr("Additional notes about the incident"),
					AdditionalAppointments:  util.StringPtr("Therapy appointment scheduled"),
					ClientID:                client.ID,
					EmployeeAbsenteeism:     "No",
					Emails:                  []string{"report@example.com", "supervisor@example.com"},
				}
			},
			checks: func(t *testing.T, incident CreateIncidentRow, err error) {
				require.NoError(t, err, "CreateIncident() with all fields should not error")
				require.NotZero(t, incident.ID, "incident ID should be set")
				require.True(t, incident.Violence, "violence field should be true")
				require.True(t, incident.Medicines, "medicines field should be true")
				require.NotNil(t, incident.IncidentExplanation, "incident explanation should not be nil")
				require.Equal(t, "Client had an altercation with another resident", *incident.IncidentExplanation)
				require.Equal(t, SeverityOfIncidentEnumFatal, incident.SeverityOfIncident)
				require.Equal(t, RecurrenceRiskEnumHigh, incident.RecurrenceRisk)
				require.NotNil(t, incident.PhysicalInjuryDesc)
				require.NotNil(t, incident.PsychologicalDamageDesc)
			},
		},
		{
			name: "creation with edge case - all boolean flags false",
			setup: func(ctx context.Context, qtx *Queries) CreateIncidentParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				client := createRandomClientDetails(ctx, qtx)

				return CreateIncidentParams{
					EmployeeID:              employee.ID,
					LocationID:              location.ID,
					ReporterInvolvement:     IncidentReporterInvolvementEnumAlarmed,
					InformWho:               []string{},
					IncidentDate:            pgtype.Date{Time: time.Now(), Valid: true},
					RuntimeIncident:         "Minor incident",
					IncidentType:            "Administrative",
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
					SeverityOfIncident:      SeverityOfIncidentEnumNearIncident,
					RecurrenceRisk:          RecurrenceRiskEnumVeryHigh,
					ClientID:                client.ID,
					EmployeeAbsenteeism:     "Yes",
					Emails:                  []string{},
				}
			},
			checks: func(t *testing.T, incident CreateIncidentRow, err error) {
				require.NoError(t, err)
				require.False(t, incident.PassingAway)
				require.False(t, incident.SelfHarm)
				require.False(t, incident.Violence)
				require.Empty(t, incident.InformWho)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			incident, err := qtx.CreateIncident(ctx, params)
			tt.checks(t, incident, err)
		})
	}
}

func TestGetIncident(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, incident GetIncidentRow, err error)
	}{
		{
			name: "successfully retrieve existing incident",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				incident := createRandomIncident(ctx, qtx)
				return incident.ID
			},
			checks: func(t *testing.T, incident GetIncidentRow, err error) {
				require.NoError(t, err, "GetIncident() should not error for existing incident")
				require.NotZero(t, incident.ID, "incident ID should be set")
				require.NotNil(t, incident.EmployeeFirstName, "employee first name should not be nil")
				require.NotNil(t, incident.ClientFirstName, "client first name should not be nil")
			},
		},
		{
			name: "error retrieving non-existent incident",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 99999
			},
			checks: func(t *testing.T, incident GetIncidentRow, err error) {
				require.Error(t, err, "GetIncident() should error for non-existent incident")
			},
		},
		{
			name: "retrieve incident with all fields populated",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				incident := createRandomIncidentWithAllFields(ctx, qtx)
				return incident.ID
			},
			checks: func(t *testing.T, incident GetIncidentRow, err error) {
				require.NoError(t, err)
				require.NotNil(t, incident.IncidentExplanation, "incident explanation should be populated")
				require.NotNil(t, incident.PhysicalInjuryDesc, "physical injury description should be populated")
				require.NotEmpty(t, incident.Technical, "technical field should not be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			incidentID := tt.setup(ctx, qtx)
			incident, err := qtx.GetIncident(ctx, incidentID)
			tt.checks(t, incident, err)
		})
	}
}

func TestListIncidents(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32)
		checks func(t *testing.T, incidents []ListIncidentsRow, err error)
	}{
		{
			name: "successfully list incidents for client",
			setup: func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32) {
				client := createRandomClientDetails(ctx, qtx)
				// Create multiple incidents for the client
				for i := 0; i < 3; i++ {
					_ = createRandomIncidentForClient(ctx, qtx, client.ID)
				}
				return client.ID, 10, 0
			},
			checks: func(t *testing.T, incidents []ListIncidentsRow, err error) {
				require.NoError(t, err, "ListIncidents() should not error")
				require.Greater(t, len(incidents), 0, "should return at least one incident")
				require.GreaterOrEqual(t, len(incidents), 3, "should return all created incidents")
				// Verify incidents are sorted by date
				for i := 0; i < len(incidents)-1; i++ {
					require.True(t, incidents[i].IncidentDate.Time.After(incidents[i+1].IncidentDate.Time) ||
						incidents[i].IncidentDate.Time.Equal(incidents[i+1].IncidentDate.Time),
						"incidents should be ordered by incident_date DESC")
				}
			},
		},
		{
			name: "list incidents with pagination",
			setup: func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32) {
				client := createRandomClientDetails(ctx, qtx)
				for i := 0; i < 5; i++ {
					_ = createRandomIncidentForClient(ctx, qtx, client.ID)
				}
				return client.ID, 2, 0 // limit to 2
			},
			checks: func(t *testing.T, incidents []ListIncidentsRow, err error) {
				require.NoError(t, err)
				require.LessOrEqual(t, len(incidents), 2, "should respect limit parameter")
			},
		},
		{
			name: "list incidents with offset",
			setup: func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32) {
				client := createRandomClientDetails(ctx, qtx)
				for i := 0; i < 4; i++ {
					_ = createRandomIncidentForClient(ctx, qtx, client.ID)
				}
				return client.ID, 2, 2 // limit 2, offset 2
			},
			checks: func(t *testing.T, incidents []ListIncidentsRow, err error) {
				require.NoError(t, err)
				require.LessOrEqual(t, len(incidents), 2)
			},
		},
		{
			name: "list incidents for client with no incidents",
			setup: func(ctx context.Context, qtx *Queries) (uuid.UUID, int32, int32) {
				client := createRandomClientDetails(ctx, qtx)
				return client.ID, 10, 0
			},
			checks: func(t *testing.T, incidents []ListIncidentsRow, err error) {
				require.NoError(t, err)
				require.Empty(t, incidents, "should return empty slice for client with no incidents")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			clientID, limit, offset := tt.setup(ctx, qtx)
			incidents, err := qtx.ListIncidents(ctx, ListIncidentsParams{
				ClientID: clientID,
				Limit:    limit,
				Offset:   offset,
			})
			tt.checks(t, incidents, err)
		})
	}
}

func TestConfirmIncident(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, result ConfirmIncidentRow, err error)
	}{
		{
			name: "successfully confirm existing incident",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				incident := createRandomIncident(ctx, qtx)
				return incident.ID
			},
			checks: func(t *testing.T, result ConfirmIncidentRow, err error) {
				require.NoError(t, err, "ConfirmIncident() should not error")
				require.True(t, result.IsConfirmed, "incident should be confirmed")
				require.Equal(t, result.ID, result.ID, "incident ID should match")
			},
		},
		{
			name: "confirm incident with file URL",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				incident := createRandomIncidentWithFileUrl(ctx, qtx)
				return incident.ID
			},
			checks: func(t *testing.T, result ConfirmIncidentRow, err error) {
				require.NoError(t, err)
				require.True(t, result.IsConfirmed)
				require.NotNil(t, result.FileUrl, "file URL should be present")
			},
		},
		{
			name: "confirm non-existent incident",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 99999
			},
			checks: func(t *testing.T, result ConfirmIncidentRow, err error) {
				require.Error(t, err, "ConfirmIncident() should error for non-existent incident")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			incidentID := tt.setup(ctx, qtx)
			result, err := qtx.ConfirmIncident(ctx, incidentID)
			tt.checks(t, result, err)
		})
	}
}

func TestDeleteIncident(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, err error)
	}{
		{
			name: "successfully delete existing incident",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				incident := createRandomIncident(ctx, qtx)
				return incident.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteIncident() should not error")
			},
		},
		{
			name: "delete non-existent incident does not error",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 99999
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteIncident() should not error even for non-existent incident")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			incidentID := tt.setup(ctx, qtx)
			err = qtx.DeleteIncident(ctx, incidentID)
			tt.checks(t, err)
		})
	}
}

// Random data generators for tests

func createRandomIncident(ctx context.Context, qtx *Queries) CreateIncidentRow {
	employee := createRandomEmployeeProfile(ctx, qtx)
	location := createRandomLocation(ctx, qtx)
	client := createRandomClientDetails(ctx, qtx)

	incident, err := qtx.CreateIncident(ctx, CreateIncidentParams{
		EmployeeID:              employee.ID,
		LocationID:              location.ID,
		ReporterInvolvement:     IncidentReporterInvolvementEnumAlarmed,
		InformWho:               []string{"Manager"},
		IncidentDate:            pgtype.Date{Time: time.Now(), Valid: true},
		RuntimeIncident:         "Standard incident for testing",
		IncidentType:            "Test Incident",
		PassingAway:             false,
		SelfHarm:                false,
		Violence:                false,
		FireWaterDamage:         false,
		Accident:                true,
		ClientAbsence:           false,
		Medicines:               false,
		Organization:            false,
		UseProhibitedSubstances: false,
		OtherNotifications:      false,
		SeverityOfIncident:      SeverityOfIncidentEnumSerious,
		RecurrenceRisk:          RecurrenceRiskEnumVeryLow,
		ClientID:                client.ID,
		EmployeeAbsenteeism:     "No",
		Emails:                  []string{"test@example.com"},
	})
	if err != nil {
		panic("failed to create random incident: " + err.Error())
	}
	return incident
}

func createRandomIncidentForClient(ctx context.Context, qtx *Queries, clientID uuid.UUID) CreateIncidentRow {
	employee := createRandomEmployeeProfile(ctx, qtx)
	location := createRandomLocation(ctx, qtx)

	incident, err := qtx.CreateIncident(ctx, CreateIncidentParams{
		EmployeeID:              employee.ID,
		LocationID:              location.ID,
		ReporterInvolvement:     IncidentReporterInvolvementEnumWitness,
		InformWho:               []string{"Manager"},
		IncidentDate:            pgtype.Date{Time: time.Now().AddDate(0, 0, -8), Valid: true},
		RuntimeIncident:         "Incident for " + clientID.String(),
		IncidentType:            "Test Incident",
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
		SeverityOfIncident:      SeverityOfIncidentEnumFatal,
		RecurrenceRisk:          RecurrenceRiskEnumVeryHigh,
		ClientID:                clientID,
		EmployeeAbsenteeism:     "No",
		Emails:                  []string{},
	})
	if err != nil {
		panic("failed to create random incident for client: " + err.Error())
	}
	return incident
}

func createRandomIncidentWithAllFields(ctx context.Context, qtx *Queries) CreateIncidentRow {
	employee := createRandomEmployeeProfile(ctx, qtx)
	location := createRandomLocation(ctx, qtx)
	client := createRandomClientDetails(ctx, qtx)

	incident, err := qtx.CreateIncident(ctx, CreateIncidentParams{
		EmployeeID:              employee.ID,
		LocationID:              location.ID,
		ReporterInvolvement:     IncidentReporterInvolvementEnumFoundAfterwards,
		InformWho:               []string{"Manager", "HR", "Safety"},
		IncidentDate:            pgtype.Date{Time: time.Now(), Valid: true},
		RuntimeIncident:         "Comprehensive incident with all fields",
		IncidentType:            "Comprehensive Test",
		PassingAway:             false,
		SelfHarm:                true,
		Violence:                true,
		FireWaterDamage:         true,
		Accident:                true,
		ClientAbsence:           true,
		Medicines:               true,
		Organization:            true,
		UseProhibitedSubstances: true,
		OtherNotifications:      true,
		SeverityOfIncident:      SeverityOfIncidentEnumNearIncident,
		IncidentExplanation:     util.StringPtr("Comprehensive explanation of incident"),
		RecurrenceRisk:          RecurrenceRiskEnumHigh,
		IncidentPreventSteps:    util.StringPtr("Steps to prevent recurrence"),
		IncidentTakenMeasures:   util.StringPtr("Measures already taken"),
		Technical:               []string{"Video recorded", "Photos taken"},
		Organizational:          []string{"Policy updated"},
		MeseWorker:              []string{"Jane Smith"},
		ClientOptions:           []string{"Therapy", "Counseling"},
		OtherCause:              util.StringPtr("Environmental factors"),
		CauseExplanation:        util.StringPtr("Detailed cause explanation"),
		PhysicalInjury:          PhysicalInjuryEnumNotNoticeableYet,
		PhysicalInjuryDesc:      util.StringPtr("Significant physical injury"),
		PsychologicalDamage:     PsychologicalDamageEnumOther,
		PsychologicalDamageDesc: util.StringPtr("Significant psychological impact"),
		NeededConsultation:      NeededConsultationEnumConsultGp,
		Succession:              []string{"Follow-up", "Review", "Debrief"},
		SuccessionDesc:          util.StringPtr("Detailed succession plan"),
		Other:                   true,
		OtherDesc:               util.StringPtr("Additional detailed notes"),
		AdditionalAppointments:  util.StringPtr(gofakeit.Sentence(5)),
		ClientID:                client.ID,
		EmployeeAbsenteeism:     "Yes",
		Emails:                  []string{gofakeit.Email(), gofakeit.Email()},
	})
	if err != nil {
		panic("failed to create random incident with all fields: " + err.Error())
	}
	return incident
}

func createRandomIncidentWithFileUrl(ctx context.Context, qtx *Queries) CreateIncidentRow {
	employee := createRandomEmployeeProfile(ctx, qtx)
	location := createRandomLocation(ctx, qtx)
	client := createRandomClientDetails(ctx, qtx)

	incident, err := qtx.CreateIncident(ctx, CreateIncidentParams{
		EmployeeID:              employee.ID,
		LocationID:              location.ID,
		ReporterInvolvement:     IncidentReporterInvolvementEnumAlarmed,
		InformWho:               []string{"Manager"},
		IncidentDate:            pgtype.Date{Time: time.Now(), Valid: true},
		RuntimeIncident:         "Incident with file URL",
		IncidentType:            "With Documentation",
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
		SeverityOfIncident:      SeverityOfIncidentEnumSerious,
		IncidentExplanation:     util.StringPtr("Incident documented with file URL"),
		RecurrenceRisk:          RecurrenceRiskEnumVeryHigh,
		ClientID:                client.ID,
		EmployeeAbsenteeism:     "No",
		Emails:                  []string{},
	})
	if err != nil {
		panic("failed to create random incident with file URL: " + err.Error())
	}
	return incident
}
