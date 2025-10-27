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
	GetCarePlanOverview(ctx *gin.Context, carePlanID int64) (*GetCarePlanOverviewResponse, error)
	UpdateCarePlanOverview(ctx context.Context, carePlanID int64, req *UpdateCarePlanOverviewRequest) (*UpdateCarePlanOverviewResponse, error)
	DeleteCarePlan(ctx context.Context, carePlanID int64) error

	// Care Plan Objectives
	CreateCarePlanObjective(ctx context.Context, carePlanID int64, req *CreateCarePlanObjectiveRequest) (*CreateCarePlanObjectiveResponse, error)
	GetCarePlanObjectivesAndActions(ctx context.Context, carePlanID int64) (*GetCarePlanObjectivesResponse, error)
	UpdateCarePlanObjective(ctx context.Context, objectiveID int64, req *UpdateCarePlanObjectiveRequest) (*UpdateCarePlanObjectiveResponse, error)
	DeleteCarePlanObjective(ctx context.Context, objectiveID int64) error

	// Care Plan Actions
	CreateCarePlanAction(ctx context.Context, objectiveID int64, req *CreateCarePlanActionsRequest) (*CreateCarePlanActionsResponse, error)
	UpdateCarePlanAction(ctx context.Context, actionID int64, req *UpdateCarePlanActionsRequest) (*UpdateCarePlanActionsResponse, error)
	DeleteCarePlanAction(ctx context.Context, actionID int64) error

	// Care Plan Interventions
	CreateCarePlanIntervention(ctx context.Context, carePlanID int64, req *CreateCarePlanInterventionRequest) (*CreateCarePlanInterventionResponse, error)
	GetCarePlanInterventions(ctx context.Context, carePlanID int64) (*GetCarePlanInterventionsResponse, error)
	UpdateCarePlanIntervention(ctx context.Context, interventionID int64, req *UpdateCarePlanInterventionRequest) (*UpdateCarePlanInterventionResponse, error)
	DeleteCarePlanIntervention(ctx context.Context, interventionID int64) error

	// Care Plan Success Metrics
	CreateCarePlanSuccessMetric(ctx context.Context, carePlanID int64, req *CreateCarePlanSuccessMetricsRequest) (*CreateCarePlanSuccessMetricsResponse, error)
	GetCarePlanSuccessMetrics(ctx context.Context, carePlanID int64) ([]GetCarePlanSuccessMetricsResponse, error)
	UpdateCarePlanSuccessMetric(ctx context.Context, metricID int64, req *UpdateCarePlanSuccessMetricsRequest) (*UpdateCarePlanSuccessMetricsResponse, error)
	DeleteCarePlanSuccessMetric(ctx context.Context, metricID int64) error

	// Care Plan Risks
	CreateCarePlanRisk(ctx context.Context, carePlanID int64, req *CreateCarePlanRisksRequest) (*CreateCarePlanRisksResponse, error)
	GetCarePlanRisks(ctx context.Context, carePlanID int64) ([]GetCarePlanRisksResponse, error)
	UpdateCarePlanRisk(ctx context.Context, riskID int64, req *UpdateCarePlanRisksRequest) (*UpdateCarePlanRisksResponse, error)
	DeleteCarePlanRisk(ctx context.Context, riskID int64) error

	// Care Plan Support Network
	CreateCarePlanSupportNetwork(ctx context.Context, carePlanID int64, req *CreateCarePlanSupportNetworkRequest) (*CreateCarePlanSupportNetworkResponse, error)
	GetCarePlanSupportNetwork(ctx context.Context, carePlanID int64) ([]GetCarePlanSupportNetworkResponse, error)
	UpdateCarePlanSupportNetwork(ctx context.Context, supportID int64, req *UpdateCarePlanSupportNetworkRequest) (*UpdateCarePlanSupportNetworkResponse, error)
	DeleteCarePlanSupportNetwork(ctx context.Context, supportID int64) error

	// Care Plan Resources
	CreateCarePlanResource(ctx context.Context, carePlanID int64, req *CreateCarePlanResourcesRequest) (*CreateCarePlanResourcesResponse, error)
	GetCarePlanResources(ctx context.Context, carePlanID int64) ([]GetCarePlanResourcesResponse, error)
	UpdateCarePlanResource(ctx context.Context, resourceID int64, req *UpdateCarePlanResourcesRequest) (*UpdateCarePlanResourcesResponse, error)
	DeleteCarePlanResource(ctx context.Context, resourceID int64) error

	// Care Plan Reports
	CreateCarePlanReport(ctx context.Context, carePlanID int64, employeeID uuid.UUID, req *CreateCarePlanReportRequest) (*CreateCarePlanReportResponse, error)
	ListCarePlanReports(ctx *gin.Context, carePlanID int64, req *ListCarePlanReportsRequest) (*pagination.Response[ListCarePlanReportsResponse], error)
	UpdateCarePlanReport(ctx context.Context, reportID int64, req *UpdateCarePlanReportRequest) (*UpdateCarePlanReportResponse, error)
	DeleteCarePlanReport(ctx context.Context, reportID int64) error
}

type carePlanService struct {
	*deps.ServiceDependencies
}

func NewCarePlanService(deps *deps.ServiceDependencies) CarePlanService {
	return &carePlanService{
		ServiceDependencies: deps,
	}
}
