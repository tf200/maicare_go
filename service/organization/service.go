package organization

import (
	"context"
	"maicare_go/service/deps"
)

type OrganizationService interface {
	// Organization methods
	CreateOrganization(ctx context.Context, req CreateOrganisationRequest) (*CreateOrganisationResponse, error)
	ListOrganizations(ctx context.Context) ([]ListOrganisationsResponse, error)
	GetOrganizationByID(ctx context.Context, organizationID int64) (*GetOrganisationResponse, error)
	GetOrganizationCounts(ctx context.Context, organizationID int64) (*GetOrganisationCountResponse, error)
	UpdateOrganization(ctx context.Context, organizationID int64, req UpdateOrganisationRequest) (*GetOrganisationResponse, error)
	DeleteOrganization(ctx context.Context, organizationID int64) (*DeleteOrganisationResponse, error)
	// Location methods
	ListOrgLocations(ctx context.Context, organizationID int64) ([]ListLocationsResponse, error)
	ListAllLocations(ctx context.Context) ([]ListLocationsResponse, error)
	CreateLocation(ctx context.Context, organizationID int64, req CreateLocationRequest) (*CreateLocationResponse, error)
	UpdateLocation(ctx context.Context, locationID int64, req UpdateLocationRequest) (*UpdateLocationResponse, error)
	DeleteLocation(ctx context.Context, locationID int64) (*DeleteLocationResponse, error)
	GetLocationByID(ctx context.Context, locationID int64) (*GetLocationResponse, error)

	// Location Shift methods
	CreateShift(ctx context.Context, req *CreateShiftApiRequest, locationID int64) (*CreateShiftApiResponse, error)
	UpdateShift(ctx context.Context, shiftID int64, req *UpdateShiftApiRequest) (*UpdateShiftApiResponse, error)
	DeleteShift(ctx context.Context, shiftID int64) error
	ListShiftsByLocationID(ctx context.Context, locationID int64) ([]ListShiftsByLocationIDResponse, error)
	
}

type organizationService struct {
	*deps.ServiceDependencies
}

func NewOrganizationService(deps *deps.ServiceDependencies) OrganizationService {
	return &organizationService{
		ServiceDependencies: deps,
	}
}
