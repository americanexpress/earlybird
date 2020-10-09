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

package scan

import (
	"fmt"
	"os"
	"regexp"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
)

//Init loads in all the Earlybird rules into the CombinedRules global variable
func Init(cfg cfgreader.EarlybirdConfig) {
	if cfg.OutputFormat != "json" && !cfg.HideMeta {
		fmt.Println("Go-EarlyBird version: ", cfg.Version)
		// Display options
		fmt.Println("Severity Fail threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.SeverityFailLevel))
		fmt.Println("Confidence Fail threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.ConfidenceFailLevel))
		fmt.Println("Severity Display threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.SeverityDisplayLevel))
		fmt.Println("Confidence Display threshold (at or above): ", cfgreader.Settings.TranslateLevelID(cfg.ConfidenceDisplayLevel))
		fmt.Println("Max file size to scan: ", cfg.MaxFileSize, " bytes")

		for _, moduleName := range cfg.EnabledModules {
			fmt.Println("Loading module: ", moduleName)
		}
	}

	// Init rule set for modules
	for _, moduleName := range cfg.EnabledModules {
		CombinedRules = append(CombinedRules, loadRuleConfigs(cfg.SeverityDisplayLevel, cfg.ConfidenceDisplayLevel, cfg.ConfigDir+moduleName+ruleSuffix)...)
	}

	//Load solutions for the rules
	if cfg.ShowSolutions {
		SolutionConfigs = loadSolutions(cfg.ConfigDir + "solutions.json")
	}

	// Init label configs
	Labels = loadLabelConfigs(cfg.ConfigDir + "labels.json")

	//Load false positive rules
	FalsePositiveRules = loadFalsePositives(cfg.ConfigDir + "false-positives.json")

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
}

//loadRuleConfigs loads the rules from the JSON config file, compiles the rules and defines the search area
func loadRuleConfigs(displaylevel int, displayconfidencelevel int, path string) []Rule {
	var rules, tmpRules Rules
	err := cfgreader.LoadConfig(&tmpRules, path)
	if err != nil {
		fmt.Println("Failed to load rules file", err)
	}

	for i := range tmpRules.Rules {
		if tmpRules.Rules[i].Severity <= displaylevel && tmpRules.Rules[i].Confidence <= displayconfidencelevel {
			tmpRules.Rules[i].Searcharea = tmpRules.Searcharea
			tmpRules.Rules[i].CompiledPattern = regexp.MustCompile(tmpRules.Rules[i].Pattern)
			rules.Rules = append(rules.Rules, tmpRules.Rules[i])
		}
	}
	return rules.Rules
}

//loadLabelConfigs loads the labels from the config file
func loadLabelConfigs(path string) (LabelConfigRules map[int]LabelConfigs) {
	LabelConfigRules = make(map[int]LabelConfigs)
	var tmpRules LabelConfigs
	err := cfgreader.LoadConfig(&tmpRules, path)
	if err != nil {
		fmt.Println("Failed to load labels file", err)
	}

	for i := range tmpRules.Labels {
		for _, code := range tmpRules.Labels[i].Codes {
			var AppendedLabelConfigRules LabelConfigs
			AppendedLabelConfigRules.Labels = append(LabelConfigRules[code].Labels, tmpRules.Labels[i])
			LabelConfigRules[code] = AppendedLabelConfigRules
		}
	}
	return LabelConfigRules
}

//loadFalsePositives loads in and compiles all the false positive rules for Earlybird
func loadFalsePositives(path string) (FalsePositiveRules map[int]FalsePositives) {
	FalsePositiveRules = make(map[int]FalsePositives)
	// init the array
	var tmpRules FalsePositives
	err := cfgreader.LoadConfig(&tmpRules, path)
	if err != nil {
		fmt.Println("Failed to load false positives file", err)
	}
	for i := range tmpRules.FalsePositives {
		tmpRules.FalsePositives[i].CompiledPattern = regexp.MustCompile(tmpRules.FalsePositives[i].Pattern)
		for _, code := range tmpRules.FalsePositives[i].Codes {
			var AppendedFalsePositives FalsePositives
			AppendedFalsePositives.FalsePositives = append(FalsePositiveRules[code].FalsePositives, tmpRules.FalsePositives[i])
			FalsePositiveRules[code] = AppendedFalsePositives
		}
	}
	return FalsePositiveRules
}

//loadSolutions loads in solutions from the json config file
func loadSolutions(path string) (solutionConfigs map[int]Solution) {
	solutionConfigs = make(map[int]Solution)
	var tmp Solutions
	err := cfgreader.LoadConfig(&tmp, path)
	if err != nil {
		fmt.Println("Failed to load solutions file", err)
	}
	for i := range tmp.Solutions {
		solutionConfigs[tmp.Solutions[i].ID] = tmp.Solutions[i]
	}
	return solutionConfigs
}
