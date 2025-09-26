package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/service/appointment"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
// @Param request body AddParticipantToAppointmentRequest true "Add participant request"
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
// @Param request body AddClientToAppointmentRequest true "Add client request"
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
// @Param request body ListAppointmentsForEmployeeInRangeRequest true "List appointments request"
// @Success 200 {object} Response[ListAppointmentsForEmployeeInRangeResponse]
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

// ListAppointmentsForClientRequest represents the request payload for listing appointments for a client in a date range
type ListAppointmentsForClientRequest struct {
	StartDate time.Time `json:"start_date" binding:"required" example:"2025-04-27T00:00:00Z"`
	EndDate   time.Time `json:"end_date" binding:"required" example:"2025-04-30T23:59:59Z"`
}

// ListAppointmentsForClientResponse represents the response payload for listing appointments for a client in a date range
type ListAppointmentsForClientResponse struct {
	ID                    uuid.UUID             `json:"id"`
	CreatorEmployeeID     *int64                `json:"creator_employee_id"`
	StartTime             time.Time             `json:"start_time"`
	EndTime               time.Time             `json:"end_time"`
	Location              *string               `json:"location"`
	Description           *string               `json:"description"`
	Color                 *string               `json:"color"`
	Status                string                `json:"status"`
	RecurrenceType        *string               `json:"recurrence_type"`
	RecurrenceInterval    *int32                `json:"recurrence_interval"`
	RecurrenceEndDate     pgtype.Date           `json:"recurrence_end_date"`
	ConfirmedByEmployeeID *int32                `json:"confirmed_by_employee_id"`
	ConfirmedAt           time.Time             `json:"confirmed_at"`
	CreatedAt             time.Time             `json:"created_at"`
	ParticipantsDetails   []ParticipantsDetails `json:"participants_details"`
	ClientsDetails        []ClientsDetails      `json:"clients_details"`
}

