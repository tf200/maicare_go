package ai

type AiHandler struct {
	OpenRouterAPIKey string
}

func NewAiHandler(openRouterAPIKey string) *AiHandler {
	return &AiHandler{
		OpenRouterAPIKey: openRouterAPIKey,
	}
}
