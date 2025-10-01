package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateOrganisationRequest represents a request to create an organisation
type CreateOrganisationRequest struct {
	Name       string  `json:"name" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	PostalCode string  `json:"postal_code" binding:"required"`
	City       string  `json:"city" binding:"required"`
	Email      *string `json:"email"`
	KvkNumber  *string `json:"kvk_number"`
	BtwNumber  *string `json:"btw_number"`
}

// CreateOrganisationResponse represents a response for CreateOrganisationApi
type CreateOrganisationResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	PostalCode string  `json:"postal_code"`
	City       string  `json:"city"`
	Email      *string `json:"email"`
	KvkNumber  *string `json:"kvk_number"`
	BtwNumber  *string `json:"btw_number"`
}

// @Summary Create an organisation
// @Description Create a new organisation
// @Tags organisations
// @Accept json
// @Produce json
// @Param input body CreateOrganisationRequest true "Create organisation"
// @Success 200 {object} Response[CreateOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations [post]
func (server *Server) CreateOrganisationApi(ctx *gin.Context) {
	var req CreateOrganisationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateOrganisationApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	organisation, err := server.store.CreateOrganisation(ctx, db.CreateOrganisationParams{
		Name:       req.Name,
		Address:    req.Address,
		PostalCode: req.PostalCode,
		City:       req.City,
		Email:      req.Email,
		KvkNumber:  req.KvkNumber,
		BtwNumber:  req.BtwNumber,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateOrganisationApi", "Failed to create organisation", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create organisation"))
		return
	}

	res := SuccessResponse(CreateOrganisationResponse{
		ID:         organisation.ID,
		Name:       organisation.Name,
		Address:    organisation.Address,
		PostalCode: organisation.PostalCode,
		City:       organisation.City,
		Email:      organisation.Email,
		KvkNumber:  organisation.KvkNumber,
		BtwNumber:  organisation.BtwNumber,
	}, "Organisation created successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListOrganisationsResponse represents an organisation in the list
type ListOrganisationsResponse struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	PostalCode    string  `json:"postal_code"`
	City          string  `json:"city"`
	Email         *string `json:"email"`
	KvkNumber     *string `json:"kvk_number"`
	BtwNumber     *string `json:"btw_number"`
	LocationCount int64   `json:"location_count"`
}

// @Summary List all organisations
// @Description Get a list of all organisations
// @Tags organisations
// @Produce json
// @Success 200 {object} Response[[]ListOrganisationsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /organisations [get]
func (server *Server) ListOrganisationsApi(ctx *gin.Context) {
	organisations, err := server.store.ListOrganisations(ctx)

	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListOrganisationsApi", "Failed to list organisations", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to list organisations"))
		return
	}
	if len(organisations) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse([]ListOrganisationsResponse{}, "No organisations found"))
		return
	}

	responseOrganisations := make([]ListOrganisationsResponse, len(organisations))
	for i, organisation := range organisations {
		responseOrganisations[i] = ListOrganisationsResponse{
			ID:            organisation.ID,
			Name:          organisation.Name,
			Address:       organisation.Address,
			PostalCode:    organisation.PostalCode,
			City:          organisation.City,
			Email:         organisation.Email,
			KvkNumber:     organisation.KvkNumber,
			BtwNumber:     organisation.BtwNumber,
			LocationCount: organisation.LocationCount,
		}
	}

	res := SuccessResponse(responseOrganisations, "Organisations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetOrganisationResponse represents a response for GetOrganisationApi
type GetOrganisationResponse struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	PostalCode    string    `json:"postal_code"`
	City          string    `json:"city"`
	Email         *string   `json:"email"`
	KvkNumber     *string   `json:"kvk_number"`
	BtwNumber     *string   `json:"btw_number"`
	LocationCount int64     `json:"location_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// @Summary Get an organisation
// @Description Get an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Success 200 {object} Response[GetOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id} [get]
func (server *Server) GetOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetOrganisationApi", "Invalid organisation ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}

	organisation, err := server.store.GetOrganisation(ctx, organisationID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetOrganisationApi", "Failed to get organisation", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get organisation"))
		return
	}

	res := SuccessResponse(GetOrganisationResponse{
		ID:            organisation.ID,
		Name:          organisation.Name,
		Address:       organisation.Address,
		PostalCode:    organisation.PostalCode,
		City:          organisation.City,
		Email:         organisation.Email,
		KvkNumber:     organisation.KvkNumber,
		BtwNumber:     organisation.BtwNumber,
		LocationCount: organisation.LocationCount,
		CreatedAt:     organisation.CreatedAt.Time,
		UpdatedAt:     organisation.UpdatedAt.Time,
	}, "Organisation retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

type GetOrganisationCountResponse struct {
	OrganisationID   int64  `json:"organisation_id"`
	OrganisationName string `json:"organisation_name"`
	LocationCount    int64  `json:"location_count"`
	ClientCount      int64  `json:"client_count"`
	EmployeeCount    int64  `json:"employee_count"`
}

// @Summary Get organisation counts
// @Description Get counts of locations, clients, and employees for an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Success 200 {object} Response[GetOrganisationCountResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id}/counts [get]
func (server *Server) GetOrganisationCountApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetOrganisationCountApi", "Invalid organisation ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}

	count, err := server.store.GetOrganisationCounts(ctx, organisationID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetOrganisationCountApi", "Failed to get organisation location count", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get organisation location count"))
		return
	}
	res := SuccessResponse(GetOrganisationCountResponse{
		OrganisationID:   count.OrganisationID,
		OrganisationName: count.OrganisationName,
		LocationCount:    count.LocationCount,
		ClientCount:      count.ClientCount,
		EmployeeCount:    count.EmployeeCount,
	}, "Organisation counts retrieved successfully")
	ctx.JSON(http.StatusOK, res)

}

