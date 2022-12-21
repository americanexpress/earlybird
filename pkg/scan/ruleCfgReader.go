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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
)

// Init loads in all the Earlybird rules into the CombinedRules global variable
func Init(cfg cfgreader.EarlybirdConfig) {
	if cfg.OutputFormat != "json" && !cfg.HideMeta {
		log.Println("Starting Go-EarlyBird version: ", cfg.Version)
		// Display options
		fmt.Println("Severity Fail threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.SeverityFailLevel))
		fmt.Println("Confidence Fail threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.ConfidenceFailLevel))
		fmt.Println("Severity Display threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.SeverityDisplayLevel))
		fmt.Println("Confidence Display threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.ConfidenceDisplayLevel))
		fmt.Println("Max file size to scan: ", cfg.MaxFileSize, " bytes")
	}

	// Init rule set for modules
	for moduleName, fileName := range cfg.EnabledModulesMap {
		log.Println("Loading module: ", moduleName)
		CombinedRules = append(CombinedRules, loadRuleConfigs(cfg, moduleName, fileName)...)
	}

	var err error
	//Load solutions for the rules
	if cfg.ShowSolutions {
		SolutionConfigs, err = loadSolutions(cfg.SolutionsConfigDir)

		if err != nil {
			log.Fatal("error loading solutions", err)
		}
	}

	// Init label configs
	Labels, err = loadLabelConfigs(cfg.LabelsConfigDir)

	if err != nil {
		log.Fatal("error loading labels file")
	}

	//Load false positive rules
	FalsePositiveRules, err = loadFalsePositives(cfg.FalsePositivesConfigDir)

	if err != nil {
		log.Fatal("error loading false positive rules", err)
	}

	// If we're only displaying the rules to be run, filter out anything that we wouldn't fail on and exit to skip the scan.
	// Only one output option since we exit before any interaction with Writers
	if cfg.RulesOnly {
		fmt.Println("\nShowing Rules Only (no scan to be executed)")
		fmt.Println()
		for _, combinedRule := range CombinedRules {
			if combinedRule.Severity <= cfg.SeverityFailLevel && combinedRule.Confidence <= cfg.ConfidenceFailLevel {
				fmt.Println("Code: ", combinedRule.Code)
				fmt.Println("Caption: ", combinedRule.Caption)
				fmt.Println("Pattern: ", combinedRule.Pattern)
				fmt.Println("Severity: ", cfgreader.Settings.TranslateLevelID(combinedRule.Severity))
				fmt.Println("Confidence: ", cfgreader.Settings.TranslateLevelID(combinedRule.Confidence))
				fmt.Println("Solution: ", combinedRule.Solution)
				fmt.Println("Example: ", combinedRule.Example)
				fmt.Println()
			}
		}
		os.Exit(0)
	}

	// Compile adjusted severity regex patterns
	if cfg.AdjustedSeverityCategories != nil {
		for i := range cfg.AdjustedSeverityCategories {
			if cfg.AdjustedSeverityCategories[i].Category == "" {
				log.Fatal("Missing required field category")
			}

			if cfg.AdjustedSeverityCategories[i].Patterns == nil {
				log.Fatal("Missing required field patterns")
			}

			if cfg.AdjustedSeverityCategories[i].AdjustedDisplaySeverity == "" {
				log.Fatal("Missing required field adjusted_display_severity")
			}

			for _, regEx := range cfg.AdjustedSeverityCategories[i].Patterns {
				compiled := regexp.MustCompile(regEx)

				cfg.AdjustedSeverityCategories[i].CompiledPatterns = append(cfg.AdjustedSeverityCategories[i].CompiledPatterns, compiled)
			}
		}
	}
}

// loadRuleConfigs loads the rules from the JSON config file, compiles the rules and defines the search area
func loadRuleConfigs(cfg cfgreader.EarlybirdConfig, moduleName, fileName string) []Rule {
	var rules, tmpRules Rules
	rulePath := filepath.Join(cfg.RulesConfigDir, fileName)

	err := cfgreader.LoadConfig(&tmpRules, rulePath)
	if err != nil {
		log.Println("Failed to load rules file", err)
	}

	for i := range tmpRules.Rules {
		if customRules, ok := cfg.ModuleConfigs.Modules[moduleName]; ok && tmpRules.Rules[i].Severity <= customRules.DisplaySeverityLevel && tmpRules.Rules[i].Confidence <= customRules.DisplayConfidenceLevel {
			tmpRules.Rules[i].Searcharea = tmpRules.Searcharea
			tmpRules.Rules[i].CompiledPattern = regexp.MustCompile(tmpRules.Rules[i].Pattern)
			rules.Rules = append(rules.Rules, tmpRules.Rules[i])
		} else if tmpRules.Rules[i].Severity <= cfg.SeverityDisplayLevel && tmpRules.Rules[i].Confidence <= cfg.ConfidenceDisplayLevel {
			tmpRules.Rules[i].Searcharea = tmpRules.Searcharea
			tmpRules.Rules[i].CompiledPattern = regexp.MustCompile(tmpRules.Rules[i].Pattern)
			rules.Rules = append(rules.Rules, tmpRules.Rules[i])
		}
	}
	return rules.Rules
}

// loadLabelConfigs loads the labels from the config file
func loadLabelConfigs(dirPath string) (LabelConfigRules map[int]LabelConfigs, err error) {
	LabelConfigRules = make(map[int]LabelConfigs)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		var tmpRules LabelConfigs
		err = cfgreader.LoadConfig(&tmpRules, path)
		if err != nil {
			log.Fatal("Failed to load labels file", err)
		}

		for i := range tmpRules.Labels {
			for _, code := range tmpRules.Labels[i].Codes {
				var AppendedLabelConfigRules LabelConfigs
				AppendedLabelConfigRules.Labels = append(LabelConfigRules[code].Labels, tmpRules.Labels[i])
				LabelConfigRules[code] = AppendedLabelConfigRules
			}
		}

		return nil
	})

	return LabelConfigRules, err
}

// loadFalsePositives loads in and compiles all the false positive rules for Earlybird
func loadFalsePositives(dirPath string) (FalsePositiveRules map[int]FalsePositives, err error) {
	FalsePositiveRules = make(map[int]FalsePositives)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		var tmpRules FalsePositives
		err = cfgreader.LoadConfig(&tmpRules, path)
		if err != nil {
			log.Fatal("Failed to load false positives file ", path, err)
		}

		for i := range tmpRules.FalsePositives {
			tmpRules.FalsePositives[i].CompiledPattern = regexp.MustCompile(tmpRules.FalsePositives[i].Pattern)
			for _, code := range tmpRules.FalsePositives[i].Codes {
				var AppendedFalsePositives FalsePositives
				AppendedFalsePositives.FalsePositives = append(FalsePositiveRules[code].FalsePositives, tmpRules.FalsePositives[i])
				FalsePositiveRules[code] = AppendedFalsePositives
			}
		}

		return nil
	})

	return FalsePositiveRules, err
}

// loadSolutions loads in solutions from the json config file
func loadSolutions(dirPath string) (solutionConfigs map[int]Solution, err error) {
	solutionConfigs = make(map[int]Solution)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		var tmp Solutions
		err = cfgreader.LoadConfig(&tmp, path)
		if err != nil {
			log.Fatal("Failed to load solutions file", err)
		}
		for i := range tmp.Solutions {
			solutionConfigs[tmp.Solutions[i].ID] = tmp.Solutions[i]
		}

		return nil
	})

	return solutionConfigs, err
}
