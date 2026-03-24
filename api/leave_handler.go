package api

import (
	"errors"
	"net/http"

	"maicare_go/service/leave"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateLeaveRequestApi creates a new leave request for the authenticated employee.
// @Summary Create leave request
// @Description Creates a new leave request for the authenticated employee.
// @Tags leave-requests
// @Accept json
// @Produce json
// @Param request body leave.CreateLeaveRequestRequest true "Leave request payload"
// @Success 201 {object} Response[leave.CreateLeaveRequestResponse]
// @Router /leave-requests [post]
func (server *Server) CreateLeaveRequestApi(ctx *gin.Context) {
	var req leave.CreateLeaveRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.CreateLeaveRequest(ctx, payload.EmployeeID, &req)
	if err != nil {
		if errors.Is(err, leave.ErrLeaveRequestInvalidRequest) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Leave request created successfully"))
}

// CreateLeaveRequestByAdminApi creates a new leave request for a selected employee.
// @Summary Admin create leave request
// @Description Creates a new leave request for the selected employee.
// @Tags leave-requests
// @Accept json
// @Produce json
// @Param request body leave.CreateLeaveRequestByAdminRequest true "Admin leave request payload"
// @Success 201 {object} Response[leave.CreateLeaveRequestResponse]
// @Router /leave-requests/admin [post]
func (server *Server) CreateLeaveRequestByAdminApi(ctx *gin.Context) {
	var req leave.CreateLeaveRequestByAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.CreateLeaveRequestByAdmin(ctx, payload.EmployeeID, &req)
	if err != nil {
		if errors.Is(err, leave.ErrLeaveRequestInvalidRequest) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Leave request created successfully"))
}

// UpdateLeaveRequestApi updates a leave request for the authenticated employee.
// @Summary Update leave request
// @Description Updates a pending future leave request owned by the authenticated employee.
// @Tags leave-requests
// @Accept json
// @Produce json
// @Param id path string true "Leave request ID"
// @Param request body leave.UpdateLeaveRequestRequest true "Leave request update payload"
// @Success 200 {object} Response[leave.UpdateLeaveRequestResponse]
// @Router /leave-requests/{id} [put]
func (server *Server) UpdateLeaveRequestApi(ctx *gin.Context) {
	leaveRequestID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req leave.UpdateLeaveRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.UpdateLeaveRequest(ctx, payload.EmployeeID, leaveRequestID, &req)
	if err != nil {
		ctx.JSON(mapLeaveRequestErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave request updated successfully"))
}

// UpdateLeaveRequestByAdminApi updates a leave request by admin/coordinator.
// @Summary Admin update leave request
// @Description Updates an editable leave request by admin/coordinator with a mandatory admin update note.
// @Tags leave-requests
// @Accept json
// @Produce json
// @Param id path string true "Leave request ID"
// @Param request body leave.UpdateLeaveRequestAdminRequest true "Admin leave request update payload"
// @Success 200 {object} Response[leave.UpdateLeaveRequestResponse]
// @Router /leave-requests/{id}/admin [put]
func (server *Server) UpdateLeaveRequestByAdminApi(ctx *gin.Context) {
	leaveRequestID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req leave.UpdateLeaveRequestAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.UpdateLeaveRequestByAdmin(ctx, payload.EmployeeID, leaveRequestID, &req)
	if err != nil {
		ctx.JSON(mapLeaveRequestErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave request updated successfully"))
}

// DecideLeaveRequestByAdminApi approves or rejects a pending leave request.
// @Summary Decide leave request
// @Description Approves or rejects a pending leave request by admin/coordinator.
// @Tags leave-requests
// @Accept json
// @Produce json
// @Param id path string true "Leave request ID"
// @Param request body leave.DecideLeaveRequestRequest true "Leave request decision payload"
// @Success 200 {object} Response[leave.DecideLeaveRequestResponse]
// @Router /leave-requests/{id}/decision [post]
func (server *Server) DecideLeaveRequestByAdminApi(ctx *gin.Context) {
	leaveRequestID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req leave.DecideLeaveRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.DecideLeaveRequestByAdmin(ctx, payload.EmployeeID, leaveRequestID, &req)
	if err != nil {
		ctx.JSON(mapLeaveRequestErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave request decided successfully"))
}

// ListMyLeaveRequestsApi lists leave requests for the authenticated employee.
// @Summary List my leave requests
// @Description Lists paginated leave requests for the authenticated employee.
// @Tags leave-requests
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param status query string false "Leave request status filter"
// @Success 200 {object} Response[pagination.Response[leave.LeaveRequestListItem]]
// @Router /leave-requests/my [get]
func (server *Server) ListMyLeaveRequestsApi(ctx *gin.Context) {
	var req leave.ListMyLeaveRequestsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.ListMyLeaveRequests(ctx, payload.EmployeeID, &req)
	if err != nil {
		if errors.Is(err, leave.ErrLeaveRequestInvalidRequest) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave requests retrieved successfully"))
}

// GetMyLeaveRequestStatsApi returns current-year leave request stats for the authenticated employee.
// @Summary Get my leave request stats
// @Description Returns current-year counts for open, approved, rejected leave requests, and approved sickness absence.
// @Tags leave-requests
// @Produce json
// @Success 200 {object} Response[leave.MyLeaveRequestStatsResponse]
// @Failure 401,403,500 {object} Response[any]
// @Router /leave-requests/my/stats [get]
func (server *Server) GetMyLeaveRequestStatsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.GetMyLeaveRequestStats(ctx, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave request stats retrieved successfully"))
}

// GetLeaveRequestStatsApi returns current-year leave request stats across all employees.
// @Summary Get leave request stats
// @Description Returns current-year counts for open, approved, rejected leave requests, and approved sickness absence across all employees.
// @Tags leave-requests
// @Produce json
// @Success 200 {object} Response[leave.LeaveRequestStatsResponse]
// @Failure 401,500 {object} Response[any]
// @Router /leave-requests/stats [get]
func (server *Server) GetLeaveRequestStatsApi(ctx *gin.Context) {
	_, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.GetLeaveRequestStats(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave request stats retrieved successfully"))
}

func mapLeaveRequestErrorStatus(err error) int {
	switch {
	case errors.Is(err, leave.ErrLeaveRequestInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, leave.ErrLeaveBalanceInvalidAdjust):
		return http.StatusBadRequest
	case errors.Is(err, leave.ErrLeaveRequestForbidden):
		return http.StatusForbidden
	case errors.Is(err, leave.ErrLeaveRequestNotFound):
		return http.StatusNotFound
	case errors.Is(err, leave.ErrLeaveRequestStateInvalid):
		return http.StatusConflict
	case errors.Is(err, leave.ErrLeaveBalanceInsufficient):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// ListLeaveRequestsApi lists leave requests for admins with filters.
// @Summary List leave requests
// @Description Lists paginated leave requests with optional status and employee filters.
// @Tags leave-requests
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param status query string false "Leave request status filter"
// @Param employee_search query string false "Employee name search (first/last name)"
// @Success 200 {object} Response[pagination.Response[leave.LeaveRequestListItem]]
// @Router /leave-requests [get]
func (server *Server) ListLeaveRequestsApi(ctx *gin.Context) {
	var req leave.ListLeaveRequestsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.ListLeaveRequests(ctx, &req)
	if err != nil {
		if errors.Is(err, leave.ErrLeaveRequestInvalidRequest) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave requests retrieved successfully"))
}

// ListLeaveBalancesApi lists leave balances for admins with filters.
// @Summary List leave balances
// @Description Lists paginated leave balances with optional employee and year filters.
// @Tags leave-balances
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param employee_search query string false "Employee name search (first/last name)"
// @Param year query int false "Year filter"
// @Success 200 {object} Response[pagination.Response[leave.LeaveBalanceListItem]]
// @Router /leave-balances [get]
func (server *Server) ListLeaveBalancesApi(ctx *gin.Context) {
	var req leave.ListLeaveBalancesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.ListLeaveBalances(ctx, &req)
	if err != nil {
		if errors.Is(err, leave.ErrLeaveRequestInvalidRequest) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave balances retrieved successfully"))
}

// ListMyLeaveBalancesApi lists leave balances for the authenticated employee.
// @Summary List my leave balances
// @Description Lists paginated leave balances for the authenticated employee with optional year filter.
// @Tags leave-balances
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param year query int false "Year filter"
// @Success 200 {object} Response[pagination.Response[leave.LeaveBalanceListItem]]
// @Router /leave-balances/my [get]
func (server *Server) ListMyLeaveBalancesApi(ctx *gin.Context) {
	var req leave.ListMyLeaveBalancesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.ListMyLeaveBalances(ctx, payload.EmployeeID, &req)
	if err != nil {
		if errors.Is(err, leave.ErrLeaveRequestInvalidRequest) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave balances retrieved successfully"))
}

// AdjustLeaveBalanceApi applies an admin adjustment to leave balance totals.
// @Summary Adjust leave balance
// @Description Applies admin correction deltas to leave balance totals and records an audit note.
// @Tags leave-balances
// @Accept json
// @Produce json
// @Param request body leave.AdjustLeaveBalanceRequest true "Leave balance adjustment payload"
// @Success 200 {object} Response[leave.AdjustLeaveBalanceResponse]
// @Router /leave-balances/adjust [post]
func (server *Server) AdjustLeaveBalanceApi(ctx *gin.Context) {
	var req leave.AdjustLeaveBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LeaveService.AdjustLeaveBalance(ctx, payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(mapLeaveRequestErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Leave balance adjusted successfully"))
}
