package api

import (
	"fmt"
	"net/http"

	db "maicare_go/db/sqlc"
	"maicare_go/service/appointment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

	// Approved appointments are immutable for non-admins.
	ev, err := server.businessService.Store.GetCalendarEventByID(ctx, eventID)
	if err == nil {
		approved := ev.WorkApprovalStatus == db.CalendarEventWorkApprovalStatusEnumApproved

		// For recurring "single" updates, also block editing an already-approved occurrence override.
		if !approved && req.Scope == appointment.MutationScopeSingle && req.RecurrenceID != nil && ev.RecurringEventID == nil {
			if occ, err := server.businessService.Store.GetCalendarEventOverrideByMasterAndRecurrence(ctx, db.GetCalendarEventOverrideByMasterAndRecurrenceParams{
				MasterEventID: &ev.ID,
				RecurrenceID:  pgtype.Timestamptz{Time: req.RecurrenceID.UTC(), Valid: true},
			}); err == nil {
				approved = occ.WorkApprovalStatus == db.CalendarEventWorkApprovalStatusEnumApproved
			}
		}

		if approved {
			hasPermission, err := server.businessService.AuthService.HasPermission(ctx, payload.UserId, "APPOINTMENT.WORK_APPROVAL.UPDATE")
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to check permissions")))
				return
			}
			if !hasPermission {
				ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("approved appointments can only be edited by admin")))
				return
			}
		}
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

	// Approved appointments are immutable for non-admins.
	ev, err := server.businessService.Store.GetCalendarEventByID(ctx, eventID)
	if err == nil {
		approved := ev.WorkApprovalStatus == db.CalendarEventWorkApprovalStatusEnumApproved
		if !approved && req.Scope == appointment.MutationScopeSingle && req.RecurrenceID != nil && ev.RecurringEventID == nil {
			if occ, err := server.businessService.Store.GetCalendarEventOverrideByMasterAndRecurrence(ctx, db.GetCalendarEventOverrideByMasterAndRecurrenceParams{
				MasterEventID: &ev.ID,
				RecurrenceID:  pgtype.Timestamptz{Time: req.RecurrenceID.UTC(), Valid: true},
			}); err == nil {
				approved = occ.WorkApprovalStatus == db.CalendarEventWorkApprovalStatusEnumApproved
			}
		}

		if approved {
			hasPermission, err := server.businessService.AuthService.HasPermission(ctx, payload.UserId, "APPOINTMENT.WORK_APPROVAL.UPDATE")
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to check permissions")))
				return
			}
			if !hasPermission {
				ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("approved appointments can only be deleted by admin")))
				return
			}
		}
	}

	err = server.businessService.AppointmentService.DeleteEvent(ctx, eventID, req, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse[any](nil, "Event deleted successfully"))
}

// @Summary Approve/reject appointment worked hours
// @Description Admin updates work approval status for an appointment (and optionally a specific recurrence occurrence)
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body appointment.SetEventWorkApprovalRequest true "Set work approval request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 403 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /events/{id}/work_approval [put]
func (server *Server) SetEventWorkApprovalApi(ctx *gin.Context) {
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

	var req appointment.SetEventWorkApprovalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	if req.Status == "rejected" {
		if req.RejectionReason == nil || len(*req.RejectionReason) == 0 {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("rejection_reason is required when status is rejected")))
			return
		}
	}

	err = server.businessService.AppointmentService.SetEventWorkApproval(ctx, eventID, &req, payload.EmployeeID, payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse[any](nil, "Work approval updated successfully"))
}

// @Summary List work approval queue items (admin)
// @Description Lists recent appointments across all employees that include at least one client attendee, with approval status.
// @Tags events
// @Accept json
// @Produce json
// @Param request body appointment.ListWorkApprovalQueueRequest true "List work approval queue request"
// @Success 200 {object} Response[appointment.ListWorkApprovalQueueResponse]
// @Failure 400 {object} Response[any]
// @Failure 401 {object} Response[any]
// @Failure 403 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /events/work_approval_queue [post]
func (server *Server) ListWorkApprovalQueueApi(ctx *gin.Context) {
	_, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	var req appointment.ListWorkApprovalQueueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	resp, err := server.businessService.AppointmentService.ListWorkApprovalQueue(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(resp, "Work approval queue listed successfully"))
}
