package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func randomCarePlanStatus() CarePlanStatusEnum {
	return CarePlanStatusEnumPending
}

func randomCarePlanTimeframe() CarePlanTimeframeEnum {
	return CarePlanTimeframeEnumShortTerm
}

func randomCarePlanObjectiveStatus() CarePlanObjectiveStatusEnum {
	return CarePlanObjectiveStatusEnumInProgress
}

func randomCarePlanInterventionFrequency() CarePlanInterventionFrequencyEnum {
	return CarePlanInterventionFrequencyEnumDaily
}

func randomCarePlanReportType() CarePlanReportTypeEnum {
	return CarePlanReportTypeEnumProgress
}

func randomCarePlanRiskLevel() CarePlanRiskLevelEnum {
	return CarePlanRiskLevelEnumMedium
}

func randomNullCarePlanTimeframe() NullCarePlanTimeframeEnum {
	return NullCarePlanTimeframeEnum{
		CarePlanTimeframeEnum: randomCarePlanTimeframe(),
		Valid:                 true,
	}
}

func randomNullCarePlanObjectiveStatus() NullCarePlanObjectiveStatusEnum {
	return NullCarePlanObjectiveStatusEnum{
		CarePlanObjectiveStatusEnum: randomCarePlanObjectiveStatus(),
		Valid:                       true,
	}
}

func randomNullCarePlanInterventionFrequency() NullCarePlanInterventionFrequencyEnum {
	return NullCarePlanInterventionFrequencyEnum{
		CarePlanInterventionFrequencyEnum: randomCarePlanInterventionFrequency(),
		Valid:                             true,
	}
}

func randomNullCarePlanReportType() NullCarePlanReportTypeEnum {
	return NullCarePlanReportTypeEnum{
		CarePlanReportTypeEnum: randomCarePlanReportType(),
		Valid:                  true,
	}
}

func randomNullCarePlanRiskLevel() NullCarePlanRiskLevelEnum {
	return NullCarePlanRiskLevelEnum{
		CarePlanRiskLevelEnum: randomCarePlanRiskLevel(),
		Valid:                 true,
	}
}

func createRandomMaturityMatrix(ctx context.Context, t *testing.T, q *Queries) MaturityMatrix {
	// Use one of the preset maturity matrix records (IDs 1-13 from seed data)
	id := int64(util.RandomInt(1, 13))
	maturityMatrix, err := q.GetMaturityMatrix(ctx, id)
	require.NoError(t, err)
	return maturityMatrix
}

func createRandomClientMaturityMatrixAssessment(ctx context.Context, t *testing.T, q *Queries) CreateClientMaturityMatrixAssessmentRow {
	client := createRandomClientDetails(ctx, q)
	maturityMatrix := createRandomMaturityMatrix(ctx, t, q)
	params := CreateClientMaturityMatrixAssessmentParams{
		ClientID:         client.ID,
		MaturityMatrixID: maturityMatrix.ID,
		StartDate:        pgtype.Date{Time: time.Now(), Valid: true},
		EndDate:          pgtype.Date{Time: time.Now().AddDate(0, 1, 0), Valid: true},
		TargetLevel:      3,
		InitialLevel:     1,
		CurrentLevel:     2,
	}
	assessment, err := q.CreateClientMaturityMatrixAssessment(ctx, params)
	require.NoError(t, err)
	return assessment
}

func createRandomCarePlan(ctx context.Context, t *testing.T, q *Queries) CarePlan {
	assessment := createRandomClientMaturityMatrixAssessment(ctx, t, q)
	employee := createRandomEmployeeProfile(ctx, q)
	params := CreateCarePlanParams{
		AssessmentID:          assessment.ID,
		GeneratedByEmployeeID: &employee.ID,
		AssessmentSummary:     util.RandomString(100),
		RawLlmResponse:        []byte(`{"response": "` + util.RandomString(200) + `"}`),
		Status:                randomCarePlanStatus(),
	}
	carePlan, err := q.CreateCarePlan(ctx, params)
	require.NoError(t, err)
	return carePlan
}

