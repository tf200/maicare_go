package db

import (
	"context"
	"maicare_go/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomProgressReport(t *testing.T, clientID int64, employeeID int64) ProgressReport {
	arg := CreateProgressReportParams{
		ClientID:       clientID,
		EmployeeID:     &employeeID,
		Title:          util.StringPtr("Test Progress Report"),
		Date:           pgtype.Timestamptz{Time: util.RandomTIme(), Valid: true},
		ReportText:     "Test Progress Report",
		Type:           "morning_report",
		EmotionalState: "normal",
	}
	progressReport, err := testQueries.CreateProgressReport(context.Background(), arg)
	require.NoError(t, err)
	return progressReport
}

func TestCreateProgressReport(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	createRandomProgressReport(t, client.ID, employee.ID)
}

func TestListProgressReport(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	for i := 0; i < 10; i++ {
		createRandomProgressReport(t, client.ID, employee.ID)
	}
	arg := ListProgressReportsParams{
		ClientID: client.ID,
		Limit:    5,
		Offset:   5,
	}
	progressReports, err := testQueries.ListProgressReports(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, progressReports, 5)

}

func TestGetProgressReport(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	progressReport1 := createRandomProgressReport(t, client.ID, employee.ID)
	progressReport2, err := testQueries.GetProgressReport(context.Background(), progressReport1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, progressReport2)
	require.Equal(t, progressReport1.ID, progressReport2.ID)
	require.Equal(t, progressReport1.ClientID, progressReport2.ClientID)
}

func TestUpdateProgressReport(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	progressReport1 := createRandomProgressReport(t, client.ID, employee.ID)
	arg := UpdateProgressReportParams{
		ID:             progressReport1.ID,
		ReportText:     util.StringPtr("Updated Progress Report"),
		EmotionalState: util.StringPtr("happy"),
	}
	progressReport2, err := testQueries.UpdateProgressReport(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, progressReport2)
	require.Equal(t, progressReport1.ID, progressReport2.ID)
	require.NotEqual(t, progressReport1.ReportText, progressReport2.ReportText)
	require.NotEqual(t, progressReport1.EmotionalState, progressReport2.EmotionalState)
}

func TestGetProgressReportsByDateRange(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	for i := 0; i < 10; i++ {
		createRandomProgressReport(t, client.ID, employee.ID)
	}
	startDate := util.RandomTIme()
	endDate := startDate.AddDate(1, 1, 5)
	arg := GetProgressReportsByDateRangeParams{
		ClientID:  client.ID,
		StartDate: pgtype.Timestamptz{Time: startDate, Valid: true},
		EndDate:   pgtype.Timestamptz{Time: endDate, Valid: true},
	}

	progressReports, err := testQueries.GetProgressReportsByDateRange(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, progressReports, 5)
	for _, report := range progressReports {
		require.NotEmpty(t, report)
		require.Equal(t, client.ID, report.ClientID)
		require.GreaterOrEqual(t, report.Date.Time.Unix(), startDate.Unix())
		require.LessOrEqual(t, report.Date.Time.Unix(), endDate.Unix())
	}
}

func createRandomAiGeneratedReport(t *testing.T, clientID int64) AiGeneratedReport {
	startdate := util.RandomTIme()
	enddate := startdate.AddDate(0, 0, 7)

	arg := CreateAiGeneratedReportParams{
		ClientID:   clientID,
		ReportText: "Test AI Generated Report",
		StartDate:  pgtype.Date{Time: startdate, Valid: true},
		EndDate:    pgtype.Date{Time: enddate, Valid: true},
	}

	aiGeneratedReport, err := testQueries.CreateAiGeneratedReport(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, aiGeneratedReport)
	require.Equal(t, arg.ClientID, aiGeneratedReport.ClientID)
	require.Equal(t, arg.ReportText, aiGeneratedReport.ReportText)
	return aiGeneratedReport
}

func TestCreateAiGeneratedReport(t *testing.T) {
	client := createRandomClientDetails(t)
	createRandomAiGeneratedReport(t, client.ID)
}

func TestListAiGeneratedReports(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 10; i++ {
		createRandomAiGeneratedReport(t, client.ID)
	}
	arg := ListAiGeneratedReportsParams{
		ClientID: client.ID,
		Limit:    5,
		Offset:   5,
	}
	aiGeneratedReports, err := testQueries.ListAiGeneratedReports(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, aiGeneratedReports, 5)
}

func TestGetAiGeneratedReport(t *testing.T) {
	client := createRandomClientDetails(t)
	aiGeneratedReport1 := createRandomAiGeneratedReport(t, client.ID)
	aiGeneratedReport2, err := testQueries.GetAiGeneratedReport(context.Background(), aiGeneratedReport1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, aiGeneratedReport2)
	require.Equal(t, aiGeneratedReport1.ID, aiGeneratedReport2.ID)
	require.Equal(t, aiGeneratedReport1.ClientID, aiGeneratedReport2.ClientID)
}


