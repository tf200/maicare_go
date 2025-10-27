package api

import (
	"net/http"

	clientp "maicare_go/service/client"

	"github.com/gin-gonic/gin"
)

// ListAllIncidentsApi handles the API request to list all incidents
// @Summary List all incidents
// @Description List all incidents with pagination and filtering options
// @Tags incidents
// @Produce json
// @Param is_confirmed query bool false "Filter by confirmation status"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} Response[pagination.Response[ListAllIncidentsResponse]]
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /incidents [get]
func (serevr *Server) ListAllIncidentsApi(ctx *gin.Context) {
	var req clientp.ListAllIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pag, err := serevr.businessService.ClientService.ListAllIncidents(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(pag, "Incidents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}