func createRandomCarePlanObjective(ctx context.Context, t *testing.T, q *Queries) CarePlanObjective {
	carePlan := createRandomCarePlan(ctx, t, q)
	params := CreateCarePlanObjectiveParams{
		CarePlanID:  carePlan.ID,
		Timeframe:   randomCarePlanTimeframe(),
		GoalTitle:   util.RandomString(50),
		Description: util.RandomString(100),
		TargetDate:  pgtype.Date{Time: time.Now().AddDate(0, 0, 30), Valid: true},
	}
	objective, err := q.CreateCarePlanObjective(ctx, params)
	require.NoError(t, err)
	return objective
}

func createRandomCarePlanAction(ctx context.Context, t *testing.T, q *Queries) CarePlanAction {
	objective := createRandomCarePlanObjective(ctx, t, q)
	params := CreateCarePlanActionParams{
		ObjectiveID:       objective.ID,
		ActionDescription: util.RandomString(100),
		SortOrder:         1,
	}
	action, err := q.CreateCarePlanAction(ctx, params)
	require.NoError(t, err)
	return action
}

func createRandomCarePlanIntervention(ctx context.Context, t *testing.T, q *Queries) CarePlanIntervention {
	carePlan := createRandomCarePlan(ctx, t, q)
	params := CreateCarePlanInterventionParams{
		CarePlanID:              carePlan.ID,
		Frequency:               randomCarePlanInterventionFrequency(),
		InterventionDescription: util.RandomString(100),
	}
	intervention, err := q.CreateCarePlanIntervention(ctx, params)
	require.NoError(t, err)
	return intervention
}

func createRandomCarePlanReport(ctx context.Context, t *testing.T, q *Queries) CarePlanReport {
	carePlan := createRandomCarePlan(ctx, t, q)
	employee := createRandomEmployeeProfile(ctx, q)
	params := CreateCarePlanReportParams{
		CarePlanID:          carePlan.ID,
		ReportType:          randomCarePlanReportType(),
		ReportContent:       util.RandomString(200),
		CreatedByEmployeeID: employee.ID,
		IsCritical:          false,
	}
	report, err := q.CreateCarePlanReport(ctx, params)
	require.NoError(t, err)
	return report
}

func createRandomCarePlanResource(ctx context.Context, t *testing.T, q *Queries) CarePlanResource {
	carePlan := createRandomCarePlan(ctx, t, q)
	params := CreateCarePlanResourcesParams{
		CarePlanID:          carePlan.ID,
		ResourceDescription: util.RandomString(100),
		IsObtained:          false,
		ObtainedDate:        pgtype.Date{Valid: false},
	}
	resource, err := q.CreateCarePlanResources(ctx, params)
	require.NoError(t, err)
	return resource
}

func createRandomCarePlanRisk(ctx context.Context, t *testing.T, q *Queries) CarePlanRisk {
	carePlan := createRandomCarePlan(ctx, t, q)
	params := CreateCarePlanRiskParams{
		CarePlanID:         carePlan.ID,
		RiskDescription:    util.RandomString(100),
		MitigationStrategy: util.RandomString(100),
		RiskLevel:          randomCarePlanRiskLevel(),
	}
	risk, err := q.CreateCarePlanRisk(ctx, params)
	require.NoError(t, err)
	return risk
}

func createRandomCarePlanSuccessMetric(ctx context.Context, t *testing.T, q *Queries) CarePlanMetric {
	carePlan := createRandomCarePlan(ctx, t, q)
	currentValue := util.RandomString(10)
	params := CreateCarePlanSuccessMetricParams{
		CarePlanID:        carePlan.ID,
		MetricName:        util.RandomString(50),
		TargetValue:       util.RandomString(10),
		MeasurementMethod: util.RandomString(50),
		CurrentValue:      &currentValue,
	}
	metric, err := q.CreateCarePlanSuccessMetric(ctx, params)
	require.NoError(t, err)
	return metric
}

func createRandomCarePlanSupportNetwork(ctx context.Context, t *testing.T, q *Queries) CarePlanSupportNetwork {
	carePlan := createRandomCarePlan(ctx, t, q)
	params := CreateCarePlanSupportNetworkParams{
		CarePlanID:                carePlan.ID,
		RoleTitle:                 util.RandomString(50),
		ResponsibilityDescription: util.RandomString(100),
	}
	network, err := q.CreateCarePlanSupportNetwork(ctx, params)
	require.NoError(t, err)
	return network
}

func TestCreateCarePlan(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			require.NotZero(t, carePlan.ID)
			require.Equal(t, randomCarePlanStatus(), carePlan.Status)
		})
	}
}

