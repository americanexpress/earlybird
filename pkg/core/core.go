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

package core

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/americanexpress/earlybird/pkg/api"
	cfgreader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/file"
	"github.com/americanexpress/earlybird/pkg/git"
	"github.com/americanexpress/earlybird/pkg/scan"
	configupdate "github.com/americanexpress/earlybird/pkg/update"
	"github.com/americanexpress/earlybird/pkg/utils"
	"github.com/americanexpress/earlybird/pkg/writers"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
)

//GitClone clones git repositories into a temporary directory
func (eb *EarlybirdCfg) GitClone(ptr PTRGitConfig) {
	var scanRepos []string
	gitPassword := os.Getenv("gitpassword")
	if *ptr.Repo != "" {
		scanRepos = []string{*ptr.Repo}
		eb.Config.Gitrepo = *ptr.Repo
	}

	if *ptr.Project != "" {
		if *ptr.RepoUser == "" {
			fmt.Println("Please use the -git-user flag to scan a Git Project or Organisation ")
			os.Exit(1)
		}

		gitPassword = utils.GetGitURL(ptr.Repo, ptr.RepoUser)
		scanRepos = git.ReposPerProject(*ptr.Project, *ptr.RepoUser, gitPassword)

		if eb.Config.OutputFormat != "json" && !(*ptrStreamInput) {
			fmt.Println("Cloning", len(scanRepos), "Repositories in", utils.GetGitProject(*ptr.Project))
		}
	}

	// Display the directory or repo being scanned
	if len(scanRepos) != 0 {
		if gitPassword == "" {
			gitPassword = utils.GetGitURL(ptr.Repo, ptr.RepoUser)
		}
		var err error
		if *ptr.RepoUser != "" { // use auth
			eb.Config.SearchDir, err = git.CloneGitRepos(scanRepos, *ptr.RepoUser, gitPassword, (eb.Config.OutputFormat == "json"))
		} else {
			eb.Config.SearchDir, err = git.CloneGitRepos(scanRepos, "", "", (eb.Config.OutputFormat == "json")) //Blank no auth
		}
		if err != nil {
			fmt.Println("Failed to clone repository:", err)
			os.Exit(1)
		}
	} else {
		if eb.Config.OutputFormat != "json" && !(*ptrStreamInput) {
			fmt.Println("Scanning directory: ", eb.Config.SearchDir)
		}
	}
}

