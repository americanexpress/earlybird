package scan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
)

func TestBuildLLMFileChunks(t *testing.T) {
	chunks := buildLLMFileChunks([]string{"one", "two", "three", "four", "five"}, 2, 1024)

	if len(chunks) != 3 {
		t.Fatalf("buildLLMFileChunks() chunk count = %d, want 3", len(chunks))
	}
	if chunks[0].StartLine != 1 || len(chunks[0].Lines) != 2 {
		t.Fatalf("unexpected first chunk: %+v", chunks[0])
	}
	if chunks[1].StartLine != 3 || len(chunks[1].Lines) != 2 {
		t.Fatalf("unexpected second chunk: %+v", chunks[1])
	}
	if chunks[2].StartLine != 5 || len(chunks[2].Lines) != 1 {
		t.Fatalf("unexpected third chunk: %+v", chunks[2])
	}
}

func TestLLMScanCallsEndpointAndParsesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("Authorization header = %q, want Bearer test-key", got)
		}

		var req llmChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		if req.Model != "gpt-test" {
			t.Fatalf("request model = %q, want gpt-test", req.Model)
		}
		if len(req.Messages) != 2 {
			t.Fatalf("message count = %d, want 2", len(req.Messages))
		}
		if !strings.Contains(req.Messages[1].Content, "2: api_key = secret-value") {
			t.Fatalf("prompt did not include numbered file lines: %q", req.Messages[1].Content)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(llmChatCompletionResponse{
			Choices: []llmChatCompletionChoice{{
				Message: llmMessage{Content: `{"findings":[{"line":2,"credential_type":"api_key","confidence":"high","candidate":"secret-value","reason":"looks like a hard-coded credential"}]}`},
			}},
		})
	}))
	defer server.Close()

	cfg := &cfgReader.EarlybirdConfig{
		EnableLLMScan:     true,
		LLMEndpoint:       server.URL,
		LLMAPIKey:         "test-key",
		LLMModel:          "gpt-test",
		LLMTimeoutSeconds: 5,
		LLMMaxLines:       10,
		LLMMaxBytes:       2048,
	}

	findings, err := llm_scan(cfg, LLMJob{
		FileName:  "sample.env",
		FilePath:  "/tmp/sample.env",
		FileLines: []string{"username = app", "api_key = secret-value"},
	})
	if err != nil {
		t.Fatalf("llm_scan() error = %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("llm_scan() findings = %d, want 1", len(findings))
	}
	if findings[0].Line != 2 {
		t.Fatalf("finding line = %d, want 2", findings[0].Line)
	}
	if findings[0].CredentialType != "api_key" {
		t.Fatalf("finding credential type = %q, want api_key", findings[0].CredentialType)
	}
}

func TestValidateLLMConfigRequiresAPIKey(t *testing.T) {
	cfg := &cfgReader.EarlybirdConfig{
		EnableLLMScan:     true,
		LLMEndpoint:       "https://api.openai.com/v1/chat/completions",
		LLMModel:          "gpt-test",
		LLMTimeoutSeconds: 5,
		LLMMaxLines:       10,
		LLMMaxBytes:       1024,
	}

	err := validateLLMConfig(cfg)
	if err == nil {
		t.Fatal("validateLLMConfig() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "api key") {
		t.Fatalf("validateLLMConfig() error = %v, want api key message", err)
	}
}