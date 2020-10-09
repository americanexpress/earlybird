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
	"testing"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/utils"
)

var (
	config     cfgreader.Configs
	configPath = utils.GetConfigDir() + "earlybird.json"
)

func init() {
	if err := cfgreader.LoadConfig(&config, configPath); err != nil {
		fmt.Println("LoadConfig() = error:", err, " want error nil")
	}
}

func Test_loadRuleConfigs(t *testing.T) {
	fmt.Println(config.ModuleConfigs)
	for _, module := range config.ModuleConfigs {
		rules := loadRuleConfigs(4, 4, utils.GetConfigDir()+module.Name+ruleSuffix)
		if len(rules) == 0 {
			t.Errorf("loadRuleConfigs() = %v, failed to load rules", rules)
		}
		for _, rule := range rules {
			matchValue := rule.CompiledPattern.FindStringSubmatch(rule.Example)
			if len(matchValue) == 0 {
				t.Errorf("Failed to match pattern to example. Module: %s Rule Code: %d", module.Name, rule.Code)
			}
		}
	}
}

func Test_loadLabelConfigs(t *testing.T) {
	if gotLabelConfigRules := loadLabelConfigs(utils.GetConfigDir() + "labels.json"); len(gotLabelConfigRules) == 0 {
		t.Errorf("loadLabelConfigs() = %v, Failed to load labels", gotLabelConfigRules)
	}
}

func Test_loadFalsePositives(t *testing.T) {
	if gotFalsePositiveRules := loadFalsePositives(utils.GetConfigDir() + "false-positives.json"); len(gotFalsePositiveRules) == 0 {
		t.Errorf("loadFalsePositives() = %v, Failed to load any false positive rules", gotFalsePositiveRules)
	}
}

func Test_loadSolutions(t *testing.T) {
	if gotSolutionConfigs := loadSolutions(utils.GetConfigDir() + "solutions.json"); len(gotSolutionConfigs) == 0 {
		t.Errorf("loadSolutions() = %v, Failed to load any solutions", gotSolutionConfigs)
	}
}
