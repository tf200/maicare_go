package clientp

import (
	"context"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Helper function to convert goals slice to JSON bytes
func goalsToJSON(goals []IntakeAssessmentGoal) ([]byte, error) {
	if goals == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(goals)
}

// Helper function to convert JSON bytes to goals slice
func jsonToGoals(data []byte) ([]IntakeAssessmentGoal, error) {
	if len(data) == 0 || string(data) == "null" {
		return []IntakeAssessmentGoal{}, nil
	}
	var goals []IntakeAssessmentGoal
	if err := json.Unmarshal(data, &goals); err != nil {
		return nil, err
	}
	return goals, nil
}

func (s *clientService) GenerateIntakeGoals(ctx context.Context, req *GenerateIntakeGoalsRequest) (*GenerateIntakeGoalsResponse, error) {
	if req.IntakeAssessmentID != uuid.Nil {
		assessment, err := s.Store.GetIntakeMaturityAssessment(ctx, req.IntakeAssessmentID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateIntakeGoals", "Failed to get intake assessment", zap.Error(err))
			return nil, err
		}
		req.TopicID = assessment.TopicID
		req.CurrentLevel = int(assessment.CurrentLevel)
		intakeForm, err := s.Store.GetIntakeForm(ctx, assessment.IntakeFormID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateIntakeGoals", "Failed to get intake form", zap.Error(err))
			return nil, err
		}
		req.RegistrationFormID = intakeForm.RegistrationFormID
	}

	registrationForm, err := s.Store.GetRegistrationForm(ctx, req.RegistrationFormID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateIntakeGoals", "Failed to get registration form", zap.Error(err))
		return nil, err
	}

	// 2. Get the registration form for client goals
	intakeForm, err := s.Store.GetIntakeFormByRegistrationFormID(ctx, registrationForm.ID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateIntakeGoals", "Failed to get intake form", zap.Error(err))
		return nil, err
	}

	// 3. Get level description from maturity matrix
	level, err := s.Store.GetLevelDescription(ctx, db.GetLevelDescriptionParams{
		ID:    req.TopicID,
		Level: fmt.Sprintf("%d", req.CurrentLevel),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "GenerateIntakeGoals", "Failed to get level description", zap.Error(err))
		// Continue with empty description if not found
	}

	// 4. Call AI Service
	clientGoals := ""
	if len(registrationForm.ClientGoals) > 0 {
		// We ignore the error here as it's a slice of strings/simple types
		goalsJSON, _ := json.Marshal(registrationForm.ClientGoals)
		clientGoals = string(goalsJSON)
	}

	riskContext := buildRiskContext(registrationForm, intakeForm)

	aiGoals, err := s.AIService.GenerateIntakeGoals(
		ctx,
		level.TopicName,
		int32(req.CurrentLevel),
		level.LevelDescription,
		req.UserDesc,
		clientGoals,
		riskContext,
	)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateIntakeGoals", "AI generation failed", zap.Error(err))
		return nil, err
	}

	// 5. Map to response
	responseGoals := make([]IntakeAssessmentGoal, 0, len(aiGoals))
	for _, g := range aiGoals {
		responseGoals = append(responseGoals, IntakeAssessmentGoal{
			Title:       g.Title,
			Description: g.Description,
			Priority:    g.Priority,
		})
	}

	return &GenerateIntakeGoalsResponse{
		Goals: responseGoals,
	}, nil
}

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
