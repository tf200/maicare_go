package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type intakeFormService struct {
	repo      domain.IntakeFormRepository
	logger    domain.Logger
	aiService domain.AIService
}

func NewIntakeFormService(repo domain.IntakeFormRepository, logger domain.Logger, aiService domain.AIService) domain.IntakeFormService {
	return &intakeFormService{repo: repo, logger: logger, aiService: aiService}
}

func (s *intakeFormService) CreateIntakeForm(ctx context.Context, params domain.CreateIntakeFormParams) (*domain.IntakeForm, error) {
	// Validate registration form exists and is processed
	regForm, err := s.repo.GetRegistrationForm(ctx, params.RegistrationFormID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.CreateIntakeForm", "failed to get registration form", err, zap.String("registration_form_id", params.RegistrationFormID.String()))
		}
		return nil, fmt.Errorf("registration form not found: %w", err)
	}

	if regForm.FormStatus != db.FormStatusEnumProcessed {
		return nil, errors.New("intake can only be created for processed registration forms")
	}

	intakeForm, err := s.repo.CreateIntakeForm(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.CreateIntakeForm", "failed to create intake form", err, zap.String("registration_form_id", params.RegistrationFormID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.CreateIntakeForm", "intake form created successfully", zap.String("intake_form_id", intakeForm.ID.String()))
	}

	return intakeForm, nil
}

func (s *intakeFormService) ListIntakeForms(ctx context.Context, params domain.ListIntakeFormsParams) (*domain.ListResult[domain.IntakeFormListItem], error) {
	items, totalCount, err := s.repo.ListIntakeForms(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.ListIntakeForms", "failed to list intake forms", err)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.ListIntakeForms", "intake forms listed successfully", zap.Int64("total_count", totalCount))
	}

	return &domain.ListResult[domain.IntakeFormListItem]{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

func (s *intakeFormService) GetIntakeFormTotals(ctx context.Context) (*domain.IntakeFormTotals, error) {
	totals, err := s.repo.GetIntakeFormTotals(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.GetIntakeFormTotals", "failed to get intake form totals", err)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.GetIntakeFormTotals", "intake form totals retrieved successfully")
	}

	return totals, nil
}

func (s *intakeFormService) GetIntakeForm(ctx context.Context, id uuid.UUID) (*domain.IntakeFormDetail, error) {
	detail, err := s.repo.GetIntakeFormDetail(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.GetIntakeForm", "failed to get intake form detail", err, zap.String("intake_form_id", id.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.GetIntakeForm", "intake form detail retrieved successfully", zap.String("intake_form_id", id.String()))
	}

	return detail, nil
}

func (s *intakeFormService) UpdateIntakeForm(ctx context.Context, params domain.UpdateIntakeFormParams) (*domain.IntakeForm, error) {
	// Validate SelfSufficiency if provided
	if params.SelfSufficiency != nil && (*params.SelfSufficiency < 0 || *params.SelfSufficiency > 5) {
		return nil, domain.ErrInvalidIntakeFormUpdate
	}

	// Validate EvaluationIntervalsWeeks if provided
	if params.EvaluationIntervalsWeeks != nil && *params.EvaluationIntervalsWeeks < 0 {
		return nil, domain.ErrInvalidIntakeFormUpdate
	}

	// Check if there are fields to update
	if params.DateOfIntake == nil && params.CareType == nil && params.IntakeParticipants == nil &&
		params.FamilySituation == nil && params.PsychologicalState == nil && params.SelfSufficiency == nil &&
		params.SenderID == nil && params.AssignedLocationID == nil && params.RiskAssessment == nil &&
		params.EvaluationIntervalsWeeks == nil && params.Signature == nil && len(params.ClearFields) == 0 {
		return nil, domain.ErrNoIntakeFormFieldsToUpdate
	}

	intakeForm, err := s.repo.UpdateIntakeForm(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.UpdateIntakeForm", "failed to update intake form", err, zap.String("intake_form_id", params.ID.String()))
		}

		// Handle domain errors
		if errors.Is(err, domain.ErrIntakeFormUpdateBlockedByActiveClient) ||
			errors.Is(err, domain.ErrIntakeFormUpdateConflict) {
			return nil, err
		}

		// Check for pgx.ErrNoRows
		if errors.Is(err, pgx.ErrNoRows) {
			hasActiveClient, activeClientErr := s.repo.HasActiveClientByIntakeFormID(ctx, &params.ID)
			if activeClientErr == nil && hasActiveClient {
				return nil, domain.ErrIntakeFormUpdateBlockedByActiveClient
			}

			_, getErr := s.repo.GetIntakeForm(ctx, params.ID)
			if getErr != nil {
				return nil, domain.ErrIntakeFormNotFound
			}

			return nil, domain.ErrIntakeFormUpdateConflict
		}

		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.UpdateIntakeForm", "intake form updated successfully", zap.String("intake_form_id", params.ID.String()))
	}

	return intakeForm, nil
}

