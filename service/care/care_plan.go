package care

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *carePlanService) CreateClientCarePlan(ctx context.Context, clientID, employeeID uuid.UUID, req *CreateClientCarePlanRequest) (*CreateClientCarePlanResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.Store.WithTx(tx)

	arg := db.CreateClientMaturityMatrixAssessmentParams{
		ClientID:         clientID,
		MaturityMatrixID: req.MaturityMatrixID,
		StartDate:        pgtype.Date{Time: time.Now(), Valid: true},
		EndDate:          pgtype.Date{Time: time.Now().Add(time.Hour * 24 * 365), Valid: true},
		InitialLevel:     req.InitialLevel,
		TargetLevel:      req.TargetLevel,
		CurrentLevel:     req.InitialLevel,
	}

	clientAssessments, err := qtx.CreateClientMaturityMatrixAssessment(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to create client maturity matrix assessment", zap.Error(err))
		return nil, fmt.Errorf("failed to create client maturity matrix assessment: %w", err)
	}

	details, err := s.getDetails(qtx, ctx, req.MaturityMatrixID, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to get details", zap.Error(err))
		return nil, fmt.Errorf("failed to get details: %w", err)
	}

	generatedCarePlan, err := s.GrpcClient.GenerateCarePlan(ctx, &grpclient.PersonalizedCarePlanRequest{
		ClientData: &grpclient.ClientData{
			Age:              details.Age,
			EducationLevel:   details.EducationLevel,
			LivingSituation:  util.DerefString(details.LivingSituation),
			DomainName:       clientAssessments.TopicName,
			CurrentLevel:     clientAssessments.CurrentLevel,
			LevelDescription: details.LevelDescription[req.InitialLevel-1].Description,
		},
		DomainDefinitions: map[string]*grpclient.DomainLevels{
			details.TopicName: {
				Levels: map[int32]string{
					details.LevelDescription[0].Level: details.LevelDescription[0].Description,
					details.LevelDescription[1].Level: details.LevelDescription[1].Description,
					details.LevelDescription[2].Level: details.LevelDescription[2].Description,
					details.LevelDescription[3].Level: details.LevelDescription[3].Description,
					details.LevelDescription[4].Level: details.LevelDescription[4].Description,
				},
			},
		},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to generate care plan", zap.Error(err))
		return nil, fmt.Errorf("failed to generate care plan: %w", err)
	}

	carePlanID, err := s.insertCarePlan(ctx, qtx, generatedCarePlan, clientAssessments.ID, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan", zap.Error(err))
		return nil, fmt.Errorf("failed to insert care plan: %w", err)
	}

	err = s.insertCarePlanObjectives(ctx, qtx, carePlanID, generatedCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan objectives", zap.Error(err))
		return nil, fmt.Errorf("failed to insert care plan objectives: %w", err)
	}

	err = s.insertCarePlanInterventions(ctx, qtx, carePlanID, generatedCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan interventions", zap.Error(err))
		return nil, fmt.Errorf("failed to insert care plan interventions: %w", err)
	}

	err = s.insertCarePlanSuccessMetrics(ctx, qtx, carePlanID, generatedCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan success metrics", zap.Error(err))
		return nil, fmt.Errorf("failed to insert care plan success metrics: %w", err)
	}

	err = s.insertCarePlanRiskFactors(ctx, qtx, carePlanID, generatedCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan risk factors", zap.Error(err))
		return nil, fmt.Errorf("failed to insert care plan risk factors: %w", err)
	}

	err = s.insertCarePlanSupportNetwork(ctx, qtx, carePlanID, generatedCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan support network", zap.Error(err))
		return nil, fmt.Errorf("failed to insert care plan support network: %w", err)
	}

	err = s.insertCarePlanResources(ctx, qtx, carePlanID, generatedCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan resources", zap.Error(err))
		return nil, fmt.Errorf("failed to insert care plan resources: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to commit transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &CreateClientCarePlanResponse{
		ClientID:   clientID,
		CarePlanID: carePlanID,
	}, nil
}

func (s *carePlanService) ListClientCarePlans(ctx *gin.Context, clientID uuid.UUID, req *ListClientCarePlansRequest) (*pagination.Response[ListClientCarePlansResponse], error) {
	params := req.GetParams()
	clientAssessments, err := s.Store.ListClientMaturityMatrixAssessments(ctx, db.ListClientMaturityMatrixAssessmentsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientCarePlans", "Failed to list client maturity matrix assessments", zap.Error(err))
		return nil, fmt.Errorf("failed to list client maturity matrix assessments: %w", err)
	}

	carePlans := []ListClientCarePlansResponse{}
	for _, assessment := range clientAssessments {
		carePlans = append(carePlans, ListClientCarePlansResponse{
			CarePlanID:   util.DerefInt64(assessment.CarePlanID),
			TopicName:    assessment.TopicName,
			ClientID:     assessment.ClientID,
			StartDate:    assessment.StartDate,
			EndDate:      assessment.EndDate,
			InitialLevel: assessment.InitialLevel,
			CurrentLevel: assessment.CurrentLevel,
			IsActive:     assessment.IsActive,
		})
	}

	paginationResponse := pagination.NewResponse(ctx, req.Request, carePlans, clientAssessments[0].TotalCount)
	return &paginationResponse, nil
}

func (s *carePlanService) GetCarePlanOverview(ctx *gin.Context, carePlanID int64) (*GetCarePlanOverviewResponse, error) {
	carePlan, err := s.Store.GetCarePlanOverview(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanOverview", "Failed to get care plan overview", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan overview: %w", err)
	}

	response := &GetCarePlanOverviewResponse{
		ID:                carePlan.ID,
		Domain:            carePlan.TopicName,
		CurrentLevel:      carePlan.CurrentLevel,
		TargetLevel:       carePlan.TargetLevel,
		Status:            string(carePlan.Status),
		GeneratedAt:       carePlan.GeneratedAt.Time,
		AssessmentSummary: carePlan.AssessmentSummary,
		RawLlmResponse:    string(carePlan.RawLlmResponse),
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanOverview(ctx context.Context, carePlanID int64, req *UpdateCarePlanOverviewRequest) (*UpdateCarePlanOverviewResponse, error) {
	carePlan, err := s.Store.UpdateCarePlanOverview(ctx, db.UpdateCarePlanOverviewParams{
		ID:                carePlanID,
		AssessmentSummary: req.AssessmentSummary,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpfateCarePlanOverveiw", "Failed to update care plan overview", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan overview: %w", err)
	}

	response := &UpdateCarePlanOverviewResponse{
		CarePlanID:        carePlan.ID,
		AssessmentID:      carePlan.AssessmentID,
		AssessmentSummary: carePlan.AssessmentSummary,
	}

	return response, nil
}

func (s *carePlanService) DeleteCarePlan(ctx context.Context, carePlanID int64) error {
	err := s.Store.DeleteCarePlan(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlan", "Failed to delete care plan", zap.Error(err))
		return fmt.Errorf("failed to delete care plan: %w", err)
	}
	return nil
}

// =========================== Care plan objectives and actions ===========================

func (s *carePlanService) CreateCarePlanObjective(ctx context.Context, carePlanID int64, req *CreateCarePlanObjectiveRequest) (*CreateCarePlanObjectiveResponse, error) {
	createdObj, err := s.Store.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
		CarePlanID:  carePlanID,
		Description: req.Description,
		Timeframe:   db.CarePlanTimeframeEnum(req.TimeFrame),
		GoalTitle:   req.GoalTitle,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanObjective", "Failed to create care plan objective", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan objective: %w", err)
	}

	response := &CreateCarePlanObjectiveResponse{
		ID:              createdObj.ID,
		CarePlanID:      carePlanID,
		Timeframe:       string(createdObj.Timeframe),
		GoalTitle:       createdObj.GoalTitle,
		Description:     createdObj.Description,
		TargetDate:      createdObj.TargetDate.Time,
		Status:          string(createdObj.Status),
		CompletionDate:  createdObj.CompletionDate.Time,
		CompletionNotes: createdObj.CompletionNotes,
		CreatedAt:       createdObj.CreatedAt.Time,
		UpdatedAt:       createdObj.UpdatedAt.Time,
	}

	return response, nil
}

func (s *carePlanService) GetCarePlanObjectivesAndActions(ctx context.Context, carePlanID int64) (*GetCarePlanObjectivesResponse, error) {
	rows, err := s.Store.GetCarePlanObjectivesWithActions(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanObjectivesAndActions", "Failed to get care plan objectives and actions", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan objectives and actions: %w", err)
	}

	objectiveMap := make(map[int64]*CarePlanObjectives)
	for _, row := range rows {
		objective, exists := objectiveMap[row.ObjectiveID]
		if !exists {
			objective = &CarePlanObjectives{
				ObjectiveID: row.ObjectiveID,
				Title:       row.ObjectiveTitle,
				Description: row.ObjectiveDescription,
				TimeFrame:   string(row.ObjectiveTimeframe),
				Status:      string(row.ObjectiveStatus),
				Actions:     []CarePlanActions{},
			}
			objectiveMap[row.ObjectiveID] = objective
		}
		if row.ActionID != nil {
			action := CarePlanActions{
				ActionID:          *row.ActionID,
				SortOrder:         util.DerefInt32(row.SortOrder),
				ActionDescription: util.DerefString(row.ActionDescription),
				IsCompleted:       util.DerefBool(row.IsCompleted),
				Notes:             util.DerefString(row.ActionNotes),
			}
			objective.Actions = append(objective.Actions, action)
		}
	}

	response := &GetCarePlanObjectivesResponse{
		ShortTermGoals:  make([]CarePlanObjectives, 0),
		MediumTermGoals: make([]CarePlanObjectives, 0),
		LongTermGoals:   make([]CarePlanObjectives, 0),
	}

	for _, objective := range objectiveMap {
		switch objective.TimeFrame {
		case "short_term":
			response.ShortTermGoals = append(response.ShortTermGoals, *objective)
		case "medium_term":
			response.MediumTermGoals = append(response.MediumTermGoals, *objective)
		case "long_term":
			response.LongTermGoals = append(response.LongTermGoals, *objective)
		}
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanObjective(ctx context.Context, objectiveID int64, req *UpdateCarePlanObjectiveRequest) (*UpdateCarePlanObjectiveResponse, error) {
	objective, err := s.Store.UpdateCarePlanObjective(ctx, db.UpdateCarePlanObjectiveParams{
		ID:          objectiveID,
		Timeframe:   db.NullCarePlanTimeframeFromPtr(req.TimeFrame),
		GoalTitle:   req.GoalTitle,
		Description: req.Description,
		Status:      db.NullCarePlanObjectiveStatusFromPtr(req.Status),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanObjective", "Failed to update care plan objective", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan objective: %w", err)
	}

	response := &UpdateCarePlanObjectiveResponse{
		ObjectiveID: objective.ID,
		CarePlanId:  objective.CarePlanID,
	}

	return response, nil
}

func (s *carePlanService) DeleteCarePlanObjective(ctx context.Context, objectiveID int64) error {
	err := s.Store.DeleteCarePlanObjective(ctx, objectiveID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanObjective", "Failed to delete care plan objective", zap.Error(err))
		return fmt.Errorf("failed to delete care plan objective: %w", err)
	}
	return nil
}

func (s *carePlanService) CreateCarePlanAction(ctx context.Context, objectiveID int64, req *CreateCarePlanActionsRequest) (*CreateCarePlanActionsResponse, error) {
	maxSortOrder, err := s.Store.GetCarePlanActionsMaxSortOrder(ctx, objectiveID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanAction", "Failed to get max sort order", zap.Error(err))
		return nil, fmt.Errorf("failed to get max sort order: %w", err)
	}

	action, err := s.Store.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
		ObjectiveID:       objectiveID,
		ActionDescription: req.ActionDescription,
		SortOrder:         maxSortOrder + 1,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanAction", "Failed to create care plan action", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan action: %w", err)
	}

	response := &CreateCarePlanActionsResponse{
		ActionID:          action.ID,
		ObjectiveID:       action.ObjectiveID,
		ActionDescription: action.ActionDescription,
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanAction(ctx context.Context, actionID int64, req *UpdateCarePlanActionsRequest) (*UpdateCarePlanActionsResponse, error) {
	action, err := s.Store.UpdateCarePlanAction(ctx, db.UpdateCarePlanActionParams{
		ID:                actionID,
		ActionDescription: req.ActionDescription,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanAction", "Failed to update care plan action", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan action: %w", err)
	}

	response := &UpdateCarePlanActionsResponse{
		ActionID:          action.ID,
		ObjectiveID:       action.ObjectiveID,
		ActionDescription: action.ActionDescription,
	}

	return response, nil
}

func (s *carePlanService) DeleteCarePlanAction(ctx context.Context, actionID int64) error {
	err := s.Store.DeleteCarePlanAction(ctx, actionID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanAction", "Failed to delete care plan action", zap.Error(err))
		return fmt.Errorf("failed to delete care plan action: %w", err)
	}
	return nil
}

// ========================== Care plan interventions ===========================

func (s *carePlanService) CreateCarePlanIntervention(ctx context.Context, carePlanID int64, req *CreateCarePlanInterventionRequest) (*CreateCarePlanInterventionResponse, error) {
	intervention, err := s.Store.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
		CarePlanID:              carePlanID,
		Frequency:               db.CarePlanInterventionFrequencyEnum(req.Frequency),
		InterventionDescription: req.InterventionDescription,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanIntervention", "Failed to create care plan intervention", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan intervention: %w", err)
	}

	response := &CreateCarePlanInterventionResponse{
		InterventionID:          intervention.ID,
		CarePlanID:              intervention.CarePlanID,
		Frequency:               string(intervention.Frequency),
		InterventionDescription: intervention.InterventionDescription,
	}
	return response, nil
}

