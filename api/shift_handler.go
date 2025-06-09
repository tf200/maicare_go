package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateShiftApi creates a new shift for a specific location
type CreateShiftApiRequest struct {
	ShiftName string `json:"shift"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// CreateShiftApiResponse represents the response structure for creating a shift
type CreateShiftApiResponse struct {
	ID         int64  `json:"id"`
	LocationID int64  `json:"location_id"`
	ShiftName  string `json:"shift"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

// CreateShiftApi handles the creation of a new shift for a specific location
// @Summary Create a new shift
// @Description Create a new shift for a specific location
// @Tags Shifts
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param request body CreateShiftApiRequest true "Shift creation request"
// @Success 201 {object} Response[CreateShiftApiResponse]
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /locations/{id}/shifts [post]
func (server *Server) CreateShiftApi(ctx *gin.Context) {
	locationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateShiftApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	startTime, err := util.StringToPgTime(req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	endTime, err := util.StringToPgTime(req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shift, err := server.store.CreateShift(ctx, db.CreateShiftParams{
		LocationID: locationID,
		ShiftName:  req.ShiftName,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateShiftApiResponse{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		ShiftName:  shift.ShiftName,
		StartTime:  util.PgTimeToString(shift.StartTime), // Convert here
		EndTime:    util.PgTimeToString(shift.EndTime),   // Convert here
	}, "Shift Created Successfully")

	ctx.JSON(http.StatusCreated, res)
}

// UpdateShiftApiRequest represents the request structure for updating a shift
type UpdateShiftApiRequest struct {
	ShiftName string `json:"shift"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// UpdateShiftApiResponse represents the response structure for updating a shift
type UpdateShiftApiResponse struct {
	ID         int64  `json:"id"`
	LocationID int64  `json:"location_id"`
	ShiftName  string `json:"shift"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

// UpdateShiftApi handles the update of an existing shift
// @Summary Update an existing shift
// @Description Update an existing shift by ID
// @Tags Shifts
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param shift_id path int true "Shift ID"
// @Param request body UpdateShiftApiRequest true "Shift update request"
// @Success 200 {object} Response[UpdateShiftApiResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /locations/{id}/shifts/{shift_id} [put]
func (server *Server) UpdateShiftApi(ctx *gin.Context) {
	shiftID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateShiftApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	startTime, err := util.StringToPgTime(req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	endTime, err := util.StringToPgTime(req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shift, err := server.store.UpdateShift(ctx, db.UpdateShiftParams{
		ID:        shiftID,
		ShiftName: req.ShiftName,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(UpdateShiftApiResponse{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		ShiftName:  shift.ShiftName,
		StartTime:  util.PgTimeToString(shift.StartTime), // Convert here
		EndTime:    util.PgTimeToString(shift.EndTime),   // Convert here
	}, "Shift Updated Successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteShiftApi handles the deletion of a shift by ID
// @Summary Delete a shift by ID
// @Description Delete a shift by its ID
// @Tags Shifts
// @Produce json
// @Param shift_id path int true "Shift ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /locations/{id}/shifts/{shift_id} [delete]
func (server *Server) DeleteShiftApi(ctx *gin.Context) {
	shiftID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteShift(ctx, shiftID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Shift Deleted Successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListShiftsByLocationIDResponse represents the response structure for listing shifts by location ID
type ListShiftsByLocationIDResponse struct {
	ID         int64  `json:"id"`
	LocationID int64  `json:"location_id"`
	ShiftName  string `json:"shift"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

// ListShiftByLocationID handles the retrieval of shifts for a specific location
// @Summary List shifts by location ID
// @Description List all shifts for a specific location
// @Tags Shifts
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} Response[[]ListShiftsByLocationIDResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /locations/{id}/shifts [get]
func (server *Server) ListShiftByLocationID(ctx *gin.Context) {
	locationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shifts, err := server.store.GetShiftsByLocationID(ctx, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := make([]ListShiftsByLocationIDResponse, len(shifts))
	for i, shift := range shifts {
		response[i] = ListShiftsByLocationIDResponse{
			ID:         shift.ID,
			LocationID: shift.LocationID,
			ShiftName:  shift.ShiftName,
			StartTime:  util.PgTimeToString(shift.StartTime), // Convert here
			EndTime:    util.PgTimeToString(shift.EndTime),   // Convert here
		}
	}
	res := SuccessResponse(response, "Shifts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}
