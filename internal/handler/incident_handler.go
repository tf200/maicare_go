package handler

import (
	"fmt"
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterIncidentRoutes(
	rg *gin.RouterGroup,
	handler *IncidentHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	// Client-specific incident routes
	clientsGroup := rg.Group("/clients")
	clientsGroup.GET("/:id/incidents", auth, requirePermission("CLIENT.INCIDENT.VIEW"), handler.ListClientIncidents)

	// Top-level incident routes
	incidentsGroup := rg.Group("/incidents")
	{
		incidentsGroup.POST("", auth, requirePermission("CLIENT.INCIDENT.CREATE"), handler.CreateIncident)
		incidentsGroup.GET("", auth, requirePermission("CLIENT.INCIDENT.VIEW"), handler.ListAllIncidents)
		incidentsGroup.GET("/counts", auth, requirePermission("CLIENT.INCIDENT.VIEW"), handler.GetIncidentCounts)
		incidentsGroup.GET("/:incident_id", auth, requirePermission("CLIENT.INCIDENT.VIEW"), handler.GetIncident)
		incidentsGroup.PUT("/:incident_id", auth, requirePermission("CLIENT.INCIDENT.UPDATE"), handler.UpdateIncident)
		incidentsGroup.DELETE("/:incident_id", auth, requirePermission("CLIENT.INCIDENT.DELETE"), handler.DeleteIncident)
		incidentsGroup.GET("/:incident_id/file", auth, requirePermission("CLIENT.INCIDENT.VIEW"), handler.GenerateIncidentFile)
		incidentsGroup.PUT("/:incident_id/confirm", auth, requirePermission("CLIENT.INCIDENT.CONFIRM"), handler.ConfirmIncident)
	}
}

type IncidentHandler struct {
	service domain.IncidentService
}

func NewIncidentHandler(service domain.IncidentService) *IncidentHandler {
	return &IncidentHandler{service: service}
}

// CreateIncident creates a new incident.
// @Summary Create an incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param request body createIncidentRequest true "Incident data"
// @Success 201 {object} httpapi.Envelope[createIncidentResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /incidents [post]
func (h *IncidentHandler) CreateIncident(ctx *gin.Context) {
	var req createIncidentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", ""))
		return
	}

	if req.ClientID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	result, err := h.service.CreateIncident(ctx.Request.Context(), toCreateIncidentParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create incident", ""))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toCreateIncidentResponse(*result), "Incident created successfully"))
}

// ListClientIncidents lists all incidents for a specific client.
// @Summary List all incidents for a client
// @Tags incidents
// @Produce json
// @Param id path uuid true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listIncidentsResponse]]
// @Router /clients/{id}/incidents [get]
func (h *IncidentHandler) ListClientIncidents(ctx *gin.Context) {
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid client ID", ""))
		return
	}

	var req listIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to bind query parameters", ""))
		return
	}

	pageParams := req.PageRequest.Params()
	result, err := h.service.ListIncidents(ctx.Request.Context(), domain.ListIncidentsParams{
		ClientID: clientID,
		Limit:    pageParams.Limit,
		Offset:   pageParams.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list incidents", ""))
		return
	}

	items := make([]listIncidentsResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toListIncidentsResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Incidents retrieved successfully"))
}

// GetIncident retrieves an incident by ID.
// @Summary Retrieve an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {object} httpapi.Envelope[getIncidentResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /incidents/{incident_id} [get]
func (h *IncidentHandler) GetIncident(ctx *gin.Context) {
	incidentID, err := uuid.Parse(ctx.Param("incident_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid incident ID", ""))
		return
	}

	result, err := h.service.GetIncident(ctx.Request.Context(), incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to retrieve incident", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetIncidentResponse(*result), "Incident retrieved successfully"))
}

// UpdateIncident updates an incident.
// @Summary Update an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Param incident body updateIncidentRequest true "Incident"
// @Success 200 {object} httpapi.Envelope[updateIncidentResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /incidents/{incident_id} [put]
func (h *IncidentHandler) UpdateIncident(ctx *gin.Context) {
	incidentID, err := uuid.Parse(ctx.Param("incident_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid incident ID", ""))
		return
	}

	var req updateIncidentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to process incident data", ""))
		return
	}

	result, err := h.service.UpdateIncident(ctx.Request.Context(), toUpdateIncidentParams(req, incidentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update incident", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toUpdateIncidentResponse(*result), "Incident updated successfully"))
}

// DeleteIncident deletes an incident.
// @Summary Delete an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {object} httpapi.Envelope[any]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /incidents/{incident_id} [delete]
func (h *IncidentHandler) DeleteIncident(ctx *gin.Context) {
	incidentID, err := uuid.Parse(ctx.Param("incident_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid incident ID", ""))
		return
	}

	if err := h.service.DeleteIncident(ctx.Request.Context(), incidentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to delete incident", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK([]string{}, "Incident deleted successfully"))
}

// GenerateIncidentFile generates a PDF file for an incident.
// @Summary Generate an incident file
// @Tags incidents
// @Produce application/pdf
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {file} file
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /incidents/{incident_id}/file [get]
func (h *IncidentHandler) GenerateIncidentFile(ctx *gin.Context) {
	incidentID, err := uuid.Parse(ctx.Param("incident_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid incident ID", ""))
		return
	}

	pdfBytes, fileName, err := h.service.GenerateIncidentFile(ctx.Request.Context(), incidentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to generate incident file", ""))
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ConfirmIncident confirms an incident.
// @Summary Confirm an incident
// @Tags incidents
// @Produce json
// @Param incident_id path uuid true "Incident ID"
// @Success 200 {object} httpapi.Envelope[confirmIncidentResponse]
// @Failure 400,404,500 {object} httpapi.Envelope[any]
// @Router /incidents/{incident_id}/confirm [put]
func (h *IncidentHandler) ConfirmIncident(ctx *gin.Context) {
	incidentID, err := uuid.Parse(ctx.Param("incident_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid incident ID", ""))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail(err.Error(), ""))
		return
	}

	result, err := h.service.ConfirmIncident(ctx.Request.Context(), incidentID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to confirm incident", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toConfirmIncidentResponse(*result), "Incident confirmed successfully"))
}

// ListAllIncidents lists all incidents with filtering.
// @Summary List all incidents
// @Tags incidents
// @Produce json
// @Param is_confirmed query bool false "Filter by confirmation status"
// @Param search query string false "Search by client first or last name"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[listAllIncidentsResponse]]
// @Failure 400,500 {object} httpapi.Envelope[any]
// @Router /incidents [get]
func (h *IncidentHandler) ListAllIncidents(ctx *gin.Context) {
	var req listAllIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	pageParams := req.PageRequest.Params()
	result, err := h.service.ListAllIncidents(ctx.Request.Context(), domain.ListAllIncidentsParams{
		Limit:       pageParams.Limit,
		Offset:      pageParams.Offset,
		IsConfirmed: req.IsConfirmed,
		Search:      req.Search,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list incidents", ""))
		return
	}

	items := make([]listAllIncidentsResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toListAllIncidentsResponse(item))
	}

	pageResp := httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount)
	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Incidents fetched successfully"))
}

// GetIncidentCounts gets aggregate incident counts.
// @Summary Get incident counts
// @Tags incidents
// @Produce json
// @Success 200 {object} httpapi.Envelope[getIncidentCountsResponse]
// @Failure 500 {object} httpapi.Envelope[any]
// @Router /incidents/counts [get]
func (h *IncidentHandler) GetIncidentCounts(ctx *gin.Context) {
	counts, err := h.service.GetIncidentCounts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get incident counts", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGetIncidentCountsResponse(*counts), "Incident counts fetched successfully"))
}

func getUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	idStr := ctx.GetString("user_id")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return uuid.Parse(idStr)
}
