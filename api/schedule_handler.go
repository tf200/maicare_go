package api

import (
	"fmt"
	"maicare_go/service/schedule"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a new schedule
// @Description Create a new schedule for an employee at a specific location. Supports both custom schedules and preset shifts.
// @Description Set is_custom=true and provide start_datetime/end_datetime for custom schedules
// @Description Set is_custom=false and provide location_shift_id/shift_date for preset shifts
// @Tags Schedule
// @Accept json
// @Produce json
// @Param request body CreateScheduleRequest true "Create Schedule Request"
// @Success 200 {object} Response[CreateScheduleResponse] "Schedule created successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules [post]
func (server *Server) CreateScheduleApi(ctx *gin.Context) {
	var req schedule.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized")))
		return
	}

	schedule, err := server.businessService.ScheduleService.CreateSchedule(ctx, payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create schedule: %w", err)))
		return
	}

	res := SuccessResponse(schedule, "Schedule created successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get monthly schedules by location
// @Description Get all schedules for a specific location for a given month and year
// @Tags Schedule
// @Produce json
// @Param id path int true "Location ID"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} Response[[]GetMonthlySchedulesByLocationResponse] "Monthly schedules retrieved successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /locations/{id}/monthly_schedules [get]
func (server *Server) GetMonthlySchedulesByLocationApi(ctx *gin.Context) {
	locationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req schedule.GetMonthlySchedulesByLocationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	response, err := server.businessService.ScheduleService.GetMonthlySchedulesByLocation(ctx, locationID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get monthly schedules: %w", err)))
		return
	}

	res := SuccessResponse(response, "Schedules retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get daily schedules by location
// @Description Get all schedules for a specific location for a given day
// @Tags Schedule
// @Produce json
// @Param id path int true "Location ID"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param day query int true "Day"
// @Success 200 {object} Response[GetDailySchedulesByLocationResponse] "Daily schedules retrieved successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /locations/{id}/daily_schedules [get]
func (server *Server) GetDailySchedulesByLocationApi(ctx *gin.Context) {
	locationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req schedule.GetDailySchedulesByLocationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	response, err := server.businessService.ScheduleService.GetDailySchedulesByLocation(ctx, locationID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get daily schedules: %w", err)))
		return
	}

	res := SuccessResponse(response, "Daily schedules retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get schedule by ID
// @Description Get a schedule by its ID
// @Tags Schedule
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} Response[GetScheduleByIdResponse] "Schedule retrieved successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules/{id} [get]
func (server *Server) GetScheduleByIDApi(ctx *gin.Context) {
	scheduleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid schedule ID format")))
		return
	}

	schedule, err := server.businessService.ScheduleService.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get schedule: %w", err)))
		return
	}

	res := SuccessResponse(schedule, "Schedule retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update an existing schedule
// @Description Update an existing schedule for an employee at a specific location. Supports both custom schedules and preset shifts.
// @Description Set is_custom=true and provide start_datetime/end_datetime for custom schedules
// @Description Set is_custom=false and provide location_shift_id/shift_date for preset shifts
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param request body UpdateScheduleRequest true "Update Schedule Request"
// @Success 200 {object} Response[UpdateScheduleResponse] "Schedule updated successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 404 {object} Response[any] "Schedule not found"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules/{id} [put]
func (server *Server) UpdateScheduleApi(ctx *gin.Context) {

	scheduleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid schedule ID format")))
		return
	}

	var req schedule.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized")))
		return
	}

	schedule, err := server.businessService.ScheduleService.UpdateSchedule(ctx, scheduleID, payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update schedule: %w", err)))
		return
	}

	res := SuccessResponse(schedule, "Schedule updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteScheduleApi deletes a schedule by its ID.
// @Summary Delete a schedule
// @Description Delete an existing schedule by its ID
// @Tags Schedule
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} Response[any] "Schedule deleted successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules/{id} [delete]
func (server *Server) DeleteScheduleApi(ctx *gin.Context) {
	scheduleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.ScheduleService.DeleteSchedule(ctx, scheduleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete schedule: %w", err)))
		return
	}

	res := SuccessResponse[any](nil, "Schedule deleted successfully")
	ctx.JSON(http.StatusOK, res)
}
