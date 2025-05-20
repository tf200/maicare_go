package api

import (
	db "maicare_go/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateScheduleRequest struct {
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
}

type CreateScheduleResponse struct {
	ID            int32     `json:"id"`
	EmployeeID    int64     `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (server *Server) CreateScheduleApi(ctx *gin.Context) {
	var req CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

type GetMonthlySchedulesByLocationRequest struct {
	LocationID int64 `form:"location_id"`
	Year       int32 `form:"year"`
	Month      int32 `form:"month"`
}

type Shift struct {
	ShiftID      int32     `json:"shift_id"`
	EmployeeID   int64     `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	LocationID   int64     `json:"location_id"`
}

type GetMonthlySchedulesByLocationResponse struct {
	Date   string  `json:"date"`
	Shifts []Shift `json:"shifts"`
}

func (server *Server) GetMonthlySchedulesByLocationApi(ctx *gin.Context) {
	var req GetMonthlySchedulesByLocationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetMonthlySchedulesByLocationParams{
		Year:       req.Year,
		Month:      req.Month,
		LocationID: req.LocationID,
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
			ShiftID:      schedule.ShiftID,
			EmployeeID:   schedule.EmployeeID,
			EmployeeName: schedule.EmployeeFirstName,
			StartTime:    schedule.StartDatetime.Time,
			EndTime:      schedule.EndDatetime.Time,
			LocationID:   schedule.LocationID,
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
