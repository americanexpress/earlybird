package scan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
)

type GptProvider struct{}

func (gpt GptProvider) setHeader(req *http.Request, value string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+value)
}

func (gpt GptProvider) buildRequestBody(cfg *cfgReader.EarlybirdConfig, scanJob LLMJob, chunk llmFileChunk) ([]byte, error) {
	if usesGPTResponsesAPI(cfg.LLMEndpoint) {
		requestBody := GptLlmResponsesRequest{
			Model: cfg.LLMModel,
			Input: []gptLlmResponsesInput{
				{
					Role: "system",
					Content: []gptLlmResponsesTextPart{{
						Type: "input_text",
						Text: llmSystemPrompt(cfg),
					}},
				},
				{
					Role: "user",
					Content: []gptLlmResponsesTextPart{{
						Type: "input_text",
						Text: buildLLMUserPrompt(scanJob, chunk),
					}},
				},
			},
			Text: gptLlmResponsesText{
				Format: gptLlmResponsesTextFormat{Type: "json_object"},
			},
		}

		return json.Marshal(requestBody)
	}

	requestBody := llmChatCompletionRequest{
		Model:       cfg.LLMModel,
		Temperature: 0,
		Messages: []llmMessage{
			{
				Role:    "system",
				Content: llmSystemPrompt(cfg),
			},
			{
				Role:    "user",
				Content: buildLLMUserPrompt(scanJob, chunk),
			},
		},
		ResponseFormat: &llmResponseFormat{Type: "json_object"},
	}

	return json.Marshal(requestBody)
}

func (gpt GptProvider) parseResponse(responseBody []byte) ([]LLMFinding, error) {
	var chatCompletion gptLlmChatCompletionResponse
	if err := json.Unmarshal(responseBody, &chatCompletion); err != nil {
		return nil, err
	}
	if chatCompletion.Error != nil {
		return nil, fmt.Errorf("gpt response error: %s", chatCompletion.Error.Message)
	}

	content := ""
	if len(chatCompletion.Choices) > 0 {
		content = chatCompletion.Choices[0].Message.Content
	} else {
		var responses gptLlmResponsesResponse
		if err := json.Unmarshal(responseBody, &responses); err != nil {
			return nil, err
		}
		if responses.Error != nil {
			return nil, fmt.Errorf("gpt response error: %s", responses.Error.Message)
		}
		if len(responses.Output) == 0 {
			return nil, fmt.Errorf("gpt response missing choices or output")
		}
		for _, part := range responses.Output[0].Content {
			if part.Text != "" {
				content += part.Text
			}
		}
	}

	var structured llmStructuredResponse
	if err := json.Unmarshal([]byte(content), &structured); err != nil {
		return nil, err
	}

	return structured.Findings, nil
}

func usesGPTResponsesAPI(endpoint string) bool {
	return strings.Contains(strings.ToLower(endpoint), "/responses")
}