func (s *carePlanService) GetCarePlanInterventions(ctx context.Context, carePlanID int64) (*GetCarePlanInterventionsResponse, error) {
	interventions, err := s.Store.GetCarePlanInterventions(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanInterventions", "Failed to get care plan interventions", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan interventions: %w", err)
	}

	response := &GetCarePlanInterventionsResponse{
		DailyActivities:   []Intervention{},
		WeeklyActivities:  []Intervention{},
		MonthlyActivities: []Intervention{},
	}

	for _, intervention := range interventions {
		switch intervention.Frequency {
		case "daily":
			response.DailyActivities = append(response.DailyActivities, Intervention{
				InterventionID:          intervention.ID,
				InterventionDescription: intervention.InterventionDescription,
			})
		case "weekly":
			response.WeeklyActivities = append(response.WeeklyActivities, Intervention{
				InterventionID:          intervention.ID,
				InterventionDescription: intervention.InterventionDescription,
			})
		case "monthly":
			response.MonthlyActivities = append(response.MonthlyActivities, Intervention{
				InterventionID:          intervention.ID,
				InterventionDescription: intervention.InterventionDescription,
			})

		default:
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "GetCarePlanInterventions", "Unknown intervention frequency", zap.String("frequency", string(intervention.Frequency)))
		}
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanIntervention(ctx context.Context, interventionID int64, req *UpdateCarePlanInterventionRequest) (*UpdateCarePlanInterventionResponse, error) {
	intervention, err := s.Store.UpdateCarePlanIntervention(ctx, db.UpdateCarePlanInterventionParams{
		ID:                      interventionID,
		Frequency:               db.NullCarePlanInterventionFrequencyFromPtr(req.Frequency),
		InterventionDescription: req.InterventionDescription,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanIntervention", "Failed to update care plan intervention", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan intervention: %w", err)
	}

	response := &UpdateCarePlanInterventionResponse{
		InterventionID:          intervention.ID,
		CarePlanID:              intervention.CarePlanID,
		Frequency:               string(intervention.Frequency),
		InterventionDescription: intervention.InterventionDescription,
	}
	return response, nil
}

