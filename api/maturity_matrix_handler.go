package api

import (
	_ "maicare_go/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary List all maturity matrix
// @Description Get a list of all maturity matrix
// @Tags maturity_matrix
// @Produce json
// @Success 200 {object} Response[[]care.ListCarePlanTopics]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /maturity_matrix [get]
func (server *Server) ListMaturityMatrixApi(ctx *gin.Context) {
	result, err := server.businessService.CarePlanService.ListCarePlanTopics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Maturity matrix retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}
