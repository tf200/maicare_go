package organization

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *organizationService) CreateOrganization(ctx context.Context, req CreateOrganisationRequest) (*CreateOrganisationResponse, error) {
	organisation, err := s.Store.CreateOrganisation(ctx, db.CreateOrganisationParams{
		Name:                req.Name,
		Street:              req.Street,
		HouseNumber:         req.HouseNumber,
		HouseNumberAddition: req.HouseNumberAddition,
		PostalCode:          req.PostalCode,
		City:                req.City,
		Email:               req.Email,
		KvkNumber:           req.KvkNumber,
		BtwNumber:           req.BtwNumber,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateOrganisationApi", "Failed to create organisation", zap.Error(err))
		return nil, fmt.Errorf("failed to create organisation")
	}

	return &CreateOrganisationResponse{
		ID:                  organisation.ID,
		Name:                organisation.Name,
		Street:              organisation.Street,
		HouseNumber:         organisation.HouseNumber,
		HouseNumberAddition: organisation.HouseNumberAddition,
		PostalCode:          organisation.PostalCode,
		City:                organisation.City,
		Email:               organisation.Email,
		KvkNumber:           organisation.KvkNumber,
		BtwNumber:           organisation.BtwNumber,
	}, nil
}

func (s *organizationService) ListOrganizations(ctx *gin.Context, req ListOrganisationsRequest) (*pagination.Response[ListOrganisationsResponse], error) {
	params := req.GetParams()

	organisations, err := s.Store.ListOrganisationsPaginated(ctx, db.ListOrganisationsPaginatedParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListOrganisationsApi", "Failed to list organisations", zap.Error(err))
		return nil, fmt.Errorf("failed to list organisations")
	}

	if len(organisations) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListOrganisationsResponse{}, 0)
		return &pag, nil
	}

	totalCount := organisations[0].TotalCount

	response := make([]ListOrganisationsResponse, len(organisations))
	for i, org := range organisations {
		response[i] = ListOrganisationsResponse{
			ID:                  org.ID,
			Name:                org.Name,
			Street:              org.Street,
			HouseNumber:         org.HouseNumber,
			HouseNumberAddition: org.HouseNumberAddition,
			PostalCode:          org.PostalCode,
			City:                org.City,
			Email:               org.Email,
			KvkNumber:           org.KvkNumber,
			BtwNumber:           org.BtwNumber,
			LocationCount:       org.LocationCount,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, response, totalCount)
	return &pag, nil
}

func (s *organizationService) GetOrganizationByID(ctx context.Context, organizationID uuid.UUID) (*GetOrganisationResponse, error) {
	organization, err := s.Store.GetOrganisation(ctx, organizationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetOrganisationApi", "Failed to get organisation by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get organisation by ID")
	}

	return &GetOrganisationResponse{
		ID:                  organization.ID,
		Name:                organization.Name,
		Street:              organization.Street,
		HouseNumber:         organization.HouseNumber,
		HouseNumberAddition: organization.HouseNumberAddition,
		PostalCode:          organization.PostalCode,
		City:                organization.City,
		Email:               organization.Email,
		KvkNumber:           organization.KvkNumber,
		BtwNumber:           organization.BtwNumber,
		LocationCount:       organization.LocationCount,
		CreatedAt:           organization.CreatedAt.Time,
		UpdatedAt:           organization.UpdatedAt.Time,
	}, nil
}

func (s *organizationService) GetOrganizationCounts(ctx context.Context, organizationID uuid.UUID) (*GetOrganisationCountResponse, error) {
	counts, err := s.Store.GetOrganisationCounts(ctx, organizationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetOrganisationCountsApi", "Failed to get organisation counts", zap.Error(err))
		return nil, fmt.Errorf("failed to get organisation counts")
	}

	return &GetOrganisationCountResponse{
		OrganisationID:   counts.OrganisationID,
		OrganisationName: counts.OrganisationName,
		LocationCount:    counts.LocationCount,
		ClientCount:      counts.ClientCount,
		EmployeeCount:    counts.EmployeeCount,
	}, nil
}

func (s *organizationService) UpdateOrganization(ctx context.Context, organizationID uuid.UUID, req UpdateOrganisationRequest) (*GetOrganisationResponse, error) {
	organization, err := s.Store.UpdateOrganisation(ctx, db.UpdateOrganisationParams{
		ID:                  organizationID,
		Name:                req.Name,
		Street:              req.Street,
		HouseNumber:         req.HouseNumber,
		HouseNumberAddition: req.HouseNumberAddition,
		PostalCode:          req.PostalCode,
		City:                req.City,
		Email:               req.Email,
		KvkNumber:           req.KvkNumber,
		BtwNumber:           req.BtwNumber,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateOrganisationApi", "Failed to update organisation", zap.Error(err))
		return nil, fmt.Errorf("failed to update organisation")
	}

	return &GetOrganisationResponse{
		ID:                  organization.ID,
		Name:                organization.Name,
		Street:              organization.Street,
		HouseNumber:         organization.HouseNumber,
		HouseNumberAddition: organization.HouseNumberAddition,
		PostalCode:          organization.PostalCode,
		City:                organization.City,
		Email:               organization.Email,
		KvkNumber:           organization.KvkNumber,
		BtwNumber:           organization.BtwNumber,
	}, nil
}

func (s *organizationService) DeleteOrganization(ctx context.Context, organizationID uuid.UUID) (*DeleteOrganisationResponse, error) {
	_, err := s.Store.DeleteOrganisation(ctx, organizationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteOrganisationApi", "Failed to delete organisation", zap.Error(err))
		return nil, fmt.Errorf("failed to delete organisation")
	}

	return &DeleteOrganisationResponse{
		ID: organizationID,
	}, nil
}