func (s *carePlanService) DeleteCarePlanIntervention(ctx context.Context, interventionID int64) error {
	err := s.Store.DeleteCarePlanIntervention(ctx, interventionID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanIntervention", "Failed to delete care plan intervention", zap.Error(err))
		return fmt.Errorf("failed to delete care plan intervention: %w", err)
	}
	return nil
}

// ==================== CarePlan Success Metrics ====================

func (s *carePlanService) CreateCarePlanSuccessMetric(ctx context.Context, carePlanID int64, req *CreateCarePlanSuccessMetricsRequest) (*CreateCarePlanSuccessMetricsResponse, error) {
	metric, err := s.Store.CreateCarePlanSuccessMetric(ctx, db.CreateCarePlanSuccessMetricParams{
		CarePlanID:        carePlanID,
		MetricName:        req.MetricName,
		TargetValue:       req.TargetValue,
		CurrentValue:      req.CurrentValue,
		MeasurementMethod: req.MeasurementMethod,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanSuccessMetric", "Failed to create care plan success metric", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan success metric: %w", err)
	}

	response := &CreateCarePlanSuccessMetricsResponse{
		MetricID:          metric.ID,
		MetricName:        metric.MetricName,
		CurrentValue:      metric.CurrentValue,
		TargetValue:       metric.TargetValue,
		MeasurementMethod: metric.MeasurementMethod,
	}
	return response, nil
}