func TestCreateCarePlanAction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan action creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			action := createRandomCarePlanAction(ctx, t, qtx)
			require.NotZero(t, action.ID)
			require.False(t, action.IsCompleted)
		})
	}
}

func TestCreateCarePlanIntervention(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan intervention creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			intervention := createRandomCarePlanIntervention(ctx, t, qtx)
			require.NotZero(t, intervention.ID)
			require.True(t, intervention.IsActive)
		})
	}
}

func TestCreateCarePlanObjective(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan objective creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			objective := createRandomCarePlanObjective(ctx, t, qtx)
			require.NotZero(t, objective.ID)
			require.Equal(t, randomCarePlanTimeframe(), objective.Timeframe)
		})
	}
}

func TestCreateCarePlanReport(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan report creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			report := createRandomCarePlanReport(ctx, t, qtx)
			require.NotZero(t, report.ID)
			require.Equal(t, randomCarePlanReportType(), report.ReportType)
		})
	}
}

func TestCreateCarePlanResources(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan resource creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			resource := createRandomCarePlanResource(ctx, t, qtx)
			require.NotZero(t, resource.ID)
			require.False(t, resource.IsObtained)
		})
	}
}

func TestCreateCarePlanRisk(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan risk creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			risk := createRandomCarePlanRisk(ctx, t, qtx)
			require.NotZero(t, risk.ID)
			require.True(t, risk.IsActive)
		})
	}
}

func TestCreateCarePlanSuccessMetric(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan success metric creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			metric := createRandomCarePlanSuccessMetric(ctx, t, qtx)
			require.NotZero(t, metric.ID)
			require.False(t, metric.IsAchieved)
		})
	}
}

func TestCreateCarePlanSupportNetwork(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan support network creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			network := createRandomCarePlanSupportNetwork(ctx, t, qtx)
			require.NotZero(t, network.ID)
			require.True(t, network.IsActive)
		})
	}
}

func TestCreateClientMaturityMatrixAssessment(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful client maturity matrix assessment creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			assessment := createRandomClientMaturityMatrixAssessment(ctx, t, qtx)
			require.NotZero(t, assessment.ID)
			require.True(t, assessment.IsActive)
		})
	}
}

func TestDeleteCarePlan(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			err = qtx.DeleteCarePlan(ctx, carePlan.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanAction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan action deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			action := createRandomCarePlanAction(ctx, t, qtx)
			err = qtx.DeleteCarePlanAction(ctx, action.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanIntervention(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan intervention deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			intervention := createRandomCarePlanIntervention(ctx, t, qtx)
			err = qtx.DeleteCarePlanIntervention(ctx, intervention.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanObjective(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan objective deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			objective := createRandomCarePlanObjective(ctx, t, qtx)
			err = qtx.DeleteCarePlanObjective(ctx, objective.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanReport(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan report deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			report := createRandomCarePlanReport(ctx, t, qtx)
			err = qtx.DeleteCarePlanReport(ctx, report.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanResource(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan resource deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			resource := createRandomCarePlanResource(ctx, t, qtx)
			err = qtx.DeleteCarePlanResource(ctx, resource.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanRisk(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan risk deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			risk := createRandomCarePlanRisk(ctx, t, qtx)
			err = qtx.DeleteCarePlanRisk(ctx, risk.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanSuccessMetric(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan success metric deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			metric := createRandomCarePlanSuccessMetric(ctx, t, qtx)
			err = qtx.DeleteCarePlanSuccessMetric(ctx, metric.ID)
			require.NoError(t, err)
		})
	}
}

func TestDeleteCarePlanSupportNetwork(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful care plan support network deletion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			network := createRandomCarePlanSupportNetwork(ctx, t, qtx)
			err = qtx.DeleteCarePlanSupportNetwork(ctx, network.ID)
			require.NoError(t, err)
		})
	}
}

func TestGetCarePlanActionsMaxSortOrder(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get max sort order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			objective := createRandomCarePlanObjective(ctx, t, qtx)
			maxOrder, err := qtx.GetCarePlanActionsMaxSortOrder(ctx, objective.ID)
			require.NoError(t, err)
			require.Equal(t, int32(0), maxOrder)
		})
	}
}

func TestGetCarePlanInterventions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get interventions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			interventions, err := qtx.GetCarePlanInterventions(ctx, carePlan.ID)
			require.NoError(t, err)
			require.Len(t, interventions, 0)
		})
	}
}