// UpdateOrganisationRequest represents a request to update an organisation
type UpdateOrganisationRequest struct {
	Name       *string `json:"name"`
	Address    *string `json:"address"`
	PostalCode *string `json:"postal_code"`
	City       *string `json:"city"`
	Email      *string `json:"email"`
	KvkNumber  *string `json:"kvk_number"`
	BtwNumber  *string `json:"btw_number"`
}

// UpdateOrganisationResponse represents a response for UpdateOrganisationApi
type UpdateOrganisationResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	PostalCode string  `json:"postal_code"`
	City       string  `json:"city"`
	Email      *string `json:"email"`
	KvkNumber  *string `json:"kvk_number"`
	BtwNumber  *string `json:"btw_number"`
}

// @Summary Update an organisation
// @Description Update an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Param input body UpdateOrganisationRequest true "Update organisation"
// @Success 200 {object} Response[UpdateOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id} [put]
func (server *Server) UpdateOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateOrganisationApi", "Invalid organisation ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}

	var req UpdateOrganisationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateOrganisationApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	organisation, err := server.store.UpdateOrganisation(ctx, db.UpdateOrganisationParams{
		ID:         organisationID,
		Name:       req.Name,
		Address:    req.Address,
		PostalCode: req.PostalCode,
		City:       req.City,
		Email:      req.Email,
		KvkNumber:  req.KvkNumber,
		BtwNumber:  req.BtwNumber,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateOrganisationApi", "Failed to update organisation", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to update organisation"))
		return
	}

	res := SuccessResponse(UpdateOrganisationResponse{
		ID:         organisation.ID,
		Name:       organisation.Name,
		Address:    organisation.Address,
		PostalCode: organisation.PostalCode,
		City:       organisation.City,
		Email:      organisation.Email,
		KvkNumber:  organisation.KvkNumber,
		BtwNumber:  organisation.BtwNumber,
	}, "Organisation updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// DeleteOrganisationResponse represents a response for DeleteOrganisationApi
type DeleteOrganisationResponse struct {
	ID int64 `json:"id"`
}

// @Summary Delete an organisation
// @Description Delete an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Success 200 {object} Response[DeleteOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id} [delete]
func (server *Server) DeleteOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteOrganisationApi", "Invalid organisation ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}

	_, err = server.store.DeleteOrganisation(ctx, organisationID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteOrganisationApi", "Failed to delete organisation", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete organisation"))
		return
	}

	res := SuccessResponse(DeleteOrganisationResponse{
		ID: organisationID,
	}, "Organisation deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// ListLocationsResponse represents a location in the list
type ListLocationsResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// @Summary List all locations
// @Description Get a list of all locations
// @Tags organisations
// @Param id path int true "Organisation ID"
// @Produce json
// @Success 200 {object} Response[[]ListLocationsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /organisations/{id}/locations [get]
func (server *Server) ListLocationsApi(ctx *gin.Context) {
	organizationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListLocationsApi", "Invalid organisation ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}
	locations, err := server.store.ListLocations(ctx, organizationID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListLocationsApi", "Failed to list locations", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to list locations"))
		return
	}
	if len(locations) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse([]ListLocationsResponse{}, "No locations found"))
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

