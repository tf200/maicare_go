package api

import (
	"errors"
	"fmt"
	_ "maicare_go/pagination" // for swagger
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreateIncidentApi creates an incident
// @Summary Create an incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body clientp.CreateIncidentRequest true "Incident data"
// @Success 201 {object} Response[clientp.CreateIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents [post]
func (server *Server) CreateIncidentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Invalid client ID", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req clientp.CreateIncidentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Invalid request body", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	incident, err := server.businessService.ClientService.CreateIncident(ctx, req, clientID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to create incident", zap.String("client_id", id), zap.Error(err))
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
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[clientp.ListIncidentsResponse]]
// @Router /clients/{id}/incidents [get]
func (server *Server) ListIncidentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Invalid client ID", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid client ID")))
		return
	}

	var req clientp.ListIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to bind query parameters", zap.String("client_id", id), zap.Error(err))
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
// @Param id path int true "Client ID"
// @Param incident_id path int true "Incident ID"
// @Success 200 {object} Response[clientp.GetIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id} [get]
func (server *Server) GetIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Client ID"
// @Param incident_id path int true "Incident ID"
// @Param incident body clientp.UpdateIncidentRequest true "Incident"
// @Success 200 {object} Response[clientp.UpdateIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id} [put]
func (server *Server) UpdateIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Client ID"
// @Param incident_id path int true "Incident ID"
// @Success 200 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id} [delete]
func (server *Server) DeleteIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
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
// @Produce json
// @Param incident_id path int true "Incident ID"
// @Param id path int true "Client ID"
// @Success 200 {object} Response[clientp.GenerateIncidentFileResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id}/file [get]
func (server *Server) GenerateIncidentFileApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}
	result, err := server.businessService.ClientService.GenerateIncidentFile(ctx, incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to generate incident file")))
		return
	}

	res := SuccessResponse(result, "Incident file generated successfully")

	ctx.JSON(http.StatusOK, res)
}

// ConfirmIncidentApi confirms an incident
// @Summary Confirm an incident
// @Tags incidents
// @Produce json
// @Param incident_id path int true "Incident ID"
// @Success 200 {object} Response[clientp.ConfirmIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents/{incident_id}/confirm [put]
func (server *Server) ConfirmIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}
	result, err := server.businessService.ClientService.ConfirmIncident(ctx, incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to confirm incident")))
		return
	}

	res := SuccessResponse(result, "Incident confirmed successfully")

	ctx.JSON(http.StatusOK, res)
}
