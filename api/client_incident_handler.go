package api

import (
	"errors"
	"fmt"
	"net/http"

	_ "maicare_go/pagination" // for swagger
	clientp "maicare_go/service/client"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateIncidentApi creates an incident
// @Summary Create an incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param request body clientp.CreateIncidentRequest true "Incident data"
// @Success 201 {object} Response[clientp.CreateIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents [post]
func (server *Server) CreateIncidentApi(ctx *gin.Context) {
	var req clientp.CreateIncidentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	if req.ClientID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	incident, err := server.businessService.ClientService.CreateIncident(ctx, req)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to create incident", zap.String("client_id", req.ClientID.String()), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create incident")))
		return
	}

	res := SuccessResponse(incident, "Incident created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListIncidentsApi lists all incidents
// @Summary List all incidents
// @Tags incidents
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListIncidentsResponse]]
// @Router /clients/{id}/incidents [get]
func (server *Server) ListIncidentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid client ID")))
		return
	}

	var req clientp.ListIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("failed to bind query parameters")))
		return
	}

	pag, err := server.businessService.ClientService.ListIncidents(ctx, req, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to list incidents")))
		return
	}

	res := SuccessResponse(pag, "Incidents retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetIncidentApi retrieves an incident
// @Summary Retrieve an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {object} Response[clientp.GetIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents/{incident_id} [get]
func (server *Server) GetIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}

	incident, err := server.businessService.ClientService.GetIncident(ctx, incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to retrieve incident")))
		return
	}

	res := SuccessResponse(incident, "Incident retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateIncidentApi updates an incident
// @Summary Update an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Param incident body clientp.UpdateIncidentRequest true "Incident"
// @Success 200 {object} Response[clientp.UpdateIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents/{incident_id} [put]
func (server *Server) UpdateIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}

	var req clientp.UpdateIncidentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("failed to process incident data")))
		return
	}

	result, err := server.businessService.ClientService.UpdateIncident(ctx, clientp.UpdateIncidentRequest(req), incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to update incident")))
		return
	}

	res := SuccessResponse(result, "Incident updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteIncidentApi deletes an incident
// @Summary Delete an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents/{incident_id} [delete]
func (server *Server) DeleteIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}
	err = server.businessService.ClientService.DeleteIncident(ctx, incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to delete incident")))
		return
	}
	res := SuccessResponse([]string{}, "Incident deleted successfully")

	ctx.JSON(http.StatusOK, res)
}

// GenerateIncidentFileApi generates an incident file
// @Summary Generate an incident file
// @Tags incidents
// @Produce application/pdf
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {file} file
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents/{incident_id}/file [get]
func (server *Server) GenerateIncidentFileApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}
	pdfBytes, fileName, err := server.businessService.ClientService.GenerateIncidentFile(ctx, incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to generate incident file")))
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ConfirmIncidentApi confirms an incident
// @Summary Confirm an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {object} Response[clientp.ConfirmIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents/{incident_id}/confirm [put]
func (server *Server) ConfirmIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	result, err := server.businessService.ClientService.ConfirmIncident(ctx, incidentID, payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to confirm incident")))
		return
	}

	res := SuccessResponse(result, "Incident confirmed successfully")

	ctx.JSON(http.StatusOK, res)
}
