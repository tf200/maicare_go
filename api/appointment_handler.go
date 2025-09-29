package api

import (
	"fmt"
	"maicare_go/service/appointment"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary Create an appointment
// @Description Create a new appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param request body appointment.CreateAppointmentRequest true "Create appointment request"
// @Success 200 {object} Response[appointment.CreateAppointmentResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - User not found"
// @Failure 409 {object} Response[any] "Conflict - Appointment already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments [post]
func (server *Server) CreateAppointmentApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	userID := payload.UserId

	var req appointment.CreateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	response, err := server.businessService.AppointmentService.CreateAppointment(&req, userID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create appointment")))
		return
	}

	res := SuccessResponse(response, "Appointment created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// AddParticipantToAppointment adds participants to an existing appointment
// @Summary Add participants to an appointment
// @Description Add participants to an existing appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Param request body appointment.AddParticipantToAppointmentRequest true "Add participant request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Participant already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{appointment_id}/participants [post]
func (server *Server) AddParticipantToAppointmentApi(ctx *gin.Context) {
	appointmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid appointment ID parameter")))
		return
	}

	var req appointment.AddParticipantToAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request payload")))
		return
	}

	err = server.businessService.AppointmentService.AddParticipantToAppointment(ctx, appointmentID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to add participants to appointment")))
		return
	}

	res := SuccessResponse[any](nil, "Participants added successfully")
	ctx.JSON(http.StatusOK, res)
}

// AddClientToAppointment adds clients to an existing appointment
// @Summary Add clients to an appointment
// @Description Add clients to an existing appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Param request body appointment.AddClientToAppointmentRequest true "Add client request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Client already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{appointment_id}/clients [post]
func (server *Server) AddClientToAppointmentApi(ctx *gin.Context) {
	appointmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req appointment.AddClientToAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.AppointmentService.AddClientToAppointment(ctx, appointmentID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Clients added successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListAppointmentsForEmployeeInRange lists appointments for an employee in a date range
// @Summary List appointments for an employee in a date range
// @Description List appointments for an employee in a date range
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param request body appointment.ListAppointmentsForEmployeeInRangeRequest true "List appointments request"
// @Success 200 {object} Response[appointment.ListAppointmentsForEmployeeInRangeResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Employee not found"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /employees/{id}/appointments [post]
func (server *Server) ListAppointmentsForEmployee(ctx *gin.Context) {
	employeeID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req appointment.ListAppointmentsForEmployeeInRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentList, err := server.businessService.AppointmentService.ListAppointmentsForEmployeeInRange(ctx, employeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(appointmentList, "Appointments retrieved successfully")
	ctx.JSON(http.StatusOK, res)

}

// @Summary List appointments for a client in a date range
// @Description List appointments for a client in a date range
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body appointment.ListAppointmentsForClientRequest true "List appointments request"
// @Success 200 {object} Response[appointment.ListAppointmentsForClientResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Client not found"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/appointments [post]
func (server *Server) ListAppointmentsForClientApi(ctx *gin.Context) {
	clientID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req appointment.ListAppointmentsForClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	appointmentList, err := server.businessService.AppointmentService.ListAppointmentsForClientInRange(ctx, clientID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(appointmentList, "Appointments retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetAppointmentResponse represents the response payload for getting an appointment
type ParticipantsDetails struct {
	EmployeeID int64  `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

// ClientsDetails represents the details of a client
type ClientsDetails struct {
	ClientID  int64  `json:"client_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetAppointmentApi retrieves an appointment by ID
// @Summary Get an appointment by ID
// @Description Get an appointment by ID
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Success 200 {object} Response[appointment.GetAppointmentResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{id} [get]
func (server *Server) GetAppointmentApi(ctx *gin.Context) {
	appointmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentData, err := server.businessService.AppointmentService.GetAppointment(ctx, appointmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(appointmentData, "Appointment retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateAppointmentApi updates an appointment
// @Summary Update an appointment
// @Description Update an appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Param request body appointment.UpdateAppointmentRequest true "Update appointment request"
// @Success 200 {object} Response[appointment.UpdateAppointmentResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Appointment already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{id} [put]
func (server *Server) UpdateAppointmentApi(ctx *gin.Context) {
	appointmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid appointment ID")))
		return
	}

	var req appointment.UpdateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.AppointmentService.UpdateAppointment(ctx, appointmentID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update appointment")))
		return
	}

	res := SuccessResponse(response, "Appointment updated successfully")
	ctx.JSON(http.StatusOK, res)

}

// DeleteAppointmentApi deletes an appointment
// @Summary Delete an appointment
// @Description Delete an appointment
// @Tags appointments
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Appointment already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{id} [delete]
func (server *Server) DeleteAppointmentApi(ctx *gin.Context) {
	appointmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid appointment ID")))
		return
	}

	err = server.businessService.AppointmentService.DeleteAppointment(ctx, appointmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete appointment")))
		return
	}

	res := SuccessResponse[any](nil, "Appointment deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// ConfirmAppointmentApi confirms an appointment
// @Summary Confirm an appointment
// @Description Confirm an appointment
// @Tags appointments
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Appointment already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{id}/confirm [post]
func (server *Server) ConfirmAppointmentApi(ctx *gin.Context) {
	appointmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ConfirmAppointmentApi", "Invalid appointment ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid appointment ID")))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("failed to get auth payload")))
		return
	}

	err = server.businessService.AppointmentService.ConfirmAppointment(ctx, appointmentID, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Appointment confirmed successfully")
	ctx.JSON(http.StatusOK, res)
}