// ListAllLocationsApi lists all locations across all organisations
// @Summary List all locations
// @Description Get a list of all locations across all organisations
// @Tags organisations
// @Produce json
// @Success 200 {object} Response[[]ListLocationsResponse]
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /locations [get]
func (server *Server) ListAllLocationsApi(ctx *gin.Context) {
	locations, err := server.store.ListAllLocations(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListAllLocationsApi", "Failed to list all locations", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to list all locations"))
		return
	}
	if len(locations) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse([]ListLocationsResponse{}, "No locations found"))
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

	res := SuccessResponse(responseLocations, "All locations retrieved successfully")

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
// @Tags organisations
// @Param id path int true "Organisation ID"
// @Accept json
// @Produce json
// @Param input body CreateLocationRequest true "Create location"
// @Success 200 {object} Response[CreateLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id}/locations [post]
func (server *Server) CreateLocationApi(ctx *gin.Context) {
	organisationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateLocationApi", "Invalid organisation ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}
	var req CreateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateLocationApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	location, err := server.store.CreateLocation(ctx, db.CreateLocationParams{
		OrganisationID: organisationID,
		Name:           req.Name,
		Address:        req.Address,
		Capacity:       req.Capacity,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateLocationApi", "Failed to create location", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create location"))
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
// @Tags organisations
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
		server.logBusinessEvent(LogLevelError, "UpdateLocationApi", "Invalid location ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid location ID"))
		return
	}

	var req UpdateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateLocationApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	location, err := server.store.UpdateLocation(ctx, db.UpdateLocationParams{
		Name:     req.Name,
		Address:  req.Address,
		Capacity: req.Capacity,
		ID:       locationID,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateLocationApi", "Failed to update location", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to update location"))
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
// @Tags organisations
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
		server.logBusinessEvent(LogLevelError, "DeleteLocationApi", "Invalid location ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid location ID"))
		return
	}
	_, err = server.store.DeleteLocation(ctx, locationID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteLocationApi", "Failed to delete location", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete location"))
		return
	}
	res := SuccessResponse(DeleteLocationResponse{
		ID: locationID,
	}, "Location deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetLocationResponse represents a response for GetLocationApi
type GetLocationResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// @Summary Get a location
// @Description Get a location
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} Response[GetLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [get]
func (server *Server) GetLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetLocationApi", "Invalid location ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid location ID"))
		return
	}
	location, err := server.store.GetLocation(ctx, locationID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetLocationApi", "Failed to retrieve location", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to retrieve location"))
		return
	}
	res := SuccessResponse(GetLocationResponse{
		ID:       location.ID,
		Name:     location.Name,
		Address:  location.Address,
		Capacity: location.Capacity,
	}, "Location retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}
