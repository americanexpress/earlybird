/*
 * Copyright 2021 American Express
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package scan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

// Rules is the exported definition of the Rules structure for Earlybird
type Rules struct {
	Rules      []Rule `json:"rules"`
	Searcharea string `json:"Searcharea"`
}

// Rule Each module config is a set of rules
type Rule struct {
	Code, Severity, Confidence, SolutionID            int
	Pattern, Caption, Category, Solution, Postprocess string
	CompiledPattern                                   *regexp.Regexp
	Searcharea                                        string
	CWE                                               []string
	Example                                           string
}

// Hit is a match in a file against a specific rule
type Hit struct {
	Code         int      `json:"code"`
	Filename     string   `json:"filename"`
	Caption      string   `json:"caption"`
	Category     string   `json:"category"`
	MatchValue   string   `json:"match_value"`
	LineValue    string   `json:"line_value"`
	Solution     string   `json:"solution"`
	Line         int      `json:"line"`
	Severity     string   `json:"severity"`
	SeverityID   int      `json:"severity_id"`
	Confidence   string   `json:"confidence"`
	ConfidenceID int      `json:"confidence_id"`
	Labels       []string `json:"labels"`
	CWE          []string `json:"cwe"`
	Time         string   `json:"time"`
}

// File to scan
type File struct {
	Name  string
	Path  string
	Lines []Line
	Raw   []byte
}

// Line in a file to scan
type Line struct {
	LineNum                       int
	LineValue, FilePath, FileName string
}

// Report is the Earlybird end output
type Report struct {
	Version       string   `json:"version"`
	Skipped       []string `json:"skipped"`
	Ignore        []string `json:"ignore"`
	Threshold     int      `json:"threshold"`
	Modules       []string `json:"modules"`
	Hits          []Hit    `json:"hits"`
	HitCount      int      `json:"hit_count"`
	FilesScanned  int      `json:"files_scanned"`
	RulesObserved int      `json:"rules_observed"`
	StartTime     string   `json:"start_time"`
	EndTime       string   `json:"end_time"`
	Duration      string   `json:"duration"`
}

// WorkJob As we add jobs to the pool, they need to contain the line being scanned and the file content (in Lines)
type WorkJob struct {
	WorkLine  Line
	FileLines []Line
}

type LLMJob struct {
	FileLines          []string
	FilePath, FileName string
}

type llmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GeminiContent struct {
	Role  string  `json:"role"`
	Parts []Parts `json:"parts"`
}

type Parts struct {
	Text string `json:"text"`
}

type llmResponseFormat struct {
	Type string `json:"type"`
}

type llmChatCompletionRequest struct {
	Model          string             `json:"model"`
	Messages       []llmMessage       `json:"messages"`
	Temperature    int                `json:"temperature"`
	ResponseFormat *llmResponseFormat `json:"response_format,omitempty"`
}

type gptLlmResponsesTextPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type gptLlmResponsesInput struct {
	Role    string                    `json:"role"`
	Content []gptLlmResponsesTextPart `json:"content"`
}

type gptLlmResponsesTextFormat struct {
	Type string `json:"type"`
}

type gptLlmResponsesText struct {
	Format gptLlmResponsesTextFormat `json:"format"`
}

type GptLlmResponsesRequest struct {
	Model string                 `json:"model"`
	Input []gptLlmResponsesInput `json:"input"`
	Text  gptLlmResponsesText    `json:"text"`
}

type GeminiRequestBody struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiResponseBody struct {
	Candidates    []llmGeminiCandidate `json:"candidates,omitempty"`
	ModelVersion  string               `json:"model_version,omitempty"`
	ResponseId    string               `json:"response_id,omitempty"`
	UsageMetadata struct {
		PromptTokenCount     int    `json:"prompt_tokens"`
		CandidatesTokenCount int    `json:"candidates_token_count"`
		TotalTokenCount      int    `json:"total_tokens"`
		ThoughtsTokenCount   int    `json:"thoughts_token_count,omitempty"`
		ServiceTier          string `json:"service_tier,omitempty"`
		PromptTokenDetails   struct {
			Modality   string `json:"modality"`
			TokenCount int    `json:"token_count"`
		} `json:"completion_tokens"`
	} `json:"usage_metadata,omitempty"`
}

type llmChatCompletionChoice struct {
	Message llmMessage `json:"message"`
}

type llmAPIError struct {
	Message string `json:"message"`
}

type gptLlmChatCompletionResponse struct {
	Choices []llmChatCompletionChoice `json:"choices"`
	Error   *llmAPIError              `json:"error,omitempty"`
}

type llmResponsesOutput struct {
	Content []gptLlmResponsesTextPart `json:"content,omitempty"`
}

type gptLlmResponsesResponse struct {
	Output []llmResponsesOutput `json:"output,omitempty"`
	Error  *llmAPIError         `json:"error,omitempty"`
}

type llmGeminiCandidate struct {
	Content json.RawMessage `json:"content"`
}

type llmAIResponse struct {
	Choices    []llmChatCompletionChoice `json:"choices,omitempty"`
	Candidates []llmGeminiCandidate      `json:"candidates,omitempty"`
	Error      *llmAPIError              `json:"error,omitempty"`
}

type LLMConfidence string

func (c *LLMConfidence) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*c = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*c = LLMConfidence(s)
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*c = LLMConfidence(strconv.FormatInt(i, 10))
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		if float64(int64(f)) == f {
			*c = LLMConfidence(strconv.FormatInt(int64(f), 10))
		} else {
			*c = LLMConfidence(strconv.FormatFloat(f, 'f', -1, 64))
		}
		return nil
	}

	return fmt.Errorf("invalid confidence value: %s", string(data))
}

type LLMFinding struct {
	Line           int           `json:"line"`
	CredentialType string        `json:"credential_type"`
	Confidence     LLMConfidence `json:"confidence"`
	Candidate      string        `json:"candidate,omitempty"`
	Reason         string        `json:"reason"`
}

type llmStructuredResponse struct {
	Findings []LLMFinding `json:"findings"`
}

type llmFileChunk struct {
	StartLine int
	Lines     []string
}

type llmHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// FalsePositives are the rules to match false positives post process
type FalsePositives struct {
	FalsePositives []FalsePositive `json:"rules"`
}

// FalsePositive is a rule to match false positives post process
type FalsePositive struct {
	Codes           []int
	Pattern         string
	CompiledPattern *regexp.Regexp
	FileExtensions  []string
	UseFullLine     bool
}

// Solutions to each rule / finding
type Solutions struct {
	Solutions []Solution `json:"solutions"`
}

// Solution display text for a solution
type Solution struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// LabelConfig Rule for applying labels to hits based on context
type LabelConfig struct {
	Label     string   `json:"label"`
	Keys      []string `json:"keys"`
	Multiline bool     `json:"multiline"`
	Category  string   `json:"category"`
	Codes     []int    `json:"codes"`
}

// LabelConfigs Rules for applying labels to hits based on context
type LabelConfigs struct {
	Labels []LabelConfig `json:"Labels"`
}
