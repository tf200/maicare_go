package service

import (
	"context"

	"maicare_go/internal/domain"
)

type OrganizationService struct {
	repository domain.OrganizationRepository
}

func NewOrganizationService(repository domain.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repository: repository}
}

func (s *OrganizationService) ListOrganizations(ctx context.Context, params domain.ListOrganizationsParams) (*domain.OrganizationPage, error) {
	return s.repository.ListOrganizations(ctx, params)
}

func (s *OrganizationService) ListOrganizationLocations(ctx context.Context, params domain.ListOrganizationLocationsParams) (*domain.OrganizationLocationPage, error) {
	return s.repository.ListOrganizationLocations(ctx, params)
}
