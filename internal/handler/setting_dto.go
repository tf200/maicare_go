package handler

import (
	"time"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

// ==================== Request DTOs ====================

type createDepartmentRequest struct {
	Name                     string     `json:"name" binding:"required"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
}

type updateDepartmentRequest struct {
	Name                     *string    `json:"name"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
}

type updateOrganizationProfileRequest struct {
	Name                  string  `json:"name" binding:"required"`
	DefaultTimezone       string  `json:"default_timezone" binding:"required"`
	Email                 *string `json:"email"`
	PhoneNumber           *string `json:"phone_number"`
	Website               *string `json:"website"`
	HQStreet              *string `json:"hq_street"`
	HQHouseNumber         *string `json:"hq_house_number"`
	HQHouseNumberAddition *string `json:"hq_house_number_addition"`
	HQPostalCode          *string `json:"hq_postal_code"`
	HQCity                *string `json:"hq_city"`
}

// ==================== Response DTOs ====================

type departmentResponse struct {
	ID                       uuid.UUID  `json:"id"`
	Name                     string     `json:"name"`
	Description              *string    `json:"description"`
	DepartmentHeadEmployeeID *uuid.UUID `json:"department_head_employee_id"`
	EmployeeCount            int64      `json:"employee_count"`
}

type organizationProfileResponse struct {
	Name                  string    `json:"name"`
	DefaultTimezone       string    `json:"default_timezone"`
	Email                 *string   `json:"email"`
	PhoneNumber           *string   `json:"phone_number"`
	Website               *string   `json:"website"`
	HQStreet              *string   `json:"hq_street"`
	HQHouseNumber         *string   `json:"hq_house_number"`
	HQHouseNumberAddition *string   `json:"hq_house_number_addition"`
	HQPostalCode          *string   `json:"hq_postal_code"`
	HQCity                *string   `json:"hq_city"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// ==================== Mappers ====================

func toDepartmentResponse(d domain.Department) departmentResponse {
	return departmentResponse{
		ID:                       d.ID,
		Name:                     d.Name,
		Description:              d.Description,
		DepartmentHeadEmployeeID: d.DepartmentHeadEmployeeID,
		EmployeeCount:            d.EmployeeCount,
	}
}

func toOrganizationProfileResponse(p *domain.OrganizationProfile) organizationProfileResponse {
	return organizationProfileResponse{
		Name:                  p.Name,
		DefaultTimezone:       p.DefaultTimezone,
		Email:                 p.Email,
		PhoneNumber:           p.PhoneNumber,
		Website:               p.Website,
		HQStreet:              p.HqStreet,
		HQHouseNumber:         p.HqHouseNumber,
		HQHouseNumberAddition: p.HqHouseNumberAddition,
		HQPostalCode:          p.HqPostalCode,
		HQCity:                p.HqCity,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}

func toCreateDepartmentParams(req createDepartmentRequest) domain.CreateDepartmentParams {
	return domain.CreateDepartmentParams{
		Name:                     req.Name,
		Description:              req.Description,
		DepartmentHeadEmployeeID: req.DepartmentHeadEmployeeID,
	}
}

func toUpdateDepartmentParams(id uuid.UUID, req updateDepartmentRequest) domain.UpdateDepartmentParams {
	return domain.UpdateDepartmentParams{
		ID:                       id,
		Name:                     req.Name,
		Description:              req.Description,
		DepartmentHeadEmployeeID: req.DepartmentHeadEmployeeID,
	}
}

func toUpdateOrganizationProfileParams(req updateOrganizationProfileRequest) domain.UpdateOrganizationProfileParams {
	return domain.UpdateOrganizationProfileParams{
		Name:                  req.Name,
		DefaultTimezone:       req.DefaultTimezone,
		Email:                 req.Email,
		PhoneNumber:           req.PhoneNumber,
		Website:               req.Website,
		HqStreet:              req.HQStreet,
		HqHouseNumber:         req.HQHouseNumber,
		HqHouseNumberAddition: req.HQHouseNumberAddition,
		HqPostalCode:          req.HQPostalCode,
		HqCity:                req.HQCity,
	}
}
