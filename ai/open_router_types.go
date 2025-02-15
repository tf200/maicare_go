package ai

// OpenRouter Request thes is the request to be sent to the open router api

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type JSONSchema[T any] struct {
	Name   string `json:"name"`
	Strict bool   `json:"strict"`
	Schema T      `json:"schema"`
}
type ResponseFormat[T any] struct {
	Type       string `json:"type"`
	JSONSchema T      `json:"json_schema"`
}

type Schema[T any] struct {
	Type                 string   `json:"type"`
	Properties           T        `json:"properties"`
	Required             []string `json:"required"`
	AdditionalProperties bool     `json:"additionalProperties"`
}

type SchemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type OpenRouterRequest[T any] struct {
	Model          string     `json:"model"`
	Messages       []Messages `json:"messages"`
	ResponseFormat T          `json:"response_format"`
}

// OpenRouter Response is the response provided by the open router api

type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Refusal *string `json:"refusal"`
}

type Choice struct {
	Logprobs           interface{} `json:"logprobs"`
	FinishReason       string      `json:"finish_reason"`
	NativeFinishReason string      `json:"native_finish_reason"`
	Index              int         `json:"index"`
	Message            Message     `json:"message"`
}

type OpenRouterResponse struct {
	ID       string     `json:"id"`
	Provider string     `json:"provider"`
	Model    string     `json:"model"`
	Object   string     `json:"object"`
	Created  int64      `json:"created"`
	Choices  []Choice   `json:"choices"`
	Usage    TokenUsage `json:"usage"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
