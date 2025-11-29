package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateAiGeneratedReport(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateAiGeneratedReportParams
		checks func(t *testing.T, report AiGeneratedReport, err error)
	}{
		{
			name: "successful AI report creation",
			setup: func(ctx context.Context, qtx *Queries) CreateAiGeneratedReportParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateAiGeneratedReportParams{
					ClientID:   client.ID,
					ReportText: util.RandomString(100),
					StartDate: pgtype.Date{
						Time:  time.Now().Add(-7 * 24 * time.Hour),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Now(),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, report AiGeneratedReport, err error) {
				require.NoError(t, err, "CreateAiGeneratedReport should not return an error")
				require.NotZero(t, report.ID)
				require.NotEmpty(t, report.ReportText)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			report, err := qtx.CreateAiGeneratedReport(ctx, params)
			tt.checks(t, report, err)
		})
	}
}

func TestCreateProgressReport(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateProgressReportParams
		checks func(t *testing.T, report ProgressReport, err error)
	}{
		{
			name: "successful progress report creation",
			setup: func(ctx context.Context, qtx *Queries) CreateProgressReportParams {
				client := createRandomClientDetails(ctx, qtx)
				employee := createRandomEmployeeProfile(ctx, qtx)
				return CreateProgressReportParams{
					ClientID:   client.ID,
					EmployeeID: &employee.ID,
					Title:      randomStringPtr(20),
					Date: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
					ReportText:     util.RandomString(200),
					Type:           ProgressReportTypeEnumEveningReport,
					EmotionalState: EmotionalStateEnumDepressed,
				}
			},
			checks: func(t *testing.T, report ProgressReport, err error) {
				require.NoError(t, err, "CreateProgressReport should not return an error")
				require.NotZero(t, report.ID)
				require.Equal(t, ProgressReportTypeEnumEveningReport, report.Type)
			},
		},
		{
			name: "progress report creation without employee",
			setup: func(ctx context.Context, qtx *Queries) CreateProgressReportParams {
				client := createRandomClientDetails(ctx, qtx)
				return CreateProgressReportParams{
					ClientID:   client.ID,
					EmployeeID: nil,
					Title:      nil,
					Date: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
					ReportText:     util.RandomString(200),
					Type:           ProgressReportTypeEnumMorningReport,
					EmotionalState: EmotionalStateEnumAnxious,
				}
			},
			checks: func(t *testing.T, report ProgressReport, err error) {
				require.NoError(t, err, "CreateProgressReport should not return an error")
				require.Nil(t, report.EmployeeID)
				require.Nil(t, report.Title)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			report, err := qtx.CreateProgressReport(ctx, params)
			tt.checks(t, report, err)
		})
	}
}

func TestGetAiGeneratedReport(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, report AiGeneratedReport, err error)
	}{
		{
			name: "get existing AI report",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				report := createRandomAiGeneratedReport(ctx, qtx)
				return report.ID
			},
			checks: func(t *testing.T, report AiGeneratedReport, err error) {
				require.NoError(t, err, "GetAiGeneratedReport should not return an error")
				require.NotZero(t, report.ID)
			},
		},
		{
			name: "get non-existent AI report",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, report AiGeneratedReport, err error) {
				require.Error(t, err, "GetAiGeneratedReport should error for non-existent report")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			reportID := tt.setup(ctx, qtx)
			report, err := qtx.GetAiGeneratedReport(ctx, reportID)
			tt.checks(t, report, err)
		})
	}
}

func TestGetProgressReport(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, report GetProgressReportRow, err error)
	}{
		{
			name: "get existing progress report",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				report := createRandomProgressReport(ctx, qtx)
				return report.ID
			},
			checks: func(t *testing.T, report GetProgressReportRow, err error) {
				require.NoError(t, err, "GetProgressReport should not return an error")
				require.NotZero(t, report.ID)
				require.NotEmpty(t, report.EmployeeFirstName)
			},
		},
		{
			name: "get non-existent progress report",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, report GetProgressReportRow, err error) {
				require.Error(t, err, "GetProgressReport should error for non-existent report")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			reportID := tt.setup(ctx, qtx)
			report, err := qtx.GetProgressReport(ctx, reportID)
			tt.checks(t, report, err)
		})
	}
}

