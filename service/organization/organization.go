package organization

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *organizationService) CreateOrganization(ctx context.Context, req CreateOrganisationRequest) (*CreateOrganisationResponse, error) {
	organisation, err := s.Store.CreateOrganisation(ctx, db.CreateOrganisationParams{
		Name:       req.Name,
		Address:    req.Address,
		PostalCode: req.PostalCode,
		City:       req.City,
		Email:      req.Email,
		KvkNumber:  req.KvkNumber,
		BtwNumber:  req.BtwNumber,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateOrganisationApi", "Failed to create organisation", zap.Error(err))
		return nil, fmt.Errorf("failed to create organisation")
	}

	return &CreateOrganisationResponse{
		ID:         organisation.ID,
		Name:       organisation.Name,
		Address:    organisation.Address,
		PostalCode: organisation.PostalCode,
		City:       organisation.City,
		Email:      organisation.Email,
		KvkNumber:  organisation.KvkNumber,
		BtwNumber:  organisation.BtwNumber,
	}, nil
}

func (s *organizationService) ListOrganizations(ctx context.Context) ([]ListOrganisationsResponse, error) {
	organisations, err := s.Store.ListOrganisations(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListOrganisationsApi", "Failed to list organisations", zap.Error(err))
		return nil, fmt.Errorf("failed to list organisations")
	}

	response := []ListOrganisationsResponse{}
	for _, org := range organisations {
		response = append(response, ListOrganisationsResponse{
			ID:            org.ID,
			Name:          org.Name,
			Address:       org.Address,
			PostalCode:    org.PostalCode,
			City:          org.City,
			Email:         org.Email,
			KvkNumber:     org.KvkNumber,
			BtwNumber:     org.BtwNumber,
			LocationCount: org.LocationCount,
		})
	}

	return response, nil
}

func (s *organizationService) GetOrganizationByID(ctx context.Context, organizationID int64) (*GetOrganisationResponse, error) {
	organization, err := s.Store.GetOrganisation(ctx, organizationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetOrganisationApi", "Failed to get organisation by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get organisation by ID")
	}

	return &GetOrganisationResponse{
		ID:            organization.ID,
		Name:          organization.Name,
		Address:       organization.Address,
		PostalCode:    organization.PostalCode,
		City:          organization.City,
		Email:         organization.Email,
		KvkNumber:     organization.KvkNumber,
		BtwNumber:     organization.BtwNumber,
		LocationCount: organization.LocationCount,
		CreatedAt:     organization.CreatedAt.Time,
		UpdatedAt:     organization.UpdatedAt.Time,
	}, nil
}

func (s *organizationService) GetOrganizationCounts(ctx context.Context, organizationID int64) (*GetOrganisationCountResponse, error) {
	counts, err := s.Store.GetOrganisationCounts(ctx, organizationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetOrganisationCountsApi", "Failed to get organisation counts", zap.Error(err))
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

func (s *organizationService) UpdateOrganization(ctx context.Context, organizationID int64, req UpdateOrganisationRequest) (*GetOrganisationResponse, error) {
	organization, err := s.Store.UpdateOrganisation(ctx, db.UpdateOrganisationParams{
		ID:         organizationID,
		Name:       req.Name,
		Address:    req.Address,
		PostalCode: req.PostalCode,
		City:       req.City,
		Email:      req.Email,
		KvkNumber:  req.KvkNumber,
		BtwNumber:  req.BtwNumber,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateOrganisationApi", "Failed to update organisation", zap.Error(err))
		return nil, fmt.Errorf("failed to update organisation")
	}

	return &GetOrganisationResponse{
		ID:         organization.ID,
		Name:       organization.Name,
		Address:    organization.Address,
		PostalCode: organization.PostalCode,
		City:       organization.City,
		Email:      organization.Email,
		KvkNumber:  organization.KvkNumber,
		BtwNumber:  organization.BtwNumber,
	}, nil
}

func (s *organizationService) DeleteOrganization(ctx context.Context, organizationID int64) (*DeleteOrganisationResponse, error) {
	_, err := s.Store.DeleteOrganisation(ctx, organizationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteOrganisationApi", "Failed to delete organisation", zap.Error(err))
		return nil, fmt.Errorf("failed to delete organisation")
	}

	return &DeleteOrganisationResponse{
		ID: organizationID,
	}, nil
}
