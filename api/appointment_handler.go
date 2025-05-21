package api

import (
	"errors"
	"maicare_go/async"
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateAppointmentRequest represents the request payload for creating an appointment
type CreateAppointmentRequest struct {
	StartTime              time.Time `json:"start_time" binding:"required" example:"2023-10-01T10:00:00Z"`
	EndTime                time.Time `json:"end_time" binding:"required"`
	Location               *string   `json:"location"`
	Description            *string   `json:"description"`
	RecurrenceType         string    `json:"recurrence_type" example:"NONE" enum:"NONE,DAILY,WEEKLY,MONTHLY"`
	RecurrenceInterval     *int32    `json:"recurrence_interval"`
	RecurrenceEndDate      time.Time `json:"recurrence_end_date" example:"2025-10-01T10:00:00Z"`
	ParticipantEmployeeIDs []int64   `json:"participant_employee_ids"`
	ClientIDs              []int64   `json:"client_ids"`
}

// CreateAppointmentResponse represents the response payload for creating an appointment
type CreateAppointmentResponse struct {
	ID                int64     `json:"id"`
	CreatorEmployeeID int64     `json:"creator_employee_id"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Location          *string   `json:"location"`
	Description       *string   `json:"description"`
}

// @Summary Create an appointment
// @Description Create a new appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param request body CreateAppointmentRequest true "Create appointment request"
// @Success 200 {object} Response[CreateAppointmentResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - User not found"
// @Failure 409 {object} Response[any] "Conflict - Appointment already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments [post]
func (server *Server) CreateAppointmentApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userID := payload.UserId

	employeeID, err := server.store.GetEmployeeIDByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req CreateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.StartTime.After(req.EndTime) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("start time must be before end time")))
		return
	}

	var response CreateAppointmentResponse

	if req.RecurrenceType == "NONE" {
		tx, err := server.store.ConnPool.Begin(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		defer tx.Rollback(ctx)
		qtx := server.store.WithTx(tx)

		appointment, err := qtx.CreateAppointment(ctx, db.CreateAppointmentParams{
			CreatorEmployeeID: &employeeID,
			StartTime:         pgtype.Timestamp{Time: req.StartTime, Valid: true},
			EndTime:           pgtype.Timestamp{Time: req.EndTime, Valid: true},
			Location:          req.Location,
			Description:       req.Description,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if len(req.ParticipantEmployeeIDs) > 0 {
			err = qtx.BulkAddAppointmentParticipants(ctx, db.BulkAddAppointmentParticipantsParams{
				AppointmentID: appointment.ID,
				EmployeeIds:   req.ParticipantEmployeeIDs,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		if len(req.ClientIDs) > 0 {
			err = qtx.BulkAddAppointmentClients(ctx, db.BulkAddAppointmentClientsParams{
				AppointmentID: appointment.ID,
				ClientIds:     req.ClientIDs,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		err = tx.Commit(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		response = CreateAppointmentResponse{
			ID:                appointment.ID,
			CreatorEmployeeID: *appointment.CreatorEmployeeID,
			StartTime:         appointment.StartTime.Time,
			EndTime:           appointment.EndTime.Time,
			Location:          appointment.Location,
			Description:       appointment.Description,
		}

	} else {
		appointmentTemp, err := server.store.CreateAppointmentTemplate(ctx, db.CreateAppointmentTemplateParams{
			CreatorEmployeeID:  employeeID,
			StartTime:          pgtype.Timestamp{Time: req.StartTime, Valid: true},
			EndTime:            pgtype.Timestamp{Time: req.EndTime, Valid: true},
			Location:           req.Location,
			Description:        req.Description,
			RecurrenceType:     &req.RecurrenceType,
			RecurrenceInterval: req.RecurrenceInterval,
			RecurrenceEndDate:  pgtype.Date{Time: req.RecurrenceEndDate, Valid: true},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		response = CreateAppointmentResponse{
			ID:                appointmentTemp.ID,
			CreatorEmployeeID: appointmentTemp.CreatorEmployeeID,
			StartTime:         appointmentTemp.StartTime.Time,
			EndTime:           appointmentTemp.EndTime.Time,
			Location:          appointmentTemp.Location,
			Description:       appointmentTemp.Description,
		}

		server.asynqClient.EnqueueAppointmentTask(ctx, async.AppointmentPayload{
			AppointmentTemplateID:  appointmentTemp.ID,
			ParticipantEmployeeIDs: req.ParticipantEmployeeIDs,
			ClientIDs:              req.ClientIDs,
		})

	}

	res := SuccessResponse(response, "Appointment created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// AddParticipantToAppointmentRequest represents the request payload for adding participants to an appointment
type AddParticipantToAppointmentRequest struct {
	ParticipantEmployeeIDs []int64 `json:"participant_employee_ids"`
}

// AddParticipantToAppointment adds participants to an existing appointment
// @Summary Add participants to an appointment
// @Description Add participants to an existing appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param appointment_id path int true "Appointment ID"
// @Param request body AddParticipantToAppointmentRequest true "Add participant request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Participant already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{appointment_id}/participants [post]
func (server *Server) AddParticipantToAppointmentApi(ctx *gin.Context) {
	appointmentID, err := strconv.ParseInt(ctx.Param("appointment_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req AddParticipantToAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.BulkAddAppointmentParticipants(ctx, db.BulkAddAppointmentParticipantsParams{
		AppointmentID: appointmentID,
		EmployeeIds:   req.ParticipantEmployeeIDs,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Participants added successfully")
	ctx.JSON(http.StatusOK, res)
}

// AddClientToAppointmentRequest represents the request payload for adding clients to an appointment
type AddClientToAppointmentRequest struct {
	ClientIDs []int64 `json:"client_ids"`
}

// AddClientToAppointment adds clients to an existing appointment
// @Summary Add clients to an appointment
// @Description Add clients to an existing appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Param appointment_id path int true "Appointment ID"
// @Param request body AddClientToAppointmentRequest true "Add client request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Client already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{appointment_id}/clients [post]
func (server *Server) AddClientToAppointmentApi(ctx *gin.Context) {
	appointmentID, err := strconv.ParseInt(ctx.Param("appointment_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req AddClientToAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.BulkAddAppointmentClients(ctx, db.BulkAddAppointmentClientsParams{
		AppointmentID: appointmentID,
		ClientIds:     req.ClientIDs,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Clients added successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListAppointmentsForEmployeeInRangeRequest represents the request payload for listing appointments for an employee in a date range
type ListAppointmentsForEmployeeInRangeRequest struct {
	StartDate time.Time `json:"start_date" binding:"required" example:"2025-04-27T00:00:00Z"`
	EndDate   time.Time `json:"end_date" binding:"required" example:"2025-04-30T23:59:59Z"`
}

// ListAppointmentsForEmployeeInRangeResponse represents the response payload for listing appointments for an employee in a date range
type ListAppointmentsForEmployeeInRangeResponse struct {
	ID                int64     `json:"id"`
	CreatorEmployeeID *int64    `json:"creator_employee_id"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Location          *string   `json:"location"`
	Description       *string   `json:"description"`
	Status            string    `json:"status"`
	IsConfirmed       bool      `json:"is_confirmed"`
	CreatedAt         time.Time `json:"created_at"`
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

	var req ListAppointmentsForEmployeeInRangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.StartDate.After(req.EndDate) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("start date must be before end date")))
		return
	}

	arg := db.ListEmployeeAppointmentsInRangeParams{
		EmployeeID: &employeeID,
		StartDate:  pgtype.Timestamp{Time: req.StartDate, Valid: true},
		EndDate:    pgtype.Timestamp{Time: req.EndDate, Valid: true},
	}

	appointments, err := server.store.ListEmployeeAppointmentsInRange(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	appointmentList := make([]ListAppointmentsForEmployeeInRangeResponse, len(appointments))
	for i, appointment := range appointments {
		appointmentList[i] = ListAppointmentsForEmployeeInRangeResponse{
			ID:                appointment.AppointmentID,
			CreatorEmployeeID: appointment.CreatorEmployeeID,
			StartTime:         appointment.StartTime.Time,
			EndTime:           appointment.EndTime.Time,
			Location:          appointment.Location,
			Description:       appointment.Description,
			Status:            appointment.Status,
			IsConfirmed:       appointment.IsConfirmed,
			CreatedAt:         appointment.CreatedAt.Time,
		}
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
type ListAppointmentsForClientResponse struct {
	ID                    int64       `json:"id"`
	CreatorEmployeeID     *int64      `json:"creator_employee_id"`
	StartTime             time.Time   `json:"start_time"`
	EndTime               time.Time   `json:"end_time"`
	Location              *string     `json:"location"`
	Description           *string     `json:"description"`
	Status                string      `json:"status"`
	RecurrenceType        *string     `json:"recurrence_type"`
	RecurrenceInterval    *int32      `json:"recurrence_interval"`
	RecurrenceEndDate     pgtype.Date `json:"recurrence_end_date"`
	ConfirmedByEmployeeID *int32      `json:"confirmed_by_employee_id"`
	ConfirmedAt           time.Time   `json:"confirmed_at"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	IsRecurringOccurrence bool        `json:"is_recurring_occurrence"`
}

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
	appointmentList := make([]ListAppointmentsForClientResponse, len(appointments))
	for i, appointment := range appointments {
		appointmentList[i] = ListAppointmentsForClientResponse{
			ID:                appointment.AppointmentID,
			CreatorEmployeeID: appointment.CreatorEmployeeID,
			StartTime:         appointment.StartTime.Time,
			EndTime:           appointment.EndTime.Time,
			Location:          appointment.Location,
			Description:       appointment.Description,
			Status:            appointment.Status,
			CreatedAt:         appointment.CreatedAt.Time,
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
	ID                     int64                 `json:"id"`
	AppointmentTemplatesID *int64                `json:"appointment_templates_id"`
	CreatorEmployeeID      *int64                `json:"creator_employee_id"`
	CreatorFirstName       *string               `json:"creator_first_name"`
	CreatorLastName        *string               `json:"creator_last_name"`
	StartTime              time.Time             `json:"start_time"`
	EndTime                time.Time             `json:"end_time"`
	Location               *string               `json:"location"`
	Description            *string               `json:"description"`
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
// @Param id path int true "Appointment ID"
// @Success 200 {object} Response[GetAppointmentResponse]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{id} [get]
func (server *Server) GetAppointmentApi(ctx *gin.Context) {
	appointmentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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
		Status:                appointment.Status,
		IsConfirmed:           appointment.IsConfirmed,
		ConfirmedByEmployeeID: appointment.ConfirmedByEmployeeID,
		ConfirmerFirstName:    appointment.ConfirmerFirstName,
		ConfirmerLastName:     appointment.ConfirmerLastName,
	}

	res := SuccessResponse(response, "Appointment retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// ConfirmAppointmentApi confirms an appointment
// @Summary Confirm an appointment
// @Description Confirm an appointment
// @Tags appointments
// @Produce json
// @Param id path int true "Appointment ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 401 {object} Response[any] "Unauthorized - Invalid credentials"
// @Failure 404 {object} Response[any] "Not found - Appointment not found"
// @Failure 409 {object} Response[any] "Conflict - Appointment already exists"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /appointments/{id}/confirm [post]
func (server *Server) ConfirmAppointmentApi(ctx *gin.Context) {
	appointmentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userID := payload.UserId

	employeeID, err := server.store.GetEmployeeIDByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.ConfirmAppointment(ctx, db.ConfirmAppointmentParams{
		ID:         appointmentID,
		EmployeeID: &employeeID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Appointment confirmed successfully")
	ctx.JSON(http.StatusOK, res)
}
