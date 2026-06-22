# LLM-Based Secret Scanning

## Overview

EarlyBird integrates Large Language Models (LLMs) to provide advanced, AI-driven secret detection alongside traditional pattern-matching rules. This module enables scanning of source code for hardcoded credentials, API keys, passwords, and other sensitive data by leveraging the contextual understanding capabilities of LLMs like OpenAI's GPT and Google's Gemini.

## Architecture

### How LLM Scanning Works

The LLM scanning system operates through the following flow:

1. **Initialization**: Configuration validation ensures all required LLM parameters are present
2. **Chunking**: Large files are split into manageable chunks based on line count and byte limits
3. **Provider Selection**: The appropriate LLM provider (GPT, Gemini, etc.) is selected based on the model name
4. **Request Building**: Each chunk is formatted as a request payload specific to the LLM provider's API
5. **API Call**: The request is sent to the configured LLM endpoint with authentication
6. **Response Parsing**: The LLM response is parsed to extract structured findings
7. **Result Aggregation**: Findings from all chunks are collected and logged

### System Prompt

EarlyBird sends a system prompt to guide the LLM's analysis. The default system prompt instructs the LLM to:

```
You analyze source code for likely hard-coded credentials or secrets. Return JSON only with the 
top-level field findings. Ex {line:3, credential_type:password, confidence: 95% , criticality: 4, 
reason: reason_for_flag, Candidate: <flagged_value>}.
```

You can customize this prompt via configuration.

### Provider Interface

All LLM providers implement the `LLMProvider` interface:

```go
type LLMProvider interface {
    setHeader(req *http.Request, apiKey string)
    buildRequestBody(cfg *cfgReader.EarlybirdConfig, scanJob LLMJob, chunk llmFileChunk) ([]byte, error)
    parseResponse(responseBody []byte) ([]LLMFinding, error)
}
```

**Methods:**

- **`setHeader()`**: Sets HTTP headers required by the LLM provider's API (e.g., authorization headers)
- **`buildRequestBody()`**: Formats the code chunk into the provider's required request payload format
- **`parseResponse()`**: Parses the provider's response and extracts findings as structured JSON

## Configuration

LLM scanning is controlled via the following configuration parameters. Fields marked as `No` are optional because EarlyBird applies a default value when they are not provided.

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `EnableLLMScan` | bool | Yes | `false` | Enables/disables LLM-based scanning |
| `LLMEndpoint` | string | Yes | None | The API endpoint URL for the LLM service |
| `LLMAPIKey` | string | Yes | None | API key for authentication (via `EARLYBIRD_LLM_API_KEY` env var) |
| `LLMModel` | string | Yes | None | Model identifier (e.g., `gpt-4`, `gpt-4o`, `gemini-2.0-flash`) |
| `LLMSystemPrompt` | string | No | `cfgReader.DefaultLLMSystemPrompt` | Custom system prompt; if omitted, EarlyBird uses the built-in default prompt |
| `LLMTimeoutSeconds` | int | No | `30` | Request timeout in seconds |
| `LLMMaxLines` | int | No | `200` | Maximum lines per chunk (e.g., 50-100) |
| `LLMMaxBytes` | int | No | `16000` | Maximum bytes per chunk (e.g., 8000-16000) |
| `LLMFailClosed` | bool | No | `false` | If true, scanning fails if LLM scan fails; if false, scanning continues |

The default values above match the CLI defaults used by EarlyBird when these options are not passed.

### Configuration File Example

```json
{
  "EnableLLMScan": true,
  "LLMEndpoint": "https://api.openai.com/v1/chat/completions",
    "LLMAPIKey": "${EARLYBIRD_LLM_API_KEY}",
  "LLMModel": "gpt-4o",
  "LLMSystemPrompt": "You analyze source code for likely hard-coded credentials or secrets...",
  "LLMTimeoutSeconds": 30,
  "LLMMaxLines": 100,
  "LLMMaxBytes": 16000,
  "LLMFailClosed": false
}
```

### Environment Variables

- `EARLYBIRD_LLM_API_KEY`: LLM API key

## Supported Models

### OpenAI (GPT)

**Provider Detection**: Model name contains `"gpt"` (case-insensitive)

**Supported Models**:
- `gpt-4`
- `gpt-4-turbo`
- `gpt-4o`
- `gpt-4o-mini`

**Configuration**:
```json
{
  "LLMEndpoint": "https://api.openai.com/v1/chat/completions",
  "LLMModel": "gpt-4o",
    "LLMAPIKey": "${EARLYBIRD_LLM_API_KEY}"
}
```