func TestGetCarePlanObjectivesWithActions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get objectives with actions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			rows, err := qtx.GetCarePlanObjectivesWithActions(ctx, carePlan.ID)
			require.NoError(t, err)
			require.Len(t, rows, 0)
		})
	}
}

func TestGetCarePlanOverview(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get care plan overview",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			overview, err := qtx.GetCarePlanOverview(ctx, carePlan.ID)
			require.NoError(t, err)
			require.Equal(t, carePlan.ID, overview.ID)
		})
	}
}

func TestGetCarePlanReport(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get care plan report",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			report := createRandomCarePlanReport(ctx, t, qtx)
			gotReport, err := qtx.GetCarePlanReport(ctx, report.ID)
			require.NoError(t, err)
			require.Equal(t, report.ID, gotReport.ID)
		})
	}
}

func TestGetCarePlanResources(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get care plan resources",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			resources, err := qtx.GetCarePlanResources(ctx, carePlan.ID)
			require.NoError(t, err)
			require.Len(t, resources, 0)
		})
	}
}

func TestGetCarePlanRisks(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get care plan risks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			risks, err := qtx.GetCarePlanRisks(ctx, carePlan.ID)
			require.NoError(t, err)
			require.Len(t, risks, 0)
		})
	}
}

func TestGetCarePlanSuccessMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get care plan success metrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			metrics, err := qtx.GetCarePlanSuccessMetrics(ctx, carePlan.ID)
			require.NoError(t, err)
			require.Len(t, metrics, 0)
		})
	}
}

func TestGetCarePlanSupportNetwork(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get care plan support network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			networks, err := qtx.GetCarePlanSupportNetwork(ctx, carePlan.ID)
			require.NoError(t, err)
			require.Len(t, networks, 0)
		})
	}
}

func TestGetClientMaturityMatrixAssessment(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get client maturity matrix assessment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			assessment := createRandomClientMaturityMatrixAssessment(ctx, t, qtx)
			gotAssessment, err := qtx.GetClientMaturityMatrixAssessment(ctx, assessment.ID)
			require.NoError(t, err)
			require.Equal(t, assessment.ID, gotAssessment.ID)
		})
	}
}

func TestGetLevelDescription(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get level description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			maturityMatrix := createRandomMaturityMatrix(ctx, t, qtx)
			params := GetLevelDescriptionParams{
				ID:    maturityMatrix.ID,
				Level: "1",
			}
			desc, err := qtx.GetLevelDescription(ctx, params)
			require.NoError(t, err)
			require.Equal(t, maturityMatrix.TopicName, desc.TopicName)
		})
	}
}

func TestGetMaturityMatrix(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful get maturity matrix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			maturityMatrix := createRandomMaturityMatrix(ctx, t, qtx)
			gotMatrix, err := qtx.GetMaturityMatrix(ctx, maturityMatrix.ID)
			require.NoError(t, err)
			require.Equal(t, maturityMatrix.ID, gotMatrix.ID)
		})
	}
}

func TestListCarePlanReports(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful list care plan reports",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			params := ListCarePlanReportsParams{
				CarePlanID: carePlan.ID,
				Limit:      10,
				Offset:     0,
			}
			reports, err := qtx.ListCarePlanReports(ctx, params)
			require.NoError(t, err)
			require.Len(t, reports, 0)
		})
	}
}

func TestListClientMaturityMatrixAssessments(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful list client maturity matrix assessments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			client := createRandomClientDetails(ctx, qtx)
			params := ListClientMaturityMatrixAssessmentsParams{
				ClientID: client.ID,
				Limit:    10,
				Offset:   0,
			}
			assessments, err := qtx.ListClientMaturityMatrixAssessments(ctx, params)
			require.NoError(t, err)
			require.Len(t, assessments, 0)
		})
	}
}

func TestListMaturityMatrix(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful list maturity matrix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			matrices, err := qtx.ListMaturityMatrix(ctx)
			require.NoError(t, err)
			require.Greater(t, len(matrices), 0) // Assuming some default data
		})
	}
}

func TestUpdateCarePlanAction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			action := createRandomCarePlanAction(ctx, t, qtx)
			newDesc := util.RandomString(100)
			params := UpdateCarePlanActionParams{
				ID:                action.ID,
				ActionDescription: &newDesc,
			}
			updatedAction, err := qtx.UpdateCarePlanAction(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newDesc, updatedAction.ActionDescription)
		})
	}
}

