package employees

import (
	"context"
	"fmt"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidCredentials  = fmt.Errorf("invalid credentials")
	ErrUserNotFound        = fmt.Errorf("user not found")
	ErrSessionNotFound     = fmt.Errorf("session not found")
	ErrUnauthorized        = fmt.Errorf("unauthorized")
	ErrTwoFaAlreadyEnabled = fmt.Errorf("two-factor authentication already enabled")
	ErrTwoFARequired       = fmt.Errorf("two-factor authentication required")
	ErrInvalidTwoFACode    = fmt.Errorf("invalid two-factor authentication code")
)

// AuthService Interface and implementation
//
//go:generate mockgen -source=service.go -destination=../mocks/mock_employee_service.go -package=mocks
type EmployeeService interface {
	CreateEmployee(req CreateEmployeeProfileRequest, ctx context.Context) (*CreateEmployeeProfileResponse, error)
	ListEmployees(req ListEmployeeRequest, ctx *gin.Context) (*pagination.Response[ListEmployeeResponse], error)
	UpdateEmployeeIsSubcontractor(req UpdateEmployeeIsSubcontractorRequest, employeeID int64, ctx context.Context) (*UpdateEmployeeIsSubcontractorResponse, error)
}

type employeeService struct {
	*deps.ServiceDependencies
}

func NewEmployeeService(deps *deps.ServiceDependencies) EmployeeService {
	return &employeeService{
		ServiceDependencies: deps,
	}
}