func TestUpdateProgressReport(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateProgressReportParams
		checks func(t *testing.T, report ProgressReport, err error)
	}{
		{
			name: "update existing progress report",
			setup: func(ctx context.Context, qtx *Queries) UpdateProgressReportParams {
				report := createRandomProgressReport(ctx, qtx)
				newEmployee := createRandomEmployeeProfile(ctx, qtx)
				return UpdateProgressReportParams{
					ID:         report.ID,
					EmployeeID: &newEmployee.ID,
					Title:      randomStringPtr(25),
					Date: pgtype.Timestamptz{
						Time:  time.Now().Add(time.Hour),
						Valid: true,
					},
					ReportText: randomStringPtr(250),
					Type: NullProgressReportTypeEnum{
						ProgressReportTypeEnum: ProgressReportTypeEnumContactJournal,
						Valid:                  true,
					},
					EmotionalState: NullEmotionalStateEnum{
						EmotionalStateEnum: EmotionalStateEnumDepressed,
						Valid:              true,
					},
				}
			},
			checks: func(t *testing.T, report ProgressReport, err error) {
				require.NoError(t, err, "UpdateProgressReport should not return an error")
				require.NotZero(t, report.ID)
			},
		},
		{
			name: "update non-existent progress report",
			setup: func(ctx context.Context, qtx *Queries) UpdateProgressReportParams {
				newEmployee := createRandomEmployeeProfile(ctx, qtx)
				return UpdateProgressReportParams{
					ID:         uuid.New(),
					EmployeeID: &newEmployee.ID,
					Title:      randomStringPtr(25),
					Date: pgtype.Timestamptz{
						Time:  time.Now().Add(time.Hour),
						Valid: true,
					},
					ReportText: randomStringPtr(250),
					Type: NullProgressReportTypeEnum{
						ProgressReportTypeEnum: ProgressReportTypeEnumContactJournal,
						Valid:                  true,
					},
					EmotionalState: NullEmotionalStateEnum{
						EmotionalStateEnum: EmotionalStateEnumDepressed,
						Valid:              true,
					},
				}
			},
			checks: func(t *testing.T, report ProgressReport, err error) {
				require.Error(t, err, "UpdateProgressReport should error for non-existent report")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			report, err := qtx.UpdateProgressReport(ctx, params)
			tt.checks(t, report, err)
		})
	}
}

func TestDeleteProgressReport(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, err error)
	}{
		{
			name: "delete existing progress report",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				report := createRandomProgressReport(ctx, qtx)
				return report.ID
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteProgressReport should not error")
			},
		},
		{
			name: "delete non-existent progress report",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "DeleteProgressReport should not error for non-existent report")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			reportID := tt.setup(ctx, qtx)
			err = qtx.DeleteProgressReport(ctx, reportID)
			tt.checks(t, err)
		})
	}
}

