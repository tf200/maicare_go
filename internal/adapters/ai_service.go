package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	"maicare_go/internal/domain"
)

var ErrAIServiceNotConfigured = errors.New("ai service is not configured")

type AIService struct {
	llm llms.Model
}

type AIServiceStub = AIService

func NewAIServiceStub() *AIService {
	return &AIService{}
}

func NewAIService(apiKey string, modelName string) (*AIService, error) {
	if apiKey == "" {
		return &AIService{}, nil
	}

	if modelName == "" {
		modelName = "openai/gpt-4o-mini"
	}

	llm, err := openai.New(
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithToken(apiKey),
		openai.WithModel(modelName),
		openai.WithHTTPClient(&http.Client{}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize openrouter ai service: %w", err)
	}

	return &AIService{
		llm: llm,
	}, nil
}

func (s *AIService) GenerateCarePlan(ctx context.Context, clientData interface{}) (*domain.CarePlanResponse, error) {
	if s == nil || s.llm == nil {
		return nil, ErrAIServiceNotConfigured
	}

	systemPrompt := `You are a professional care planner. Generate a structured care plan in JSON format based on the client data provided.
The output MUST be a valid JSON object matching this schema:
{
    "assessment_summary": "string",
    "objectives": [
        {
            "title": "string",
            "description": "string",
            "timeframe": "short_term | medium_term | long_term",
            "actions": ["string"]
        }
    ],
    "interventions": [
        {
            "description": "string",
            "frequency": "daily | weekly | monthly"
        }
    ]
}`
	userPromptBytes, err := json.Marshal(clientData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal client data: %w", err)
	}
	userPrompt := string(userPromptBytes)

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	resp, err := s.llm.GenerateContent(ctx, content, llms.WithJSONMode())
	if err != nil {
		return nil, err
	}

	var result domain.CarePlanResponse
	if len(resp.Choices) > 0 {
		err = json.Unmarshal([]byte(resp.Choices[0].Content), &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal care plan response: %w", err)
		}
	}

	return &result, nil
}

func (s *AIService) GenerateIntakeGoals(ctx context.Context, topicName string, currentLevel int32, levelDescription string, userDesc string, clientGoals string, riskContext string) ([]domain.AIGeneratedIntakeGoal, error) {
	if s == nil || s.llm == nil {
		return nil, ErrAIServiceNotConfigured
	}

	systemPrompt := `You are a professional care planner. Based on the client's current maturity level in a specific topic, generate 3-5 concrete Care Plan objectives.
Each objective MUST have a title, a brief description, and a priority (high, medium, or low).
Use the following priority order of inputs:
1) User description (highest priority)
2) Risk context (must be considered for safety and urgency)
3) Relevant desired goals only (ignore goals not related to the topic)
4) Topic level description
When desired goals are not relevant to the topic, exclude them.
The output MUST be a valid JSON array of objects matching this schema:
[
  {
    "title": "string",
    "description": "string",
    "priority": "high | medium | low"
  }
]`

	userPrompt := fmt.Sprintf(`Topic: %s
Current Level: %d
Level Description: %s
User Description: %s
Client's Desired Goals: %s
Risk Context: %s

Suggest structured Care Plan objectives for this client.`, topicName, currentLevel, levelDescription, userDesc, clientGoals, riskContext)

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	resp, err := s.llm.GenerateContent(ctx, content, llms.WithJSONMode())
	if err != nil {
		return nil, err
	}

	var result []domain.AIGeneratedIntakeGoal
	if len(resp.Choices) > 0 {
		err = json.Unmarshal([]byte(resp.Choices[0].Content), &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal intake goals response: %w", err)
		}
	}

	return result, nil
}

func (s *AIService) GenerateAutoReports(ctx context.Context, pastReportsText string) (string, error) {
	if s == nil || s.llm == nil {
		return "", ErrAIServiceNotConfigured
	}
	return "", ErrAIServiceNotConfigured
}

var _ domain.AIService = (*AIService)(nil)
var _ domain.AutoReportGenerator = (*AIService)(nil)

