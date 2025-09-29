package api

import (
	"errors"
	"fmt"
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateAppointmentCardApi creates a new appointment card
// @Summary Create a new appointment card
// @Description Create a new appointment card
// @Tags appointment_cards
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.CreateAppointmentCardRequest true "Request body"
// @Success 201 {object} Response[clientp.CreateAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [post]
func (server *Server) CreateAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}
	var req clientp.CreateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	appointmentCard, err := server.businessService.ClientService.CreateAppointmentCard(req, clientID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create appointment card")))
		return
	}
	res := SuccessResponse(appointmentCard, "Appointment card created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// GetAppointmentCardApi retrieves an appointment card by client ID
// @Summary Get an appointment card by client ID
// @Description Get an appointment card by client ID
// @Tags appointment_cards
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GetAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [get]
func (server *Server) GetAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	appointmentCard, err := server.businessService.ClientService.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if errors.Is(err, fmt.Errorf("appointment card not found")) {
			res := SuccessResponse[any](nil, "Appointment card not found")
			ctx.JSON(http.StatusNotFound, res)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(appointmentCard, "Appointment card retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateAppointmentCardApi updates an appointment card by client ID
// @Summary Update an appointment card by client ID
// @Description Update an appointment card by client ID
// @Tags appointment_cards
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.UpdateAppointmentCardRequest true "Request body"
// @Success 200 {object} Response[clientp.UpdateAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [put]
func (server *Server) UpdateAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req clientp.UpdateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentCard, err := server.businessService.ClientService.UpdateAppointmentCard(req, clientID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(appointmentCard, "Appointment card updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// GenerateAppointmentCardDocument generates an appointment card document by client ID
// @Summary Generate an appointment card document by client ID
// @Description Generate an appointment card document by client ID
// @Tags appointment_cards
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GenerateAppointmentCardDocumentApiResponse]
// @Router /clients/{id}/appointment_cards/generate_document [post]
func (server *Server) GenerateAppointmentCardDocumentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateAppointmentCardDocumentApi", "Invalid client ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	response, err := server.businessService.ClientService.GenerateAppointmentCardDocumentApi(ctx, clientID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateAppointmentCardDocumentApi", "Failed to generate appointment card document", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to generate appointment card document")))
		return
	}

	res := SuccessResponse(response, "Appointment card document generated successfully")
	ctx.JSON(http.StatusOK, res)

}
