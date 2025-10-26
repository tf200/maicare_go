package api

import (
	"maicare_go/service/organization"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	var req organization.CreateShiftApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shift, err := server.businessService.OrganizationService.CreateShift(ctx, &req, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(shift, "Shift Created Successfully")

	ctx.JSON(http.StatusCreated, res)
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

	var req organization.UpdateShiftApiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shift, err := server.businessService.OrganizationService.UpdateShift(ctx, shiftID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(shift, "Shift Updated Successfully")

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

	err = server.businessService.OrganizationService.DeleteShift(ctx, shiftID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Shift Deleted Successfully")
	ctx.JSON(http.StatusOK, res)
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

	shifts, err := server.businessService.OrganizationService.ListShiftsByLocationID(ctx, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(shifts, "Shifts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}
