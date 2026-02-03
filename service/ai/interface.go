package ai

import (
	"context"
)

type CarePlanResponse struct {
	AssessmentSummary string `json:"assessment_summary"`
	Objectives        []struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Timeframe   string   `json:"timeframe"` // short_term, medium_term, long_term
		Actions     []string `json:"actions"`
	} `json:"objectives"`
	Interventions []struct {
		Description string `json:"description"`
		Frequency   string `json:"frequency"` // daily, weekly, monthly
	} `json:"interventions"`
	SuccessMetrics []struct {
		Metric            string `json:"metric"`
		Target            string `json:"target"`
		MeasurementMethod string `json:"measurement_method"`
	} `json:"success_metrics"`
	RiskFactors []struct {
		Risk       string `json:"risk"`
		Mitigation string `json:"mitigation"`
		RiskLevel  string `json:"risk_level"` // low, medium, high
	} `json:"risk_factors"`
	SupportNetwork []struct {
		Role           string `json:"role"`
		Responsibility string `json:"responsibility"`
	} `json:"support_network"`
	ResourcesRequired []string `json:"resources_required"`
}

type AIService interface {
	GenerateCarePlan(ctx context.Context, clientData interface{}) (*CarePlanResponse, error)
	GenerateIntakeGoals(ctx context.Context, topicName string, currentLevel int32, levelDescription string, userDesc string, clientGoals string, riskContext string) ([]struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
	}, error)
}