func (s *carePlanService) GetCarePlanSuccessMetrics(ctx context.Context, carePlanID int64) ([]GetCarePlanSuccessMetricsResponse, error) {
	metrics, err := s.Store.GetCarePlanSuccessMetrics(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanSuccessMetrics", "Failed to get care plan success metrics", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan success metrics: %w", err)
	}

	response := []GetCarePlanSuccessMetricsResponse{}
	for _, metric := range metrics {
		response = append(response, GetCarePlanSuccessMetricsResponse{
			MetricID:          metric.ID,
			MetricName:        metric.MetricName,
			CurrentValue:      metric.CurrentValue,
			TargetValue:       metric.TargetValue,
			MeasurementMethod: metric.MeasurementMethod,
		})
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanSuccessMetric(ctx context.Context, metricID int64, req *UpdateCarePlanSuccessMetricsRequest) (*UpdateCarePlanSuccessMetricsResponse, error) {
	metric, err := s.Store.UpdateCarePlanSuccessMetric(ctx, db.UpdateCarePlanSuccessMetricParams{
		ID:                metricID,
		MetricName:        req.MetricName,
		TargetValue:       req.TargetValue,
		MeasurementMethod: req.MeasurementMethod,
		CurrentValue:      req.CurrentValue,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanSuccessMetric", "Failed to update care plan success metric", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan success metric: %w", err)
	}

	response := &UpdateCarePlanSuccessMetricsResponse{
		MetricID:          metric.ID,
		MetricName:        metric.MetricName,
		CurrentValue:      metric.CurrentValue,
		TargetValue:       metric.TargetValue,
		MeasurementMethod: metric.MeasurementMethod,
	}
	return response, nil
}

func (s *carePlanService) DeleteCarePlanSuccessMetric(ctx context.Context, metricID int64) error {
	err := s.Store.DeleteCarePlanSuccessMetric(ctx, metricID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanSuccessMetric", "Failed to delete care plan success metric", zap.Error(err))
		return fmt.Errorf("failed to delete care plan success metric: %w", err)
	}
	return nil
}

// ========================== Care plan risk factors ===========================

func (s *carePlanService) CreateCarePlanRisk(ctx context.Context, carePlanID int64, req *CreateCarePlanRisksRequest) (*CreateCarePlanRisksResponse, error) {
	riskFactor, err := s.Store.CreateCarePlanRisk(ctx, db.CreateCarePlanRiskParams{
		CarePlanID:         carePlanID,
		RiskDescription:    req.RiskDescription,
		MitigationStrategy: req.MitigationStrategy,
		RiskLevel:          db.CarePlanRiskLevelEnum(req.RiskLevel),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanRiskFactor", "Failed to create care plan risk factor", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan risk factor: %w", err)
	}

	response := &CreateCarePlanRisksResponse{
		RiskID:             riskFactor.ID,
		RiskDescription:    riskFactor.RiskDescription,
		MitigationStrategy: riskFactor.MitigationStrategy,
		RiskLevel:          string(riskFactor.RiskLevel),
	}
	return response, nil
}

func (s *carePlanService) GetCarePlanRisks(ctx context.Context, carePlanID int64) ([]GetCarePlanRisksResponse, error) {
	riskFactors, err := s.Store.GetCarePlanRisks(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanRiskFactors", "Failed to get care plan risk factors", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan risk factors: %w", err)
	}

	response := []GetCarePlanRisksResponse{}
	for _, riskFactor := range riskFactors {
		response = append(response, GetCarePlanRisksResponse{
			RiskID:             riskFactor.ID,
			RiskDescription:    riskFactor.RiskDescription,
			MitigationStrategy: riskFactor.MitigationStrategy,
			RiskLevel:          string(riskFactor.RiskLevel),
		})
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanRisk(ctx context.Context, riskID int64, req *UpdateCarePlanRisksRequest) (*UpdateCarePlanRisksResponse, error) {
	riskFactor, err := s.Store.UpdateCarePlanRisk(ctx, db.UpdateCarePlanRiskParams{
		ID:                 riskID,
		RiskDescription:    req.RiskDescription,
		MitigationStrategy: req.MitigationStrategy,
		RiskLevel:          db.NullCarePlanRiskLevelFromPtr(req.RiskLevel),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanRiskFactor", "Failed to update care plan risk factor", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan risk factor: %w", err)
	}

	response := &UpdateCarePlanRisksResponse{
		RiskID:             riskFactor.ID,
		RiskDescription:    riskFactor.RiskDescription,
		MitigationStrategy: riskFactor.MitigationStrategy,
		RiskLevel:          string(riskFactor.RiskLevel),
	}
	return response, nil
}

func (s *carePlanService) DeleteCarePlanRisk(ctx context.Context, riskID int64) error {
	err := s.Store.DeleteCarePlanRisk(ctx, riskID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanRiskFactor", "Failed to delete care plan risk factor", zap.Error(err))
		return fmt.Errorf("failed to delete care plan risk factor: %w", err)
	}
	return nil
}

// ========================== Care plan support network ===========================

func (s *carePlanService) CreateCarePlanSupportNetwork(ctx context.Context, carePlanID int64, req *CreateCarePlanSupportNetworkRequest) (*CreateCarePlanSupportNetworkResponse, error) {
	support, err := s.Store.CreateCarePlanSupportNetwork(ctx, db.CreateCarePlanSupportNetworkParams{
		CarePlanID:                carePlanID,
		RoleTitle:                 req.RoleTitle,
		ResponsibilityDescription: req.ResponsibilityDescription,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanSupportNetwork", "Failed to create care plan support network", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan support network: %w", err)
	}

	response := &CreateCarePlanSupportNetworkResponse{
		SupportNetworkID:          support.ID,
		RoleTitle:                 support.RoleTitle,
		ResponsibilityDescription: support.ResponsibilityDescription,
	}
	return response, nil
}

func (s *carePlanService) GetCarePlanSupportNetwork(ctx context.Context, carePlanID int64) ([]GetCarePlanSupportNetworkResponse, error) {
	supportNetworks, err := s.Store.GetCarePlanSupportNetwork(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanSupportNetwork", "Failed to get care plan support network", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan support network: %w", err)
	}

	response := []GetCarePlanSupportNetworkResponse{}
	for _, support := range supportNetworks {
		response = append(response, GetCarePlanSupportNetworkResponse{
			SupportNetworkID:          support.ID,
			RoleTitle:                 support.RoleTitle,
			ResponsibilityDescription: &support.ResponsibilityDescription,
		})
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanSupportNetwork(ctx context.Context, supportNetworkID int64, req *UpdateCarePlanSupportNetworkRequest) (*UpdateCarePlanSupportNetworkResponse, error) {
	support, err := s.Store.UpdateCarePlanSupportNetwork(ctx, db.UpdateCarePlanSupportNetworkParams{
		ID:                        supportNetworkID,
		RoleTitle:                 req.RoleTitle,
		ResponsibilityDescription: req.ResponsibilityDescription,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanSupportNetwork", "Failed to update care plan support network", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan support network: %w", err)
	}

	response := &UpdateCarePlanSupportNetworkResponse{
		SupportNetworkID:          support.ID,
		RoleTitle:                 support.RoleTitle,
		ResponsibilityDescription: support.ResponsibilityDescription,
	}
	return response, nil
}

func (s *carePlanService) DeleteCarePlanSupportNetwork(ctx context.Context, supportNetworkID int64) error {
	err := s.Store.DeleteCarePlanSupportNetwork(ctx, supportNetworkID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanSupportNetwork", "Failed to delete care plan support network", zap.Error(err))
		return fmt.Errorf("failed to delete care plan support network: %w", err)
	}
	return nil
}

// ========================= Care plan resources ===========================

func (s *carePlanService) CreateCarePlanResource(ctx context.Context, carePlanID int64, req *CreateCarePlanResourcesRequest) (*CreateCarePlanResourcesResponse, error) {
	resource, err := s.Store.CreateCarePlanResources(ctx, db.CreateCarePlanResourcesParams{
		CarePlanID:          carePlanID,
		ResourceDescription: req.ResourceDescription,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanResource", "Failed to create care plan resource", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan resource: %w", err)
	}

	response := &CreateCarePlanResourcesResponse{
		ID:                  resource.ID,
		ResourceDescription: resource.ResourceDescription,
		IsObtained:          resource.IsObtained,
		ObtainedDate:        &resource.ObtainedDate.Time,
	}
	return response, nil
}

func (s *carePlanService) GetCarePlanResources(ctx context.Context, carePlanID int64) ([]GetCarePlanResourcesResponse, error) {
	resources, err := s.Store.GetCarePlanResources(ctx, carePlanID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanResources", "Failed to get care plan resources", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan resources: %w", err)
	}

	response := []GetCarePlanResourcesResponse{}
	for _, resource := range resources {
		response = append(response, GetCarePlanResourcesResponse{
			ID:                  resource.ID,
			ResourceDescription: resource.ResourceDescription,
			IsObtained:          resource.IsObtained,
			ObtainedDate:        &resource.ObtainedDate.Time,
		})
	}
	return response, nil
}

func (s *carePlanService) UpdateCarePlanResource(ctx context.Context, resourceID int64, req *UpdateCarePlanResourcesRequest) (*UpdateCarePlanResourcesResponse, error) {
	resource, err := s.Store.UpdateCarePlanResource(ctx, db.UpdateCarePlanResourceParams{
		ID:                  resourceID,
		ResourceDescription: req.ResourceDescription,
		IsObtained:          req.IsObtained,
		ObtainedDate:        pgtype.Date{Time: req.ObtainedDate, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanResource", "Failed to update care plan resource", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan resource: %w", err)
	}

	response := &UpdateCarePlanResourcesResponse{
		ID:                  resource.ID,
		ResourceDescription: resource.ResourceDescription,
		IsObtained:          resource.IsObtained,
		ObtainedDate:        resource.ObtainedDate.Time,
	}
	return response, nil
}

func (s *carePlanService) DeleteCarePlanResource(ctx context.Context, resourceID int64) error {
	err := s.Store.DeleteCarePlanResource(ctx, resourceID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanResource", "Failed to delete care plan resource", zap.Error(err))
		return fmt.Errorf("failed to delete care plan resource: %w", err)
	}
	return nil
}

// ========================== Care plan Reports ===========================

func (s *carePlanService) CreateCarePlanReport(ctx context.Context, carePlanID int64, employeeID uuid.UUID, req *CreateCarePlanReportRequest) (*CreateCarePlanReportResponse, error) {
	report, err := s.Store.CreateCarePlanReport(ctx, db.CreateCarePlanReportParams{
		CarePlanID:          carePlanID,
		ReportType:          db.CarePlanReportTypeEnum(req.ReportType),
		ReportContent:       req.ReportContent,
		IsCritical:          req.IsCritical,
		CreatedByEmployeeID: employeeID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanReport", "Failed to create care plan report", zap.Error(err))
		return nil, fmt.Errorf("failed to create care plan report: %w", err)
	}

	response := &CreateCarePlanReportResponse{
		ID:            report.ID,
		CarePlanID:    report.CarePlanID,
		ReportType:    string(report.ReportType),
		ReportContent: report.ReportContent,
		IsCritical:    report.IsCritical,
		CreatedAt:     report.CreatedAt.Time,
	}
	return response, nil
}

func (s *carePlanService) ListCarePlanReports(ctx *gin.Context, carePlanID int64, req *ListCarePlanReportsRequest) (*pagination.Response[ListCarePlanReportsResponse], error) {
	params := req.GetParams()
	reports, err := s.Store.ListCarePlanReports(ctx, db.ListCarePlanReportsParams{
		CarePlanID: carePlanID,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListCarePlanReports", "Failed to list care plan reports", zap.Error(err))
		return nil, fmt.Errorf("failed to list care plan reports: %w", err)
	}

	responseData := []ListCarePlanReportsResponse{}
	for _, report := range reports {
		responseData = append(responseData, ListCarePlanReportsResponse{
			ID:                 report.ID,
			CarePlanID:         report.CarePlanID,
			ReportType:         string(report.ReportType),
			ReportContent:      report.ReportContent,
			CreatedByFirstName: report.CreatedByFirstName,
			CreatedByLastName:  report.CreatedByLastName,
			IsCritical:         report.IsCritical,
			CreatedAt:          report.CreatedAt.Time,
		})
	}

	paginationResponse := pagination.NewResponse(ctx, req.Request, responseData, reports[0].TotalCount)
	return &paginationResponse, nil
}

func (s *carePlanService) UpdateCarePlanReport(ctx context.Context, reportID int64, req *UpdateCarePlanReportRequest) (*UpdateCarePlanReportResponse, error) {
	report, err := s.Store.UpdateCarePlanReport(ctx, db.UpdateCarePlanReportParams{
		ID:            reportID,
		ReportType:    db.NullCarePlanReportTypeFromPtr(req.ReportType),
		ReportContent: req.ReportContent,
		IsCritical:    req.IsCritical,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateCarePlanReport", "Failed to update care plan report", zap.Error(err))
		return nil, fmt.Errorf("failed to update care plan report: %w", err)
	}

	response := &UpdateCarePlanReportResponse{
		ID:            report.ID,
		CarePlanID:    report.CarePlanID,
		ReportType:    string(report.ReportType),
		ReportContent: report.ReportContent,
		IsCritical:    report.IsCritical,
	}
	return response, nil
}

func (s *carePlanService) DeleteCarePlanReport(ctx context.Context, reportID int64) error {
	err := s.Store.DeleteCarePlanReport(ctx, reportID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanReport", "Failed to delete care plan report", zap.Error(err))
		return fmt.Errorf("failed to delete care plan report: %w", err)
	}
	return nil
}

// ========================== Private helper functions ===========================

func (s *carePlanService) getDetails(
	qtx *db.Queries,
	ctx context.Context,
	topicID int64,
	clientID uuid.UUID,
) (*Details, error) {
	topicDescription, err := qtx.GetMaturityMatrix(ctx, topicID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "getDetails", "Failed to get maturity matrix", zap.Error(err))
		return nil, fmt.Errorf("failed to get maturity matrix: %w", err)
	}

	var levelDescription []Level
	err = json.Unmarshal(topicDescription.LevelDescription, &levelDescription)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "getDetails", "Failed to unmarshal level description", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal level description: %w", err)
	}

	clientDetails, err := qtx.GetClientDetails(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "getDetails", "Failed to get client details", zap.Error(err))
		return nil, fmt.Errorf("failed to get client details: %w", err)
	}
	details := &Details{
		TopicName: topicDescription.TopicName,
		LivingSituation: func() *string {
			if clientDetails.LivingSituation.Valid {
				str := string(clientDetails.LivingSituation.ClientLivingSituationEnum)
				return &str
			}
			return nil
		}(),
		EducationLevel:   string(clientDetails.EducationLevel),
		Age:              int32(time.Since(clientDetails.DateOfBirth.Time).Hours() / 24 / 365),
		LevelDescription: levelDescription,
	}

	return details, nil
}

func (s *carePlanService) insertCarePlan(
	ctx context.Context,
	qtx *db.Queries,
	genCarePlan *grpclient.PersonalizedCarePlanResponse,
	assessmentID int64,
	emplpoyeeID uuid.UUID,
) (carePlanID int64, err error) {
	rawllmResp, err := json.Marshal(genCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlan", "Failed to marshal care plan", zap.Error(err))
		return 0, err
	}
	arg := db.CreateCarePlanParams{
		AssessmentID:          assessmentID,
		GeneratedByEmployeeID: &emplpoyeeID,
		AssessmentSummary:     genCarePlan.AssessmentSummary,
		RawLlmResponse:        rawllmResp,
		Status:                "draft",
	}

	carePlan, err := qtx.CreateCarePlan(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlan", "Failed to create care plan", zap.Error(err))
		return 0, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "insertCarePlan", "Care plan created successfully", zap.Int64("carePlanID", carePlan.ID))
	return carePlan.ID, nil
}

func (s *carePlanService) insertCarePlanObjectives(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID int64,
	genCarePlan *grpclient.PersonalizedCarePlanResponse,
) error {
	// short term goals
	for _, obj := range genCarePlan.CarePlanObjectives.ShortTermGoals {
		createdObj, err := qtx.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlanID,
			GoalTitle:   obj.GoalTitle,
			Description: obj.Description,
			Timeframe:   "short_term",
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true},
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create short term goal", zap.Error(err))
			return err
		}
		for i, action := range obj.SpecificActions {
			_, err := qtx.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1),
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create action step for short term goal", zap.Error(err))
				return err
			}
		}
	}
	// Meduim term goals

	for _, obj := range genCarePlan.CarePlanObjectives.MediumTermGoals {
		createdObj, err := qtx.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlanID,
			GoalTitle:   obj.GoalTitle,
			Description: obj.Description,
			Timeframe:   "medium_term",
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true},
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create medium term goal", zap.Error(err))
			return err
		}
		for i, action := range obj.SpecificActions {
			_, err := qtx.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1),
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create action step for medium term goal", zap.Error(err))
				return err
			}
		}
	}
	// Long term goals
	for _, obj := range genCarePlan.CarePlanObjectives.LongTermGoals {
		createdObj, err := qtx.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlanID,
			GoalTitle:   obj.GoalTitle,
			Description: obj.Description,
			Timeframe:   "long_term",
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true},
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create long term goal", zap.Error(err))
			return err
		}
		for i, action := range obj.SpecificActions {
			_, err := qtx.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1),
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create action step for long term goal", zap.Error(err))
				return err
			}
		}
	}
	return nil
}

