package scan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
	)

type llmHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var defaultLLMHTTPClient llmHTTPClient = &http.Client{}

type llmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmResponseFormat struct {
	Type string `json:"type"`
}

type llmChatCompletionRequest struct {
	Model          string              `json:"model"`
	Messages       []llmMessage        `json:"messages"`
	Temperature    int                 `json:"temperature"`
	ResponseFormat *llmResponseFormat  `json:"response_format,omitempty"`
}

type llmChatCompletionChoice struct {
	Message llmMessage `json:"message"`
}

type llmAPIError struct {
	Message string `json:"message"`
}

type llmChatCompletionResponse struct {
	Choices []llmChatCompletionChoice `json:"choices"`
	Error   *llmAPIError              `json:"error,omitempty"`
}

type LLMFinding struct {
	Line           int    `json:"line"`
	CredentialType string `json:"credential_type"`
	Confidence     string `json:"confidence"`
	Candidate      string `json:"candidate"`
	Reason         string `json:"reason"`
}

type llmStructuredResponse struct {
	Findings []LLMFinding `json:"findings"`
}

type llmFileChunk struct {
	StartLine int
	Lines     []string
}

func llm_scan(cfg *cfgReader.EarlybirdConfig, scanjob LLMJob) ([]LLMFinding, error) {
	if err := validateLLMConfig(cfg); err != nil {
		return nil, err
	}
	if len(scanjob.FileLines) == 0 {
		return nil, nil
	}

	chunks := buildLLMFileChunks(scanjob.FileLines, cfg.LLMMaxLines, cfg.LLMMaxBytes)
	findings := make([]LLMFinding, 0)

	for _, chunk := range chunks {
		chunkFindings, err := callLLMChunk(defaultLLMHTTPClient, cfg, scanjob, chunk)
		if err != nil {
			return findings, err
		}
		findings = append(findings, chunkFindings...)
	}

	logLLMFindings(cfg, scanjob, findings)
	return findings, nil
}

func validateLLMConfig(cfg *cfgReader.EarlybirdConfig) error {
	if cfg == nil {
		return fmt.Errorf("llm scan config is nil")
	}
	if !cfg.EnableLLMScan {
		return fmt.Errorf("llm scan is disabled")
	}
	if cfg.LLMEndpoint == "" {
		return fmt.Errorf("llm endpoint is required")
	}
	if cfg.LLMAPIKey == "" {
		return fmt.Errorf("llm api key is required: set EARLYBIRD_LLM_API_KEY or OPENAI_API_KEY")
	}
	if cfg.LLMModel == "" {
		return fmt.Errorf("llm model is required")
	}
	if cfg.LLMTimeoutSeconds <= 0 {
		return fmt.Errorf("llm timeout must be greater than zero")
	}
	if cfg.LLMMaxLines <= 0 {
		return fmt.Errorf("llm max lines must be greater than zero")
	}
	if cfg.LLMMaxBytes <= 0 {
		return fmt.Errorf("llm max bytes must be greater than zero")
	}
	return nil
}

func buildLLMFileChunks(fileLines []string, maxLines int, maxBytes int) []llmFileChunk {
	if len(fileLines) == 0 {
		return nil
	}

	chunks := make([]llmFileChunk, 0, (len(fileLines)/maxLines)+1)
	current := llmFileChunk{StartLine: 1}
	currentBytes := 0

	for index, line := range fileLines {
		lineNumber := index + 1
		approxBytes := len(strconv.Itoa(lineNumber)) + len(line) + 3

		if len(current.Lines) > 0 && (len(current.Lines) >= maxLines || currentBytes+approxBytes > maxBytes) {
			chunks = append(chunks, current)
			current = llmFileChunk{StartLine: lineNumber}
			currentBytes = 0
		}

		current.Lines = append(current.Lines, line)
		currentBytes += approxBytes
	}

	if len(current.Lines) > 0 {
		chunks = append(chunks, current)
	}

	return chunks
}

func callLLMChunk(client llmHTTPClient, cfg *cfgReader.EarlybirdConfig, scanjob LLMJob, chunk llmFileChunk) ([]LLMFinding, error) {
	requestBody := llmChatCompletionRequest{
		Model:       cfg.LLMModel,
		Temperature: 0,
		Messages: []llmMessage{
			{
				Role: "system",
				Content: "You review source code for likely hard-coded credentials or secrets. Return JSON only with the top-level field findings. Each finding must contain line, credential_type, confidence, candidate, and reason. Confidence must be one of low, medium, or high. Ignore placeholders, obvious examples, comments describing documentation samples, and non-secret identifiers.",
			},
			{
				Role: "user",
				Content: buildLLMUserPrompt(scanjob, chunk),
			},
		},
		ResponseFormat: &llmResponseFormat{Type: "json_object"},
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal llm request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.LLMTimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.LLMEndpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build llm request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.LLMAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send llm request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read llm response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("llm request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var completion llmChatCompletionResponse
	if err := json.Unmarshal(body, &completion); err != nil {
		return nil, fmt.Errorf("decode llm response: %w", err)
	}
	if completion.Error != nil {
		return nil, fmt.Errorf("llm api error: %s", completion.Error.Message)
	}
	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("llm response did not include any choices")
	}

	var structured llmStructuredResponse
	content := completion.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &structured); err != nil {
		return nil, fmt.Errorf("decode llm structured response: %w", err)
	}

	return structured.Findings, nil
}

func buildLLMUserPrompt(scanjob LLMJob, chunk llmFileChunk) string {
	var builder strings.Builder
	builder.WriteString("File name: ")
	builder.WriteString(scanjob.FileName)
	builder.WriteString("\nFile path: ")
	builder.WriteString(scanjob.FilePath)
	builder.WriteString("\nChunk start line: ")
	builder.WriteString(strconv.Itoa(chunk.StartLine))
	builder.WriteString("\nCode:\n")

	for index, line := range chunk.Lines {
		builder.WriteString(strconv.Itoa(chunk.StartLine + index))
		builder.WriteString(": ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return builder.String()
}

func logLLMFindings(cfg *cfgReader.EarlybirdConfig, scanjob LLMJob, findings []LLMFinding) {
	if len(findings) == 0 {
		log.Printf("LLM scan found no additional credentials in %s", scanjob.FilePath)
		return
	}

	for _, finding := range findings {
		candidate := finding.Candidate
		if cfg.Suppress {
			candidate = maskValue(candidate)
		}
		log.Printf(
			"LLM scan finding file=%s line=%d type=%s confidence=%s candidate=%s reason=%s",
			scanjob.FilePath,
			finding.Line,
			finding.CredentialType,
			finding.Confidence,
			candidate,
			finding.Reason,
		)
	}
}
