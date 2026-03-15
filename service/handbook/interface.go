package handbook

import (
	"context"

	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HandbookService interface {
	// Employee self-service
	GetMyActiveHandbook(ctx context.Context, employeeID uuid.UUID) (*GetMyActiveHandbookResponse, error)
	StartMyHandbook(ctx context.Context, employeeID uuid.UUID) (*StartMyHandbookResponse, error)
	CompleteMyHandbookStep(ctx context.Context, employeeID, stepID uuid.UUID, response any) (*CompleteMyHandbookStepResponse, error)

	// Admin/manager
	CreateTemplateForDepartment(ctx context.Context, actorEmployeeID uuid.UUID, req CreateTemplateForDepartmentRequest) (*HandbookTemplateAPI, error)
	CloneTemplateToDraft(ctx context.Context, actorEmployeeID uuid.UUID, req CloneTemplateToDraftRequest) (*HandbookTemplateAPI, error)
	UpdateTemplate(ctx context.Context, req UpdateTemplateRequest) (*HandbookTemplateAPI, error)
	PublishTemplate(ctx context.Context, actorEmployeeID uuid.UUID, req PublishTemplateRequest) (*HandbookTemplateAPI, error)
	ListTemplatesByDepartment(ctx context.Context, departmentID uuid.UUID) (*pagination.Response[HandbookTemplateAPI], error)
	CreateStep(ctx context.Context, req CreateStepRequest) (*CreateStepResponse, error)
	UpdateStep(ctx context.Context, req UpdateStepRequest) (*UpdateStepResponse, error)
	DeleteStep(ctx context.Context, req DeleteStepRequest) (*DeleteStepResponse, error)
	ReorderTemplateSteps(ctx context.Context, req ReorderStepsRequest) (*ReorderStepsResponse, error)
	ListStepsByTemplate(ctx context.Context, templateID uuid.UUID) ([]ListStepResponse, error)
	AssignTemplateToEmployee(ctx context.Context, actorEmployeeID uuid.UUID, req AssignTemplateToEmployeeRequest) (*AssignTemplateToEmployeeResponse, error)
	WaiveEmployeeHandbook(ctx context.Context, actorEmployeeID uuid.UUID, req WaiveEmployeeHandbookRequest) (*WaiveEmployeeHandbookResponse, error)
	ListEmployeeHandbookHistory(ctx context.Context, employeeID uuid.UUID) ([]HandbookAssignmentHistoryEntry, error)
	ListEligibleEmployees(ctx *gin.Context, actorEmployeeID uuid.UUID, req ListEligibleEmployeesRequest) (*pagination.Response[ListEligibleEmployeesResponse], error)
	ListEmployeeHandbookAssignments(ctx *gin.Context, req ListEmployeeHandbookAssignmentsRequest) (*pagination.Response[EmployeeHandbookAssignmentSummary], error)
	GetEmployeeHandbookDetails(ctx context.Context, handbookID uuid.UUID) (*GetEmployeeHandbookDetailsResponse, error)
}
