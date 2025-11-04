package api

import (
	"net/http"
	"strconv"

	_ "maicare_go/pagination"
	clientp "maicare_go/service/client"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	clientID, err := uuid.Parse(id)
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
	clientID, err := uuid.Parse(id)
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
// @Success 200 {object} Response[clientp.GetProgressReportResponse]
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
// @Param request body clientp.UpdateProgressReportRequest true "Progress Report Request"
// @Success 200 {object} Response[clientp.UpdateProgressReportResponse]
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

// GenerateAutoReportsApi is the handler for the auto reports generation API
// @Summary Generate auto reports
// @Description Generate auto reports
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.GenerateAutoReportsRequest true "Request body"
// @Success 200 {object} Response[clientp.GenerateAutoReportsResponse]
// @Router /clients/{id}/ai_progress_reports [post]
func (server *Server) GenerateAutoReportsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.GenerateAutoReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	autoReports, err := server.businessService.ClientService.GenerateAutoReports(ctx, &req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(autoReports, "Auto reports generated successfully")

	ctx.JSON(http.StatusOK, res)
}

// ConfirmProgressReportApi creates a new progress report for a client
// @Summary Confirm a progress report for a client
// @Tags progress_reports
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.ConfirmProgressReportRequest true "Progress Report Request"
// @Success 201 {object} Response[clientp.ConfirmProgressReportResponse]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/ai_progress_reports/confirm [post]
func (server *Server) ConfirmProgressReportApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ConfirmProgressReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.ConfirmAiProgressReport(ctx, clientID, &req, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Progress Report created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListAiGeneratedReportsApi lists all AI generated reports for a client
// @Summary List all AI generated reports for a client
// @Tags progress_reports
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[[]clientp.ListAiGeneratedReportsResponse]]
// @Failure 400,404 {object} Response[any]
// @Router /clients/{id}/ai_progress_reports [get]
func (server *Server) ListAiGeneratedReportsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.ListAiGeneratedReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.ClientService.ListAiGeneratedReports(ctx, &req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(pag, "Progress reports retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}
