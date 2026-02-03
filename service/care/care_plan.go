package care

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/ai"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *carePlanService) CreateClientCarePlan(ctx context.Context, clientID, employeeID uuid.UUID, req *CreateClientCarePlanRequest) (*CreateClientCarePlanResponse, error) {
	var carePlanID uuid.UUID
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		arg := db.CreateClientMaturityMatrixAssessmentParams{
			ClientID:         clientID,
			MaturityMatrixID: req.MaturityMatrixID,
			StartDate:        pgtype.Date{Time: time.Now(), Valid: true},
			EndDate:          pgtype.Date{Time: time.Now().Add(time.Hour * 24 * 365), Valid: true},
			InitialLevel:     req.InitialLevel,
			TargetLevel:      req.TargetLevel,
			CurrentLevel:     req.InitialLevel,
		}

		clientAssessments, err := q.CreateClientMaturityMatrixAssessment(ctx, arg)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to create client maturity matrix assessment", zap.Error(err))
			return fmt.Errorf("failed to create client maturity matrix assessment: %w", err)
		}

		details, err := s.getDetails(q, ctx, req.MaturityMatrixID, clientID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to get details", zap.Error(err))
			return fmt.Errorf("failed to get details: %w", err)
		}

		// Construct client data for AI
		clientData := map[string]interface{}{
			"age":               details.Age,
			"education_level":   details.EducationLevel,
			"living_situation":  util.DerefString(details.LivingSituation),
			"domain_name":       clientAssessments.TopicName,
			"current_level":     clientAssessments.CurrentLevel,
			"level_description": details.LevelDescription[req.InitialLevel-1].Description,
			"domain_levels": map[string]interface{}{
				"1": details.LevelDescription[0].Description,
				"2": details.LevelDescription[1].Description,
				"3": details.LevelDescription[2].Description,
				"4": details.LevelDescription[3].Description,
				"5": details.LevelDescription[4].Description,
			},
		}

		generatedCarePlan, err := s.AIService.GenerateCarePlan(ctx, clientData)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to generate care plan", zap.Error(err))
			return fmt.Errorf("failed to generate care plan: %w", err)
		}

		carePlanID, err = s.insertCarePlan(ctx, q, generatedCarePlan, clientAssessments.ID, employeeID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan", zap.Error(err))
			return fmt.Errorf("failed to insert care plan: %w", err)
		}

		err = s.insertCarePlanObjectives(ctx, q, carePlanID, generatedCarePlan)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan objectives", zap.Error(err))
			return fmt.Errorf("failed to insert care plan objectives: %w", err)
		}

		err = s.insertCarePlanInterventions(ctx, q, carePlanID, generatedCarePlan)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan interventions", zap.Error(err))
			return fmt.Errorf("failed to insert care plan interventions: %w", err)
		}

		err = s.insertCarePlanSuccessMetrics(ctx, q, carePlanID, generatedCarePlan)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan success metrics", zap.Error(err))
			return fmt.Errorf("failed to insert care plan success metrics: %w", err)
		}

		err = s.insertCarePlanRiskFactors(ctx, q, carePlanID, generatedCarePlan)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan risk factors", zap.Error(err))
			return fmt.Errorf("failed to insert care plan risk factors: %w", err)
		}

		err = s.insertCarePlanSupportNetwork(ctx, q, carePlanID, generatedCarePlan)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan support network", zap.Error(err))
			return fmt.Errorf("failed to insert care plan support network: %w", err)
		}

		err = s.insertCarePlanResources(ctx, q, carePlanID, generatedCarePlan)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientCarePlan", "Failed to insert care plan resources", zap.Error(err))
			return fmt.Errorf("failed to insert care plan resources: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CreateClientCarePlanResponse{
		ClientID:   clientID,
		CarePlanID: carePlanID,
	}, nil
}

