package handbook

import (
	"context"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

type HandbookService interface {
	// Employee self-service
	GetMyActiveHandbook(ctx context.Context, employeeID uuid.UUID) (*GetMyActiveHandbookResponse, error)
	StartMyHandbook(ctx context.Context, employeeID uuid.UUID) (*StartMyHandbookResponse, error)
	CompleteMyHandbookStep(ctx context.Context, employeeID, stepID uuid.UUID, response any) (*CompleteMyHandbookStepResponse, error)

	// Admin/manager
	CreateTemplateForDepartment(ctx context.Context, actorEmployeeID uuid.UUID, req CreateTemplateForDepartmentRequest) (*CreateTemplateForDepartmentResponse, error)
	ListTemplatesByDepartment(ctx context.Context, departmentID uuid.UUID) (*pagination.Response[ListTemplateResponse], error)
	CreateStep(ctx context.Context, req CreateStepRequest) (*CreateStepResponse, error)
	ListStepsByTemplate(ctx context.Context, templateID uuid.UUID) ([]ListStepResponse, error)
	AssignTemplateToEmployee(ctx context.Context, actorEmployeeID uuid.UUID, req AssignTemplateToEmployeeRequest) (*AssignTemplateToEmployeeResponse, error)
}
