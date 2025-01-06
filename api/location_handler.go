package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListLocationsResponse represents a location in the list
type ListLocationsResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// @Summary List all locations
// @Description Get a list of all locations
// @Tags locations
// @Produce json
// @Success 200 {object} Response[[]ListLocationsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /locations [get]
func (server *Server) ListLocationsApi(ctx *gin.Context) {
	locations, err := server.store.ListLocations(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	list := []ListLocationsResponse{}
	for _, location := range locations {
		list = append(list, ListLocationsResponse{
			ID:       location.ID,
			Name:     location.Name,
			Address:  location.Address,
			Capacity: location.Capacity,
		})

		res := SuccessResponse(list, "Locations retrieved successfully")

		ctx.JSON(http.StatusOK, res)
	}
}
