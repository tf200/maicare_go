package api

import (
	"errors"
	"net/http"

	latearrival "maicare_go/service/late_arrival"

	"github.com/gin-gonic/gin"
)

// CreateLateArrivalApi creates a late-arrival record for the authenticated employee.
// @Summary Create late arrival
// @Description Creates a late-arrival record linked to the employee assigned shift for that date.
// @Tags late-arrivals
// @Accept json
// @Produce json
// @Param request body late_arrival.CreateLateArrivalRequest true "Late arrival payload"
// @Success 201 {object} Response[late_arrival.CreateLateArrivalResponse]
// @Router /late-arrivals [post]
func (server *Server) CreateLateArrivalApi(ctx *gin.Context) {
	var req latearrival.CreateLateArrivalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LateArrivalService.CreateLateArrival(ctx, payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(mapLateArrivalErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Late arrival created successfully"))
}

// CreateLateArrivalByAdminApi creates a late-arrival record for any employee.
// @Summary Admin create late arrival
// @Description Creates a late-arrival record for a selected employee linked to the assigned shift.
// @Tags late-arrivals
// @Accept json
// @Produce json
// @Param request body late_arrival.CreateLateArrivalByAdminRequest true "Admin late arrival payload"
// @Success 201 {object} Response[late_arrival.CreateLateArrivalResponse]
// @Router /late-arrivals/admin [post]
func (server *Server) CreateLateArrivalByAdminApi(ctx *gin.Context) {
	var req latearrival.CreateLateArrivalByAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LateArrivalService.CreateLateArrivalByAdmin(ctx, payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(mapLateArrivalErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Late arrival created successfully"))
}

// ListMyLateArrivalsApi lists late-arrival records for the authenticated employee.
// @Summary List my late arrivals
// @Description Lists paginated late-arrival records for the authenticated employee.
// @Tags late-arrivals
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param date_from query string false "Start date filter (YYYY-MM-DD)"
// @Param date_to query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} Response[pagination.Response[late_arrival.LateArrivalListItem]]
// @Router /late-arrivals/my [get]
func (server *Server) ListMyLateArrivalsApi(ctx *gin.Context) {
	var req latearrival.ListMyLateArrivalsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	res, err := server.businessService.LateArrivalService.ListMyLateArrivals(ctx, payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(mapLateArrivalErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Late arrivals retrieved successfully"))
}

// ListLateArrivalsApi lists late-arrival records for admins with filters.
// @Summary List late arrivals
// @Description Lists paginated late-arrival records with optional employee and date filters.
// @Tags late-arrivals
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param employee_search query string false "Employee name search (first/last name)"
// @Param date_from query string false "Start date filter (YYYY-MM-DD)"
// @Param date_to query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} Response[pagination.Response[late_arrival.LateArrivalListItem]]
// @Router /late-arrivals [get]
func (server *Server) ListLateArrivalsApi(ctx *gin.Context) {
	var req latearrival.ListLateArrivalsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := server.businessService.LateArrivalService.ListLateArrivals(ctx, &req)
	if err != nil {
		ctx.JSON(mapLateArrivalErrorStatus(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Late arrivals retrieved successfully"))
}

func mapLateArrivalErrorStatus(err error) int {
	switch {
	case errors.Is(err, latearrival.ErrLateArrivalInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, latearrival.ErrLateArrivalConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
