package api

import (
	"fmt"
	_ "maicare_go/pagination"
	"maicare_go/service/organization"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create an organization
// @Description Create a new organization
// @Tags organizations
// @Accept json
// @Produce json
// @Param input body organization.CreateOrganisationRequest true "Create organization"
// @Success 200 {object} Response[organization.CreateOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organizations [post]
func (server *Server) CreateOrganisationApi(ctx *gin.Context) {
	var req organization.CreateOrganisationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	organization, err := server.businessService.OrganizationService.CreateOrganization(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create organization"))
		return
	}

	res := SuccessResponse(organization, "Organization created successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary List all organizations
// @Description Get a list of all organizations with pagination and optional name search
// @Tags organizations
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param name query string false "Search by organization name"
// @Success 200 {object} Response[pagination.Response[organization.ListOrganisationsResponse]]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /organizations [get]
func (server *Server) ListOrganisationsApi(ctx *gin.Context) {
	var req organization.ListOrganisationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	organizations, err := server.businessService.OrganizationService.ListOrganizations(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(organizations, "Organizations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Get an organization
// @Description Get an organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} Response[organization.GetOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organizations/{id} [get]
func (server *Server) GetOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organizationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organization ID"))
		return
	}

	organization, err := server.businessService.OrganizationService.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get organization by ID"))
		return
	}

	res := SuccessResponse(organization, "Organization retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get organization counts
// @Description Get counts of locations, clients, and employees for an organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} Response[organization.GetOrganisationCountResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organizations/{id}/counts [get]
func (server *Server) GetOrganisationCountApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organizationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organization ID"))
		return
	}
	organization, err := server.businessService.OrganizationService.GetOrganizationCounts(ctx, organizationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get organization counts"))
		return
	}
	res := SuccessResponse(organization, "Organization counts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update an organization
// @Description Update an organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param input body organization.UpdateOrganisationRequest true "Update organization"
// @Success 200 {object} Response[organization.UpdateOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organizations/{id} [put]
func (server *Server) UpdateOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organizationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organization ID"))
		return
	}

	var req organization.UpdateOrganisationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	organization, err := server.businessService.OrganizationService.UpdateOrganization(ctx, organizationID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to update organization"))
		return
	}

	res := SuccessResponse(organization, "Organization updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete an organization
// @Description Delete an organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} Response[organization.DeleteOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organizations/{id} [delete]
func (server *Server) DeleteOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organizationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organization ID"))
		return
	}

	result, err := server.businessService.OrganizationService.DeleteOrganization(ctx, organizationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete organization"))
		return
	}

	res := SuccessResponse(result, "Organization deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary List all locations
// @Description Get a list of all locations
// @Tags organizations
// @Param id path int true "Organization ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param name query string false "Search by location name"
// @Produce json
// @Success 200 {object} Response[pagination.Response[organization.ListLocationsResponse]]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /organizations/{id}/locations [get]
func (server *Server) ListLocationsApi(ctx *gin.Context) {
	organizationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid organization ID")))
		return
	}

	var req organization.ListLocationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	locations, err := server.businessService.OrganizationService.ListOrgLocations(ctx, organizationID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(locations, "Locations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// ListAllLocationsApi lists all locations across all organizations
// @Summary List all locations
// @Description Get a list of all locations across all organizations with pagination and optional search
// @Tags organizations
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search by location name"
// @Success 200 {object} Response[pagination.Response[organization.ListLocationsResponse]]
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /locations [get]
func (server *Server) ListAllLocationsApi(ctx *gin.Context) {
	var req organization.ListAllLocationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	locations, err := server.businessService.OrganizationService.ListAllLocations(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to list all locations"))
		return
	}

	res := SuccessResponse(locations, "All locations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Create a location
// @Description Create a new location
// @Tags organizations
// @Param id path int true "Organization ID"
// @Accept json
// @Produce json
// @Param input body organization.CreateLocationRequest true "Create location"
// @Success 200 {object} Response[organization.CreateLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id}/locations [post]
func (server *Server) CreateLocationApi(ctx *gin.Context) {
	organizationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organization ID"))
		return
	}
	var req organization.CreateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	location, err := server.businessService.OrganizationService.CreateLocation(ctx, organizationID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create location"))
		return
	}
	res := SuccessResponse(location, "Location created successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update a location
// @Description Update a location
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param input body organization.UpdateLocationRequest true "Update location"
// @Success 200 {object} Response[organization.UpdateLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [put]
func (server *Server) UpdateLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid location ID"))
		return
	}

	var req organization.UpdateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	location, err := server.businessService.OrganizationService.UpdateLocation(ctx, locationID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to update location"))
		return
	}
	res := SuccessResponse(location, "Location updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete a location
// @Description Delete a location
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} Response[organization.DeleteLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [delete]
func (server *Server) DeleteLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid location ID"))
		return
	}
	result, err := server.businessService.OrganizationService.DeleteLocation(ctx, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete location"))
		return
	}
	res := SuccessResponse(result, "Location deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get a location
// @Description Get a location
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} Response[organization.GetLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [get]
func (server *Server) GetLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid location ID"))
		return
	}
	location, err := server.businessService.OrganizationService.GetLocationByID(ctx, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get location by ID"))
		return
	}
	res := SuccessResponse(location, "Location retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}
