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
	"flag"
	"os"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/utils"
)

const (
	rulesDir          = "rules"
	falsePositivesDir = "falsepositives"
	labelsDir         = "labels"
	solutionsDir      = "solutions"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

//Define our static CLI flags
var (
	userHomeDir, _                = os.UserHomeDir()
	levelOptions                  = utils.GetDisplayList(cfgreader.Settings.GetLevelNames())
	ptrStreamInput                = flag.Bool("stream", false, "Use stream IO as input instead of file(s)")
	enableFlags                   arrayFlags
	ptrUpdateFlag                 = flag.Bool("update", false, "Update module configurations")
	ptrGitStreamInput             = flag.Bool("git-commit-stream", false, "Use stream IO of Git commit log as input instead of file(s) -- e.g., 'cat secrets.text > go-earlybird'")
	ptrVerbose                    = flag.Bool("verbose", false, "Reports details about file reads")
	ptrSuppressSecret             = flag.Bool("suppress", false, "Suppress reporting of the secret found (important if output is going to Slack or other logs)")
	ptrWorkerCount                = flag.Int("workers", 100, "Set number of workers.")
	ptrWorkLength                 = flag.Int("worksize", 2500, "Set Line Wrap Length.")
	ptrMaxFileSize                = flag.Int64("max-file-size", 10240000, "Maximum file size to scan (in bytes)")
	ptrShowFullLine               = flag.Bool("show-full-line", false, "Display the full line where the pattern match was found (warning: this can be dangerous with minified script files)")
	ptrConfigDir                  = flag.String("config", utils.GetConfigDir(), "Directory where configuration files are stored")
	ptrRulesOnly                  = flag.Bool("show-rules-only", false, "Display rules that would be run, but do not execute a scan")
	ptrSkipComments               = flag.Bool("skip-comments", false, "Skip scanning comments in files -- applies only to the 'content' module")
	ptrIgnoreFPRules              = flag.Bool("ignore-fp-rules", false, "Ignore the false positive post-process rules")
	ptrShowSolutions              = flag.Bool("show-solutions", false, "Display recommended solution for each finding")
	ptrGitStagedFlag              = flag.Bool("git-staged", false, "Scan only git staged files")
	ptrGitTrackedFlag             = flag.Bool("git-tracked", false, "Scan only git tracked files")
	ptrPath                       = flag.String("path", utils.MustGetWD(), "Directory to scan (defaults to CWD) -- ABSOLUTE PATH ONLY")
	ptrOutputFormat               = flag.String("format", "console", "Output format [ console | json | csv ]")
	ptrOutputFile                 = flag.String("file", "", "Output file -- e.g., 'go-earlybird --file=/home/jdoe/myfile.csv'")
	ptrIgnoreFile                 = flag.String("ignorefile", userHomeDir+string(os.PathSeparator)+".ge_ignore", "Patterns File (including wildcards) for files to ignore.  (e.g. *.jpg)")
	ptrFailSeverityThreshold      = flag.String("fail-severity", cfgreader.Settings.TranslateLevelID(cfgreader.Settings.FailThreshold), "Lowest severity level at which to fail "+levelOptions)
	ptrDisplaySeverityThreshold   = flag.String("display-severity", cfgreader.Settings.TranslateLevelID(cfgreader.Settings.DisplayThreshold), "Lowest severity level to display "+levelOptions)
	ptrDisplayConfidenceThreshold = flag.String("display-confidence", cfgreader.Settings.TranslateLevelID(cfgreader.Settings.DisplayConfidenceThreshold), "Lowest confidence level to display "+levelOptions)
	ptrFailConfidenceThreshold    = flag.String("fail-confidence", cfgreader.Settings.TranslateLevelID(cfgreader.Settings.FailThreshold), "Lowest confidence level at which to fail "+levelOptions)
	ptrModuleConfigFile           = flag.String("module-config-file", "", "Path to file with per module config settings")
)
