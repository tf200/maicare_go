package clientp

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/logger"
	"maicare_go/pagination"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateProgressReport(ctx context.Context, req *CreateProgressReportRequest, clientID int64) (*CreateProgressReportResponse, error) {
	arg := db.CreateProgressReportParams{
		ClientID:       clientID,
		EmployeeID:     req.EmployeeID,
		Title:          req.Title,
		Date:           pgtype.Timestamptz{Time: req.Date, Valid: true},
		ReportText:     req.ReportText,
		Type:           req.Type,
		EmotionalState: req.EmotionalState,
	}

	report, err := s.Store.CreateProgressReport(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateProgressReport", "Failed to create progress report", zap.Int64("client_id", clientID), zap.Error(err))
		return nil, err
	}
	return &CreateProgressReportResponse{
		ID:             report.ID,
		ClientID:       report.ClientID,
		Date:           report.Date.Time,
		Title:          report.Title,
		ReportText:     report.ReportText,
		EmployeeID:     report.EmployeeID,
		Type:           report.Type,
		EmotionalState: report.EmotionalState,
		CreatedAt:      report.CreatedAt.Time,
	}, nil
}

func (s *clientService) ListProgressReports(ctx *gin.Context, req *ListProgressReportsRequest, clientID int64) (*pagination.Response[ListProgressReportsResponse], error) {
	params := req.GetParams()

	arg := db.ListProgressReportsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}
	reports, err := s.Store.ListProgressReports(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListProgressReports", "Failed to list progress reports", zap.Int64("client_id", clientID), zap.Error(err))
		return nil, err
	}
	if len(reports) == 0 {
		emptyResp := pagination.NewResponse(ctx, req.Request, []ListProgressReportsResponse{}, 0)
		return &emptyResp, nil
	}

	totalCount := reports[0].TotalCount

	var resp []ListProgressReportsResponse
	for _, report := range reports {
		resp = append(resp, ListProgressReportsResponse{
			ID:                     report.ID,
			ClientID:               report.ClientID,
			Date:                   report.Date.Time,
			Title:                  report.Title,
			ReportText:             report.ReportText,
			EmployeeID:             report.EmployeeID,
			Type:                   report.Type,
			EmotionalState:         report.EmotionalState,
			CreatedAt:              report.CreatedAt.Time,
			EmployeeFirstName:      report.EmployeeFirstName,
			EmployeeLastName:       report.EmployeeLastName,
			EmployeeProfilePicture: report.EmployeeProfilePicture,
		})
	}
	pag := pagination.NewResponse(ctx, req.Request, resp, totalCount)
	return &pag, nil
}

func (s *clientService) GetProgressReport(ctx context.Context, reportID int64) (*GetProgressReportResponse, error) {
	report, err := s.Store.GetProgressReport(ctx, reportID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetProgressReport", "Failed to get progress report", zap.Int64("report_id", reportID), zap.Error(err))
		return nil, err
	}
	return &GetProgressReportResponse{
		ID:                     report.ID,
		ClientID:               report.ClientID,
		Date:                   report.Date.Time,
		Title:                  report.Title,
		ReportText:             report.ReportText,
		EmployeeID:             report.EmployeeID,
		Type:                   report.Type,
		EmotionalState:         report.EmotionalState,
		CreatedAt:              report.CreatedAt.Time,
		EmployeeFirstName:      report.EmployeeFirstName,
		EmployeeLastName:       report.EmployeeLastName,
		EmployeeProfilePicture: report.EmployeeProfilePicture,
	}, nil
}

func (s *clientService) UpdateProgressReport(ctx context.Context, req *UpdateProgressReportRequest, reportID int64) (*GetProgressReportResponse, error) {
	arg := db.UpdateProgressReportParams{
		ID:             reportID,
		EmployeeID:     req.EmployeeID,
		Title:          req.Title,
		Date:           pgtype.Timestamptz{Time: req.Date, Valid: true},
		ReportText:     req.ReportText,
		Type:           req.Type,
		EmotionalState: req.EmotionalState,
	}

	report, err := s.Store.UpdateProgressReport(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateProgressReport", "Failed to update progress report", zap.Int64("report_id", reportID), zap.Error(err))
		return nil, err
	}
	return &GetProgressReportResponse{
		ID:             report.ID,
		ClientID:       report.ClientID,
		Date:           report.Date.Time,
		Title:          report.Title,
		ReportText:     report.ReportText,
		EmployeeID:     report.EmployeeID,
		Type:           report.Type,
		EmotionalState: report.EmotionalState,
		CreatedAt:      report.CreatedAt.Time,
	}, nil
}