func (s *intakeFormService) ReplaceIntakeFormGoals(ctx context.Context, intakeFormID uuid.UUID, params domain.ReplaceIntakeFormGoalsParams) (*domain.IntakeFormGoalsResult, error) {
	result, err := s.repo.ReplaceIntakeFormGoals(ctx, intakeFormID, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.ReplaceIntakeFormGoals", "failed to replace intake form goals", err, zap.String("intake_form_id", intakeFormID.String()))
		}

		if errors.Is(err, domain.ErrIntakeGoalsUpdateBlockedByActiveClient) {
			return nil, err
		}

		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.ReplaceIntakeFormGoals", "intake form goals replaced successfully", zap.String("intake_form_id", intakeFormID.String()))
	}

	return result, nil
}

func (s *intakeFormService) GenerateIntakeGoals(ctx context.Context, params domain.GenerateIntakeGoalsParams) (*domain.GenerateIntakeGoalsResult, error) {
	// If IntakeAssessmentID is provided, fetch assessment and derive TopicID/CurrentLevel and RegistrationFormID
	if params.IntakeAssessmentID != uuid.Nil {
		assessment, err := s.repo.GetIntakeMaturityAssessment(ctx, params.IntakeAssessmentID)
		if err != nil {
			if s.logger != nil {
				s.logger.LogError(ctx, "IntakeFormService.GenerateIntakeGoals", "failed to get intake assessment", err, zap.String("assessment_id", params.IntakeAssessmentID.String()))
			}
			return nil, err
		}
		params.TopicID = assessment.TopicID
		params.CurrentLevel = int(assessment.CurrentLevel)

		intakeForm, err := s.repo.GetIntakeForm(ctx, assessment.IntakeFormID)
		if err != nil {
			if s.logger != nil {
				s.logger.LogError(ctx, "IntakeFormService.GenerateIntakeGoals", "failed to get intake form", err)
			}
			return nil, err
		}
		params.RegistrationFormID = intakeForm.RegistrationFormID
	}

	// If IntakeFormID is provided, fetch intake form and derive RegistrationFormID
	if params.IntakeFormID != uuid.Nil {
		intakeForm, err := s.repo.GetIntakeForm(ctx, params.IntakeFormID)
		if err != nil {
			if s.logger != nil {
				s.logger.LogError(ctx, "IntakeFormService.GenerateIntakeGoals", "failed to get intake form", err, zap.String("intake_form_id", params.IntakeFormID.String()))
			}
			return nil, err
		}
		params.RegistrationFormID = intakeForm.RegistrationFormID
	}

	// Get registration form
	registrationForm, err := s.repo.GetRegistrationForm(ctx, params.RegistrationFormID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.GenerateIntakeGoals", "failed to get registration form", err, zap.String("registration_form_id", params.RegistrationFormID.String()))
		}
		return nil, err
	}

	// Get intake form by registration form ID
	intakeForm, err := s.repo.GetIntakeFormByRegistrationFormID(ctx, registrationForm.ID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.GenerateIntakeGoals", "failed to get intake form", err, zap.String("registration_form_id", params.RegistrationFormID.String()))
		}
		return nil, err
	}

	// Get topic level
	topicName := ""
	levelDesc := ""
	level, err := s.repo.GetTopicLevel(ctx, db.GetTopicLevelParams{
		ID:    params.TopicID,
		Level: int32(params.CurrentLevel),
	})
	if err == nil {
		topicName = level.TopicName
		levelDesc = level.LevelDescription
	}

	// Marshal client goals to JSON string if any
	clientGoals := ""
	if len(registrationForm.ClientGoals) > 0 {
		goalsJSON, _ := json.Marshal(registrationForm.ClientGoals)
		clientGoals = string(goalsJSON)
	}

	// Build risk context
	riskContext := buildRiskContext(registrationForm, intakeForm)

	// Call AI service
	aiGoals, err := s.aiService.GenerateIntakeGoals(
		ctx,
		topicName,
		int32(params.CurrentLevel),
		levelDesc,
		params.UserDesc,
		clientGoals,
		riskContext,
	)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.GenerateIntakeGoals", "AI generation failed", err)
		}
		return nil, err
	}

	// Map AI goals to domain.IntakeAssessmentGoal
	responseGoals := make([]domain.IntakeAssessmentGoal, 0, len(aiGoals))
	for _, g := range aiGoals {
		responseGoals = append(responseGoals, domain.IntakeAssessmentGoal{
			Title:       g.Title,
			Description: g.Description,
			Priority:    g.Priority,
		})
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.GenerateIntakeGoals", "goals generated successfully", zap.Int("goals_count", len(responseGoals)))
	}

	return &domain.GenerateIntakeGoalsResult{Goals: responseGoals}, nil
}