//StartHTTP spins up the Earlybird REST API server
func (eb *EarlybirdCfg) StartHTTP(ptr PTRHTTPConfig) {
	// Set up http server
	r := mux.NewRouter()
	r.HandleFunc("/scan/git", api.GITScan(eb.Config)).Methods("GET")
	r.HandleFunc("/scan", api.Scan(eb.Config)).Methods("POST")
	r.HandleFunc("/labels", api.Labels(eb.Config.Version, scan.Labels)).Methods("GET")
	r.HandleFunc("/categorylabels", api.LabelsPerCategory(eb.Config.Version, scan.Labels)).Methods("GET")
	r.HandleFunc("/categories", api.Categories(eb.Config.Version, scan.CombinedRules)).Methods("GET")
	// Catch-all: Serve our JavaScript application's entry-point (index.html) and static assets directly.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(userHomeDir + string(os.PathSeparator) + ".eb-wa-build" + string(os.PathSeparator))))
	//Default time out settings
	serverconfig := cfgreader.ServerConfig{
		WriteTimeout: 60,
		ReadTimeout:  60,
		IdleTimeout:  120,
	}

	if *ptr.HTTPConfig != "" {
		var serverconfig cfgreader.ServerConfig
		err := cfgreader.LoadConfig(&serverconfig, *ptr.HTTPConfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	srv := &http.Server{
		Addr: *ptr.HTTP,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(serverconfig.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(serverconfig.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(serverconfig.IdleTimeout),
		Handler:      r,
	}

	if *ptr.HTTPS != "" {
		srv.Addr = *ptr.HTTPS
		err := http2.ConfigureServer(srv, &http2.Server{})
		if err != nil {
			log.Fatal("Failed to configure HTTP server", err)
		}
		fmt.Println("go-earlybird HTTPS/2 API Listening on", *ptr.HTTPS)
		log.Fatal(srv.ListenAndServeTLS(*ptr.HTTPSCert, *ptr.HTTPSKey))
	} else {
		fmt.Println("go-earlybird HTTP API Listening on", *ptr.HTTP)
		log.Fatal(srv.ListenAndServe())
	}
}

//ConfigInit loads in the earlybird configuration and CLI flags
func (eb *EarlybirdCfg) ConfigInit() {
	eb.Config.ConfigDir = utils.GetConfigDir()
	eb.Config.LevelMap = cfgreader.Settings.GetLevelMap()
	// Build the string to display available modules for the CLI flags
	availableModules := cfgreader.Settings.GetAvailableModules()

	//Load CLI arguments and parse
	flag.Var(&enableFlags, "enable", "Enable individual scanning modules "+utils.GetDisplayList(availableModules))
	flag.Parse()

	//Assign CLI arguments to our global configuration
	eb.Config.WorkerCount = *ptrWorkerCount
	eb.Config.WorkLength = *ptrWorkLength
	eb.Config.ShowFullLine = *ptrShowFullLine
	eb.Config.MaxFileSize = *ptrMaxFileSize
	eb.Config.VerboseEnabled = *ptrVerbose
	eb.Config.Suppress = *ptrSuppressSecret
	eb.Config.OutputFormat = *ptrOutputFormat
	eb.Config.OutputFile = *ptrOutputFile
	eb.Config.SearchDir = *ptrPath
	eb.Config.ConfigDir = *ptrConfigDir
	eb.Config.IgnoreFile = *ptrIgnoreFile
	eb.Config.GitStream = *ptrGitStreamInput
	eb.Config.RulesOnly = *ptrRulesOnly
	eb.Config.SkipComments = *ptrSkipComments
	eb.Config.IgnoreFPRules = *ptrIgnoreFPRules
	eb.Config.ShowSolutions = *ptrShowSolutions

	// If the streaming IO flag was specified, accept the streaming input
	if *ptrStreamInput || eb.Config.GitStream {
		eb.Config.SearchDir = ""
	}

	// Check to see if the user opted to update.  If they choose this option
	// the configuration files will be updated and the program will exit.
	if *ptrUpdateFlag {
		configPath := utils.GetConfigDir() + "earlybird.json"
		doUpdate(eb.Config.ConfigDir, configPath, cfgreader.Settings.ConfigFileURL)
	}

	eb.Config.Version = cfgreader.Settings.Version
	// Set the skip options (what not to scan) from configs
	eb.Config.AnnotationsToSkipLine = cfgreader.Settings.AnnotationsToSkip
	eb.Config.ExtensionsToSkipScan = cfgreader.Settings.ExtensionsToSkipTextScan
	// Determine which results to show and which to fail on
	eb.Config.SeverityDisplayLevel = cfgreader.Settings.TranslateLevelName(*ptrDisplaySeverityThreshold)
	eb.Config.SeverityFailLevel = cfgreader.Settings.TranslateLevelName(*ptrFailSeverityThreshold)
	// Determine which results to show and which to fail on based on confidence
	eb.Config.ConfidenceDisplayLevel = cfgreader.Settings.TranslateLevelName(*ptrDisplayConfidenceThreshold)
	eb.Config.ConfidenceFailLevel = cfgreader.Settings.TranslateLevelName(*ptrFailConfidenceThreshold)
	// Let's see if we have specified git tracked/staged files
	eb.Config.TargetType = utils.GetTargetType(*ptrGitStagedFlag, *ptrGitTrackedFlag)
	eb.Config.EnabledModules = utils.GetEnabledModules(enableFlags, availableModules)
}

//Scan Runs the scan by kicking off the different modules as go routines
func (eb *EarlybirdCfg) Scan() {
	// Validate the path passed in as the target directory to scan
	start := time.Now()
	fileContext, err := eb.FileContext()
	if err != nil {
		log.Fatal("Failed to get FileContext: ", err)
	}
	HitChannel := make(chan scan.Hit)
	go scan.SearchFiles(&eb.Config, fileContext.Files, fileContext.CompressPaths, HitChannel)

	// Send output to a writer
	eb.WriteResults(start, HitChannel, fileContext)

	utils.DeleteGit(eb.Config.Gitrepo, eb.Config.SearchDir)
	if eb.Config.FailScan {
		if eb.Config.OutputFormat == "console" {
			fmt.Fprintln(os.Stderr, "Scan detected findings above the accepted threshold -- Failing.")
		}
		os.Exit(1)
	}
}

//FileContext provides an inclusive file system context of our scan
func (eb *EarlybirdCfg) FileContext() (fileContext file.Context, err error) {
	cfg := eb.Config
	if cfg.SearchDir != "" {
		// We're going to load a 'files' slice based on the CLI args
		switch cfg.TargetType {
		case utils.Tracked:
			return file.GetGitFiles(utils.Tracked, &cfg)
		case utils.Staged:
			return file.GetGitFiles(utils.Staged, &cfg)
		default:
			return file.GetFiles(cfg.SearchDir, cfg.IgnoreFile, cfg.VerboseEnabled, cfg.MaxFileSize)
		}
	}
	if cfg.GitStream {
		var err error
		fileContext.Files, err = git.ParseGitLog(bufio.NewReader(os.Stdin))
		if err != nil {
			return fileContext, err
		}
	}
	fileContext.Files = file.GetFileFromStream(&cfg)
	return fileContext, nil
}

//WriteResults reads hits from the channel to the console or target file
func (eb *EarlybirdCfg) WriteResults(start time.Time, HitChannel chan scan.Hit, fileContext file.Context) {
	// Send output to a writer
	var err error
	switch {
	case eb.Config.OutputFormat == "json":
		var Hits []scan.Hit
		for hit := range HitChannel {
			Hits = append(Hits, hit)
		}
		report := scan.Report{
			Hits:          Hits,
			HitCount:      len(Hits),
			Skipped:       fileContext.SkippedFiles,
			Ignore:        fileContext.IgnorePatterns,
			Version:       eb.Config.Version,
			Modules:       eb.Config.EnabledModules,
			Threshold:     eb.Config.SeverityDisplayLevel,
			FilesScanned:  len(fileContext.Files),
			RulesObserved: len(scan.CombinedRules),
			StartTime:     start.UTC().Format(time.RFC3339),
			EndTime:       time.Now().UTC().Format(time.RFC3339),
			Duration:      string(time.Since(start)),
		}
		_, err = writers.WriteJSON(report, eb.Config.OutputFile)
	case eb.Config.OutputFormat == "csv":
		err = writers.WriteCSV(HitChannel, eb.Config.OutputFile)
	default:
		err = writers.WriteConsole(HitChannel, eb.Config.OutputFile, eb.Config.ShowFullLine)
		fmt.Printf("\n%d files scanned in %s", len(fileContext.Files), time.Since(start))
		fmt.Printf("\n%d rules observed\n", len(scan.CombinedRules))
	}
	if err != nil {
		fmt.Println("Writing Results failed:", err)
	}
}

// Update configs from the latest in the repo
func doUpdate(configDir string, configPath string, appConfigURL string) {
	err := configupdate.UpdateConfigFiles(configDir, configPath, appConfigURL)
	if err != nil {
		log.Fatal("Failed to update config:", err)
	}
	fmt.Println("Configurations updated.  Exiting")
	os.Exit(0)
}
