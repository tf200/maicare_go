package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListLocationsResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// @Summary List all locations
// @Description Get a list of all locations
// @Tags locations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} ListLocationsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /location [get]
func (server *Server) ListLocationsApi(ctx *gin.Context) {
	locations, err := server.store.ListLocations(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := []ListLocationsResponse{}
	for _, location := range locations {
		res = append(res, ListLocationsResponse{
			ID:       location.ID,
			Name:     location.Name,
			Address:  location.Address,
			Capacity: location.Capacity,
		})

		ctx.JSON(http.StatusOK, res)
	}
}
