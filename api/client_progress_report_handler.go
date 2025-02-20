package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateProgressReportRequest defines the request payload for CreateProgressReport API
type CreateProgressReportRequest struct {
	EmployeeID     *int64    `json:"employee_id"`
	Title          *string   `json:"title"`
	Date           time.Time `json:"date"`
	ReportText     string    `json:"report_text" binding:"required"`
	Type           string    `json:"type" binding:"required,oneof=morning_report evening_report night_report shift_report one_to_one_report process_report contact_journal other"`
	EmotionalState string    `json:"emotional_state" binding:"required,oneof=normal excited happy sad angry anxious depressed"`
}

// CreateProgressReportResponse defines the response payload for CreateProgressReport API
type CreateProgressReportResponse struct {
	ID             int64     `json:"id"`
	ClientID       int64     `json:"client_id"`
	Date           time.Time `json:"date"`
	Title          *string   `json:"title"`
	ReportText     string    `json:"report_text"`
	EmployeeID     *int64    `json:"employee_id"`
	Type           string    `json:"type"`
	EmotionalState string    `json:"emotional_state"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateProgressReportApi creates a new progress report for a client
// @Summary Create a new progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body CreateProgressReportRequest true "Progress Report Request"
// @Success 201 {object} Response[CreateProgressReportResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/progress_reports [post]
func (server *Server) CreateProgressReportApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var request CreateProgressReportRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	progressReport := db.CreateProgressReportParams{
		ClientID:       clientID,
		EmployeeID:     request.EmployeeID,
		Title:          request.Title,
		Date:           pgtype.Timestamptz{Time: request.Date, Valid: true},
		ReportText:     request.ReportText,
		Type:           request.Type,
		EmotionalState: request.EmotionalState,
	}

	createdProgressReport, err := server.store.CreateProgressReport(ctx, progressReport)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateProgressReportResponse{
		ID:             createdProgressReport.ID,
		ClientID:       createdProgressReport.ClientID,
		Date:           createdProgressReport.Date.Time,
		Title:          createdProgressReport.Title,
		ReportText:     createdProgressReport.ReportText,
		EmployeeID:     createdProgressReport.EmployeeID,
		Type:           createdProgressReport.Type,
		EmotionalState: createdProgressReport.EmotionalState,
		CreatedAt:      createdProgressReport.CreatedAt.Time,
	}, "Progress Report created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListProgressReportsRequest defines the request payload for ListProgressReports API
type ListProgressReportsRequest struct {
	pagination.Request
}

// ListProgressReportsResponse defines the response payload for ListProgressReports API
type ListProgressReportsResponse struct {
	ID                int64     `json:"id"`
	ClientID          int64     `json:"client_id"`
	Date              time.Time `json:"date"`
	Title             *string   `json:"title"`
	ReportText        string    `json:"report_text"`
	EmployeeID        *int64    `json:"employee_id"`
	Type              string    `json:"type"`
	EmotionalState    string    `json:"emotional_state"`
	CreatedAt         time.Time `json:"created_at"`
	EmployeeFirstName string    `json:"employee_first_name"`
	EmployeeLastName  string    `json:"employee_last_name"`
}

// ListProgressReportsApi lists all progress reports for a client
// @Summary List all progress reports for a client
// @Tags progress_reports
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[[]ListProgressReportsResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/progress_reports [get]
func (server *Server) ListProgressReportsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListProgressReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	arg := db.ListProgressReportsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}

	progressReports, err := server.store.ListProgressReports(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(progressReports) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse([]ListProgressReportsResponse{}, "No progress reports found"))
		return
	}

	repList := make([]ListProgressReportsResponse, len(progressReports))
	for i, progressReport := range progressReports {
		repList[i] = ListProgressReportsResponse{
			ID:                progressReport.ID,
			ClientID:          progressReport.ClientID,
			Date:              progressReport.Date.Time,
			Title:             progressReport.Title,
			ReportText:        progressReport.ReportText,
			EmployeeID:        progressReport.EmployeeID,
			Type:              progressReport.Type,
			EmotionalState:    progressReport.EmotionalState,
			CreatedAt:         progressReport.CreatedAt.Time,
			EmployeeFirstName: progressReport.EmployeeFirstName,
			EmployeeLastName:  progressReport.EmployeeLastName,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, repList, progressReports[0].TotalCount)
	res := SuccessResponse(pag, "Progress reports retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetProgressReportResponse defines the response payload for GetProgressReport API
type GetProgressReportResponse struct {
	ID                int64     `json:"id"`
	ClientID          int64     `json:"client_id"`
	Date              time.Time `json:"date"`
	Title             *string   `json:"title"`
	ReportText        string    `json:"report_text"`
	EmployeeID        *int64    `json:"employee_id"`
	Type              string    `json:"type"`
	EmotionalState    string    `json:"emotional_state"`
	CreatedAt         time.Time `json:"created_at"`
	EmployeeFirstName string    `json:"employee_first_name"`
	EmployeeLastName  string    `json:"employee_last_name"`
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

	progressReport, err := server.store.GetProgressReport(ctx, reportID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetProgressReportResponse{
		ID:                progressReport.ID,
		ClientID:          progressReport.ClientID,
		Date:              progressReport.Date.Time,
		Title:             progressReport.Title,
		ReportText:        progressReport.ReportText,
		EmployeeID:        progressReport.EmployeeID,
		Type:              progressReport.Type,
		EmotionalState:    progressReport.EmotionalState,
		CreatedAt:         progressReport.CreatedAt.Time,
		EmployeeFirstName: progressReport.EmployeeFirstName,
		EmployeeLastName:  progressReport.EmployeeLastName,
	}, "Progress Report retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateProgressReportRequest defines the request payload for UpdateProgressReport API
type UpdateProgressReportRequest struct {
	ClientID       int64     `json:"client_id"`
	EmployeeID     *int64    `json:"employee_id"`
	Title          *string   `json:"title"`
	Date           time.Time `json:"date"`
	ReportText     *string   `json:"report_text"`
	Type           *string   `json:"type"`
	EmotionalState *string   `json:"emotional_state"`
}

// UpdateProgressReportResponse defines the response payload for UpdateProgressReport API
type UpdateProgressReportResponse struct {
	ID             int64     `json:"id"`
	ClientID       int64     `json:"client_id"`
	Date           time.Time `json:"date"`
	Title          *string   `json:"title"`
	ReportText     string    `json:"report_text"`
	EmployeeID     *int64    `json:"employee_id"`
	Type           string    `json:"type"`
	EmotionalState string    `json:"emotional_state"`
	CreatedAt      time.Time `json:"created_at"`
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

	var req UpdateProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	progressReport := db.UpdateProgressReportParams{
		ID:             reportID,
		EmployeeID:     req.EmployeeID,
		Title:          req.Title,
		Date:           pgtype.Timestamptz{Time: req.Date, Valid: true},
		ReportText:     req.ReportText,
		Type:           req.Type,
		EmotionalState: req.EmotionalState,
	}

	updatedProgressReport, err := server.store.UpdateProgressReport(ctx, progressReport)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateProgressReportResponse{
		ID:             updatedProgressReport.ID,
		ClientID:       updatedProgressReport.ClientID,
		Date:           updatedProgressReport.Date.Time,
		Title:          updatedProgressReport.Title,
		ReportText:     updatedProgressReport.ReportText,
		EmployeeID:     updatedProgressReport.EmployeeID,
		Type:           updatedProgressReport.Type,
		EmotionalState: updatedProgressReport.EmotionalState,
		CreatedAt:      updatedProgressReport.CreatedAt.Time,
	}, "Progress Report updated successfully")

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
