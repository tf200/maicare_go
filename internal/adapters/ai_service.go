package adapters

import (
	"context"
	"errors"

	"maicare_go/internal/domain"
)

var ErrAIServiceNotConfigured = errors.New("ai service is not configured")

// AIServiceStub is a compile-time adapter used until the OpenRouter-backed
// implementation is ported from the legacy service/ai package.
type AIServiceStub struct{}

func NewAIServiceStub() *AIServiceStub {
	return &AIServiceStub{}
}

func (s *AIServiceStub) GenerateCarePlan(ctx context.Context, clientData interface{}) (*domain.CarePlanResponse, error) {
	return nil, ErrAIServiceNotConfigured
}

func (s *AIServiceStub) GenerateIntakeGoals(ctx context.Context, topicName string, currentLevel int32, levelDescription string, userDesc string, clientGoals string, riskContext string) ([]domain.AIGeneratedIntakeGoal, error) {
	return nil, ErrAIServiceNotConfigured
}

func (s *AIServiceStub) GenerateAutoReports(ctx context.Context, pastReportsText string) (string, error) {
	return "", ErrAIServiceNotConfigured
}

var _ domain.AIService = (*AIServiceStub)(nil)
var _ domain.AutoReportGenerator = (*AIServiceStub)(nil)
