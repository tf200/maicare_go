package employees

import (
	"context"
	"fmt"
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
type EmployeeService interface {
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
