package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
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
