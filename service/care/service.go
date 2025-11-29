package care

import (
	"context"

	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CarePlanService interface {
	ListCarePlanTopics(ctx context.Context) ([]ListCarePlanTopics, error)
	CreateClientCarePlan(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, req *CreateClientCarePlanRequest) (*CreateClientCarePlanResponse, error)
	ListClientCarePlans(ctx *gin.Context, clientID uuid.UUID, req *ListClientCarePlansRequest) (*pagination.Response[ListClientCarePlansResponse], error)
	GetCarePlanOverview(ctx *gin.Context, carePlanID uuid.UUID) (*GetCarePlanOverviewResponse, error)
	UpdateCarePlanOverview(ctx context.Context, carePlanID uuid.UUID, req *UpdateCarePlanOverviewRequest) (*UpdateCarePlanOverviewResponse, error)
	DeleteCarePlan(ctx context.Context, carePlanID uuid.UUID) error

	// Care Plan Objectives
	CreateCarePlanObjective(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanObjectiveRequest) (*CreateCarePlanObjectiveResponse, error)
	GetCarePlanObjectivesAndActions(ctx context.Context, carePlanID uuid.UUID) (*GetCarePlanObjectivesResponse, error)
	UpdateCarePlanObjective(ctx context.Context, objectiveID uuid.UUID, req *UpdateCarePlanObjectiveRequest) (*UpdateCarePlanObjectiveResponse, error)
	DeleteCarePlanObjective(ctx context.Context, objectiveID uuid.UUID) error

	// Care Plan Actions
	CreateCarePlanAction(ctx context.Context, objectiveID uuid.UUID, req *CreateCarePlanActionsRequest) (*CreateCarePlanActionsResponse, error)
	UpdateCarePlanAction(ctx context.Context, actionID uuid.UUID, req *UpdateCarePlanActionsRequest) (*UpdateCarePlanActionsResponse, error)
	DeleteCarePlanAction(ctx context.Context, actionID uuid.UUID) error

	// Care Plan Interventions
	CreateCarePlanIntervention(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanInterventionRequest) (*CreateCarePlanInterventionResponse, error)
	GetCarePlanInterventions(ctx context.Context, carePlanID uuid.UUID) (*GetCarePlanInterventionsResponse, error)
	UpdateCarePlanIntervention(ctx context.Context, interventionID uuid.UUID, req *UpdateCarePlanInterventionRequest) (*UpdateCarePlanInterventionResponse, error)
	DeleteCarePlanIntervention(ctx context.Context, interventionID uuid.UUID) error

	// Care Plan Success Metrics
	CreateCarePlanSuccessMetric(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanSuccessMetricsRequest) (*CreateCarePlanSuccessMetricsResponse, error)
	GetCarePlanSuccessMetrics(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanSuccessMetricsResponse, error)
	UpdateCarePlanSuccessMetric(ctx context.Context, metricID uuid.UUID, req *UpdateCarePlanSuccessMetricsRequest) (*UpdateCarePlanSuccessMetricsResponse, error)
	DeleteCarePlanSuccessMetric(ctx context.Context, metricID uuid.UUID) error

	// Care Plan Risks
	CreateCarePlanRisk(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanRisksRequest) (*CreateCarePlanRisksResponse, error)
	GetCarePlanRisks(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanRisksResponse, error)
	UpdateCarePlanRisk(ctx context.Context, riskID uuid.UUID, req *UpdateCarePlanRisksRequest) (*UpdateCarePlanRisksResponse, error)
	DeleteCarePlanRisk(ctx context.Context, riskID uuid.UUID) error

	// Care Plan Support Network
	CreateCarePlanSupportNetwork(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanSupportNetworkRequest) (*CreateCarePlanSupportNetworkResponse, error)
	GetCarePlanSupportNetwork(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanSupportNetworkResponse, error)
	UpdateCarePlanSupportNetwork(ctx context.Context, supportID uuid.UUID, req *UpdateCarePlanSupportNetworkRequest) (*UpdateCarePlanSupportNetworkResponse, error)
	DeleteCarePlanSupportNetwork(ctx context.Context, supportID uuid.UUID) error

	// Care Plan Resources
	CreateCarePlanResource(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanResourcesRequest) (*CreateCarePlanResourcesResponse, error)
	GetCarePlanResources(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanResourcesResponse, error)
	UpdateCarePlanResource(ctx context.Context, resourceID uuid.UUID, req *UpdateCarePlanResourcesRequest) (*UpdateCarePlanResourcesResponse, error)
	DeleteCarePlanResource(ctx context.Context, resourceID uuid.UUID) error

	// Care Plan Reports
	CreateCarePlanReport(ctx context.Context, carePlanID uuid.UUID, employeeID uuid.UUID, req *CreateCarePlanReportRequest) (*CreateCarePlanReportResponse, error)
	ListCarePlanReports(ctx *gin.Context, carePlanID uuid.UUID, req *ListCarePlanReportsRequest) (*pagination.Response[ListCarePlanReportsResponse], error)
	UpdateCarePlanReport(ctx context.Context, reportID uuid.UUID, req *UpdateCarePlanReportRequest) (*UpdateCarePlanReportResponse, error)
	DeleteCarePlanReport(ctx context.Context, reportID uuid.UUID) error
}

type carePlanService struct {
	*deps.ServiceDependencies
}

func NewCarePlanService(deps *deps.ServiceDependencies) CarePlanService {
	return &carePlanService{
		ServiceDependencies: deps,
	}
}
