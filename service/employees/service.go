package employees

import (
	"context"
	"fmt"

	"maicare_go/async/aclient"
	"maicare_go/pagination"
	"maicare_go/service/appointment"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	UpdateEmployeeIsSubcontractor(req UpdateEmployeeIsSubcontractorRequest, employeeID uuid.UUID, ctx context.Context) (*UpdateEmployeeIsSubcontractorResponse, error)
	GetEmployeeProfile(userID uuid.UUID, ctx context.Context) (*GetEmployeeProfileResponse, error)
	GetEmployeeProfileDetails(userID uuid.UUID, ctx context.Context) (*GetEmployeeProfileDetailsResponse, error)
	GetEmployeeProfileByID(employeeID, currentUserID uuid.UUID, ctx context.Context) (*GetEmployeeProfileByIDResponse, error)
	UpdateEmployeeProfile(req UpdateEmployeeProfileRequest, employeeID uuid.UUID, ctx context.Context) (*UpdateEmployeeProfileResponse, error)
	SetEmployeeProfilePicture(req SetEmployeeProfilePictureRequest, employeeID uuid.UUID, ctx context.Context) (*SetEmployeeProfilePictureResponse, error)
	GetEmployeeCounts(ctx context.Context) (*GetEmployeeCountsResponse, error)
	SearchEmployeesByNameOrEmail(req SearchEmployeesByNameOrEmailRequest, ctx context.Context) ([]SearchEmployeesByNameOrEmailResponse, error)

	// Contract methods
	AddEmployeeContractDetails(req AddEmployeeContractDetailsRequest, employeeID uuid.UUID, ctx context.Context) (*AddEmployeeContractDetailsResponse, error)
	GetEmployeeContractDetails(employeeID uuid.UUID, ctx context.Context) (*GetEmployeeContractDetailsResponse, error)

	// Education methods
	AddEducationToEmployeeProfile(req AddEducationToEmployeeProfileRequest, employeeID uuid.UUID, ctx context.Context) (*AddEducationToEmployeeProfileResponse, error)
	ListEmployeeEducation(employeeID uuid.UUID, ctx context.Context) ([]ListEmployeeEducationResponse, error)
	UpdateEmployeeEducation(req UpdateEmployeeEducationRequest, educationID uuid.UUID, ctx context.Context) (*UpdateEmployeeEducationResponse, error)
	DeleteEmployeeEducation(educationID uuid.UUID, ctx context.Context) (*DeleteEmployeeEducationResponse, error)

	// Experience methods
	AddEmployeeExperience(req AddEmployeeExperienceRequest, employeeID uuid.UUID, ctx context.Context) (*AddEmployeeExperienceResponse, error)
	ListEmployeeExperience(employeeID uuid.UUID, ctx context.Context) ([]ListEmployeeExperienceResponse, error)
	UpdateEmployeeExperience(req UpdateEmployeeExperienceRequest, experienceID uuid.UUID, ctx context.Context) (*UpdateEmployeeExperienceResponse, error)
	DeleteEmployeeExperience(experienceID uuid.UUID, ctx context.Context) (*DeleteEmployeeExperienceResponse, error)

	// Certification methods
	AddEmployeeCertification(req AddEmployeeCertificationRequest, employeeID uuid.UUID, ctx context.Context) (*AddEmployeeCertificationResponse, error)
	ListEmployeeCertification(employeeID uuid.UUID, ctx context.Context) ([]ListEmployeeCertificationResponse, error)
	UpdateEmployeeCertification(req UpdateEmployeeCertificationRequest, certificationID uuid.UUID, ctx context.Context) (*UpdateEmployeeCertificationResponse, error)
	DeleteEmployeeCertification(certificationID uuid.UUID, ctx context.Context) (*DeleteEmployeeCertificationResponse, error)

	// Working hours methods
	ListWorkingHours(ctx context.Context, employeeID uuid.UUID, req *ListWorkingHoursRequest) (*ListWorkingHoursResponse, error)
	GetMyScheduleTimeline(ctx context.Context, employeeID uuid.UUID, req *GetMyScheduleTimelineRequest) ([]GetMyScheduleTimelineDayResponse, error)
}

type employeeService struct {
	*deps.ServiceDependencies
	asynqClient        aclient.AsynqClientInterface
	appointmentService appointment.AppointmentService
}

func NewEmployeeService(deps *deps.ServiceDependencies, asynqClient aclient.AsynqClientInterface, appointmentService appointment.AppointmentService) EmployeeService {
	return &employeeService{
		ServiceDependencies: deps,
		asynqClient:         asynqClient,
		appointmentService:  appointmentService,
	}
}
