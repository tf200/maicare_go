package api

import (
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"

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
	responseLocations := make([]ListLocationsResponse, len(locations))
	for i, location := range locations {
		responseLocations[i] = ListLocationsResponse{
			ID:       location.ID,
			Name:     location.Name,
			Address:  location.Address,
			Capacity: location.Capacity,
		}
	}

	res := SuccessResponse(responseLocations, "Locations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// CreateLocationRequest represents a request to create a location
type CreateLocationRequest struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Capacity *int32 `json:"capacity"`
}

// CreateLocationResponse represents a response for CreateLocationApi
type CreateLocationResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// @Summary Create a location
// @Description Create a new location
// @Tags locations
// @Accept json
// @Produce json
// @Param input body CreateLocationRequest true "Create location"
// @Success 200 {object} Response[CreateLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations [post]
func (server *Server) CreateLocationApi(ctx *gin.Context) {
	var req CreateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	location, err := server.store.CreateLocation(ctx, db.CreateLocationParams{
		Name:     req.Name,
		Address:  req.Address,
		Capacity: req.Capacity,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(CreateLocationResponse{
		ID:       location.ID,
		Name:     location.Name,
		Address:  location.Address,
		Capacity: location.Capacity,
	}, "Location created successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateLocationRequest represents a request to update a location
type UpdateLocationRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Capacity *int32  `json:"capacity"`
}

// UpdateLocationResponse represents a response for UpdateLocationApi
type UpdateLocationResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// @Summary Update a location
// @Description Update a location
// @Tags locations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param input body UpdateLocationRequest true "Update location"
// @Success 200 {object} Response[UpdateLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [put]
func (server *Server) UpdateLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	location, err := server.store.UpdateLocation(ctx, db.UpdateLocationParams{
		Name:     req.Name,
		Address:  req.Address,
		Capacity: req.Capacity,
		ID:       locationID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(UpdateLocationResponse{
		ID:       location.ID,
		Name:     location.Name,
		Address:  location.Address,
		Capacity: location.Capacity,
	}, "Location updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteLocationResponse represents a response for DeleteLocationApi
type DeleteLocationResponse struct {
	ID int64 `json:"id"`
}

// @Summary Delete a location
// @Description Delete a location
// @Tags locations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} Response[DeleteLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [delete]
func (server *Server) DeleteLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.DeleteLocation(ctx, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(DeleteLocationResponse{
		ID: locationID,
	}, "Location deleted successfully")
	ctx.JSON(http.StatusOK, res)
}
