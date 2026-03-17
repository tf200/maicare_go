package api

import (
	"errors"
	"net/http"

	"maicare_go/service/schedule"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) CreateShiftSwapRequestApi(ctx *gin.Context) {
	var req schedule.CreateShiftSwapRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.ScheduleService.CreateShiftSwapRequest(ctx, payload.EmployeeID, &req)
	if err != nil {
		status := mapShiftSwapErrorStatus(err)
		ctx.JSON(status, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Shift swap request created successfully"))
}

func (server *Server) RespondShiftSwapRequestApi(ctx *gin.Context) {
	swapID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req schedule.RespondShiftSwapRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.ScheduleService.RespondToShiftSwapRequest(ctx, payload.EmployeeID, swapID, &req)
	if err != nil {
		status := mapShiftSwapErrorStatus(err)
		ctx.JSON(status, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Shift swap request response recorded successfully"))
}

func (server *Server) AdminDecisionShiftSwapRequestApi(ctx *gin.Context) {
	swapID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req schedule.AdminDecisionShiftSwapRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.ScheduleService.AdminDecisionShiftSwapRequest(ctx, payload.EmployeeID, swapID, &req)
	if err != nil {
		status := mapShiftSwapErrorStatus(err)
		ctx.JSON(status, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Shift swap request admin decision recorded successfully"))
}

func (server *Server) ListMyShiftSwapRequestsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.ScheduleService.ListMyShiftSwapRequests(ctx, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Shift swap requests retrieved successfully"))
}

func (server *Server) ListShiftSwapRequestsApi(ctx *gin.Context) {
	var req schedule.ListShiftSwapRequestsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := server.businessService.ScheduleService.ListShiftSwapRequests(ctx, &req)
	if err != nil {
		if errors.Is(err, schedule.ErrShiftSwapInvalidRequest) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Shift swap requests retrieved successfully"))
}

func mapShiftSwapErrorStatus(err error) int {
	switch {
	case errors.Is(err, schedule.ErrShiftSwapInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, schedule.ErrScheduleNotFound),
		errors.Is(err, schedule.ErrShiftSwapNotFound):
		return http.StatusNotFound
	case errors.Is(err, schedule.ErrShiftSwapStateInvalid),
		errors.Is(err, schedule.ErrShiftSwapDuplicateActiveRequest),
		errors.Is(err, schedule.ErrShiftSwapExpired),
		errors.Is(err, schedule.ErrShiftSwapScheduleOwnership),
		errors.Is(err, schedule.ErrShiftSwapConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
