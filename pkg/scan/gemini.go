package scan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
)

type GeminiProvider struct{}

func isGeminiModel(model string) bool {
	return strings.Contains(strings.ToLower(model), "gemini")
}

func (gemini GeminiProvider) setHeader(req *http.Request, value string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", value)
}

func (gemini GeminiProvider) buildRequestBody(cfg *cfgReader.EarlybirdConfig, scanJob LLMJob, chunk llmFileChunk) ([]byte, error) {

	requestBody := GeminiRequestBody{
		Contents: []GeminiContent{
			{
				Role:  "user",
				Parts: []Parts{{Text: llmSystemPrompt(cfg)}},
			},
			{
				Role:  "user",
				Parts: []Parts{{Text: buildLLMUserPrompt(scanJob, chunk)}},
			},
		},
	}
	// fmt.Println("request body", requestBody)
	return json.Marshal(requestBody)
}

func (gemini GeminiProvider) parseResponse(responseBody []byte) ([]LLMFinding, error) {
	// var completion any
	var completion GeminiResponseBody
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return nil, err
	}
	if len(completion.Candidates) == 0 {
		return nil, fmt.Errorf("gemini response missing candidates")
	}
	// fmt.Print("parsed gemini response: ", string(completion.Candidates[0].Content))

	content, err := extractGeminiCandidateText(completion.Candidates[0].Content)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Gemini response: ", content)

	var structured llmStructuredResponse
	if err := json.Unmarshal([]byte(content), &structured); err != nil {
		return nil, err
	}
	// fmt.Println("structured response: %+v", structured)
	// return nil, nil

	return structured.Findings, nil
}

func extractGeminiCandidateText(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}

	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return normalizeGeminiJSONPayload(text), nil
	}

	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "", err
	}

	if t, ok := extractGeminiTextFromObject(obj); ok {
		return normalizeGeminiJSONPayload(t), nil
	}

	return "", fmt.Errorf("unsupported gemini candidate content format")
}

func normalizeGeminiJSONPayload(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return content
	}

	if strings.HasPrefix(content, "```") {
		newline := strings.Index(content, "\n")
		if newline >= 0 {
			content = strings.TrimSpace(content[newline+1:])
		}
		if strings.HasSuffix(content, "```") {
			content = strings.TrimSpace(strings.TrimSuffix(content, "```"))
		}
	}

	if strings.HasPrefix(strings.ToLower(content), "json") {
		newline := strings.Index(content, "\n")
		if newline >= 0 {
			firstLine := strings.TrimSpace(content[:newline])
			if strings.EqualFold(firstLine, "json") {
				content = strings.TrimSpace(content[newline+1:])
			}
		}
	}

	if strings.HasPrefix(content, "`") && strings.HasSuffix(content, "`") {
		content = strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(content, "`"), "`"))
	}

	return strings.TrimSpace(content)
}

func extractGeminiTextFromObject(obj map[string]any) (string, bool) {
	if text, ok := obj["text"]; ok {
		if s, ok := text.(string); ok && s != "" {
			return s, true
		}
	}

	if content, ok := obj["content"]; ok {
		switch c := content.(type) {
		case string:
			return c, true
		case map[string]any:
			if text, ok := c["text"].(string); ok && text != "" {
				return text, true
			}
			if parts, ok := c["parts"]; ok {
				return extractGeminiTextFromParts(parts)
			}
		}
	}

	if parts, ok := obj["parts"]; ok {
		return extractGeminiTextFromParts(parts)
	}

	return "", false
}

func extractGeminiTextFromParts(raw any) (string, bool) {
	parts, ok := raw.([]any)
	if !ok {
		return "", false
	}

	var builder strings.Builder
	for _, part := range parts {
		if p, ok := part.(map[string]any); ok {
			if text, ok := p["text"].(string); ok {
				builder.WriteString(text)
			}
		}
	}

	if builder.Len() == 0 {
		return "", false
	}

	return builder.String(), true
}
