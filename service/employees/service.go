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
	GetEmployeeProfile(userID int64, ctx context.Context) (*GetEmployeeProfileResponse, error)
	GetEmployeeProfileByID(employeeID, currentUserID int64, ctx context.Context) (*GetEmployeeProfileByIDResponse, error)
	UpdateEmployeeProfile(req UpdateEmployeeProfileRequest, employeeID int64, ctx context.Context) (*UpdateEmployeeProfileResponse, error)
	SetEmployeeProfilePicture(req SetEmployeeProfilePictureRequest, employeeID int64, ctx context.Context) (*SetEmployeeProfilePictureResponse, error)
	GetEmployeeCounts(ctx context.Context) (*GetEmployeeCountsResponse, error)
	SearchEmployeesByNameOrEmail(req SearchEmployeesByNameOrEmailRequest, ctx context.Context) ([]SearchEmployeesByNameOrEmailResponse, error)

	// Contract methods
	AddEmployeeContractDetails(req AddEmployeeContractDetailsRequest, employeeID int64, ctx context.Context) (*AddEmployeeContractDetailsResponse, error)
	GetEmployeeContractDetails(employeeID int64, ctx context.Context) (*GetEmployeeContractDetailsResponse, error)

	// Education methods
	AddEducationToEmployeeProfile(req AddEducationToEmployeeProfileRequest, employeeID int64, ctx context.Context) (*AddEducationToEmployeeProfileResponse, error)
	ListEmployeeEducation(employeeID int64, ctx context.Context) ([]ListEmployeeEducationResponse, error)
	UpdateEmployeeEducation(req UpdateEmployeeEducationRequest, educationID int64, ctx context.Context) (*UpdateEmployeeEducationResponse, error)
	DeleteEmployeeEducation(educationID int64, ctx context.Context) (*DeleteEmployeeEducationResponse, error)

	// Experience methods
	AddEmployeeExperience(req AddEmployeeExperienceRequest, employeeID int64, ctx context.Context) (*AddEmployeeExperienceResponse, error)
	ListEmployeeExperience(employeeID int64, ctx context.Context) ([]ListEmployeeExperienceResponse, error)
	UpdateEmployeeExperience(req UpdateEmployeeExperienceRequest, experienceID int64, ctx context.Context) (*UpdateEmployeeExperienceResponse, error)
	DeleteEmployeeExperience(experienceID int64, ctx context.Context) (*DeleteEmployeeExperienceResponse, error)

	// Certification methods
	AddEmployeeCertification(req AddEmployeeCertificationRequest, employeeID int64, ctx context.Context) (*AddEmployeeCertificationResponse, error)
	ListEmployeeCertification(employeeID int64, ctx context.Context) ([]ListEmployeeCertificationResponse, error)
	UpdateEmployeeCertification(req UpdateEmployeeCertificationRequest, certificationID int64, ctx context.Context) (*UpdateEmployeeCertificationResponse, error)
	DeleteEmployeeCertification(certificationID int64, ctx context.Context) (*DeleteEmployeeCertificationResponse, error)
}

type employeeService struct {
	*deps.ServiceDependencies
}

func NewEmployeeService(deps *deps.ServiceDependencies) EmployeeService {
	return &employeeService{
		ServiceDependencies: deps,
	}
}
