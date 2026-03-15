package settings

import (
	"context"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/service/deps"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type settingsService struct {
	*deps.ServiceDependencies
}

func NewSettingsService(deps *deps.ServiceDependencies) SettingsService {
	return &settingsService{ServiceDependencies: deps}
}

func (s *settingsService) ListDepartments(ctx *gin.Context) (*pagination.Response[ListDepartmentResponse], error) {
	depts, err := s.Store.ListDepartments(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ListDepartmentResponse, 0, len(depts))
	for _, d := range depts {
		out = append(out, ListDepartmentResponse{
			ID:                       d.ID,
			Name:                     d.Name,
			Description:              d.Description,
			DepartmentHeadEmployeeID: d.DepartmentHeadEmployeeID,
			EmployeeCount:            d.EmployeeCount,
		})
	}

	resp := &pagination.Response[ListDepartmentResponse]{
		Next:     nil,
		Previous: nil,
		Count:    int64(len(out)),
		PageSize: int32(len(out)),
		Results:  out,
	}
	return resp, nil
}

func (s *settingsService) CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*CreateDepartmentResponse, error) {
	dept, err := s.Store.CreateDepartment(ctx, db.CreateDepartmentParams{
		Name:                     req.Name,
		Description:              util.OtpString(req.Description),
		DepartmentHeadEmployeeID: req.DepartmentHeadEmployeeID,
	})
	if err != nil {
		return nil, err
	}

	return &CreateDepartmentResponse{
		ID:                       dept.ID,
		Name:                     dept.Name,
		Description:              dept.Description,
		DepartmentHeadEmployeeID: dept.DepartmentHeadEmployeeID,
	}, nil
}

func (s *settingsService) UpdateDepartment(ctx context.Context, departmentID uuid.UUID, req UpdateDepartmentRequest) (*UpdateDepartmentResponse, error) {
	name := util.OtpString(req.Name)
	if req.Name != nil && name == nil {
		return nil, fmt.Errorf("name cannot be empty")
	}

	dept, err := s.Store.UpdateDepartment(ctx, db.UpdateDepartmentParams{
		ID:                       departmentID,
		Name:                     name,
		Description:              util.OtpString(req.Description),
		DepartmentHeadEmployeeID: req.DepartmentHeadEmployeeID,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateDepartmentResponse{
		ID:                       dept.ID,
		Name:                     dept.Name,
		Description:              dept.Description,
		DepartmentHeadEmployeeID: dept.DepartmentHeadEmployeeID,
	}, nil
}

func (s *settingsService) GetOrganizationProfile(ctx context.Context) (*GetOrganizationProfileResponse, error) {
	profile, err := s.Store.GetAppOrganizationProfile(ctx)
	if err != nil {
		return nil, err
	}

	return mapOrganizationProfile(profile), nil
}

func (s *settingsService) UpdateOrganizationProfile(ctx context.Context, req UpdateOrganizationProfileRequest) (*GetOrganizationProfileResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	defaultTimezone := strings.TrimSpace(req.DefaultTimezone)
	if defaultTimezone == "" {
		return nil, fmt.Errorf("default_timezone is required")
	}

	if _, err := time.LoadLocation(defaultTimezone); err != nil {
		return nil, fmt.Errorf("invalid default_timezone: %w", err)
	}

	email := util.OtpString(req.Email)
	if email != nil {
		if _, err := mail.ParseAddress(*email); err != nil {
			return nil, fmt.Errorf("invalid email")
		}
	}

	website := util.OtpString(req.Website)
	if website != nil {
		parsed, err := url.ParseRequestURI(*website)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return nil, fmt.Errorf("invalid website")
		}
	}

	profile, err := s.Store.UpdateAppOrganizationProfile(ctx, db.UpdateAppOrganizationProfileParams{
		Name:                  name,
		DefaultTimezone:       defaultTimezone,
		Email:                 email,
		PhoneNumber:           util.OtpString(req.PhoneNumber),
		Website:               website,
		HqStreet:              util.OtpString(req.HQStreet),
		HqHouseNumber:         util.OtpString(req.HQHouseNumber),
		HqHouseNumberAddition: util.OtpString(req.HQHouseNumberAddition),
		HqPostalCode:          util.OtpString(req.HQPostalCode),
		HqCity:                util.OtpString(req.HQCity),
	})
	if err != nil {
		return nil, err
	}

	return mapOrganizationProfile(profile), nil
}

func mapOrganizationProfile(profile db.AppOrganizationProfile) *GetOrganizationProfileResponse {
	return &GetOrganizationProfileResponse{
		Name:                  profile.Name,
		DefaultTimezone:       profile.DefaultTimezone,
		Email:                 profile.Email,
		PhoneNumber:           profile.PhoneNumber,
		Website:               profile.Website,
		HQStreet:              profile.HqStreet,
		HQHouseNumber:         profile.HqHouseNumber,
		HQHouseNumberAddition: profile.HqHouseNumberAddition,
		HQPostalCode:          profile.HqPostalCode,
		HQCity:                profile.HqCity,
		CreatedAt:             profile.CreatedAt.Time,
		UpdatedAt:             profile.UpdatedAt.Time,
	}
}
