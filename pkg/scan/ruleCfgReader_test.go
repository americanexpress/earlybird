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
	cfgreader "github.com/americanexpress/earlybird/v4/pkg/config"
	"github.com/americanexpress/earlybird/v4/pkg/utils"
	"path"
	"testing"
)

var (
	config cfgreader.EarlybirdConfig
)

func init() {
	config.ConfigDir = path.Join(utils.MustGetWD(), utils.GetConfigDir())
	config.RulesConfigDir = path.Join(config.ConfigDir, "rules")
	// low severity
	config.SeverityDisplayLevel = 4
	// low confidence
	config.ConfidenceDisplayLevel = 4
}

func Test_loadRuleConfigs(t *testing.T) {
	rules := loadRuleConfigs(config, "content", "content.yaml")

	if len(rules) == 0 {
		t.Errorf("loadRuleConfigs() = %v, failed to load rules", rules)
	}
}

func Test_loadLabelConfigs(t *testing.T) {
	gotLabelConfigRules, err := loadLabelConfigs(path.Join(config.ConfigDir, "labels"))

	if err != nil {
		t.Fatal(err)
	}

	if len(gotLabelConfigRules) == 0 {
		t.Errorf("loadLabelConfigs() = %v, Failed to load labels", gotLabelConfigRules)
	}
}

func Test_loadFalsePositives(t *testing.T) {
	gotFalsePositiveRules, err := loadFalsePositives(path.Join(config.ConfigDir, "falsepositives"))

	if err != nil {
		t.Fatal(err)
	}

	if len(gotFalsePositiveRules) == 0 {
		t.Errorf("loadFalsePositives() = %v, Failed to load any false positive rules", gotFalsePositiveRules)
	}
}

func Test_loadSolutions(t *testing.T) {
	gotSolutionConfigs, err := loadSolutions(path.Join(config.ConfigDir, "solutions"))

	if err != nil {
		t.Fatal(err)
	}

	if len(gotSolutionConfigs) == 0 {
		t.Errorf("loadSolutions() = %v, Failed to load any solutions", gotSolutionConfigs)
	}
}
