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

package cfgreader

import "regexp"

// ServerConfig is the timeout configuration for the Earlybird REST API server
type ServerConfig struct {
	WriteTimeout int `json:"write-timeout"`
	ReadTimeout  int `json:"read-timeout"`
	IdleTimeout  int `json:"idle-timeout"`
}

type AdjustedSeverityCategory struct {
	Category                string   `json:"category"`
	Patterns                []string `json:"patterns"`
	CompiledPatterns        []*regexp.Regexp
	AdjustedDisplaySeverity string `json:"adjusted_display_severity"`
	// If UseFilename is true the match will be based on hit.Filename
	UseFilename bool `json:"use_filename"`
	// If UseLineValue is true the match will be based on hit.LineValue
	UseLineValue bool `json:"use_line_value"`
	// If UseFileName and UseLine value are either both false or undefined the match defaults to hit.MatchValue
}

// Configs is the result of earlybird.json
type Configs struct {
	LevelConfigs []struct {
		Name string `json:"level_name"`
		ID   int    `json:"level_id"`
	} `json:"finding_levels"`
	AnnotationsToSkip          []string                   `json:"text_ignore_patterns"`
	ConfigBaseUrl              string                     `json:"config_base_url"`
	ExtensionsToSkipTextScan   []string                   `json:"filename_skip_text_scanning_extensions"`
	FailThreshold              int                        `json:"fail_threshold_level"`
	DisplayThreshold           int                        `json:"display_threshold_level"`
	DisplayConfidenceThreshold int                        `json:"display_confidence_threshold_level"`
	ConfigFileURL              string                     `json:"earlybird_config_url"`
	Version                    string                     `json:"version"`
	AdjustedSeverityCategories []AdjustedSeverityCategory `json:"adjusted_severity_categories_patterns"`
}

// Config from -module-config-file flag
type ModuleConfig struct {
	DisplaySeverity        string `json:"display_severity"`
	DisplayConfidence      string `json:"display_confidence"`
	DisplaySeverityLevel   int
	DisplayConfidenceLevel int
}

type ModuleConfigs struct {
	Modules map[string]ModuleConfig `json:"modules"`
}

// EarlybirdConfig is the overall scan configs from config file and cli params
type EarlybirdConfig struct {
	AvailableModules           []string
	RuleModulesFilenameMap     map[string]string
	SearchDir                  string
	Gitrepo                    string
	TargetType                 string
	EnabledModulesMap          map[string]string
	EnabledModules             []string
	WithConsole                bool
	OutputFormat               string
	OutputFile                 string
	IgnoreFile                 string
	IgnoreFailure              bool
	SeverityFailLevel          int
	SeverityDisplayLevel       int
	ConfidenceFailLevel        int
	ConfidenceDisplayLevel     int
	ConfigDir                  string
	RulesConfigDir             string
	FalsePositivesConfigDir    string
	SolutionsConfigDir         string
	LabelsConfigDir            string
	LevelMap                   map[string]int
	Suppress                   bool
	VerboseEnabled             bool
	GitStream                  bool
	MaxFileSize                int64
	ShowFullLine               bool
	FailScan                   bool
	RulesOnly                  bool
	ExtensionsToSkipScan       []string
	AnnotationsToSkipLine      []string
	SkipComments               bool
	IgnoreFPRules              bool
	ShowSolutions              bool
	Version                    string
	WorkerCount                int
	WorkLength                 int
	HideMeta                   bool
	StrictJKS                  bool
	ModuleConfigs              ModuleConfigs
	AdjustedSeverityCategories []AdjustedSeverityCategory
}
