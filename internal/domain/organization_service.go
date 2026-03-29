package domain

import "context"

type OrganizationService interface {
	ListOrganizations(ctx context.Context, params ListOrganizationsParams) (*OrganizationPage, error)
	ListOrganizationLocations(ctx context.Context, params ListOrganizationLocationsParams) (*OrganizationLocationPage, error)
}
