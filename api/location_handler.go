package api

import (
	"fmt"
	"net/http"
	"strconv"

	"maicare_go/service/organization"

	"github.com/gin-gonic/gin"
)

// @Summary Create an organisation
// @Description Create a new organisation
// @Tags organisations
// @Accept json
// @Produce json
// @Param input body organization.CreateOrganisationRequest true "Create organisation"
// @Success 200 {object} Response[organization.CreateOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations [post]
func (server *Server) CreateOrganisationApi(ctx *gin.Context) {
	var req organization.CreateOrganisationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	organisation, err := server.businessService.OrganizationService.CreateOrganization(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create organisation"))
		return
	}

	res := SuccessResponse(organisation, "Organisation created successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary List all organisations
// @Description Get a list of all organisations
// @Tags organisations
// @Produce json
// @Success 200 {object} Response[[]organization.ListOrganisationsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /organisations [get]
func (server *Server) ListOrganisationsApi(ctx *gin.Context) {
	organisations, err := server.businessService.OrganizationService.ListOrganizations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to list organisations"))
		return
	}

	res := SuccessResponse(organisations, "Organisations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Get an organisation
// @Description Get an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Success 200 {object} Response[organization.GetOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id} [get]
func (server *Server) GetOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}

	organisation, err := server.businessService.OrganizationService.GetOrganizationByID(ctx, organisationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get organisation by ID"))
		return
	}

	res := SuccessResponse(organisation, "Organisation retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get organisation counts
// @Description Get counts of locations, clients, and employees for an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Success 200 {object} Response[organization.GetOrganisationCountResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id}/counts [get]
func (server *Server) GetOrganisationCountApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}
	organisation, err := server.businessService.OrganizationService.GetOrganizationCounts(ctx, organisationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get organisation counts"))
		return
	}
	res := SuccessResponse(organisation, "Organisation counts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update an organisation
// @Description Update an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Param input body organization.UpdateOrganisationRequest true "Update organisation"
// @Success 200 {object} Response[organization.UpdateOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id} [put]
func (server *Server) UpdateOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}

	var req organization.UpdateOrganisationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	organisation, err := server.businessService.OrganizationService.UpdateOrganization(ctx, organisationID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to update organisation"))
		return
	}

	res := SuccessResponse(organisation, "Organisation updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete an organisation
// @Description Delete an organisation by ID
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Organisation ID"
// @Success 200 {object} Response[organization.DeleteOrganisationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id} [delete]
func (server *Server) DeleteOrganisationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	organisationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}

	result, err := server.businessService.OrganizationService.DeleteOrganization(ctx, organisationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete organisation"))
		return
	}

	res := SuccessResponse(result, "Organisation deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary List all locations
// @Description Get a list of all locations
// @Tags organisations
// @Param id path int true "Organisation ID"
// @Produce json
// @Success 200 {object} Response[[]organization.ListLocationsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /organisations/{id}/locations [get]
func (server *Server) ListLocationsApi(ctx *gin.Context) {
	organizationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}
	locations, err := server.businessService.OrganizationService.ListOrgLocations(ctx, organizationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to list locations"))
		return
	}

	res := SuccessResponse(locations, "Locations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// ListAllLocationsApi lists all locations across all organisations
// @Summary List all locations
// @Description Get a list of all locations across all organisations
// @Tags organisations
// @Produce json
// @Success 200 {object} Response[[]organization.ListLocationsResponse]
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /locations [get]
func (server *Server) ListAllLocationsApi(ctx *gin.Context) {
	locations, err := server.businessService.OrganizationService.ListAllLocations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to list all locations"))
		return
	}

	res := SuccessResponse(locations, "All locations retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Create a location
// @Description Create a new location
// @Tags organisations
// @Param id path int true "Organisation ID"
// @Accept json
// @Produce json
// @Param input body organization.CreateLocationRequest true "Create location"
// @Success 200 {object} Response[organization.CreateLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /organisations/{id}/locations [post]
func (server *Server) CreateLocationApi(ctx *gin.Context) {
	organisationID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid organisation ID"))
		return
	}
	var req organization.CreateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	location, err := server.businessService.OrganizationService.CreateLocation(ctx, organisationID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create location"))
		return
	}
	res := SuccessResponse(location, "Location created successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update a location
// @Description Update a location
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param input body organization.UpdateLocationRequest true "Update location"
// @Success 200 {object} Response[organization.UpdateLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [put]
func (server *Server) UpdateLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := strconv.ParseInt(id, 10, 64)
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
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} Response[organization.DeleteLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [delete]
func (server *Server) DeleteLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := strconv.ParseInt(id, 10, 64)
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
// @Tags organisations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} Response[organization.GetLocationResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /locations/{id} [get]
func (server *Server) GetLocationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := strconv.ParseInt(id, 10, 64)
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
