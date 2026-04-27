package repository

import (
	"context"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/util"
)

// ==================== Department Repository ====================

type DepartmentRepository struct {
	queries *db.Queries
}

func NewDepartmentRepository(queries *db.Queries) *DepartmentRepository {
	return &DepartmentRepository{queries: queries}
}

func (r *DepartmentRepository) List(ctx context.Context) ([]domain.Department, error) {
	rows, err := r.queries.ListDepartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list departments: %w", err)
	}

	result := make([]domain.Department, len(rows))
	for i, row := range rows {
		result[i] = domain.Department{
			ID:                       row.ID,
			Name:                     row.Name,
			Description:              row.Description,
			DepartmentHeadEmployeeID: row.DepartmentHeadEmployeeID,
			EmployeeCount:            row.EmployeeCount,
			CreatedAt:                row.CreatedAt.Time,
			UpdatedAt:                row.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (r *DepartmentRepository) Create(ctx context.Context, params domain.CreateDepartmentParams) (*domain.Department, error) {
	dept, err := r.queries.CreateDepartment(ctx, db.CreateDepartmentParams{
		Name:                     params.Name,
		Description:              util.OtpString(params.Description),
		DepartmentHeadEmployeeID: params.DepartmentHeadEmployeeID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return &domain.Department{
		ID:                       dept.ID,
		Name:                     dept.Name,
		Description:              dept.Description,
		DepartmentHeadEmployeeID: dept.DepartmentHeadEmployeeID,
		CreatedAt:                dept.CreatedAt.Time,
		UpdatedAt:                dept.UpdatedAt.Time,
	}, nil
}

func (r *DepartmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	dept, err := r.queries.GetDepartment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return &domain.Department{
		ID:                       dept.ID,
		Name:                     dept.Name,
		Description:              dept.Description,
		DepartmentHeadEmployeeID: dept.DepartmentHeadEmployeeID,
		CreatedAt:                dept.CreatedAt.Time,
		UpdatedAt:                dept.UpdatedAt.Time,
	}, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, params domain.UpdateDepartmentParams) (*domain.Department, error) {
	name := util.OtpString(params.Name)
	if params.Name != nil && name == nil {
		return nil, fmt.Errorf("name cannot be empty")
	}

	dept, err := r.queries.UpdateDepartment(ctx, db.UpdateDepartmentParams{
		ID:                       params.ID,
		Name:                     name,
		Description:              util.OtpString(params.Description),
		DepartmentHeadEmployeeID: params.DepartmentHeadEmployeeID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return &domain.Department{
		ID:                       dept.ID,
		Name:                     dept.Name,
		Description:              dept.Description,
		DepartmentHeadEmployeeID: dept.DepartmentHeadEmployeeID,
		CreatedAt:                dept.CreatedAt.Time,
		UpdatedAt:                dept.UpdatedAt.Time,
	}, nil
}

// ==================== Organization Profile Repository ====================

type OrganizationProfileRepository struct {
	queries *db.Queries
}

func NewOrganizationProfileRepository(queries *db.Queries) *OrganizationProfileRepository {
	return &OrganizationProfileRepository{queries: queries}
}

func (r *OrganizationProfileRepository) Get(ctx context.Context) (*domain.OrganizationProfile, error) {
	profile, err := r.queries.GetAppOrganizationProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization profile: %w", err)
	}

	return mapOrganizationProfileDBToDomain(profile), nil
}

func (r *OrganizationProfileRepository) Update(ctx context.Context, params domain.UpdateOrganizationProfileParams) (*domain.OrganizationProfile, error) {
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	defaultTimezone := strings.TrimSpace(params.DefaultTimezone)
	if defaultTimezone == "" {
		return nil, fmt.Errorf("default_timezone is required")
	}

	if _, err := time.LoadLocation(defaultTimezone); err != nil {
		return nil, fmt.Errorf("invalid default_timezone: %w", err)
	}

	email := util.OtpString(params.Email)
	if email != nil {
		if _, err := mail.ParseAddress(*email); err != nil {
			return nil, fmt.Errorf("invalid email")
		}
	}

	website := util.OtpString(params.Website)
	if website != nil {
		parsed, err := url.ParseRequestURI(*website)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return nil, fmt.Errorf("invalid website")
		}
	}

	profile, err := r.queries.UpdateAppOrganizationProfile(ctx, db.UpdateAppOrganizationProfileParams{
		Name:                  name,
		DefaultTimezone:       defaultTimezone,
		Email:                 email,
		PhoneNumber:           util.OtpString(params.PhoneNumber),
		Website:               website,
		HqStreet:              util.OtpString(params.HqStreet),
		HqHouseNumber:         util.OtpString(params.HqHouseNumber),
		HqHouseNumberAddition: util.OtpString(params.HqHouseNumberAddition),
		HqPostalCode:          util.OtpString(params.HqPostalCode),
		HqCity:                util.OtpString(params.HqCity),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update organization profile: %w", err)
	}

	return mapOrganizationProfileDBToDomain(profile), nil
}

// ==================== Helpers ====================

func mapOrganizationProfileDBToDomain(p db.AppOrganizationProfile) *domain.OrganizationProfile {
	return &domain.OrganizationProfile{
		Name:                  p.Name,
		DefaultTimezone:       p.DefaultTimezone,
		Email:                 p.Email,
		PhoneNumber:           p.PhoneNumber,
		Website:               p.Website,
		HqStreet:              p.HqStreet,
		HqHouseNumber:         p.HqHouseNumber,
		HqHouseNumberAddition: p.HqHouseNumberAddition,
		HqPostalCode:          p.HqPostalCode,
		HqCity:                p.HqCity,
		CreatedAt:             p.CreatedAt.Time,
		UpdatedAt:             p.UpdatedAt.Time,
	}
}
