package ai

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"encoding/json"
)

/*
These types represent the desired response format from the llm

	set Strict to true for the for lower chance of failure

	LLMSpellingResponseFormat is the format of the response from the LLM

	this will be sent in the request to the open router api
*/

type LLMSpellingResponseFormat struct { // here
	CorrectedText  SchemaProperty `json:"corrected_text"`
	CorrectedWords SchemaProperty `json:"corrected_words"`
}

/* CorrectedContent Represents the final output taken after parsing the llm response */

type CorrectedContent struct {
	CorrectedText string `json:"corrected_text"`
}

func (ai *AiHandler) SpellingCheck(text string, model string) (*CorrectedContent, error) {
	request := OpenRouterRequest[ResponseFormat[JSONSchema[Schema[LLMSpellingResponseFormat]]]]{
		Model: model,
		Messages: []Messages{
			{
				Role:    "user",
				Content: fmt.Sprintf("correct the spelling : %s"+"the response need to in json format in this exact structure {\"corrected_text\" : \"text here\", \"corrected_words\" : \"changed words here\" }", text),
			},
		},
		ResponseFormat: ResponseFormat[JSONSchema[Schema[LLMSpellingResponseFormat]]]{
			Type: "json_schema",
			JSONSchema: JSONSchema[Schema[LLMSpellingResponseFormat]]{
				Name:   "correct_spelling",
				Strict: true,
				Schema: Schema[LLMSpellingResponseFormat]{
					Type: "object",
					Properties: LLMSpellingResponseFormat{
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
	body, err := io.ReadAll(resp.Body)
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

type LLMAutoReportsResponseFormat struct { // here
	GeneratedReport SchemaProperty `json:"generated_report"`
}

type AutoReportsContent struct {
	GeneratedReport string `json:"generated_report"`
}

func (ai *AiHandler) GenerateAutoReports(text string, model string) (*AutoReportsContent, error) {
	request := OpenRouterRequest[ResponseFormat[JSONSchema[Schema[LLMAutoReportsResponseFormat]]]]{
		Model: model,
		Messages: []Messages{
			{
				Role: "user",
				Content: fmt.Sprintf(`you are given some past reports of one of our clients your Job is to make from 
										them a report that summarizes all of them make sure the response is in json format
									, here are the prevois reports  %s`, text),
			},
		},
		ResponseFormat: ResponseFormat[JSONSchema[Schema[LLMAutoReportsResponseFormat]]]{
			Type: "json_schema",
			JSONSchema: JSONSchema[Schema[LLMAutoReportsResponseFormat]]{
				Name:   "correct_spelling",
				Strict: true,
				Schema: Schema[LLMAutoReportsResponseFormat]{
					Type: "object",
					Properties: LLMAutoReportsResponseFormat{
						GeneratedReport: SchemaProperty{
							Type:        "string",
							Description: "The generated report",
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
	body, err := io.ReadAll(resp.Body)
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

	var report AutoReportsContent
	err = json.Unmarshal([]byte(openRouterResponse.Choices[0].Message.Content), &report)
	if err != nil {
		log.Printf("Error unmarshalling corrected content: %v", err)
		return nil, err
	}

	return &report, nil
}
