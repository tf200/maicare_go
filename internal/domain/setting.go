package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Status: scaffolded, unwired

// ==================== Domain Types ====================

type Department struct {
	ID                       uuid.UUID
	Name                     string
	Description              *string
	DepartmentHeadEmployeeID *uuid.UUID
	EmployeeCount            int64
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type OrganizationProfile struct {
	Name                  string
	DefaultTimezone       string
	Email                 *string
	PhoneNumber           *string
	Website               *string
	HqStreet              *string
	HqHouseNumber         *string
	HqHouseNumberAddition *string
	HqPostalCode          *string
	HqCity                *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// ==================== Parameter Structs ====================

type CreateDepartmentParams struct {
	Name                     string
	Description              *string
	DepartmentHeadEmployeeID *uuid.UUID
}

type UpdateDepartmentParams struct {
	ID                       uuid.UUID
	Name                     *string
	Description              *string
	DepartmentHeadEmployeeID *uuid.UUID
}

type UpdateOrganizationProfileParams struct {
	Name                  string
	DefaultTimezone       string
	Email                 *string
	PhoneNumber           *string
	Website               *string
	HqStreet              *string
	HqHouseNumber         *string
	HqHouseNumberAddition *string
	HqPostalCode          *string
	HqCity                *string
}

// ==================== Repository Interfaces ====================

type DepartmentRepository interface {
	List(ctx context.Context) ([]Department, error)
	Create(ctx context.Context, params CreateDepartmentParams) (*Department, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Department, error)
	Update(ctx context.Context, params UpdateDepartmentParams) (*Department, error)
}

type OrganizationProfileRepository interface {
	Get(ctx context.Context) (*OrganizationProfile, error)
	Update(ctx context.Context, params UpdateOrganizationProfileParams) (*OrganizationProfile, error)
}

// ==================== Service Interfaces ====================

type DepartmentService interface {
	ListDepartments(ctx context.Context) ([]Department, error)
	CreateDepartment(ctx context.Context, params CreateDepartmentParams) (*Department, error)
	UpdateDepartment(ctx context.Context, params UpdateDepartmentParams) (*Department, error)
}

type OrganizationProfileService interface {
	GetOrganizationProfile(ctx context.Context) (*OrganizationProfile, error)
	UpdateOrganizationProfile(ctx context.Context, params UpdateOrganizationProfileParams) (*OrganizationProfile, error)
}