**API Features**:
- Supports both standard Chat Completions API and newer Responses API
- Requests must include authorization header: `Authorization: Bearer <api-key>`
- Returns JSON response with `choices[0].message.content`

### Google Gemini

**Provider Detection**: Model name contains `"gemini"` (case-insensitive)

**Supported Models**:
- `gemini-2.0-flash`
- `gemini-1.5-pro`
- `gemini-1.5-flash`

**Configuration**:
```json
{
  "LLMEndpoint": "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent",
  "LLMModel": "gemini-2.0-flash",
  "LLMAPIKey": "${GEMINI_API_KEY}"
}
```

**API Features**:
- Uses header `X-goog-api-key` for authentication instead of Authorization
- Request body uses `contents` array format
- Responses are wrapped in `candidates` array

## Adding a New LLM Provider

### Step 1: Create Provider Implementation File

Create a new file `pkg/scan/<provider_name>.go` that implements the `LLMProvider` interface.

**File Structure**:
```go
package scan

import (
    "net/http"
    cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
)

// Define detection function
func is<ProviderName>Model(model string) bool {
    return strings.Contains(strings.ToLower(model), "<provider_keyword>")
}

// Define provider struct
type <ProviderName>Provider struct{}

// Implement setHeader()
func (p <ProviderName>Provider) setHeader(req *http.Request, apiKey string) {
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)
    // Add provider-specific headers as needed
}

// Implement buildRequestBody()
func (p <ProviderName>Provider) buildRequestBody(
    cfg *cfgReader.EarlybirdConfig,
    scanJob LLMJob,
    chunk llmFileChunk,
) ([]byte, error) {
    // Transform the code chunk into the provider's request format
    requestBody := YourRequestStruct{
        Model: cfg.LLMModel,
        Messages: []YourMessage{
            {Role: "system", Content: llmSystemPrompt(cfg)},
            {Role: "user", Content: buildLLMUserPrompt(scanJob, chunk)},
        },
    }
    return json.Marshal(requestBody)
}

// Implement parseResponse()
func (p <ProviderName>Provider) parseResponse(responseBody []byte) ([]LLMFinding, error) {
    var response YourResponseStruct
    if err := json.Unmarshal(responseBody, &response); err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }
    
    // Extract findings from provider's response format
    var structured llmStructuredResponse
    if err := json.Unmarshal([]byte(response.Content), &structured); err != nil {
        return nil, fmt.Errorf("parse findings: %w", err)
    }
    
    return structured.Findings, nil
}
```

### Step 2: Register Provider in Provider Selector

Update `pkg/scan/llm_Scan.go` in the `getLLMProvider()` function to include detection for your new provider:

```go
func getLLMProvider(cfg *cfgReader.EarlybirdConfig) LLMProvider {
    switch {
    case isGeminiModel(cfg.LLMModel):
        return GeminiProvider{}
    case is<ProviderName>Model(cfg.LLMModel):
        return <ProviderName>Provider{}
    default:
        return GptProvider{}
    }
}
```

### Step 3: Define Request/Response Structures

In `pkg/scan/structures.go`, add the necessary structs to map your provider's API:

```go
type <ProviderName>RequestBody struct {
    Model    string                    `json:"model"`
    Messages []<ProviderName>Message   `json:"messages"`
    // Add other fields as required by the provider
}

type <ProviderName>Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type <ProviderName>ResponseBody struct {
    // Map the provider's response structure
    // Ensure it includes the LLM findings as JSON
}
```

### Step 4: Write Tests

Create `pkg/scan/<provider_name>_test.go` to test your implementation:

```go
package scan

import (
    "net/http"
    "net/http/httptest"
    "testing"
    cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
)

func TestLLMScan<ProviderName>BuildsRequestAndParsesResponse(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify headers
        if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
            t.Fatalf("Authorization header = %q, want Bearer test-key", got)
        }
        
        // Verify request body format
        var req <ProviderName>RequestBody
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            t.Fatalf("Decode() error = %v", err)
        }
        
        // Send mock response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "result": `{"findings":[...]}`,
        })
    }))
    defer server.Close()
    
    cfg := &cfgReader.EarlybirdConfig{
        EnableLLMScan:     true,
        LLMEndpoint:       server.URL,
        LLMAPIKey:         "test-key",
        LLMModel:          "<provider-model>",
        LLMTimeoutSeconds: 5,
        LLMMaxLines:       100,
        LLMMaxBytes:       16000,
    }
    
    findings, err := llm_scan(cfg, LLMJob{
        FileName:  "test.py",
        FilePath:  "/tmp/test.py",
        FileLines: []string{"password = 'secret'"},
    })
    
    if err != nil {
        t.Fatalf("llm_scan() error = %v", err)
    }
    // Assert findings
}
```

