package clientp

import (
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"

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


  