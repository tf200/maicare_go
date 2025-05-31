package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateScheduleRequest represents the request body for creating a schedule.
type CreateScheduleRequest struct {
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime" example:"2023-10-01T09:00:00Z"`
	EndDatetime   time.Time `json:"end_datetime" example:"2023-10-01T17:00:00Z"`
}

// CreateScheduleResponse represents the response body after creating a schedule.
type CreateScheduleResponse struct {
	ID            int64     `json:"id"`
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// @Summary Create a new schedule
// @Description Create a new schedule for an employee at a specific location
// @Tags Schedule
// @Accept json
// @Produce json
// @Param request body CreateScheduleRequest true "Create Schedule Request"
// @Success 200 {object} Response[CreateScheduleRequest] "Schedule created successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules [post]
func (server *Server) CreateScheduleApi(ctx *gin.Context) {
	var req CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate that StartDatetime is before EndDatetime
	if req.StartDatetime.After(req.EndDatetime) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("start_datetime must be before end_datetime")))
		return
	}

	arg := db.CreateScheduleParams{
		EmployeeID:    req.EmployeeID,
		LocationID:    req.LocationID,
		StartDatetime: pgtype.Timestamp{Time: req.StartDatetime, Valid: true},
		EndDatetime:   pgtype.Timestamp{Time: req.EndDatetime, Valid: true},
	}
	schedule, err := server.store.CreateSchedule(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateScheduleResponse{
		ID:            schedule.ID,
		EmployeeID:    schedule.EmployeeID,
		LocationID:    schedule.LocationID,
		StartDatetime: schedule.StartDatetime.Time,
		EndDatetime:   schedule.EndDatetime.Time,
		CreatedAt:     schedule.CreatedAt.Time,
		UpdatedAt:     schedule.UpdatedAt.Time,
	}, "Schedule created successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetMonthlySchedulesByLocationApi retrieves the monthly schedules for a specific location.
type GetMonthlySchedulesByLocationRequest struct {
	Year  int32 `form:"year"`
	Month int32 `form:"month"`
}

// Shift represents a work shift for an employee.
type Shift struct {
	ShiftID           int64     `json:"shift_id"`
	EmployeeID        int64     `json:"employee_id"`
	EmployeeFirstName string    `json:"employee_first_name"`
	EmployeeLastName  string    `json:"employee_last_name"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	LocationID        int64     `json:"location_id"`
}

// GetMonthlySchedulesByLocationResponse represents the response body for monthly schedules.
type GetMonthlySchedulesByLocationResponse struct {
	Date   string  `json:"date"`
	Shifts []Shift `json:"shifts"`
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
	var req GetMonthlySchedulesByLocationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetMonthlySchedulesByLocationParams{
		Year:       req.Year,
		Month:      req.Month,
		LocationID: locationID,
	}

	schedules, err := server.store.GetMonthlySchedulesByLocation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	calendar := make(map[string][]Shift)

	for _, schedule := range schedules {
		day := schedule.Day.Time.Format("2006-01-02")
		shift := Shift{
			ShiftID:           schedule.ShiftID,
			EmployeeID:        schedule.EmployeeID,
			EmployeeFirstName: schedule.EmployeeFirstName,
			EmployeeLastName:  schedule.EmployeeLastName,
			StartTime:         schedule.StartDatetime.Time,
			EndTime:           schedule.EndDatetime.Time,
			LocationID:        schedule.LocationID,
		}
		calendar[day] = append(calendar[day], shift)
	}

	var response []GetMonthlySchedulesByLocationResponse
	for date, shifts := range calendar {
		response = append(response, GetMonthlySchedulesByLocationResponse{
			Date:   date,
			Shifts: shifts,
		})
	}
	res := SuccessResponse(response, "Schedules retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetDailySchedulesByLocationApi retrieves the daily schedules for a specific location.
type GetDailySchedulesByLocationRequest struct {
	Year  int32 `form:"year" binding:"required"`
	Month int32 `form:"month" binding:"required"`
	Day   int32 `form:"day" binding:"required"`
}

// GetDailySchedulesByLocationResponse represents the response body for daily schedules.
type GetDailySchedulesByLocationResponse struct {
	Date   string  `json:"date"`
	Shifts []Shift `json:"shifts"`
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

	var req GetDailySchedulesByLocationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetDailySchedulesByLocationParams{
		Year:       req.Year,
		Month:      req.Month,
		Day:        req.Day,
		LocationID: locationID,
	}

	schedules, err := server.store.GetDailySchedulesByLocation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var shifts []Shift
	var targetDate string

	for _, schedule := range schedules {
		if targetDate == "" {
			targetDate = schedule.Day.Time.Format("2006-01-02")
		}
		shift := Shift{
			ShiftID:           schedule.ShiftID,
			EmployeeID:        schedule.EmployeeID,
			EmployeeFirstName: schedule.EmployeeFirstName,
			EmployeeLastName:  schedule.EmployeeLastName,
			StartTime:         schedule.StartDatetime.Time,
			EndTime:           schedule.EndDatetime.Time,
			LocationID:        schedule.LocationID,
		}
		shifts = append(shifts, shift)
	}

	// If no schedules found, still create the target date
	if targetDate == "" {
		targetDate = time.Date(int(req.Year), time.Month(req.Month), int(req.Day), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}

	response := GetDailySchedulesByLocationResponse{
		Date:   targetDate,
		Shifts: shifts,
	}

	res := SuccessResponse(response, "Daily schedules retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetScheduleByIdResponse represents the response body for retrieving a schedule by ID.
type GetScheduleByIdResponse struct {
	ID                int64     `json:"id"`
	EmployeeID        int64     `json:"employee_id"`
	EmployeeFirstName string    `json:"employee_first_name"`
	EmployeeLastName  string    `json:"employee_last_name"`
	LocationID        int64     `json:"location_id"`
	LocationName      string    `json:"location_name"`
	StartDatetime     time.Time `json:"start_datetime"`
	EndDatetime       time.Time `json:"end_datetime"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// @Summary Get schedule by ID
// @Description Get a schedule by its ID
// @Tags Schedule
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} Response[GetScheduleByIdResponse] "Schedule retrieved successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules/{id} [get]
func (server *Server) GetScheduleByIDApi(ctx *gin.Context) {
	scheduleID, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	schedule, err := server.store.GetScheduleById(ctx, scheduleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetScheduleByIdResponse{
		ID:                schedule.ID,
		EmployeeID:        schedule.EmployeeID,
		EmployeeFirstName: schedule.EmployeeFirstName,
		EmployeeLastName:  schedule.EmployeeLastName,
		LocationID:        schedule.LocationID,
		LocationName:      schedule.LocationName,
		StartDatetime:     schedule.StartDatetime.Time,
		EndDatetime:       schedule.EndDatetime.Time,
		CreatedAt:         schedule.CreatedAt.Time,
		UpdatedAt:         schedule.UpdatedAt.Time,
	}, "Schedule retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateScheduleRequest represents the request body for updating a schedule.
type UpdateScheduleRequest struct {
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
}

// UpdateScheduleResponse represents the response body after updating a schedule.
type UpdateScheduleResponse struct {
	ID            int64     `json:"id"`
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// @Summary Update a schedule
// @Description Update an existing schedule by its ID
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body UpdateScheduleRequest true "Update Schedule Request"
// @Success 200 {object} Response[UpdateScheduleResponse] "Schedule updated successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules/{id} [put]
func (server *Server) UpdateScheduleApi(ctx *gin.Context) {
	scheduleID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateScheduleParams{
		ID:            scheduleID,
		EmployeeID:    req.EmployeeID,
		LocationID:    req.LocationID,
		StartDatetime: pgtype.Timestamp{Time: req.StartDatetime, Valid: true},
		EndDatetime:   pgtype.Timestamp{Time: req.EndDatetime, Valid: true},
	}

	schedule, err := server.store.UpdateSchedule(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateScheduleResponse{
		ID:            schedule.ID,
		EmployeeID:    schedule.EmployeeID,
		LocationID:    schedule.LocationID,
		StartDatetime: schedule.StartDatetime.Time,
		EndDatetime:   schedule.EndDatetime.Time,
		CreatedAt:     schedule.CreatedAt.Time,
		UpdatedAt:     schedule.UpdatedAt.Time,
	}, "Schedule updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteScheduleApi deletes a schedule by its ID.
// @Summary Delete a schedule
// @Description Delete an existing schedule by its ID
// @Tags Schedule
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} Response[any] "Schedule deleted successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules/{id} [delete]
func (server *Server) DeleteScheduleApi(ctx *gin.Context) {
	scheduleID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteSchedule(ctx, scheduleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Schedule deleted successfully")
	ctx.JSON(http.StatusOK, res)
}
