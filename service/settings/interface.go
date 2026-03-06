package settings

import (
	"context"

	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsService interface {
	ListDepartments(ctx *gin.Context) (*pagination.Response[ListDepartmentResponse], error)
	CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*CreateDepartmentResponse, error)
	UpdateDepartment(ctx context.Context, departmentID uuid.UUID, req UpdateDepartmentRequest) (*UpdateDepartmentResponse, error)
	GetOrganizationProfile(ctx context.Context) (*GetOrganizationProfileResponse, error)
	UpdateOrganizationProfile(ctx context.Context, req UpdateOrganizationProfileRequest) (*GetOrganizationProfileResponse, error)
	ListRoles(ctx context.Context) ([]ListRolesApiResponse, error)
	ListAllPermissions(ctx context.Context) ([]PermissionGroupResponse, error)
	ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]ListAllRolePermissionsApiResponse, error)
	AssignRoleToEmployee(ctx context.Context, employeeID uuid.UUID, req *AssignRoleToEmployeeParams) (*AssignRoleToEmployeeApiResponse, error)
	ListUserRolesAndPermissionsApi(ctx context.Context, employeeID uuid.UUID) (*ListUserRolesAndPermissionsApiResponse, error)
	ReplaceUserPermissionOverrides(ctx context.Context, employeeID uuid.UUID, req *ReplaceUserPermissionOverridesRequest) (*ReplaceUserPermissionOverridesResponse, error)
	AddPermissionsToRole(ctx context.Context, roleID uuid.UUID, req *AddPermissionsToRoleRequest) (*AddPermissionsToRoleResponse, error)
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error)
}