// @Summary List appointments for a client in a date range
// @Description List appointments for a client in a date range
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body ListAppointmentsForClientRequest true "List appointments request"
// @Success 200 {object} Response[ListAppointmentsForClientResponse]
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

	var req ListAppointmentsForClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.StartDate.After(req.EndDate) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("start date must be before end date")))
		return
	}

	arg := db.ListClientAppointmentsInRangeParams{
		ClientID:  clientID,
		StartDate: pgtype.Timestamp{Time: req.StartDate, Valid: true},
		EndDate:   pgtype.Timestamp{Time: req.EndDate, Valid: true},
	}

	appointments, err := server.store.ListClientAppointmentsInRange(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(appointments) == 0 {
		res := SuccessResponse([]ListAppointmentsForClientResponse{}, "No appointments found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	appointmentList := make([]ListAppointmentsForClientResponse, len(appointments))
	for i, appointment := range appointments {
		participants, err := server.store.GetAppointmentParticipants(ctx, appointment.AppointmentID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		clientDetails, err := server.store.GetAppointmentClients(ctx, appointment.AppointmentID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		participantsDetails := make([]ParticipantsDetails, len(participants))
		for j, participant := range participants {
			participantsDetails[j] = ParticipantsDetails{
				EmployeeID: participant.EmployeeID,
				FirstName:  participant.FirstName,
				LastName:   participant.LastName,
			}
		}
		clientsDetails := make([]ClientsDetails, len(clientDetails))
		for j, client := range clientDetails {
			clientsDetails[j] = ClientsDetails{
				ClientID:  client.ClientID,
				FirstName: client.FirstName,
				LastName:  client.LastName,
			}
		}

		appointmentList[i] = ListAppointmentsForClientResponse{
			ID:                  appointment.AppointmentID,
			CreatorEmployeeID:   appointment.CreatorEmployeeID,
			StartTime:           appointment.StartTime.Time,
			EndTime:             appointment.EndTime.Time,
			Location:            appointment.Location,
			Description:         appointment.Description,
			Color:               appointment.Color,
			Status:              appointment.Status,
			CreatedAt:           appointment.CreatedAt.Time,
			ParticipantsDetails: participantsDetails,
			ClientsDetails:      clientsDetails,
		}
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

// GetAppointmentResponse represents the response payload for getting an appointment
type GetAppointmentResponse struct {
	ID                     uuid.UUID             `json:"id"`
	AppointmentTemplatesID *int64                `json:"appointment_templates_id"`
	CreatorEmployeeID      *int64                `json:"creator_employee_id"`
	CreatorFirstName       *string               `json:"creator_first_name"`
	CreatorLastName        *string               `json:"creator_last_name"`
	StartTime              time.Time             `json:"start_time"`
	EndTime                time.Time             `json:"end_time"`
	Location               *string               `json:"location"`
	Description            *string               `json:"description"`
	Color                  *string               `json:"color"`
	Status                 string                `json:"status"`
	IsConfirmed            bool                  `json:"is_confirmed"`
	ConfirmedByEmployeeID  *int64                `json:"confirmed_by_employee_id"`
	ConfirmerFirstName     *string               `json:"confirmer_first_name"`
	ConfirmerLastName      *string               `json:"confirmer_last_name"`
	ConfirmedAt            time.Time             `json:"confirmed_at"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	ParticipantsDetails    []ParticipantsDetails `json:"participants_details"`
	ClientsDetails         []ClientsDetails      `json:"clients_details"`
}

// GetAppointmentApi retrieves an appointment by ID
// @Summary Get an appointment by ID
// @Description Get an appointment by ID
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Success 200 {object} Response[GetAppointmentResponse]
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

	appointment, err := server.store.GetScheduledAppointmentByID(ctx, appointmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	participants, err := server.store.GetAppointmentParticipants(ctx, appointment.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	clientDetails, err := server.store.GetAppointmentClients(ctx, appointment.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	participantsDetails := make([]ParticipantsDetails, len(participants))
	for i, participant := range participants {
		participantsDetails[i] = ParticipantsDetails{
			EmployeeID: participant.EmployeeID,
			FirstName:  participant.FirstName,
			LastName:   participant.LastName,
		}
	}

	clientsDetails := make([]ClientsDetails, len(clientDetails))
	for i, client := range clientDetails {
		clientsDetails[i] = ClientsDetails{
			ClientID:  client.ClientID,
			FirstName: client.FirstName,
			LastName:  client.LastName,
		}
	}

	response := GetAppointmentResponse{
		ID:                    appointment.ID,
		CreatorEmployeeID:     appointment.CreatorEmployeeID,
		CreatorFirstName:      appointment.CreatorFirstName,
		CreatorLastName:       appointment.CreatorLastName,
		StartTime:             appointment.StartTime.Time,
		EndTime:               appointment.EndTime.Time,
		Location:              appointment.Location,
		Description:           appointment.Description,
		Color:                 appointment.Color,
		Status:                appointment.Status,
		IsConfirmed:           appointment.IsConfirmed,
		ConfirmedByEmployeeID: appointment.ConfirmedByEmployeeID,
		ConfirmerFirstName:    appointment.ConfirmerFirstName,
		ConfirmerLastName:     appointment.ConfirmerLastName,
		ParticipantsDetails:   participantsDetails,
		ClientsDetails:        clientsDetails,
	}

	res := SuccessResponse(response, "Appointment retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// ConfirmAppointmentRequest represents the request payload for confirming an appointment
type UpdateAppointmentRequest struct {
	StartTime              time.Time `json:"start_time" binding:"required" example:"2023-10-01T10:00:00Z"`
	EndTime                time.Time `json:"end_time" binding:"required"`
	Location               *string   `json:"location"`
	Color                  *string   `json:"color" example:"#FF5733"`
	Description            *string   `json:"description"`
	ClientIDs              *[]int64  `json:"client_ids"`
	ParticipantEmployeeIDs *[]int64  `json:"participant_employee_ids"`
}

// UpdateAppointmentResponse represents the response payload for updating an appointment
type UpdateAppointmentResponse struct {
	ID                     uuid.UUID        `json:"id"`
	AppointmentTemplatesID *uuid.UUID       `json:"appointment_templates_id"`
	CreatorEmployeeID      *int64           `json:"creator_employee_id"`
	StartTime              pgtype.Timestamp `json:"start_time"`
	EndTime                pgtype.Timestamp `json:"end_time"`
	Location               *string          `json:"location"`
	Description            *string          `json:"description"`
	Color                  *string          `json:"color"`
	Status                 string           `json:"status"`
	IsConfirmed            bool             `json:"is_confirmed"`
	ConfirmedByEmployeeID  *int64           `json:"confirmed_by_employee_id"`
	ConfirmedAt            pgtype.Timestamp `json:"confirmed_at"`
	CreatedAt              pgtype.Timestamp `json:"created_at"`
	UpdatedAt              pgtype.Timestamp `json:"updated_at"`
}

// UpdateAppointmentApi updates an appointment
// @Summary Update an appointment
// @Description Update an appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID (UUID)"
// @Param request body UpdateAppointmentRequest true "Update appointment request"
// @Success 200 {object} Response[UpdateAppointmentResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Appointment already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{id} [put]
func (server *Server) UpdateAppointmentApi(ctx *gin.Context) {
	appointmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Invalid appointment ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid appointment ID")))
		return
	}

	var req UpdateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	if req.StartTime.After(req.EndTime) {
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Start time must be before end time", zap.Time("start_time", req.StartTime), zap.Time("end_time", req.EndTime))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("start time must be before end time")))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to begin transaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to begin transaction")))
		return
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to rollback appointment update", zap.Error(rollbackErr))
		}
	}()
	qtx := server.store.WithTx(tx)

	appointment, err := qtx.UpdateAppointment(ctx, db.UpdateAppointmentParams{
		ID:          appointmentID,
		StartTime:   pgtype.Timestamp{Time: req.StartTime, Valid: true},
		EndTime:     pgtype.Timestamp{Time: req.EndTime, Valid: true},
		Location:    req.Location,
		Color:       req.Color,
		Description: req.Description,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to update appointment", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update appointment")))
		return
	}

	if req.ParticipantEmployeeIDs != nil {
		err = qtx.DeleteAppointmentParticipants(ctx, appointment.ID)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to delete appointment participants", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete appointment participants")))
			return
		}

		if len(*req.ParticipantEmployeeIDs) > 0 {
			err = qtx.BulkAddAppointmentParticipants(ctx, db.BulkAddAppointmentParticipantsParams{
				AppointmentID: appointment.ID,
				EmployeeIds:   *req.ParticipantEmployeeIDs,
			})
			if err != nil {
				server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to add appointment participants", zap.Error(err))
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to add appointment participants")))
				return
			}
		}
	}

	if req.ClientIDs != nil {
		err = qtx.DeleteAppointmentClients(ctx, appointment.ID)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to delete appointment clients", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete appointment clients")))
			return
		}

		if len(*req.ClientIDs) > 0 {
			err = qtx.BulkAddAppointmentClients(ctx, db.BulkAddAppointmentClientsParams{
				AppointmentID: appointment.ID,
				ClientIds:     *req.ClientIDs,
			})
			if err != nil {
				server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to add appointment clients", zap.Error(err))
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to add appointment clients")))
				return
			}
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateAppointmentApi", "Failed to commit transaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to commit transaction")))
		return
	}

	response := UpdateAppointmentResponse{
		ID:                     appointment.ID,
		AppointmentTemplatesID: appointment.AppointmentTemplatesID,
		CreatorEmployeeID:      appointment.CreatorEmployeeID,
		StartTime:              appointment.StartTime,
		EndTime:                appointment.EndTime,
		Location:               appointment.Location,
		Description:            appointment.Description,
		Color:                  appointment.Color,
		Status:                 appointment.Status,
		IsConfirmed:            appointment.IsConfirmed,
		ConfirmedByEmployeeID:  appointment.ConfirmedByEmployeeID,
		ConfirmedAt:            appointment.ConfirmedAt,
		CreatedAt:              appointment.CreatedAt,
		UpdatedAt:              appointment.UpdatedAt,
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
		server.logBusinessEvent(LogLevelError, "DeleteAppointmentApi", "Invalid appointment ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid appointment ID")))
		return
	}

	err = server.store.DeleteAppointment(ctx, appointmentID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteAppointmentApi", "Failed to delete appointment", zap.Error(err))
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
		server.logBusinessEvent(LogLevelError, "ConfirmAppointmentApi", "Failed to get auth payload", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("failed to get auth payload")))
		return
	}

	err = server.store.ConfirmAppointment(ctx, db.ConfirmAppointmentParams{
		ID:         appointmentID,
		EmployeeID: &payload.EmployeeID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Appointment confirmed successfully")
	ctx.JSON(http.StatusOK, res)
}