func TestUpdateCarePlanIntervention(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan intervention",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			intervention := createRandomCarePlanIntervention(ctx, t, qtx)
			newDesc := util.RandomString(100)
			params := UpdateCarePlanInterventionParams{
				ID:                      intervention.ID,
				Frequency:               randomNullCarePlanInterventionFrequency(),
				InterventionDescription: &newDesc,
			}
			updatedIntervention, err := qtx.UpdateCarePlanIntervention(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newDesc, updatedIntervention.InterventionDescription)
		})
	}
}

func TestUpdateCarePlanObjective(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan objective",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			objective := createRandomCarePlanObjective(ctx, t, qtx)
			newTitle := util.RandomString(50)
			params := UpdateCarePlanObjectiveParams{
				ID:          objective.ID,
				Timeframe:   randomNullCarePlanTimeframe(),
				GoalTitle:   &newTitle,
				Description: nil,
				Status:      randomNullCarePlanObjectiveStatus(),
			}
			updatedObjective, err := qtx.UpdateCarePlanObjective(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newTitle, updatedObjective.GoalTitle)
		})
	}
}

func TestUpdateCarePlanOverview(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan overview",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			carePlan := createRandomCarePlan(ctx, t, qtx)
			newSummary := util.RandomString(100)
			params := UpdateCarePlanOverviewParams{
				ID:                carePlan.ID,
				AssessmentSummary: &newSummary,
			}
			updatedCarePlan, err := qtx.UpdateCarePlanOverview(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newSummary, updatedCarePlan.AssessmentSummary)
		})
	}
}

func TestUpdateCarePlanReport(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan report",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			report := createRandomCarePlanReport(ctx, t, qtx)
			newContent := util.RandomString(200)
			params := UpdateCarePlanReportParams{
				ID:            report.ID,
				ReportType:    randomNullCarePlanReportType(),
				ReportContent: &newContent,
				IsCritical:    nil,
			}
			updatedReport, err := qtx.UpdateCarePlanReport(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newContent, updatedReport.ReportContent)
		})
	}
}

func TestUpdateCarePlanResource(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			resource := createRandomCarePlanResource(ctx, t, qtx)
			newDesc := util.RandomString(100)
			params := UpdateCarePlanResourceParams{
				ID:                  resource.ID,
				ResourceDescription: &newDesc,
				IsObtained:          nil,
				ObtainedDate:        pgtype.Date{Valid: false},
			}
			updatedResource, err := qtx.UpdateCarePlanResource(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newDesc, updatedResource.ResourceDescription)
		})
	}
}

func TestUpdateCarePlanRisk(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan risk",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			risk := createRandomCarePlanRisk(ctx, t, qtx)
			newDesc := util.RandomString(100)
			params := UpdateCarePlanRiskParams{
				ID:                 risk.ID,
				RiskDescription:    &newDesc,
				MitigationStrategy: nil,
				RiskLevel:          randomNullCarePlanRiskLevel(),
			}
			updatedRisk, err := qtx.UpdateCarePlanRisk(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newDesc, updatedRisk.RiskDescription)
		})
	}
}

func TestUpdateCarePlanSuccessMetric(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan success metric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			metric := createRandomCarePlanSuccessMetric(ctx, t, qtx)
			newName := util.RandomString(50)
			params := UpdateCarePlanSuccessMetricParams{
				ID:                metric.ID,
				MetricName:        &newName,
				TargetValue:       nil,
				MeasurementMethod: nil,
				CurrentValue:      nil,
			}
			updatedMetric, err := qtx.UpdateCarePlanSuccessMetric(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newName, updatedMetric.MetricName)
		})
	}
}

func TestUpdateCarePlanSupportNetwork(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "successful update care plan support network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			network := createRandomCarePlanSupportNetwork(ctx, t, qtx)
			newTitle := util.RandomString(50)
			params := UpdateCarePlanSupportNetworkParams{
				ID:                        network.ID,
				RoleTitle:                 &newTitle,
				ResponsibilityDescription: nil,
			}
			updatedNetwork, err := qtx.UpdateCarePlanSupportNetwork(ctx, params)
			require.NoError(t, err)
			require.Equal(t, newTitle, updatedNetwork.RoleTitle)
		})
	}
}
