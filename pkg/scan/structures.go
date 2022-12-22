/*
 * Copyright 2023 American Express
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
	"regexp"
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
