package api

import (
	"fmt"
	"net/http"

	"maicare_go/service/appointment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create calendar event
// @Description Create appointment or reminder event (optionally recurring)
// @Tags events
// @Accept json
// @Produce json
// @Param request body appointment.CreateEventRequest true "Create event request"
// @Success 201 {object} Response[appointment.EventResponse]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /events [post]
func (server *Server) CreateEventApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	var req appointment.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.AppointmentService.CreateEvent(ctx, &req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(response, "Event created successfully"))
}

// @Summary List events
// @Description List calendar events for current employee in a date range
// @Tags events
// @Accept json
// @Produce json
// @Param request body appointment.ListEventsRequest true "List events request"
// @Success 200 {object} Response[[]appointment.EventOccurrenceResponse]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /events/list [post]
func (server *Server) ListEventsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	var req appointment.ListEventsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	// Check if filtering by a different employee
	if req.EmployeeID != nil && *req.EmployeeID != payload.EmployeeID {
		// Require APPOINTMENT.VIEW_ALL permission to view other employees' events
		hasPermission, err := server.businessService.AuthService.HasPermission(ctx, payload.UserId, "APPOINTMENT.VIEW_ALL")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to check permissions")))
			return
		}
		if !hasPermission {
			ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("insufficient permissions to view other employees' events")))
			return
		}
	}

	response, err := server.businessService.AppointmentService.ListEvents(ctx, req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Events listed successfully"))
}

// @Summary Get event by ID
// @Description Get event by ID
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} Response[appointment.EventResponse]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /events/{id} [get]
func (server *Server) GetEventApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid event id")))
		return
	}

	response, err := server.businessService.AppointmentService.GetEvent(ctx, eventID, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Event retrieved successfully"))
}

// @Summary Update event
// @Description Update event with mutation scope (single/series/future)
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body appointment.UpdateEventRequest true "Update event request"
// @Success 200 {object} Response[appointment.EventResponse]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /events/{id} [patch]
func (server *Server) UpdateEventApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid event id")))
		return
	}

	var req appointment.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.AppointmentService.UpdateEvent(ctx, eventID, &req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Event updated successfully"))
}

// @Summary Delete event
// @Description Cancel/delete event with mutation scope
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body appointment.DeleteEventRequest true "Delete event request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /events/{id} [delete]
func (server *Server) DeleteEventApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid event id")))
		return
	}

	var req appointment.DeleteEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	err = server.businessService.AppointmentService.DeleteEvent(ctx, eventID, req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse[any](nil, "Event deleted successfully"))
}