func (s *carePlanService) insertCarePlanInterventions(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID int64,
	genCarePlan *grpclient.PersonalizedCarePlanResponse,
) error {
	for _, intervention := range genCarePlan.Interventions.DailyActivities {
		_, err := qtx.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlanID,
			Frequency:               "daily",
			InterventionDescription: intervention,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanInterventions", "Failed to create care plan intervention", zap.Error(err))
			return err
		}
	}

	for _, intervention := range genCarePlan.Interventions.WeeklyActivities {
		_, err := qtx.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlanID,
			Frequency:               "weekly",
			InterventionDescription: intervention,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanInterventions", "Failed to create care plan intervention", zap.Error(err))
			return err
		}
	}

	for _, intervention := range genCarePlan.Interventions.MonthlyActivities {
		_, err := qtx.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlanID,
			Frequency:               "monthly",
			InterventionDescription: intervention,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanInterventions", "Failed to create care plan intervention", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *carePlanService) insertCarePlanSuccessMetrics(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID int64,
	genCarePlan *grpclient.PersonalizedCarePlanResponse,
) error {
	for _, metric := range genCarePlan.SuccessMetrics {
		_, err := qtx.CreateCarePlanSuccessMetric(ctx, db.CreateCarePlanSuccessMetricParams{
			CarePlanID:        carePlanID,
			MetricName:        metric.Metric,
			TargetValue:       metric.Target,
			MeasurementMethod: metric.MeasurementMethod,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanSuccessMetrics", "Failed to create care plan success metric", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *carePlanService) insertCarePlanRiskFactors(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID int64,
	genCarePlan *grpclient.PersonalizedCarePlanResponse,
) error {
	for _, risk := range genCarePlan.RiskFactors {
		_, err := qtx.CreateCarePlanRisk(ctx, db.CreateCarePlanRiskParams{
			CarePlanID:         carePlanID,
			RiskDescription:    risk.Risk,
			MitigationStrategy: risk.Mitigation,
			RiskLevel:          db.CarePlanRiskLevelEnum(risk.RiskLevel), // Use pointer to allow NULL values
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanRiskFactors", "Failed to create care plan risk factor", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *carePlanService) insertCarePlanSupportNetwork(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID int64,
	genCarePlan *grpclient.PersonalizedCarePlanResponse,
) error {
	for _, network := range genCarePlan.SupportNetwork {
		_, err := qtx.CreateCarePlanSupportNetwork(ctx, db.CreateCarePlanSupportNetworkParams{
			CarePlanID:                carePlanID,
			RoleTitle:                 network.Role,
			ResponsibilityDescription: network.Responsibility,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanSupportNetworks", "Failed to create care plan support network", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *carePlanService) insertCarePlanResources(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID int64,
	genCarePlan *grpclient.PersonalizedCarePlanResponse,
) error {
	for _, resource := range genCarePlan.ResourcesRequired {
		_, err := qtx.CreateCarePlanResources(ctx, db.CreateCarePlanResourcesParams{
			CarePlanID:          carePlanID,
			ResourceDescription: resource,
			IsObtained:          false,
			ObtainedDate:        pgtype.Date{Time: time.Now(), Valid: false},
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanResources", "Failed to create care plan resource", zap.Error(err))
			return err
		}
	}

	return nil
}
