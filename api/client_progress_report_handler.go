package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateProgressReportApi creates a new progress report for a client
// @Summary Create a new progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.CreateProgressReportRequest true "Progress Report Request"
// @Success 201 {object} Response[clientp.CreateProgressReportResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/progress_reports [post]
func (server *Server) CreateProgressReportApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.CreateProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	progressReport, err := server.businessService.ClientService.CreateProgressReport(ctx, &req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(progressReport, "Progress Report created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListProgressReportsApi lists all progress reports for a client
// @Summary List all progress reports for a client
// @Tags progress_reports
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[[]clientp.ListProgressReportsResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/progress_reports [get]
func (server *Server) ListProgressReportsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ListProgressReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	pag, err := server.businessService.ClientService.ListProgressReports(ctx, &req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(pag, "Progress reports retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetProgressReportApi retrieves a progress report for a client
// @Summary Retrieve a progress report for a client
// @Tags progress_reports
// @Produce json
// @Param id path int true "Client ID"
// @Param report_id path int true "Progress Report ID"
// @Success 200 {object} Response[GetProgressReportResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/progress_reports/{report_id} [get]
func (server *Server) GetProgressReportApi(ctx *gin.Context) {
	id := ctx.Param("report_id")
	reportID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	progressReport, err := server.businessService.ClientService.GetProgressReport(ctx, reportID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(progressReport, "Progress Report retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateProgressReportApi updates a progress report for a client
// @Summary Update a progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param report_id path int true "Progress Report ID"
// @Param request body UpdateProgressReportRequest true "Progress Report Request"
// @Success 200 {object} Response[UpdateProgressReportResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/progress_reports/{report_id} [put]
func (server *Server) UpdateProgressReportApi(ctx *gin.Context) {
	id := ctx.Param("report_id")
	reportID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	updatedReport, err := server.businessService.ClientService.UpdateProgressReport(ctx, &req, reportID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(updatedReport, "Progress Report updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteProgressReportApi deletes a progress report for a client
// @Summary Delete a progress report for a client
// @Tags progress_reports
// @Produce json
// @Param id path int true "Client ID"
// @Param report_id path int true "Progress Report ID"
// @Success 200 {object} Response[any]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/progress_reports/{report_id} [delete]
func (server *Server) DeleteProgressReportApi(ctx *gin.Context) {
	id := ctx.Param("report_id")
	reportID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.ClientService.DeleteProgressReport(ctx, reportID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse[any](nil, "Progress Report deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// GenerateAutoReportsRequest is the request format for the auto reports generation API
type GenerateAutoReportsRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// GenerateAutoReportsResponse is the response format for the auto reports generation API
type GenerateAutoReportsResponse struct {
	Report string `json:"report"`
}

// GenerateAutoReportsApi is the handler for the auto reports generation API
// @Summary Generate auto reports
// @Description Generate auto reports
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body GenerateAutoReportsRequest true "Request body"
// @Success 200 {object} Response[GenerateAutoReportsResponse]
// @Router /clients/{id}/ai_progress_reports [post]
func (server *Server) GenerateAutoReportsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req GenerateAutoReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetProgressReportsByDateRangeParams{
		ClientID:  clientID,
		StartDate: pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		EndDate:   pgtype.Timestamptz{Time: req.EndDate, Valid: true},
	}
	progressReports, err := server.store.GetProgressReportsByDateRange(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var builder strings.Builder
	for _, report := range progressReports {
		fmt.Fprintf(&builder,
			"Date: %s\nType: %s\nEmotional State: %s\nReport Text: %s\n\n",
			report.Date.Time.GoString(),
			report.Type,
			report.EmotionalState,
			report.ReportText)
	}
	text := builder.String()

	autoRep, err := server.aiHandler.GenerateAutoReports(text, "google/gemini-2.0-flash-001")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GenerateAutoReportsResponse{
		Report: autoRep.GeneratedReport,
	}, "Auto reports generated successfully")

	ctx.JSON(http.StatusOK, res)

}

// ConfirmProgressReportRequest defines the request payload for ConfirmProgressReport API
type ConfirmProgressReportRequest struct {
	ReportText string    `json:"report_text" binding:"required"`
	Startdate  time.Time `json:"start_date" binding:"required"`
	Enddate    time.Time `json:"end_date" binding:"required"`
}

// ConfirmProgressReportResponse defines the response payload for ConfirmProgressReport API
type ConfirmProgressReportResponse struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	ReportText string    `json:"report_text"`
	CreatedAt  time.Time `json:"created_at"`
}

// ConfirmProgressReportApi creates a new progress report for a client
// @Summary Confirm a progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body ConfirmProgressReportRequest true "Progress Report Request"
// @Success 201 {object} Response[ConfirmProgressReportResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/ai_progress_reports/confirm [post]
func (server *Server) ConfirmProgressReportApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ConfirmProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	progressReport := db.CreateAiGeneratedReportParams{
		ClientID:   clientID,
		ReportText: req.ReportText,
		StartDate:  pgtype.Date{Time: req.Startdate, Valid: true},
		EndDate:    pgtype.Date{Time: req.Enddate, Valid: true},
	}

	createdProgressReport, err := server.store.CreateAiGeneratedReport(ctx, progressReport)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(ConfirmProgressReportResponse{
		ID:         createdProgressReport.ID,
		ClientID:   createdProgressReport.ClientID,
		StartDate:  createdProgressReport.StartDate.Time,
		EndDate:    createdProgressReport.EndDate.Time,
		ReportText: createdProgressReport.ReportText,
		CreatedAt:  createdProgressReport.CreatedAt.Time,
	}, "Progress Report created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListAiGeneratedReportsRequest defines the request payload for ListAiGeneratedReports API
type ListAiGeneratedReportsRequest struct {
	pagination.Request
}

// ListAiGeneratedReportsResponse defines the response payload for ListAiGeneratedReports API
type ListAiGeneratedReportsResponse struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	ReportText string    `json:"report_text"`
	CreatedAt  time.Time `json:"created_at"`
}

// ListAiGeneratedReportsApi lists all AI generated reports for a client
// @Summary List all AI generated reports for a client
// @Tags progress_reports
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[[]ListAiGeneratedReportsResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/ai_progress_reports [get]
func (server *Server) ListAiGeneratedReportsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListAiGeneratedReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	arg := db.ListAiGeneratedReportsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}

	progressReports, err := server.store.ListAiGeneratedReports(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(progressReports) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListAiGeneratedReportsResponse{}, 0)
		res := SuccessResponse(pag, "No progress reports found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	repList := make([]ListAiGeneratedReportsResponse, len(progressReports))
	for i, progressReport := range progressReports {
		repList[i] = ListAiGeneratedReportsResponse{
			ID:         progressReport.ID,
			ClientID:   progressReport.ClientID,
			StartDate:  progressReport.StartDate.Time,
			EndDate:    progressReport.EndDate.Time,
			ReportText: progressReport.ReportText,
			CreatedAt:  progressReport.CreatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, repList, progressReports[0].TotalCount)
	res := SuccessResponse(pag, "Progress reports retrieved successfully")
	ctx.JSON(http.StatusOK, res)

}
