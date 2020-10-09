/*
 * Copyright 2020 American Express
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

//ServerConfig is the timeout configuration for the Earlybird REST API server
type ServerConfig struct {
	WriteTimeout int `json:"write-timeout"`
	ReadTimeout  int `json:"read-timeout"`
	IdleTimeout  int `json:"idle-timeout"`
}

//Configs is the result of earlybird.json
type Configs struct {
	ModuleConfigs []struct {
		Name           string `json:"name"`
		DefaultEnabled bool   `json:"default_enabled"`
		ConfigURL      string `json:"config_url"`
	} `json:"modules"`
	LevelConfigs []struct {
		Name string `json:"level_name"`
		ID   int    `json:"level_id"`
	} `json:"finding_levels"`
	AnnotationsToSkip          []string `json:"text_ignore_patterns"`
	ExtensionsToSkipTextScan   []string `json:"filename_skip_text_scanning_extensions"`
	FailThreshold              int      `json:"fail_threshold_level"`
	DisplayThreshold           int      `json:"display_threshold_level"`
	DisplayConfidenceThreshold int      `json:"display_confidence_threshold_level"`
	ConfigFileURL              string   `json:"earlybird_config_url"`
	Version                    string   `json:"version"`
}

//EarlybirdConfig is the overall scan configs from config file and cli params
type EarlybirdConfig struct {
	SearchDir              string
	Gitrepo                string
	TargetType             string
	EnabledModules         []string
	OutputFormat           string
	OutputFile             string
	IgnoreFile             string
	SeverityFailLevel      int
	SeverityDisplayLevel   int
	ConfidenceFailLevel    int
	ConfidenceDisplayLevel int
	ConfigDir              string
	LevelMap               map[string]int
	Suppress               bool
	VerboseEnabled         bool
	GitStream              bool
	MaxFileSize            int64
	ShowFullLine           bool
	FailScan               bool
	RulesOnly              bool
	ExtensionsToSkipScan   []string
	AnnotationsToSkipLine  []string
	SkipComments           bool
	IgnoreFPRules          bool
	ShowSolutions          bool
	Version                string
	WorkerCount            int
	WorkLength             int
	HideMeta               bool
	Configs                Configs
}
