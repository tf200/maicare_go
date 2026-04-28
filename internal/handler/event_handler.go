package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"
	"maicare_go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterEventRoutes(
	rg *gin.RouterGroup,
	handler *EventHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	events := rg.Group("/events")
	{
		events.POST("", auth, requirePermission("APPOINTMENT.CREATE"), handler.CreateEvent)
		events.POST("/list", auth, requirePermission("APPOINTMENT.VIEW"), handler.ListEvents)
		events.GET("/:id", auth, requirePermission("APPOINTMENT.VIEW"), handler.GetEvent)
		events.PATCH("/:id", auth, requirePermission("APPOINTMENT.UPDATE"), handler.UpdateEvent)
		events.DELETE("/:id", auth, requirePermission("APPOINTMENT.DELETE"), handler.DeleteEvent)
		events.PUT("/:id/work_approval", auth, requirePermission("APPOINTMENT.WORK_APPROVAL.UPDATE"), handler.SetEventWorkApproval)
		events.POST("/work_approval_queue", auth, requirePermission("APPOINTMENT.WORK_APPROVAL.UPDATE"), handler.ListWorkApprovalQueue)
	}
}

type EventHandler struct {
	service domain.EventService
}

func NewEventHandler(service domain.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// CreateEvent creates a new calendar event.
// @Summary Create calendar event
// @Description Create appointment or reminder event (optionally recurring)
// @Tags events
// @Accept json
// @Produce json
// @Param request body domain.CreateEventRequest true "Create event request"
// @Success 201 {object} httpapi.Envelope[domain.EventResponse]
// @Failure 400 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[any]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /events [post]
func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	employeeID := middleware.EmployeeIDFromContext(ctx.Request.Context())
	if employeeID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", ""))
		return
	}

	var req domain.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	response, err := h.service.CreateEvent(ctx.Request.Context(), &req, employeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(response, "Event created successfully"))
}

// ListEvents lists calendar events in a date range.
// @Summary List events
// @Description List calendar events for current employee in a date range
// @Tags events
// @Accept json
// @Produce json
// @Param request body domain.ListEventsRequest true "List events request"
// @Success 200 {object} httpapi.Envelope[[]domain.EventOccurrenceResponse]
// @Failure 400 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[any]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /events/list [post]
func (h *EventHandler) ListEvents(ctx *gin.Context) {
	employeeID := middleware.EmployeeIDFromContext(ctx.Request.Context())
	if employeeID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", ""))
		return
	}

	var req domain.ListEventsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	response, err := h.service.ListEvents(ctx.Request.Context(), req, employeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(response, "Events listed successfully"))
}

// GetEvent retrieves a single event by ID.
// @Summary Get event by ID
// @Description Get event by ID
// @Tags events
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} httpapi.Envelope[domain.EventResponse]
// @Failure 400 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[any]
// @Failure 404 {object} httpapi.Envelope[any]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(ctx *gin.Context) {
	employeeID := middleware.EmployeeIDFromContext(ctx.Request.Context())
	if employeeID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", ""))
		return
	}

	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid event id", ""))
		return
	}

	response, err := h.service.GetEvent(ctx.Request.Context(), eventID, employeeID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(response, "Event retrieved successfully"))
}

// UpdateEvent updates an event with mutation scope.
// @Summary Update event
// @Description Update event with mutation scope (single/series/future)
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body domain.UpdateEventRequest true "Update event request"
// @Success 200 {object} httpapi.Envelope[domain.EventResponse]
// @Failure 400 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[any]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /events/{id} [patch]
func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	employeeID := middleware.EmployeeIDFromContext(ctx.Request.Context())
	if employeeID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", ""))
		return
	}

	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid event id", ""))
		return
	}

	var req domain.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	response, err := h.service.UpdateEvent(ctx.Request.Context(), eventID, &req, employeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(response, "Event updated successfully"))
}

// DeleteEvent cancels/deletes an event.
// @Summary Delete event
// @Description Cancel/delete event with mutation scope
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body domain.DeleteEventRequest true "Delete event request"
// @Success 200 {object} httpapi.Envelope[any]
// @Failure 400 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[any]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(ctx *gin.Context) {
	employeeID := middleware.EmployeeIDFromContext(ctx.Request.Context())
	if employeeID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", ""))
		return
	}

	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid event id", ""))
		return
	}

	var req domain.DeleteEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	err = h.service.DeleteEvent(ctx.Request.Context(), eventID, req, employeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Event deleted successfully"))
}

// SetEventWorkApproval approves/rejects work approval for an event.
// @Summary Set event work approval
// @Description Admin updates work approval status for an appointment
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body domain.SetEventWorkApprovalRequest true "Set work approval request"
// @Success 200 {object} httpapi.Envelope[any]
// @Failure 400 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[any]
// @Failure 403 {object} httpapi.Envelope[any]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /events/{id}/work_approval [put]
func (h *EventHandler) SetEventWorkApproval(ctx *gin.Context) {
	employeeID := middleware.EmployeeIDFromContext(ctx.Request.Context())
	if employeeID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", ""))
		return
	}

	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid event id", ""))
		return
	}

	var req domain.SetEventWorkApprovalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	// For handler validation, user ID is obtained from auth context
	userID := middleware.UserIDFromContext(ctx.Request.Context())

	err = h.service.SetEventWorkApproval(ctx.Request.Context(), eventID, &req, employeeID, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK[any](nil, "Work approval updated successfully"))
}

// ListWorkApprovalQueue lists the work approval queue.
// @Summary List work approval queue
// @Description Lists recent appointments across all employees that include at least one client attendee, with approval status
// @Tags events
// @Accept json
// @Produce json
// @Param request body domain.ListWorkApprovalQueueRequest true "List work approval queue request"
// @Success 200 {object} httpapi.Envelope[domain.ListWorkApprovalQueueResponse]
// @Failure 400 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[any]
// @Failure 403 {object} httpapi.Envelope[any]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /events/work_approval_queue [post]
func (h *EventHandler) ListWorkApprovalQueue(ctx *gin.Context) {
	employeeID := middleware.EmployeeIDFromContext(ctx.Request.Context())
	if employeeID == uuid.Nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", ""))
		return
	}

	var req domain.ListWorkApprovalQueueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	resp, err := h.service.ListWorkApprovalQueue(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(resp, "Work approval queue listed successfully"))
}


