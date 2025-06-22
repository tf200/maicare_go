package api

import (
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// WorkingHourItem represents a single working hour item, which can be either a schedule or an appointment.
type WorkingHourItem struct {
	ID            int64     `json:"id"`
	Type          string    `json:"type"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	DurationHours float64   `json:"duration_hours"`
	Location      string    `json:"location"`
	LocationID    *int64    `json:"location_id,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Status        *string   `json:"status,omitempty"`
	Color         string    `json:"color"`
}

// Summary contains the summary of working hours for an employee in a given period.
type Summary struct {
	TotalHours       float64 `json:"total_hours"`
	AppointmentHours float64 `json:"appointment_hours"`
	ShiftHours       float64 `json:"shift_hours"`
	TotalDaysWorked  int     `json:"total_days_worked"`
}

// Period information
type Period struct {
	Year           int32  `json:"year"`
	Month          int32  `json:"month"`
	MonthName      string `json:"month_name"`
	IsCurrentMonth bool   `json:"is_current_month"`
	DateRange      struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"date_range"`
}

// ListWorkingHoursRequest represents the request parameters for listing working hours.
type ListWorkingHoursRequest struct {
	Year  int32 `form:"year" binding:"required"`
	Month int32 `form:"month" binding:"required"`
}

// ListWorkingHoursResponse represents the response structure for listing working hours.
type ListWorkingHoursResponse struct {
	EmployeeID   int64             `json:"employee_id"`
	Period       Period            `json:"period"`
	Summary      Summary           `json:"summary"`
	WorkingHours []WorkingHourItem `json:"working_hours"`
}

// @Summary List working hours for an employee
// @Description List working hours for an employee in a given month and year
// @Tags Working Hours
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} Response[ListWorkingHoursResponse] "
// @Failure 400 {object} Response[any] "Invalid request parameters"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /employees/{id}/working_hours [get]
func (server *Server) ListWorkingHours(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	var req ListWorkingHoursRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	periodStart := time.Date(int(req.Year), time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)

	var periodEnd time.Time

	if req.Year == int32(time.Now().Year()) && req.Month == int32(time.Now().Month()) {
		periodEnd = time.Now().AddDate(0, 0, 1)
	} else {
		periodEnd = periodStart.AddDate(0, 1, 0)
	}

	employeeAppointments, err := server.store.ListEmployeeAppointmentsInRange(ctx,
		db.ListEmployeeAppointmentsInRangeParams{
			StartDate: pgtype.Timestamp{
				Time:  periodStart,
				Valid: true,
			},
			EndDate: pgtype.Timestamp{
				Time:  periodEnd,
				Valid: true,
			},
			EmployeeID: &employeeID,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	employeeSchedules, err := server.store.GetEmployeeSchedules(
		ctx,
		db.GetEmployeeSchedulesParams{
			PeriodStart: pgtype.Timestamp{
				Time:  periodStart,
				Valid: true,
			},
			PeriodEnd: pgtype.Timestamp{
				Time:  periodEnd,
				Valid: true,
			},
			EmployeeID: employeeID,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	workingHours := make([]WorkingHourItem, len(employeeSchedules)+len(employeeAppointments))
	var appointmentHours, shiftHours float64
	uniqueDays := make(map[string]bool)

	for i, schedule := range employeeSchedules {
		duration := schedule.EndDatetime.Time.Sub(schedule.StartDatetime.Time).Hours()

		shiftHours += duration
		dayKey := schedule.StartDatetime.Time.Format("2006-01-02")
		uniqueDays[dayKey] = true

		workingHours[i] = WorkingHourItem{
			ID:            schedule.ID,
			Type:          "schedule",
			StartTime:     schedule.StartDatetime.Time,
			EndTime:       schedule.EndDatetime.Time,
			DurationHours: duration,
			Location:      schedule.LocationName,
			LocationID:    &schedule.LocationID,
			Description:   nil,
			Status:        nil,
		}
	}

	for i, appointment := range employeeAppointments {
		duration := appointment.EndTime.Time.Sub(appointment.StartTime.Time).Hours()

		appointmentHours += duration
		dayKey := appointment.StartTime.Time.Format("2006-01-02")
		uniqueDays[dayKey] = true

		workingHours[len(employeeSchedules)+i] = WorkingHourItem{
			ID:            appointment.AppointmentID,
			Type:          "appointment",
			StartTime:     appointment.StartTime.Time,
			EndTime:       appointment.EndTime.Time,
			DurationHours: duration,
			Location:      *appointment.Location,
			LocationID:    nil,
			Description:   appointment.Description,
			Status:        &appointment.Status,
		}
	}

	totalHours := appointmentHours + shiftHours

	summary := Summary{
		TotalHours:       totalHours,
		AppointmentHours: appointmentHours,
		ShiftHours:       shiftHours,
		TotalDaysWorked:  len(uniqueDays),
	}

	period := Period{
		Year:           req.Year,
		Month:          req.Month,
		MonthName:      periodStart.Month().String(),
		IsCurrentMonth: req.Year == int32(time.Now().Year()) && req.Month == int32(time.Now().Month()),
		DateRange: struct {
			Start string `json:"start"`
			End   string `json:"end"`
		}{
			Start: periodStart.Format("2006-01-02"),
			End:   periodEnd.Format("2006-01-02"),
		},
	}

	res := SuccessResponse(ListWorkingHoursResponse{
		EmployeeID:   employeeID,
		Period:       period,
		Summary:      summary,
		WorkingHours: workingHours,
	}, "List working hours successfully")
	ctx.JSON(http.StatusOK, res)

}
