package organization

import (
	"context"

	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationService interface {
	// Organization methods
	CreateOrganization(ctx context.Context, req CreateOrganisationRequest) (*CreateOrganisationResponse, error)
	ListOrganizations(ctx *gin.Context, req ListOrganisationsRequest) (*pagination.Response[ListOrganisationsResponse], error)
	GetOrganizationByID(ctx context.Context, organizationID uuid.UUID) (*GetOrganisationResponse, error)
	GetOrganizationCounts(ctx context.Context, organizationID uuid.UUID) (*GetOrganisationCountResponse, error)
	UpdateOrganization(ctx context.Context, organizationID uuid.UUID, req UpdateOrganisationRequest) (*GetOrganisationResponse, error)
	DeleteOrganization(ctx context.Context, organizationID uuid.UUID) (*DeleteOrganisationResponse, error)
	// Location methods
	ListOrgLocations(ctx *gin.Context, organizationID uuid.UUID, req ListLocationsRequest) (*pagination.Response[ListLocationsResponse], error)
	ListAllLocations(ctx *gin.Context, req ListAllLocationsRequest) (*pagination.Response[ListLocationsResponse], error)
	CreateLocation(ctx context.Context, organizationID uuid.UUID, req CreateLocationRequest) (*CreateLocationResponse, error)
	UpdateLocation(ctx context.Context, locationID uuid.UUID, req UpdateLocationRequest) (*UpdateLocationResponse, error)
	DeleteLocation(ctx context.Context, locationID uuid.UUID) (*DeleteLocationResponse, error)
	GetLocationByID(ctx context.Context, locationID uuid.UUID) (*GetLocationResponse, error)

	// Location Shift methods
	CreateShift(ctx context.Context, req *CreateShiftApiRequest, locationID uuid.UUID) (*CreateShiftApiResponse, error)
	UpdateShift(ctx context.Context, shiftID uuid.UUID, req *UpdateShiftApiRequest) (*UpdateShiftApiResponse, error)
	DeleteShift(ctx context.Context, shiftID uuid.UUID) error
	ListShiftsByLocationID(ctx context.Context, locationID uuid.UUID) ([]ListShiftsByLocationIDResponse, error)
}

type organizationService struct {
	*deps.ServiceDependencies
}

func NewOrganizationService(deps *deps.ServiceDependencies) OrganizationService {
	return &organizationService{
		ServiceDependencies: deps,
	}
}