func (s *clientService) DeleteProgressReport(ctx context.Context, reportID int64) error {
	err := s.Store.DeleteProgressReport(ctx, reportID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteProgressReport", "Failed to delete progress report", zap.Int64("report_id", reportID), zap.Error(err))
		return err
	}
	return nil
}

// TO DO GENERATE AUTO REPORTS

func (s *clientService) GenerateAutoReports(ctx context.Context, req *GenerateAutoReportsRequest, clientID int64) (*GenerateAutoReportsResponse, error) {
	arg := db.GetProgressReportsByDateRangeParams{
		ClientID:  clientID,
		StartDate: pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		EndDate:   pgtype.Timestamptz{Time: req.EndDate, Valid: true},
	}
	reports, err := s.Store.GetProgressReportsByDateRange(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateAutoReports", "Failed to get progress reports for auto report generation", zap.Int64("client_id", clientID), zap.Error(err))
		return nil, err
	}

	var builder strings.Builder
	for _, report := range reports {
		fmt.Fprintf(&builder,
			"Date: %s\nType: %s\nEmotional State: %s\nReport Text: %s\n\n",
			report.Date.Time.GoString(),
			report.Type,
			report.EmotionalState,
			report.ReportText)
	}
	text := builder.String()

	autoRep, err := s.GrpcClient.GenerateAutoReports(ctx, &grpclient.PastReports{
		Text: text,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateAutoReports", "Failed to generate auto reports via gRPC", zap.Int64("client_id", clientID), zap.Error(err))
		return nil, err
	}

	return &GenerateAutoReportsResponse{
		Report: autoRep.GetReport(),
	}, nil

}

func (s *clientService) ConfirmAiProgressReport(ctx context.Context, clientID int64, req *ConfirmProgressReportRequest, reportID int64) (*ConfirmProgressReportResponse, error) {
	progressReport := db.CreateAiGeneratedReportParams{
		ClientID:   clientID,
		ReportText: req.ReportText,
		StartDate:  pgtype.Date{Time: req.Startdate, Valid: true},
		EndDate:    pgtype.Date{Time: req.Enddate, Valid: true},
	}
	createdProgressReport, err := s.Store.CreateAiGeneratedReport(ctx, progressReport)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ConfirmAiProgressReport", "Failed to confirm AI generated progress report", zap.Int64("report_id", reportID), zap.Error(err))
		return nil, err
	}
	return &ConfirmProgressReportResponse{
		ID:         createdProgressReport.ID,
		ClientID:   createdProgressReport.ClientID,
		StartDate:  createdProgressReport.StartDate.Time,
		EndDate:    createdProgressReport.EndDate.Time,
		ReportText: createdProgressReport.ReportText,
		CreatedAt:  createdProgressReport.CreatedAt.Time,
	}, nil
}

func (s *clientService) ListAiGeneratedReports(ctx *gin.Context, req *ListAiGeneratedReportsRequest, clientID int64) (*pagination.Response[ListAiGeneratedReportsResponse], error) {
	params := req.GetParams()

	arg := db.ListAiGeneratedReportsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}
	reports, err := s.Store.ListAiGeneratedReports(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAiGeneratedReports", "Failed to list AI generated reports", zap.Int64("client_id", clientID), zap.Error(err))
		return nil, err
	}

	result := []ListAiGeneratedReportsResponse{}
	for _, report := range reports {
		result = append(result, ListAiGeneratedReportsResponse{
			ID:         report.ID,
			ClientID:   report.ClientID,
			StartDate:  report.StartDate.Time,
			EndDate:    report.EndDate.Time,
			ReportText: report.ReportText,
			CreatedAt:  report.CreatedAt.Time,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, result, reports[0].TotalCount)
	return &pag, nil
}
