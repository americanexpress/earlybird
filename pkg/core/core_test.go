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

package core

import (
	"os"
	"path/filepath"
	"testing"

	cfgReader "github.com/americanexpress/earlybird/v4/pkg/config"
	"github.com/americanexpress/earlybird/v4/pkg/scan"
	"github.com/americanexpress/earlybird/v4/pkg/utils"
)

var eb EarlybirdCfg

func setup() {
	wd := utils.MustGetWD()
	eb.Config.LabelsConfigDir = filepath.Join(wd, "../../config/labels")
	eb.Config.SolutionsConfigDir = filepath.Join(wd, "../../config/solutions")
	eb.Config.FalsePositivesConfigDir = filepath.Join(wd, "../../config/falsepositives")
}

func cleanup() {
	eb = EarlybirdCfg{}
}

func TestMain(m *testing.M) {
	setup()
	scan.Init(eb.Config)

	os.Exit(m.Run())
}

// Program will exit with error if config init fails
func TestEarlybirdCfg_ConfigInit(t *testing.T) {
	eb.ConfigInit()
}

func TestEarlybirdCfg_Scan(t *testing.T) {
	eb.Config.SearchDir = utils.MustGetWD()
	eb.Scan()
}

func TestEarlybirdCfg_GitClone(t *testing.T) {
	if os.Getenv("local") == "" {
		t.Skip("If test cases not running locally, skip cloning external repositories for CI/CD purposes.")
	}

	var (
		FakeRepo = "https://github.com/carnal0wnage/fake_commited_secrets"
		RepoUser string
		Project  string
	)
	ptr := PTRGitConfig{
		Repo:     &FakeRepo,
		RepoUser: &RepoUser,
		Project:  &Project,
	}

	eb.GitClone(ptr)

	//Delete temporary cloned repository directory
	utils.DeleteGit(FakeRepo, eb.Config.SearchDir)
}

func TestEarlybirdCfg_getDefaultModuleSettings(t *testing.T) {
	modules := map[string]cfgReader.ModuleConfig{
		"inclusivity": {
			DisplaySeverity: "high",
		},
	}
	eb.Config.ModuleConfigs.Modules = modules

	eb.getDefaultModuleSettings()

	// we didn't explicitly configure  DisplayConfidence, make sure it got set to global default
	if got, want := eb.Config.ModuleConfigs.Modules["inclusivity"].DisplayConfidenceLevel, eb.Config.ConfidenceDisplayLevel; got != want {
		t.Fatalf("Unexpected default value set, got: %d, want: %d", got, want)
	}

	cleanup()
}

func TestEarlybirdCfg_BaselineWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	baselinePath := filepath.Join(tempDir, "earlybird_baseline.json")

	// Create test file with a secret
	secretFile := filepath.Join(tempDir, "credentials.txt")
	err := os.WriteFile(secretFile, []byte("AKIAIOSFODNN7EXAMPLE\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	testEb := EarlybirdCfg{}
	testEb.Config = eb.Config
	testEb.Config.SearchDir = tempDir
	testEb.Config.BaselineFile = baselinePath
	testEb.Config.OutputFormat = "json"
	testEb.Config.OutputFile = filepath.Join(tempDir, "scan_out.json")

	// 1. Generate baseline
	testEb.Scan()

	if _, err := os.Stat(baselinePath); os.IsNotExist(err) {
		t.Fatalf("Expected baseline file to be created at %s", baselinePath)
	}

	// 2. Scan with ignore-baseline: should exclude baseline findings
	ignoreEb := EarlybirdCfg{}
	ignoreEb.Config = eb.Config
	ignoreEb.Config.SearchDir = tempDir
	ignoreEb.Config.IgnoreBaselineFile = baselinePath
	ignoreEb.Config.OutputFormat = "json"
	ignoreEb.Config.OutputFile = filepath.Join(tempDir, "ignored_out.json")

	ignoreEb.Scan()

	if ignoreEb.Config.FailScan {
		t.Errorf("Expected FailScan to be false when all findings are in baseline")
	}

	// 3. Add a new secret: scan should detect new finding
	newSecretFile := filepath.Join(tempDir, "new_secret.txt")
	err = os.WriteFile(newSecretFile, []byte("ghp_111111111111111111111111111111111111\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	newScanEb := EarlybirdCfg{}
	newScanEb.Config = eb.Config
	newScanEb.Config.SearchDir = tempDir
	newScanEb.Config.IgnoreBaselineFile = baselinePath
	newScanEb.Config.OutputFormat = "json"
	newScanEb.Config.OutputFile = filepath.Join(tempDir, "new_out.json")

	newScanEb.Scan()
}