func TestGetProgressReportsByDateRange(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) GetProgressReportsByDateRangeParams
		checks func(t *testing.T, reports []ProgressReport, err error)
	}{
		{
			name: "get reports in date range",
			setup: func(ctx context.Context, qtx *Queries) GetProgressReportsByDateRangeParams {
				client := createRandomClientDetails(ctx, qtx)
				// Create a report in the range
				createRandomProgressReportForClient(ctx, qtx, client.ID)
				startDate := time.Now().Add(-24 * time.Hour)
				endDate := time.Now().Add(24 * time.Hour)
				return GetProgressReportsByDateRangeParams{
					ClientID: client.ID,
					StartDate: pgtype.Timestamptz{
						Time:  startDate,
						Valid: true,
					},
					EndDate: pgtype.Timestamptz{
						Time:  endDate,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, reports []ProgressReport, err error) {
				require.NoError(t, err, "GetProgressReportsByDateRange should not error")
				require.GreaterOrEqual(t, len(reports), 1)
			},
		},
		{
			name: "get reports outside date range",
			setup: func(ctx context.Context, qtx *Queries) GetProgressReportsByDateRangeParams {
				client := createRandomClientDetails(ctx, qtx)
				startDate := time.Now().Add(48 * time.Hour)
				endDate := time.Now().Add(72 * time.Hour)
				return GetProgressReportsByDateRangeParams{
					ClientID: client.ID,
					StartDate: pgtype.Timestamptz{
						Time:  startDate,
						Valid: true,
					},
					EndDate: pgtype.Timestamptz{
						Time:  endDate,
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, reports []ProgressReport, err error) {
				require.NoError(t, err, "GetProgressReportsByDateRange should not error")
				require.Len(t, reports, 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			reports, err := qtx.GetProgressReportsByDateRange(ctx, params)
			tt.checks(t, reports, err)
		})
	}
}

func TestListAiGeneratedReports(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListAiGeneratedReportsParams
		checks func(t *testing.T, reports []ListAiGeneratedReportsRow, err error)
	}{
		{
			name: "list AI reports for client",
			setup: func(ctx context.Context, qtx *Queries) ListAiGeneratedReportsParams {
				client := createRandomClientDetails(ctx, qtx)
				// Create a report
				createRandomAiGeneratedReportForClient(ctx, qtx, client.ID)
				return ListAiGeneratedReportsParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, reports []ListAiGeneratedReportsRow, err error) {
				require.NoError(t, err, "ListAiGeneratedReports should not error")
				require.GreaterOrEqual(t, len(reports), 1)
			},
		},
		{
			name: "list AI reports with pagination",
			setup: func(ctx context.Context, qtx *Queries) ListAiGeneratedReportsParams {
				client := createRandomClientDetails(ctx, qtx)
				return ListAiGeneratedReportsParams{
					ClientID: client.ID,
					Limit:    5,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, reports []ListAiGeneratedReportsRow, err error) {
				require.NoError(t, err, "ListAiGeneratedReports should not error")
				require.LessOrEqual(t, len(reports), 5)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			reports, err := qtx.ListAiGeneratedReports(ctx, params)
			tt.checks(t, reports, err)
		})
	}
}

func TestListProgressReports(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListProgressReportsParams
		checks func(t *testing.T, reports []ListProgressReportsRow, err error)
	}{
		{
			name: "list progress reports for client",
			setup: func(ctx context.Context, qtx *Queries) ListProgressReportsParams {
				client := createRandomClientDetails(ctx, qtx)
				// Create a report
				createRandomProgressReportForClient(ctx, qtx, client.ID)
				return ListProgressReportsParams{
					ClientID: client.ID,
					Limit:    10,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, reports []ListProgressReportsRow, err error) {
				require.NoError(t, err, "ListProgressReports should not error")
				require.GreaterOrEqual(t, len(reports), 1)
			},
		},
		{
			name: "list progress reports with pagination",
			setup: func(ctx context.Context, qtx *Queries) ListProgressReportsParams {
				client := createRandomClientDetails(ctx, qtx)
				return ListProgressReportsParams{
					ClientID: client.ID,
					Limit:    5,
					Offset:   0,
				}
			},
			checks: func(t *testing.T, reports []ListProgressReportsRow, err error) {
				require.NoError(t, err, "ListProgressReports should not error")
				require.LessOrEqual(t, len(reports), 5)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			reports, err := qtx.ListProgressReports(ctx, params)
			tt.checks(t, reports, err)
		})
	}
}

// Helpers
func createRandomAiGeneratedReport(ctx context.Context, qtx *Queries) AiGeneratedReport {
	client := createRandomClientDetails(ctx, qtx)
	params := CreateAiGeneratedReportParams{
		ClientID:   client.ID,
		ReportText: util.RandomString(100),
		StartDate: pgtype.Date{
			Time:  time.Now().Add(-7 * 24 * time.Hour),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
	}
	report, err := qtx.CreateAiGeneratedReport(ctx, params)
	if err != nil {
		panic(err)
	}
	return report
}

func createRandomAiGeneratedReportForClient(ctx context.Context, qtx *Queries, clientID uuid.UUID) AiGeneratedReport {
	params := CreateAiGeneratedReportParams{
		ClientID:   clientID,
		ReportText: util.RandomString(100),
		StartDate: pgtype.Date{
			Time:  time.Now().Add(-7 * 24 * time.Hour),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
	}
	report, err := qtx.CreateAiGeneratedReport(ctx, params)
	if err != nil {
		panic(err)
	}
	return report
}

func createRandomProgressReport(ctx context.Context, qtx *Queries) ProgressReport {
	client := createRandomClientDetails(ctx, qtx)
	employee := createRandomEmployeeProfile(ctx, qtx)
	params := CreateProgressReportParams{
		ClientID:   client.ID,
		EmployeeID: &employee.ID,
		Title:      randomStringPtr(20),
		Date: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ReportText:     util.RandomString(200),
		Type:           ProgressReportTypeEnumMorningReport,
		EmotionalState: EmotionalStateEnumDepressed,
	}
	report, err := qtx.CreateProgressReport(ctx, params)
	if err != nil {
		panic(err)
	}
	return report
}

func createRandomProgressReportForClient(ctx context.Context, qtx *Queries, clientID uuid.UUID) ProgressReport {
	employee := createRandomEmployeeProfile(ctx, qtx)
	params := CreateProgressReportParams{
		ClientID:   clientID,
		EmployeeID: &employee.ID,
		Title:      randomStringPtr(20),
		Date: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ReportText:     util.RandomString(200),
		Type:           ProgressReportTypeEnumEveningReport,
		EmotionalState: EmotionalStateEnumDepressed,
	}
	report, err := qtx.CreateProgressReport(ctx, params)
	if err != nil {
		panic(err)
	}
	return report
}
