package api

import (
	"net/http"

	_ "maicare_go/pagination"
	clientp "maicare_go/service/client"

	"github.com/gin-gonic/gin"
)

// ListAllIncidentsApi handles the API request to list all incidents
// @Summary List all incidents
// @Description List all incidents with pagination and filtering options
// @Tags incidents
// @Produce json
// @Param is_confirmed query bool false "Filter by confirmation status"
// @Param search query string false "Search by client first or last name"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} Response[pagination.Response[clientp.ListAllIncidentsResponse]]
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /incidents [get]
func (server *Server) ListAllIncidentsApi(ctx *gin.Context) {
	var req clientp.ListAllIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := server.businessService.ClientService.ListAllIncidents(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(pag, "Incidents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetIncidentCountsApi gets incident aggregate counts
// @Summary Get incident counts
// @Description Get counts for serious/fatal incidents, pending confirmations, and incidents in the past 24 hours
// @Tags incidents
// @Produce json
// @Success 200 {object} Response[clientp.GetIncidentCountsResponse]
// @Failure 500 {object} Response[any]
// @Router /incidents/counts [get]
func (server *Server) GetIncidentCountsApi(ctx *gin.Context) {
	counts, err := server.businessService.ClientService.GetIncidentCounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(counts, "Incident counts fetched successfully")
	ctx.JSON(http.StatusOK, res)
}
