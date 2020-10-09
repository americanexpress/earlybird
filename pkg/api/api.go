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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	cfgreader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/file"
	"github.com/americanexpress/earlybird/pkg/git"
	"github.com/americanexpress/earlybird/pkg/scan"
	"github.com/americanexpress/earlybird/pkg/utils"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

//Scan uses the Earlybird config to search uploaded multipart files for secrets
func Scan(cfg cfgreader.EarlybirdConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(1024 << 20) // 1GB upload limit
		if err != nil {
			http.Error(w, "File upload too large: "+err.Error(), http.StatusInternalServerError)
			return
		}
		start := time.Now()

		//Get files from req
		formdata := r.MultipartForm
		fileList, err := file.MultipartToScanFiles(formdata.File["scan"], cfg)
		if err != nil {
			http.Error(w, "Failed to parse file upload: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Define our result objects and start scan process
		var Hits []scan.Hit
		HitChannel := make(chan scan.Hit)
		go scan.SearchFiles(&cfg, fileList, []string{}, HitChannel)

		for hit := range HitChannel {
			Hits = append(Hits, hit)
		}

		//Format our results into an Earlybird report
		report := scan.Report{
			Hits:          Hits,
			HitCount:      len(Hits),
			Version:       cfg.Version,
			Modules:       cfg.EnabledModules,
			Threshold:     cfg.SeverityDisplayLevel,
			FilesScanned:  len(fileList),
			RulesObserved: len(scan.CombinedRules),
			StartTime:     start.UTC().Format(time.RFC3339),
			EndTime:       time.Now().UTC().Format(time.RFC3339),
			Duration:      string(time.Since(start)),
		}

		//Encode and send JSON response
		response, err := json.MarshalIndent(report, "", "\t")
		if err != nil {
			http.Error(w, "Failed to encode JSON response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(response))
	}
}

//GITScan searches for secrets in git repositories based off the Earlybird config, supports authentication via env variables "gituser" and "gitpassword"
func GITScan(cfg cfgreader.EarlybirdConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var blank string
		start := time.Now()
		mycfg := cfg

		giturls, ok := r.URL.Query()["url"]
		if !ok || len(giturls[0]) < 1 {
			http.Error(w, "GIT URL Parameter is missing", http.StatusInternalServerError)
			return
		}

		giturl := giturls[0]
		utils.GetGitURL(&giturl, &blank)
		mycfg.SearchDir, err = git.CloneGitRepos([]string{giturl}, os.Getenv("gituser"), os.Getenv("gitpassword"), (cfg.OutputFormat == "json"))
		if err != nil {
			if err == transport.ErrAuthenticationRequired {
				http.Error(w, "Failed to clone, repository is private. Please enter a public repository URL.", http.StatusInternalServerError)
			} else {
				http.Error(w, "Failed to clone, please verify your repository is available", http.StatusInternalServerError)
			}
			utils.DeleteGit(giturl, mycfg.SearchDir)
			return
		}

		fileContext, err := file.GetFiles(mycfg.SearchDir, mycfg.IgnoreFile, mycfg.VerboseEnabled, cfg.MaxFileSize)
		if err != nil {
			http.Error(w, "Failed to load scan files: "+err.Error(), http.StatusInternalServerError)
			return
		}
		//Delete our tmp directory when done
		defer utils.DeleteGit(giturl, mycfg.SearchDir)
		// Start building a list of hits.  The module go routines will all dump back to this
		var Hits []scan.Hit
		HitChannel := make(chan scan.Hit)
		//Create pointer to reduce memory overhead
		go scan.SearchFiles(&mycfg, fileContext.Files, fileContext.CompressPaths, HitChannel)

		for hit := range HitChannel {
			Hits = append(Hits, hit)
		}

		report := scan.Report{
			Hits:          Hits,
			HitCount:      len(Hits),
			Skipped:       fileContext.SkippedFiles,
			Ignore:        fileContext.IgnorePatterns,
			Version:       cfg.Version,
			Modules:       mycfg.EnabledModules,
			Threshold:     mycfg.SeverityDisplayLevel,
			FilesScanned:  len(fileContext.Files),
			RulesObserved: len(scan.CombinedRules),
			StartTime:     start.UTC().Format(time.RFC3339),
			EndTime:       time.Now().UTC().Format(time.RFC3339),
			Duration:      string(time.Since(start)),
		}

		response, err := json.MarshalIndent(report, "", "\t")
		if err != nil {
			http.Error(w, "Failed to encode JSON response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(response))
	}
}

//Labels returns all the available Earlybird labels in "LabelsReponse" format
func Labels(version string, scanLabels map[int]scan.LabelConfigs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Build label list
		var labels []string
		for _, lblcfgs := range scanLabels { //Cycle through all the label config arrays
			for _, lblcfg := range lblcfgs.Labels { //Cycle through label configs
				if !utils.Contains(labels, lblcfg.Label) { //Check if labels is in list before appending
					labels = append(labels, lblcfg.Label) //Append unique label to list
				}
			}
		}

		resp := LabelsResponse{
			Version: version,
			Labels:  labels,
		}

		response, err := json.MarshalIndent(resp, "", "\t")
		if err != nil {
			http.Error(w, "Failed to encode JSON response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(response))
	}
}

//Categories returns all the available Earlybird categories in "CategoriesReponse" format
func Categories(version string, rules []scan.Rule) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categories []string
		for _, rule := range rules { //Cycle through all the rules
			if !utils.Contains(categories, rule.Category) { //Check if catgory is in list before appending
				categories = append(categories, rule.Category) //Append unique category to list
			}
		}

		resp := CategoriesResponse{
			Version:    version,
			Categories: categories,
		}

		response, err := json.MarshalIndent(resp, "", "\t")
		if err != nil {
			http.Error(w, "Failed to encode JSON response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(response))
	}
}

//LabelsPerCategory returns all the available Earlybird labels per categories in "CategoryLabelsResponse" format
func LabelsPerCategory(version string, scanLabels map[int]scan.LabelConfigs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categorylabels := labelsToLabelsPerCategory(scanLabels)
		resp := CategoryLabelsResponse{
			Version:        version,
			CategoryLabels: categorylabels,
		}

		response, err := json.MarshalIndent(resp, "", "\t")
		if err != nil {
			http.Error(w, "Failed to encode JSON response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(response))
	}
}

func labelsToLabelsPerCategory(scanLabels map[int]scan.LabelConfigs) (categorylabels map[string][]string) {
	categorylabels = make(map[string][]string)
	for _, lblcfgs := range scanLabels { //Cycle through all the label config arrays
		for _, lblcfg := range lblcfgs.Labels { //Cycle through label configs
			if list, ok := categorylabels[lblcfg.Category]; ok { //Check if category exists
				if !utils.Contains(categorylabels[lblcfg.Category], lblcfg.Label) { //Check if labels is in category list before appending
					categorylabels[lblcfg.Category] = append(list, lblcfg.Label) //Append unique label to category list
				}
			} else { // If category doesn't exist, create it and assign first label
				categorylabels[lblcfg.Category] = []string{lblcfg.Label}
			}
		}
	}
	return categorylabels
}
