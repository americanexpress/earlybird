package scan

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
)

var defaultLLMHTTPClient llmHTTPClient = &http.Client{}

type LLMProvider interface {
	setHeader(req *http.Request, apiKey string)
	buildRequestBody(cfg *cfgReader.EarlybirdConfig, scanJob LLMJob, chunk llmFileChunk) ([]byte, error)
	parseResponse(responseBody []byte) ([]LLMFinding, error)
}

func getLLMProvider(cfg *cfgReader.EarlybirdConfig) LLMProvider {
	switch {
	case isGeminiModel(cfg.LLMModel):
		return GeminiProvider{}
	default:
		return GptProvider{}
	}
}

func llmSystemPrompt(cfg *cfgReader.EarlybirdConfig) string {
	if cfg != nil && cfg.LLMSystemPrompt != "" {
		return cfg.LLMSystemPrompt
	}
	return cfgReader.DefaultLLMSystemPrompt
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

// validateLLMConfig checks that the necessary configuration for LLM scan is present and valid.
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
		return fmt.Errorf("llm api key is required: set EARLYBIRD_LLM_API_KEY")
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

func callLLMChunk(client llmHTTPClient, cfg *cfgReader.EarlybirdConfig, scanJob LLMJob, chunk llmFileChunk) ([]LLMFinding, error) {
	provider := getLLMProvider(cfg)
	payload, err := provider.buildRequestBody(cfg, scanJob, chunk)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.LLMTimeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.LLMEndpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build llm request: %w", err)
	}
	provider.setHeader(req, cfg.LLMAPIKey)
	var responseBody []byte
	responseBody, err = sendLLMRequest(client, req)
	if err != nil {
		return nil, err
	}
	return provider.parseResponse(responseBody)
}

func sendLLMRequest(client llmHTTPClient, req *http.Request) ([]byte, error) {
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

	return body, nil
}

func buildLLMUserPrompt(scanJob LLMJob, chunk llmFileChunk) string {
	var builder strings.Builder
	builder.WriteString("File name: ")
	builder.WriteString(scanJob.FileName)
	builder.WriteString("\nFile path: ")
	builder.WriteString(scanJob.FilePath)
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

func logLLMFindings(cfg *cfgReader.EarlybirdConfig, scanJob LLMJob, findings []LLMFinding) {
	if len(findings) == 0 {
		fmt.Printf("LLM scan found no additional credentials in %s", scanJob.FilePath)
		return
	}
	fmt.Println("LLM scan finding.")

	for i, finding := range findings {
		fmt.Println("Finding #:", i+1)
		fmt.Println("\tFilename=", scanJob.FilePath)
		fmt.Println("\tLine=", finding.Line)
		fmt.Println("\tType=", finding.CredentialType)
		fmt.Println("\tConfidence=", finding.Confidence)
		fmt.Println("\tDescription=", finding.Reason)
	}
}
