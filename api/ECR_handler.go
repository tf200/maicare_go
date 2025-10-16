package api

import (
	"fmt"
	_"maicare_go/pagination"
	"maicare_go/service/ecr"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ======================================================
// ADMIN DASHBOARD
// ======================================================

// DischargeOverviewApi returns a list of discharge overview.
// @Summary Returns a list of discharge overview.
// @Tags ECR
// @Produce json
// @Param page query integer false "page number"
// @Param page_size query integer false "number of items per page"
// @Param filter_type query string false "filter type" Enums(urgent, contract, status_change, all)
// @Success 200 {object} Response[pagination.Response[ecr.DischargeOverviewResponse]]
// @Failure 400 {object} Response[any]
// @Router /ecr/discharge_overview [get]
func (server *Server) DischargeOverviewApi(ctx *gin.Context) {
	var req ecr.DischargeOverviewRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "DischargeOverviewApi", "Failed to bind query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("failed to bind query parameters")))
		return
	}

	overview, err := server.businessService.ECRService.DischargeOverview(ctx, req)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DischargeOverviewApi", "Failed to get discharge overview", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get discharge overview")))
		return
	}

	res := SuccessResponse(overview, "Discharge overview Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)
}

// TotalDischargeCount returns the total discharge count.
// @Summary Returns the total discharge count.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[ecr.TotalDischargeCountResponse]
// @Failure 400 {object} Response[any]
// @Router /ecr/total_discharge_count [get]
func (server *Server) TotalDischargeCountApi(ctx *gin.Context) {
	count, err := server.businessService.ECRService.TotalDischargeCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(count, "Total Discharge Count Retrieved Successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListEmployeesByContractEndDateApi handles the API request for listing employees by contract end date.
// @Summary Lists employees whose contracts are ending soon.
// @Description This endpoint retrieves a list of employees whose contracts are approaching their end date.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[[]ecr.ListEmployeesByContractEndDateResponse]
// @Failure 500 {object} Response[any]
// @Router /ecr/employee_ending_contract [get]
func (server *Server) ListEmployeesByContractEndDateApi(ctx *gin.Context) {
	employees, err := server.businessService.ECRService.ListEmployeesByContractEndDate(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListEmployeesByContractEndDateApi", "Failed to list employees by contract end date", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(employees, "Employees listed successfully"))
}

// ListLatestPaymentsApi handles the API request for listing the latest payments.
// @Summary Lists the latest payments.
// @Description This endpoint retrieves a list of the latest payments.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[[]ecr.ListLatestPaymentsResponse]
// @Failure 500 {object} Response[any]
// @Router /ecr/latest_payments [get]
func (server *Server) ListLatestPaymentsApi(ctx *gin.Context) {
	payments, err := server.businessService.ECRService.ListLatestPayments(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListLatestPaymentsApi", "Failed to list latest payments", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(payments, "Latest payments listed successfully"))
}

// ListUpcomingAppointmentsApi handles the API request for listing upcoming appointments.
// @Summary Lists upcoming appointments.
// @Description This endpoint retrieves a list of upcoming appointments.
// @Tags ECR
// @Produce json
// @Success 200 {object} Response[[]ecr.ListUpcomingAppointmentsResponse]
// @Failure 500 {object} Response[any]
// @Router /ecr/upcoming_appointments [get]
func (server *Server) ListUpcomingAppointmentsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListUpcomingAppointmentsApi", "Failed to get auth payload", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("failed to get auth payload")))
		return
	}

	appointments, err := server.businessService.ECRService.ListUpcomingAppointments(ctx, payload.EmployeeID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListUpcomingAppointmentsApi", "Failed to list upcoming appointments", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(appointments, "Upcoming appointments listed successfully"))
}
