package api

import (
	"errors"
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateAppointmentRequest represents the request payload for creating an appointment
type CreateAppointmentRequest struct {
	StartTime              time.Time `json:"start_time" binding:"required"`
	EndTime                time.Time `json:"end_time" binding:"required"`
	Location               *string   `json:"location"`
	Description            *string   `json:"description"`
	Status                 string    `json:"status" binding:"required"`
	RecurrenceType         *string   `json:"recurrence_type"`
	RecurrenceInterval     *int32    `json:"recurrence_interval"`
	RecurrenceEndDate      time.Time `json:"recurrence_end_date"`
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

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)
	qtx := server.store.WithTx(tx)

	appointment, err := qtx.CreateAppointment(ctx, db.CreateAppointmentParams{
		CreatorEmployeeID:  employeeID,
		StartTime:          pgtype.Timestamp{Time: req.StartTime, Valid: true},
		EndTime:            pgtype.Timestamp{Time: req.EndTime, Valid: true},
		Location:           req.Location,
		Description:        req.Description,
		Status:             req.Status,
		RecurrenceType:     req.RecurrenceType,
		RecurrenceInterval: req.RecurrenceInterval,
		RecurrenceEndDate:  pgtype.Date{Time: req.RecurrenceEndDate, Valid: true},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, participantID := range req.ParticipantEmployeeIDs {
		err = qtx.AddAppointmentParticipant(ctx, db.AddAppointmentParticipantParams{
			AppointmentID: appointment.ID,
			EmployeeID:    participantID,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	for _, clientID := range req.ClientIDs {
		err = qtx.AddAppointmentClient(ctx, db.AddAppointmentClientParams{
			AppointmentID: appointment.ID,
			ClientID:      clientID,
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

	res := SuccessResponse(CreateAppointmentResponse{
		ID:                appointment.ID,
		CreatorEmployeeID: appointment.CreatorEmployeeID,
		StartTime:         appointment.StartTime.Time,
		EndTime:           appointment.EndTime.Time,
		Location:          appointment.Location,
		Description:       appointment.Description,
	}, "Appointment created successfully")

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

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)
	qtx := server.store.WithTx(tx)

	for _, participantID := range req.ParticipantEmployeeIDs {
		err = qtx.AddAppointmentParticipant(ctx, db.AddAppointmentParticipantParams{
			AppointmentID: appointmentID,
			EmployeeID:    participantID,
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

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)
	qtx := server.store.WithTx(tx)

	for _, clientID := range req.ClientIDs {
		err = qtx.AddAppointmentClient(ctx, db.AddAppointmentClientParams{
			AppointmentID: appointmentID,
			ClientID:      clientID,
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

	res := SuccessResponse[any](nil, "Clients added successfully")
	ctx.JSON(http.StatusOK, res)
}


