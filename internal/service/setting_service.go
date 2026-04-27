package service

import (
	"context"

	"maicare_go/internal/domain"
)

// ==================== Department Service ====================

type DepartmentService struct {
	repo domain.DepartmentRepository
}

func NewDepartmentService(repo domain.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) ListDepartments(ctx context.Context) ([]domain.Department, error) {
	return s.repo.List(ctx)
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, params domain.CreateDepartmentParams) (*domain.Department, error) {
	return s.repo.Create(ctx, params)
}

func (s *DepartmentService) UpdateDepartment(ctx context.Context, params domain.UpdateDepartmentParams) (*domain.Department, error) {
	return s.repo.Update(ctx, params)
}

// ==================== Organization Profile Service ====================

type OrganizationProfileService struct {
	repo domain.OrganizationProfileRepository
}

func NewOrganizationProfileService(repo domain.OrganizationProfileRepository) *OrganizationProfileService {
	return &OrganizationProfileService{repo: repo}
}

func (s *OrganizationProfileService) GetOrganizationProfile(ctx context.Context) (*domain.OrganizationProfile, error) {
	return s.repo.Get(ctx)
}

func (s *OrganizationProfileService) UpdateOrganizationProfile(ctx context.Context, params domain.UpdateOrganizationProfileParams) (*domain.OrganizationProfile, error) {
	return s.repo.Update(ctx, params)
}