func (s *intakeFormService) UpdateIntakeConclusion(ctx context.Context, id uuid.UUID, params domain.UpdateIntakeConclusionParams) (*domain.IntakeFormConclusion, error) {
	conclusion, err := s.repo.UpdateIntakeConclusion(ctx, id, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.UpdateIntakeConclusion", "failed to update intake conclusion", err, zap.String("intake_form_id", id.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.UpdateIntakeConclusion", "intake conclusion updated successfully", zap.String("intake_form_id", id.String()))
	}

	return conclusion, nil
}

func (s *intakeFormService) PromoteIntakeToClient(ctx context.Context, params domain.PromoteIntakeToClientParams) (*domain.PromoteIntakeToClientResult, error) {
	result, err := s.repo.PromoteIntakeToClient(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "IntakeFormService.PromoteIntakeToClient", "failed to promote intake to client", err, zap.String("intake_form_id", params.IntakeFormID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "IntakeFormService.PromoteIntakeToClient", "intake promoted to client successfully", zap.String("intake_form_id", params.IntakeFormID.String()), zap.String("client_id", result.ClientID.String()))
	}

	return result, nil
}

// Helper functions

func buildRiskContext(registrationForm db.GetRegistrationFormRow, intakeForm db.IntakeForm) string {
	risks := []string{}

	if isTruePtr(registrationForm.RiskAggressiveBehavior) {
		risks = append(risks, "aggressive behavior")
	}
	if isTruePtr(registrationForm.RiskSuicidalSelfharm) {
		risks = append(risks, "suicidal/self-harm")
	}
	if isTruePtr(registrationForm.RiskSubstanceAbuse) {
		risks = append(risks, "substance abuse")
	}
	if isTruePtr(registrationForm.RiskPsychiatricIssues) {
		risks = append(risks, "psychiatric issues")
	}
	if isTruePtr(registrationForm.RiskCriminalHistory) {
		risks = append(risks, "criminal history")
	}
	if isTruePtr(registrationForm.RiskFlightBehavior) {
		risks = append(risks, "flight behavior")
	}
	if isTruePtr(registrationForm.RiskWeaponPossession) {
		risks = append(risks, "weapon possession")
	}
	if isTruePtr(registrationForm.RiskSexualBehavior) {
		risks = append(risks, "sexual behavior")
	}
	if isTruePtr(registrationForm.RiskDayNightRhythm) {
		risks = append(risks, "day-night rhythm")
	}
	if isTruePtr(registrationForm.RiskOther) {
		otherDesc := strings.TrimSpace(ptrString(registrationForm.RiskOtherDescription))
		if otherDesc != "" {
			risks = append(risks, fmt.Sprintf("other: %s", otherDesc))
		} else {
			risks = append(risks, "other")
		}
	}

	additionalNotes := strings.TrimSpace(ptrString(registrationForm.RiskAdditionalNotes))
	intakeRisk := strings.TrimSpace(ptrString(intakeForm.RiskAssessment))

	if len(risks) == 0 && additionalNotes == "" && intakeRisk == "" {
		return "none reported"
	}

	parts := []string{}
	if len(risks) > 0 {
		parts = append(parts, "flags: "+strings.Join(risks, ", "))
	}
	if additionalNotes != "" {
		parts = append(parts, "additional notes: "+additionalNotes)
	}
	if intakeRisk != "" {
		parts = append(parts, "risk assessment: "+intakeRisk)
	}

	return strings.Join(parts, " | ")
}

func isTruePtr(value *bool) bool {
	return value != nil && *value
}

func ptrString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}