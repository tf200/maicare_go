package ai

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
	CorrectedText string `json:"generated_objectives"`
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
		return nil, fmt.Errorf("no choices returned: %v", string(body))
	}

	var correctedContent CorrectedContent
	err = json.Unmarshal([]byte(openRouterResponse.Choices[0].Message.Content), &correctedContent)
	if err != nil {
		log.Printf("Error unmarshalling corrected content: %v", err)
		return nil, err
	}

	return &correctedContent, nil
}

type LLMAutoReportsResponseFormat struct {
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
				Name:   "report_text",
				Strict: true,
				Schema: Schema[LLMAutoReportsResponseFormat]{
					Type: "object",
					Properties: LLMAutoReportsResponseFormat{
						GeneratedReport: SchemaProperty{
							Type:        "string",
							Description: "The generated report",
						},
					},
					Required:             []string{"generated_report"},
					AdditionalProperties: false,
				},
			},
		},
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+ai.OpenRouterAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, string(body))
	}

	var openRouterResponse OpenRouterResponse
	err = json.Unmarshal(body, &openRouterResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if len(openRouterResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned : %v req : %v", string(body), string(jsonRequest))
	}

	var report AutoReportsContent
	err = json.Unmarshal([]byte(openRouterResponse.Choices[0].Message.Content), &report)
	if err != nil {
		log.Printf("Error unmarshalling llm json content: %v", err)
		return nil, fmt.Errorf("error unmarshalling llm json content: %v, %v", err, string(body))
	}

	return &report, nil
}

// ========================================================================================================
type Objectives struct {
	ObjectiveDescription string `json:"objective_description"`
	DueDate              string `json:"due_date"`
}

type GeneratedObjectivesContent struct {
	GeneratedObjectives []Objectives `json:"generated_objectives"`
}

func (ai *AiHandler) GenerateObjectives(levelDescription, goal, description, startDate, endDate string, model string) (*GeneratedObjectivesContent, error) {
	generalDescriptionEscaped := strings.ReplaceAll(levelDescription, "\"", "\\\"")
	generalDescriptionEscaped = strings.ReplaceAll(generalDescriptionEscaped, "\n", "\\n")

	goalEscaped := strings.ReplaceAll(goal, "\"", "\\\"")
	goalEscaped = strings.ReplaceAll(goalEscaped, "\n", "\\n")

	descriptionEscaped := strings.ReplaceAll(description, "\"", "\\\"")
	descriptionEscaped = strings.ReplaceAll(descriptionEscaped, "\n", "\\n")

	startDateEscaped := strings.ReplaceAll(startDate, "\"", "\\\"")
	endDateEscaped := strings.ReplaceAll(endDate, "\"", "\\\"")

	// Now construct your request with the escaped strings
	request := fmt.Sprintf(`{
		"model": "%s",
		"messages": [
		  {
			"role": "user",
			"content": "Ok, in our company, we manage our clients. One thing we do is assess their problems. The way we do this is by first assigning them a level from 1 to 5 based on their situation, corresponding to a domain.\n\nYou will be given:\n\nA description of their current general level\nA short description of their current situation\nA short description of what we are trying to achieve\nA start date and an end date for the goal to be achieved\n\nI need you to generate a set of objectives they can work toward to achieve that goal.\n\n{\n  \\\"general_description\\\": \\\"%s\\\",\n  \\\"goal\\\": \\\"%s\\\",\n  \\\"description\\\": \\\"%s\\\",\n  \\\"start_date\\\": \\\"%s\\\",\n  \\\"end_date\\\": \\\"%s\\\"\n}\n\nExpected Response Format: The response should be in JSON format, providing a structured list of objectives to help the client progress toward the goal. if the given data does not allow you to achieve this respond with empty array \n\nExample response format:\n{\n  \\\"generated_objectives\\\": [\n    {\n      \\\"objective_description\\\": \\\"Start paying off debt\\\",\n      \\\"due_date\\\": \\\"2023-01-01\\\"\n    },\n    {\n      \\\"objective_description\\\": \\\"Find a job\\\",\n      \\\"due_date\\\": \\\"2023-01-01\\\"\n    }\n  ]\n}"
		  }
		]
	}`, model, generalDescriptionEscaped, goalEscaped, descriptionEscaped, startDateEscaped, endDateEscaped)

	jsonRequest := []byte(request)

	req, err := http.NewRequest(http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+ai.OpenRouterAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, string(body))
	}

	var openRouterResponse OpenRouterResponse
	err = json.Unmarshal(body, &openRouterResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v, %v ", err, string(body))
	}

	if len(openRouterResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned : %v req : %v", string(body), string(jsonRequest))
	}

	// Get the content from the response
	content := openRouterResponse.Choices[0].Message.Content

	// Remove markdown code block indicators if they exist
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Now unmarshal the cleaned content
	var objectives GeneratedObjectivesContent
	err = json.Unmarshal([]byte(content), &objectives)
	if err != nil {
		log.Printf("Error unmarshalling llm json content: %v", err)
		return nil, fmt.Errorf("error unmarshalling llm json content: %v, content: %v", err, content)
	}

	return &objectives, nil

}
