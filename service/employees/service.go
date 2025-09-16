package employees

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/service/deps"
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
	CreateEmployee(req CreateEmployeeRequest, ctx context.Context) (*CreateEmployeeResult, error)
	ListEmployees(req ListEmployeesRequest, ctx context.Context) ([]db.ListEmployeeProfileRow, *int64, error)
	UpdateEmployeeIsSubcontractor(req UpdateEmployeeIsSubcontractorRequest, ctx context.Context) (*UpdateEmployeeIsSubcontractorResult, error)
}

type employeeService struct {
	*deps.ServiceDependencies
}

func NewEmployeeService(deps *deps.ServiceDependencies) EmployeeService {
	return &employeeService{
		ServiceDependencies: deps,
	}
}