## Response Format

All LLM providers must return findings in the following JSON structure:

```json
{
  "findings": [
    {
      "line": 42,
      "credential_type": "api_key",
      "confidence": "high",
      "candidate": "sk-1234567890...",
      "reason": "OpenAI API key pattern detected"
    }
  ]
}
```

**Field Descriptions**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `line` | int | Yes | Line number in the file where the credential was found |
| `credential_type` | string | Yes | Type of credential (e.g., `api_key`, `password`, `token`) |
| `confidence` | string | Yes | Confidence level: `low`, `medium`, `high` |
| `candidate` | string | No | The actual credential value or pattern found |
| `reason` | string | Yes | Human-readable explanation of why this was flagged |

## Performance Considerations

### Chunking Strategy

Files are split into chunks to manage API payload limits and costs:

- **`LLMMaxLines`**: Controls how many lines are included per API call (default: 50-100)
- **`LLMMaxBytes`**: Prevents chunks from exceeding byte limits (default: 8000-16000)

**Chunking Algorithm**:
1. Iterate through file lines sequentially
2. Accumulate lines until either max line count or byte limit is reached
3. Create a new chunk and continue
4. Each chunk has a `StartLine` to track line numbers in the file

### Timeout Configuration

Set `LLMTimeoutSeconds` appropriate for your LLM provider:
- OpenAI: 30-60 seconds
- Google Gemini: 30-60 seconds
- Custom providers: depends on latency and model complexity

### Cost Optimization

- Use smaller models for non-critical scans
- Increase `LLMMaxLines` and `LLMMaxBytes` to reduce API calls
- Configure `LLMFailClosed: false` to prevent scanning failures from blocking CI/CD

## Error Handling

The LLM scan validates configuration before execution:

```
Validation Checks:
- Config must not be nil
- EnableLLMScan must be true
- LLMEndpoint must be set
- LLMAPIKey must be set (via env var or config)
- LLMModel must be set
- LLMTimeoutSeconds must be > 0
- LLMMaxLines must be > 0
- LLMMaxBytes must be > 0
```

If validation fails, scanning returns an error. If `LLMFailClosed` is `false`, the scan continues without LLM analysis.

## Integration with Rule-Based Scanning

LLM findings are independent of rule-based findings. Both are reported separately in the scan results:

- **Rule-based findings**: Pattern matches from YAML configuration files
- **LLM findings**: AI-detected findings from the LLM provider

LLM findings do NOT apply false positive filters, labels, or post-processing rules defined for rule-based scanning.

## Troubleshooting

### Common Issues

**Issue**: LLM scan returns no findings
- Check `EnableLLMScan` is `true`
- Verify API key is set correctly via environment variable
- Confirm endpoint URL is correct for your provider
- Check network connectivity to the LLM endpoint

**Issue**: "llm scan is disabled"
- Set `EnableLLMScan: true` in configuration

**Issue**: "llm api key is required"
- Set `EARLYBIRD_LLM_API_KEY` environment variable

**Issue**: Timeout errors
- Increase `LLMTimeoutSeconds`
- Reduce `LLMMaxLines` to send smaller chunks
- Verify network connectivity

**Issue**: Response parsing errors
- Confirm the LLM is returning JSON with the expected structure
- Update system prompt to ensure JSON format compliance
- Check provider-specific response format in implementation

## Example: Running LLM Scan

```bash
export EARLYBIRD_LLM_API_KEY=sk-your-key-here

go-earlybird --path=/project/src \
  --config=config/earlybird.json
```

With configuration:
```json
{
  "EnableLLMScan": true,
  "LLMEndpoint": "https://api.openai.com/v1/chat/completions",
  "LLMModel": "gpt-4o",
  "LLMTimeoutSeconds": 30,
  "LLMMaxLines": 100,
  "LLMMaxBytes": 16000
}
```

## References

- [OpenAI Chat Completions API](https://platform.openai.com/docs/api-reference/chat/create)
- [Google Gemini API](https://ai.google.dev/api/rest)
- [CWE-798: Use of Hardcoded Credentials](https://cwe.mitre.org/data/definitions/798.html)
```