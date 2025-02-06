package ai

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"
)

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type SpellCheckFormat struct {
	Type       string     `json:"type"`
	JSONSchema JSONSchema `json:"json_schema"`
}

type JSONSchema struct {
	Name   string `json:"name"`
	Strict bool   `json:"strict"`
	Schema Schema `json:"schema"`
}

type Schema struct {
	Type                 string           `json:"type"`
	Properties           SchemaProperties `json:"properties"`
	Required             []string         `json:"required"`
	AdditionalProperties bool             `json:"additionalProperties"`
}

type SchemaProperties struct {
	CorrectedText  SchemaProperty `json:"corrected_text"`
	CorrectedWords SchemaProperty `json:"corrected_words"`
}

type SchemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type OpenRouterRequest struct {
	Model          string           `json:"model"`
	Messages       []Messages       `json:"messages"`
	ResponseFormat SpellCheckFormat `json:"response_format"`
}

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
type CorrectedContent struct {
	CorrectedText string `json:"corrected_text"`
}

func (ai *AiHandler) SpellingCheck(text string, model string) (*CorrectedContent, error) {
	request := OpenRouterRequest{
		Model: model,
		Messages: []Messages{
			{
				Role:    "user",
				Content: fmt.Sprintf("correct the spelling : %s"+"the response need to in json format in this exact structure {\"corrected_text\" : \"text here\", \"corrected_words\" : \"changed words here\" }", text),
			},
		},
		ResponseFormat: SpellCheckFormat{
			Type: "json_schema",
			JSONSchema: JSONSchema{
				Name:   "correct_spelling",
				Strict: true,
				Schema: Schema{
					Type: "object",
					Properties: SchemaProperties{
						CorrectedText: SchemaProperty{
							Type:        "string",
							Description: "The corrected text",
						},
						CorrectedWords: SchemaProperty{
							Type:        "array",
							Description: "The words that were corrected",
						},
					},
					Required:             []string{"corrected_text", "corrected_words"},
					AdditionalProperties: false,
				},
			},
		},
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ai.OpenRouterAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, string(body))
	}

	var openRouterResponse OpenRouterResponse
	err = json.Unmarshal(body, &openRouterResponse)
	if err != nil {
		return nil, err
	}

	if len(openRouterResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	var correctedContent CorrectedContent
	err = json.Unmarshal([]byte(openRouterResponse.Choices[0].Message.Content), &correctedContent)
	if err != nil {
		log.Printf("Error unmarshalling corrected content: %v", err)
		return nil, err
	}

	return &correctedContent, nil
}
