package api

import (
	"database/sql"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateScheduleRequest represents the request body for creating a schedule.
type CreateScheduleRequest struct {
	EmployeeID int64 `json:"employee_id"`
	LocationID int64 `json:"location_id"`
	IsCustom   bool  `json:"is_custom" example:"true"` // true for custom schedule, false for preset shift

	// For custom schedules (required when is_custom = true)
	StartDatetime *time.Time `json:"start_datetime,omitempty" example:"2023-10-01T09:00:00Z"`
	EndDatetime   *time.Time `json:"end_datetime,omitempty" example:"2023-10-01T17:00:00Z"`

	// For preset shift-based schedules (required when is_custom = false)
	LocationShiftID *int64  `json:"location_shift_id,omitempty" example:"1"`
	ShiftDate       *string `json:"shift_date,omitempty" example:"2023-10-01"` // Date to apply the shift
}

// CreateScheduleResponse represents the response body after creating a schedule.
type CreateScheduleResponse struct {
	ID            int64     `json:"id"`
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Color         *string   `json:"color"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Additional info if created from preset shift
	LocationShiftID *int64  `json:"location_shift_id,omitempty"`
	ShiftName       *string `json:"shift_name,omitempty"`
}

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
	var req CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate request based on is_custom flag
	if req.IsCustom {
		// Custom schedule validation
		if req.StartDatetime == nil || req.EndDatetime == nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("start_datetime and end_datetime are required for custom schedules")))
			return
		}
		if req.LocationShiftID != nil || req.ShiftDate != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("location_shift_id and shift_date should not be provided for custom schedules")))
			return
		}
	} else {
		// Preset shift validation
		if req.LocationShiftID == nil || req.ShiftDate == nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("location_shift_id and shift_date are required for preset shift schedules")))
			return
		}
		if req.StartDatetime != nil || req.EndDatetime != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("start_datetime and end_datetime should not be provided for preset shift schedules")))
			return
		}
	}

	var startDatetime, endDatetime time.Time
	var locationShiftID *int64
	var shiftName *string

	if req.IsCustom {
		// Handle custom schedule
		if req.StartDatetime.After(*req.EndDatetime) {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("start_datetime must be before end_datetime")))
			return
		}
		startDatetime = *req.StartDatetime
		endDatetime = *req.EndDatetime
	} else {
		// Handle preset shift
		// First, get the location_shift details
		locationShift, err := server.store.GetShiftByID(ctx, *req.LocationShiftID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid location_shift_id: %v", err)))
			return
		}

		// Verify the shift belongs to the specified location
		if locationShift.LocationID != req.LocationID {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("location_shift_id does not belong to the specified location")))
			return
		}

		// Parse the shift date
		shiftDate, err := time.Parse("2006-01-02", *req.ShiftDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid shift_date format, expected YYYY-MM-DD: %v", err)))
			return
		}

		// Convert pgtype.Time (microseconds since midnight) to time components
		startHour, startMin, startSec, startNano := util.MicrosecondsToTimeComponents(locationShift.StartTime.Microseconds)
		endHour, endMin, endSec, endNano := util.MicrosecondsToTimeComponents(locationShift.EndTime.Microseconds)

		// Combine date with shift times to create full datetime
		startDatetime = time.Date(
			shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
			startHour, startMin, startSec, startNano,
			shiftDate.Location(),
		)

		endDatetime = time.Date(
			shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
			endHour, endMin, endSec, endNano,
			shiftDate.Location(),
		)

		// Handle shifts that cross midnight (end time is before start time)
		if locationShift.EndTime.Microseconds < locationShift.StartTime.Microseconds {
			endDatetime = endDatetime.AddDate(0, 0, 1)
		}

		locationShiftID = req.LocationShiftID
		shiftName = &locationShift.ShiftName
	}

	// Create the schedule
	arg := db.CreateScheduleParams{
		EmployeeID:      req.EmployeeID,
		LocationID:      req.LocationID,
		LocationShiftID: locationShiftID,
		StartDatetime:   pgtype.Timestamp{Time: startDatetime, Valid: true},
		EndDatetime:     pgtype.Timestamp{Time: endDatetime, Valid: true},
	}

	schedule, err := server.store.CreateSchedule(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		Color:           schedule.Color,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: locationShiftID,
		ShiftName:       shiftName,
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
	Color             *string   `json:"color"`                // Optional field for color coding
	ShiftName         *string   `json:"shift_name,omitempty"` // Optional field for shift name
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
			Color:             schedule.Color,
		}
		if schedule.ShiftName != nil {
			// If shift name is provided, add it to the shift
			shift.ShiftName = schedule.ShiftName
		} else {

			shift.ShiftName = util.StringPtr("Custom Shift")
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
			Color:             schedule.Color,
		}
		if schedule.ShiftName != nil {
			// If shift name is provided, add it to the shift
			shift.ShiftName = schedule.ShiftName
		} else {

			shift.ShiftName = util.StringPtr("Custom Shift")
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
	LocationShiftID   *int64    `json:"location_shift_id,omitempty"` // Optional field for preset shift
	LocationShiftName *string   `json:"shift_name,omitempty"`        // Optional field for shift name
	Color             *string   `json:"color"`                       // Optional field for color coding
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
		Color:             schedule.Color,
		LocationShiftID:   schedule.LocationShiftID,
		LocationShiftName: schedule.LocationShiftName,
	}, "Schedule retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateScheduleRequest represents the request body for updating a schedule.
type UpdateScheduleRequest struct {
	EmployeeID *int64 `json:"employee_id,omitempty"`
	LocationID *int64 `json:"location_id,omitempty"`
	IsCustom   *bool  `json:"is_custom,omitempty" example:"true"` // true for custom schedule, false for preset shift

	// For custom schedules (required when is_custom = true)
	StartDatetime *time.Time `json:"start_datetime,omitempty" example:"2023-10-01T09:00:00Z"`
	EndDatetime   *time.Time `json:"end_datetime,omitempty" example:"2023-10-01T17:00:00Z"`

	// For preset shift-based schedules (required when is_custom = false)
	LocationShiftID *int64  `json:"location_shift_id,omitempty" example:"1"`
	ShiftDate       *string `json:"shift_date,omitempty" example:"2023-10-01"` // Date to apply the shift

	Color *string `json:"color,omitempty" example:"#FF5733"`
}

// UpdateScheduleResponse represents the response body after updating a schedule.
type UpdateScheduleResponse struct {
	ID            int64     `json:"id"`
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Color         *string   `json:"color"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Additional info if updated from preset shift
	LocationShiftID *int64  `json:"location_shift_id,omitempty"`
	ShiftName       *string `json:"shift_name,omitempty"`
}

// @Summary Update an existing schedule
// @Description Update an existing schedule for an employee at a specific location. Supports both custom schedules and preset shifts.
// @Description Set is_custom=true and provide start_datetime/end_datetime for custom schedules
// @Description Set is_custom=false and provide location_shift_id/shift_date for preset shifts
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path int64 true "Schedule ID"
// @Param request body UpdateScheduleRequest true "Update Schedule Request"
// @Success 200 {object} Response[UpdateScheduleResponse] "Schedule updated successfully"
// @Failure 400 {object} Response[any] "Bad Request"
// @Failure 404 {object} Response[any] "Schedule not found"
// @Failure 500 {object} Response[any] "Internal Server Error"
// @Router /schedules/{id} [put]
func (server *Server) UpdateScheduleApi(ctx *gin.Context) {
	// Get schedule ID from URL parameter
	scheduleIDStr := ctx.Param("id")
	scheduleID, err := strconv.ParseInt(scheduleIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid schedule ID")))
		return
	}

	// Get existing schedule
	existingSchedule, err := server.store.GetScheduleById(ctx, scheduleID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("schedule not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Determine if this is a custom schedule update
	// is_custom must be explicitly provided to use start_datetime and end_datetime
	isCustom := req.IsCustom != nil && *req.IsCustom

	// If is_custom is not provided or false, ignore start_datetime and end_datetime
	if req.IsCustom == nil || !*req.IsCustom {
		req.StartDatetime = nil
		req.EndDatetime = nil

		// Determine schedule type based on other fields or existing schedule
		if req.IsCustom == nil {
			if req.LocationShiftID != nil || req.ShiftDate != nil {
				isCustom = false
			} else {
				// If no schedule type fields are provided, keep existing type
				// Check if existing schedule has location_shift_id to determine type
				isCustom = existingSchedule.LocationShiftID == nil
			}
		}
	}

	// Validate request based on schedule type
	if isCustom {
		// Custom schedule validation
		if req.LocationShiftID != nil || req.ShiftDate != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("location_shift_id and shift_date should not be provided for custom schedules")))
			return
		}
	} else {
		// For preset shift schedules, ignore start_datetime and end_datetime if provided
		// No validation error, just ignore these fields
		req.StartDatetime = nil
		req.EndDatetime = nil
	}

	// Prepare update parameters with existing values as defaults
	var startDatetime, endDatetime time.Time
	var locationShiftID *int64
	var shiftName *string
	var employeeID int64 = existingSchedule.EmployeeID
	var locationID int64 = existingSchedule.LocationID
	var color *string = existingSchedule.Color

	// Update fields if provided
	if req.EmployeeID != nil {
		employeeID = *req.EmployeeID
	}
	if req.LocationID != nil {
		locationID = *req.LocationID
	}
	if req.Color != nil {
		color = req.Color
	}

	if isCustom {
		// Handle custom schedule
		startDatetime = existingSchedule.StartDatetime.Time
		endDatetime = existingSchedule.EndDatetime.Time

		if req.StartDatetime != nil {
			startDatetime = *req.StartDatetime
		}
		if req.EndDatetime != nil {
			endDatetime = *req.EndDatetime
		}

		// Validate datetime order
		if startDatetime.After(endDatetime) {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("start_datetime must be before end_datetime")))
			return
		}

		locationShiftID = nil // Clear location_shift_id for custom schedules
	} else {
		// Handle preset shift
		var shiftIDToUse int64
		var shiftDateToUse string

		// Use existing values if not provided
		if req.LocationShiftID != nil {
			shiftIDToUse = *req.LocationShiftID
		} else if existingSchedule.LocationShiftID != nil {
			shiftIDToUse = *existingSchedule.LocationShiftID
		} else {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("location_shift_id is required for preset shift schedules")))
			return
		}

		if req.ShiftDate != nil {
			shiftDateToUse = *req.ShiftDate
		} else {
			// Extract date from existing start_datetime
			shiftDateToUse = existingSchedule.StartDatetime.Time.Format("2006-01-02")
		}

		// Get the location_shift details
		locationShift, err := server.store.GetShiftByID(ctx, shiftIDToUse)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid location_shift_id: %v", err)))
			return
		}

		// Verify the shift belongs to the specified location
		if locationShift.LocationID != locationID {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("location_shift_id does not belong to the specified location")))
			return
		}

		// Parse the shift date
		shiftDate, err := time.Parse("2006-01-02", shiftDateToUse)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid shift_date format, expected YYYY-MM-DD: %v", err)))
			return
		}

		// Convert pgtype.Time (microseconds since midnight) to time components
		startHour, startMin, startSec, startNano := util.MicrosecondsToTimeComponents(locationShift.StartTime.Microseconds)
		endHour, endMin, endSec, endNano := util.MicrosecondsToTimeComponents(locationShift.EndTime.Microseconds)

		// Combine date with shift times to create full datetime
		startDatetime = time.Date(
			shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
			startHour, startMin, startSec, startNano,
			shiftDate.Location(),
		)

		endDatetime = time.Date(
			shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
			endHour, endMin, endSec, endNano,
			shiftDate.Location(),
		)

		// Handle shifts that cross midnight (end time is before start time)
		if locationShift.EndTime.Microseconds < locationShift.StartTime.Microseconds {
			endDatetime = endDatetime.AddDate(0, 0, 1)
		}

		locationShiftID = &shiftIDToUse
		shiftName = &locationShift.ShiftName
	}

	// Update the schedule
	arg := db.UpdateScheduleParams{
		ID:              scheduleID,
		EmployeeID:      employeeID,
		LocationID:      locationID,
		LocationShiftID: locationShiftID,
		StartDatetime:   pgtype.Timestamp{Time: startDatetime, Valid: true},
		EndDatetime:     pgtype.Timestamp{Time: endDatetime, Valid: true},
		Color:           color,
	}

	schedule, err := server.store.UpdateSchedule(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		Color:           schedule.Color,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: locationShiftID,
		ShiftName:       shiftName,
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
