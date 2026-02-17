package api

import (
	"fmt"
	"net/http"

	clientp "maicare_go/service/client"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAppointmentCardApi retrieves an appointment card by client ID
// @Summary Get an appointment card by client ID
// @Description Get an appointment card by client ID
// @Tags appointment_cards
// @Produce json
// @Param id path uuid true "Client ID"
// @Success 200 {object} Response[clientp.GetAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [get]
func (server *Server) GetAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	appointmentCard, err := server.businessService.ClientService.GetAppointmentCard(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if appointmentCard == nil {
		res := SuccessResponse[any](nil, "Appointment card not found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := SuccessResponse(appointmentCard, "Appointment card retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateAppointmentCardApi creates or updates an appointment card by client ID
// @Summary Create or update an appointment card by client ID
// @Description Create or update an appointment card by client ID
// @Tags appointment_cards
// @Accept json
// @Produce json
// @Param id path uuid true "Client ID"
// @Param request body clientp.UpdateAppointmentCardRequest true "Request body"
// @Success 200 {object} Response[clientp.UpdateAppointmentCardResponse]
// @Router /clients/{id}/appointment_cards [put]
func (server *Server) UpdateAppointmentCardApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req clientp.UpdateAppointmentCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	appointmentCard, err := server.businessService.ClientService.UpdateAppointmentCard(
		req,
		clientID,
		ctx,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(appointmentCard, "Appointment card updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// GenerateAppointmentCardDocument generates an appointment card document by client ID
// @Summary Generate an appointment card document by client ID
// @Description Generate and download an appointment card PDF by client ID
// @Tags appointment_cards
// @Produce application/pdf
// @Param id path uuid true "Client ID"
// @Success 200 {file} file
// @Router /clients/{id}/appointment_cards/generate_document [post]
func (server *Server) GenerateAppointmentCardDocumentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	pdfBytes, fileName, err := server.businessService.ClientService.GenerateAppointmentCardDocumentApi(
		ctx,
		clientID,
	)
	if err != nil {
		if err.Error() == "appointment card not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(fmt.Errorf("failed to generate appointment card document")),
		)
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