func (s *carePlanService) ListClientCarePlans(ctx *gin.Context, clientID uuid.UUID, req *ListClientCarePlansRequest) (*pagination.Response[ListClientCarePlansResponse], error) {
	params := req.GetParams()
	var clientAssessments []db.ListClientMaturityMatrixAssessmentsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		clientAssessments, txErr = q.ListClientMaturityMatrixAssessments(ctx, db.ListClientMaturityMatrixAssessmentsParams{
			ClientID: clientID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
		return txErr
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientCarePlans", "Failed to list client maturity matrix assessments", zap.Error(err))
		return nil, fmt.Errorf("failed to list client maturity matrix assessments: %w", err)
	}

	carePlans := []ListClientCarePlansResponse{}
	for _, assessment := range clientAssessments {
		carePlans = append(carePlans, ListClientCarePlansResponse{
			CarePlanID:   util.DerefUUID(assessment.CarePlanID),
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

func (s *carePlanService) GetCarePlanOverview(ctx *gin.Context, carePlanID uuid.UUID) (*GetCarePlanOverviewResponse, error) {
	var carePlan db.GetCarePlanOverviewRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		carePlan, txErr = q.GetCarePlanOverview(ctx, carePlanID)
		return txErr
	})
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

func (s *carePlanService) UpdateCarePlanOverview(ctx context.Context, carePlanID uuid.UUID, req *UpdateCarePlanOverviewRequest) (*UpdateCarePlanOverviewResponse, error) {
	var carePlan db.CarePlan
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		carePlan, txErr = q.UpdateCarePlanOverview(ctx, db.UpdateCarePlanOverviewParams{
			ID:                carePlanID,
			AssessmentSummary: req.AssessmentSummary,
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlan(ctx context.Context, carePlanID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlan(ctx, carePlanID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlan", "Failed to delete care plan", zap.Error(err))
		return fmt.Errorf("failed to delete care plan: %w", err)
	}
	return nil
}

// =========================== Care plan objectives and actions ===========================

func (s *carePlanService) CreateCarePlanObjective(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanObjectiveRequest) (*CreateCarePlanObjectiveResponse, error) {
	var createdObj db.CarePlanObjective
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		createdObj, txErr = q.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlanID,
			Description: req.Description,
			Timeframe:   db.CarePlanTimeframeEnum(req.TimeFrame),
			GoalTitle:   req.GoalTitle,
		})
		return txErr
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

func (s *carePlanService) GetCarePlanObjectivesAndActions(ctx context.Context, carePlanID uuid.UUID) (*GetCarePlanObjectivesResponse, error) {
	var rows []db.GetCarePlanObjectivesWithActionsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		rows, txErr = q.GetCarePlanObjectivesWithActions(ctx, carePlanID)
		return txErr
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetCarePlanObjectivesAndActions", "Failed to get care plan objectives and actions", zap.Error(err))
		return nil, fmt.Errorf("failed to get care plan objectives and actions: %w", err)
	}

	objectiveMap := make(map[uuid.UUID]*CarePlanObjectives)
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

func (s *carePlanService) UpdateCarePlanObjective(ctx context.Context, objectiveID uuid.UUID, req *UpdateCarePlanObjectiveRequest) (*UpdateCarePlanObjectiveResponse, error) {
	var objective db.CarePlanObjective
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		objective, txErr = q.UpdateCarePlanObjective(ctx, db.UpdateCarePlanObjectiveParams{
			ID:          objectiveID,
			Timeframe:   db.NullCarePlanTimeframeFromPtr(req.TimeFrame),
			GoalTitle:   req.GoalTitle,
			Description: req.Description,
			Status:      db.NullCarePlanObjectiveStatusFromPtr(req.Status),
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanObjective(ctx context.Context, objectiveID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanObjective(ctx, objectiveID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanObjective", "Failed to delete care plan objective", zap.Error(err))
		return fmt.Errorf("failed to delete care plan objective: %w", err)
	}
	return nil
}

func (s *carePlanService) CreateCarePlanAction(ctx context.Context, objectiveID uuid.UUID, req *CreateCarePlanActionsRequest) (*CreateCarePlanActionsResponse, error) {
	var action db.CarePlanAction
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		maxSortOrder, err := q.GetCarePlanActionsMaxSortOrder(ctx, objectiveID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanAction", "Failed to get max sort order", zap.Error(err))
			return fmt.Errorf("failed to get max sort order: %w", err)
		}

		action, err = q.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
			ObjectiveID:       objectiveID,
			ActionDescription: req.ActionDescription,
			SortOrder:         maxSortOrder + 1,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateCarePlanAction", "Failed to create care plan action", zap.Error(err))
			return fmt.Errorf("failed to create care plan action: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	response := &CreateCarePlanActionsResponse{
		ActionID:          action.ID,
		ObjectiveID:       action.ObjectiveID,
		ActionDescription: action.ActionDescription,
	}

	return response, nil
}

func (s *carePlanService) UpdateCarePlanAction(ctx context.Context, actionID uuid.UUID, req *UpdateCarePlanActionsRequest) (*UpdateCarePlanActionsResponse, error) {
	var action db.CarePlanAction
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		action, txErr = q.UpdateCarePlanAction(ctx, db.UpdateCarePlanActionParams{
			ID:                actionID,
			ActionDescription: req.ActionDescription,
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanAction(ctx context.Context, actionID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanAction(ctx, actionID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanAction", "Failed to delete care plan action", zap.Error(err))
		return fmt.Errorf("failed to delete care plan action: %w", err)
	}
	return nil
}

// ========================== Care plan interventions ===========================

func (s *carePlanService) CreateCarePlanIntervention(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanInterventionRequest) (*CreateCarePlanInterventionResponse, error) {
	var intervention db.CarePlanIntervention
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		intervention, txErr = q.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlanID,
			Frequency:               db.CarePlanInterventionFrequencyEnum(req.Frequency),
			InterventionDescription: req.InterventionDescription,
		})
		return txErr
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

func (s *carePlanService) GetCarePlanInterventions(ctx context.Context, carePlanID uuid.UUID) (*GetCarePlanInterventionsResponse, error) {
	var interventions []db.CarePlanIntervention
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		interventions, txErr = q.GetCarePlanInterventions(ctx, carePlanID)
		return txErr
	})
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

func (s *carePlanService) UpdateCarePlanIntervention(ctx context.Context, interventionID uuid.UUID, req *UpdateCarePlanInterventionRequest) (*UpdateCarePlanInterventionResponse, error) {
	var intervention db.CarePlanIntervention
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		intervention, txErr = q.UpdateCarePlanIntervention(ctx, db.UpdateCarePlanInterventionParams{
			ID:                      interventionID,
			Frequency:               db.NullCarePlanInterventionFrequencyFromPtr(req.Frequency),
			InterventionDescription: req.InterventionDescription,
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanIntervention(ctx context.Context, interventionID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanIntervention(ctx, interventionID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanIntervention", "Failed to delete care plan intervention", zap.Error(err))
		return fmt.Errorf("failed to delete care plan intervention: %w", err)
	}
	return nil
}

// ==================== CarePlan Success Metrics ====================

func (s *carePlanService) CreateCarePlanSuccessMetric(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanSuccessMetricsRequest) (*CreateCarePlanSuccessMetricsResponse, error) {
	var metric db.CarePlanMetric
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		metric, txErr = q.CreateCarePlanSuccessMetric(ctx, db.CreateCarePlanSuccessMetricParams{
			CarePlanID:        carePlanID,
			MetricName:        req.MetricName,
			TargetValue:       req.TargetValue,
			CurrentValue:      req.CurrentValue,
			MeasurementMethod: req.MeasurementMethod,
		})
		return txErr
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

func (s *carePlanService) GetCarePlanSuccessMetrics(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanSuccessMetricsResponse, error) {
	var metrics []db.CarePlanMetric
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		metrics, txErr = q.GetCarePlanSuccessMetrics(ctx, carePlanID)
		return txErr
	})
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

func (s *carePlanService) UpdateCarePlanSuccessMetric(ctx context.Context, metricID uuid.UUID, req *UpdateCarePlanSuccessMetricsRequest) (*UpdateCarePlanSuccessMetricsResponse, error) {
	var metric db.CarePlanMetric
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		metric, txErr = q.UpdateCarePlanSuccessMetric(ctx, db.UpdateCarePlanSuccessMetricParams{
			ID:                metricID,
			MetricName:        req.MetricName,
			TargetValue:       req.TargetValue,
			MeasurementMethod: req.MeasurementMethod,
			CurrentValue:      req.CurrentValue,
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanSuccessMetric(ctx context.Context, metricID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanSuccessMetric(ctx, metricID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanSuccessMetric", "Failed to delete care plan success metric", zap.Error(err))
		return fmt.Errorf("failed to delete care plan success metric: %w", err)
	}
	return nil
}

// ========================== Care plan risk factors ===========================

func (s *carePlanService) CreateCarePlanRisk(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanRisksRequest) (*CreateCarePlanRisksResponse, error) {
	var riskFactor db.CarePlanRisk
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		riskFactor, txErr = q.CreateCarePlanRisk(ctx, db.CreateCarePlanRiskParams{
			CarePlanID:         carePlanID,
			RiskDescription:    req.RiskDescription,
			MitigationStrategy: req.MitigationStrategy,
			RiskLevel:          db.CarePlanRiskLevelEnum(req.RiskLevel),
		})
		return txErr
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

func (s *carePlanService) GetCarePlanRisks(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanRisksResponse, error) {
	var riskFactors []db.CarePlanRisk
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		riskFactors, txErr = q.GetCarePlanRisks(ctx, carePlanID)
		return txErr
	})
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

func (s *carePlanService) UpdateCarePlanRisk(ctx context.Context, riskID uuid.UUID, req *UpdateCarePlanRisksRequest) (*UpdateCarePlanRisksResponse, error) {
	var riskFactor db.CarePlanRisk
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		riskFactor, txErr = q.UpdateCarePlanRisk(ctx, db.UpdateCarePlanRiskParams{
			ID:                 riskID,
			RiskDescription:    req.RiskDescription,
			MitigationStrategy: req.MitigationStrategy,
			RiskLevel:          db.NullCarePlanRiskLevelFromPtr(req.RiskLevel),
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanRisk(ctx context.Context, riskID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanRisk(ctx, riskID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanRiskFactor", "Failed to delete care plan risk factor", zap.Error(err))
		return fmt.Errorf("failed to delete care plan risk factor: %w", err)
	}
	return nil
}

// ========================== Care plan support network ===========================

func (s *carePlanService) CreateCarePlanSupportNetwork(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanSupportNetworkRequest) (*CreateCarePlanSupportNetworkResponse, error) {
	var support db.CarePlanSupportNetwork
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		support, txErr = q.CreateCarePlanSupportNetwork(ctx, db.CreateCarePlanSupportNetworkParams{
			CarePlanID:                carePlanID,
			RoleTitle:                 req.RoleTitle,
			ResponsibilityDescription: req.ResponsibilityDescription,
		})
		return txErr
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

func (s *carePlanService) GetCarePlanSupportNetwork(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanSupportNetworkResponse, error) {
	var supportNetworks []db.CarePlanSupportNetwork
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		supportNetworks, txErr = q.GetCarePlanSupportNetwork(ctx, carePlanID)
		return txErr
	})
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

func (s *carePlanService) UpdateCarePlanSupportNetwork(ctx context.Context, supportNetworkID uuid.UUID, req *UpdateCarePlanSupportNetworkRequest) (*UpdateCarePlanSupportNetworkResponse, error) {
	var support db.CarePlanSupportNetwork
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		support, txErr = q.UpdateCarePlanSupportNetwork(ctx, db.UpdateCarePlanSupportNetworkParams{
			ID:                        supportNetworkID,
			RoleTitle:                 req.RoleTitle,
			ResponsibilityDescription: req.ResponsibilityDescription,
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanSupportNetwork(ctx context.Context, supportNetworkID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanSupportNetwork(ctx, supportNetworkID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanSupportNetwork", "Failed to delete care plan support network", zap.Error(err))
		return fmt.Errorf("failed to delete care plan support network: %w", err)
	}
	return nil
}

// ========================= Care plan resources ===========================

func (s *carePlanService) CreateCarePlanResource(ctx context.Context, carePlanID uuid.UUID, req *CreateCarePlanResourcesRequest) (*CreateCarePlanResourcesResponse, error) {
	var resource db.CarePlanResource
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		resource, txErr = q.CreateCarePlanResources(ctx, db.CreateCarePlanResourcesParams{
			CarePlanID:          carePlanID,
			ResourceDescription: req.ResourceDescription,
		})
		return txErr
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

func (s *carePlanService) GetCarePlanResources(ctx context.Context, carePlanID uuid.UUID) ([]GetCarePlanResourcesResponse, error) {
	var resources []db.CarePlanResource
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		resources, txErr = q.GetCarePlanResources(ctx, carePlanID)
		return txErr
	})
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

func (s *carePlanService) UpdateCarePlanResource(ctx context.Context, resourceID uuid.UUID, req *UpdateCarePlanResourcesRequest) (*UpdateCarePlanResourcesResponse, error) {
	var resource db.CarePlanResource
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		resource, txErr = q.UpdateCarePlanResource(ctx, db.UpdateCarePlanResourceParams{
			ID:                  resourceID,
			ResourceDescription: req.ResourceDescription,
			IsObtained:          req.IsObtained,
			ObtainedDate:        pgtype.Date{Time: req.ObtainedDate, Valid: true},
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanResource(ctx context.Context, resourceID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanResource(ctx, resourceID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteCarePlanResource", "Failed to delete care plan resource", zap.Error(err))
		return fmt.Errorf("failed to delete care plan resource: %w", err)
	}
	return nil
}

// ========================== Care plan Reports ===========================

func (s *carePlanService) CreateCarePlanReport(ctx context.Context, carePlanID uuid.UUID, employeeID uuid.UUID, req *CreateCarePlanReportRequest) (*CreateCarePlanReportResponse, error) {
	var report db.CarePlanReport
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		report, txErr = q.CreateCarePlanReport(ctx, db.CreateCarePlanReportParams{
			CarePlanID:          carePlanID,
			ReportType:          db.CarePlanReportTypeEnum(req.ReportType),
			ReportContent:       req.ReportContent,
			IsCritical:          req.IsCritical,
			CreatedByEmployeeID: employeeID,
		})
		return txErr
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

func (s *carePlanService) ListCarePlanReports(ctx *gin.Context, carePlanID uuid.UUID, req *ListCarePlanReportsRequest) (*pagination.Response[ListCarePlanReportsResponse], error) {
	params := req.GetParams()
	var reports []db.ListCarePlanReportsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		reports, txErr = q.ListCarePlanReports(ctx, db.ListCarePlanReportsParams{
			CarePlanID: carePlanID,
			Limit:      params.Limit,
			Offset:     params.Offset,
		})
		return txErr
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

func (s *carePlanService) UpdateCarePlanReport(ctx context.Context, reportID uuid.UUID, req *UpdateCarePlanReportRequest) (*UpdateCarePlanReportResponse, error) {
	var report db.CarePlanReport
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var txErr error
		report, txErr = q.UpdateCarePlanReport(ctx, db.UpdateCarePlanReportParams{
			ID:            reportID,
			ReportType:    db.NullCarePlanReportTypeFromPtr(req.ReportType),
			ReportContent: req.ReportContent,
			IsCritical:    req.IsCritical,
		})
		return txErr
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

func (s *carePlanService) DeleteCarePlanReport(ctx context.Context, reportID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteCarePlanReport(ctx, reportID)
	})
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
	topicID uuid.UUID,
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
	genCarePlan *ai.CarePlanResponse,
	assessmentID uuid.UUID,
	emplpoyeeID uuid.UUID,
) (carePlanID uuid.UUID, err error) {
	rawllmResp, err := json.Marshal(genCarePlan)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlan", "Failed to marshal care plan", zap.Error(err))
		return uuid.Nil, err
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
		return uuid.Nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "insertCarePlan", "Care plan created successfully", zap.String("carePlanID", carePlan.ID.String()))
	return carePlan.ID, nil
}

func (s *carePlanService) insertCarePlanObjectives(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID uuid.UUID,
	genCarePlan *ai.CarePlanResponse,
) error {
	for _, obj := range genCarePlan.Objectives {
		createdObj, err := qtx.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlanID,
			GoalTitle:   obj.Title,
			Description: obj.Description,
			Timeframe:   db.CarePlanTimeframeEnum(obj.Timeframe),
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true},
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create objective", zap.Error(err))
			return err
		}
		for i, action := range obj.Actions {
			_, err := qtx.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1),
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "insertCarePlanObjectives", "Failed to create action step", zap.Error(err))
				return err
			}
		}
	}
	return nil
}

func (s *carePlanService) insertCarePlanInterventions(
	ctx context.Context,
	qtx *db.Queries,
	carePlanID uuid.UUID,
	genCarePlan *ai.CarePlanResponse,
) error {
	for _, intervention := range genCarePlan.Interventions {
		_, err := qtx.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlanID,
			Frequency:               db.CarePlanInterventionFrequencyEnum(intervention.Frequency),
			InterventionDescription: intervention.Description,
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
	carePlanID uuid.UUID,
	genCarePlan *ai.CarePlanResponse,
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
	carePlanID uuid.UUID,
	genCarePlan *ai.CarePlanResponse,
) error {
	for _, risk := range genCarePlan.RiskFactors {
		_, err := qtx.CreateCarePlanRisk(ctx, db.CreateCarePlanRiskParams{
			CarePlanID:         carePlanID,
			RiskDescription:    risk.Risk,
			MitigationStrategy: risk.Mitigation,
			RiskLevel:          db.CarePlanRiskLevelEnum(risk.RiskLevel),
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
	carePlanID uuid.UUID,
	genCarePlan *ai.CarePlanResponse,
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
	carePlanID uuid.UUID,
	genCarePlan *ai.CarePlanResponse,
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
