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
	"bufio"
	"crypto/sha1"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	cfgReader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/postprocess"
)

var (
	//Labels is a map of all our labels, accessible by the rule unique code
	Labels map[int]LabelConfigs
	//CombinedRules is a global array where we load all our precompiled rules
	CombinedRules []Rule
	//FalsePositiveRules is a map of our false positive rules sorted by the rule unique code
	FalsePositiveRules map[int]FalsePositives
	//SolutionConfigs is a map of our solutions sorted by the rule unique code
	SolutionConfigs map[int]Solution
	//CompressPattern is a pattern used to identify compressed zip files
	CompressPattern = regexp.MustCompile(compressRegex)
	tempPattern     = regexp.MustCompile(tempRegex)
)

//SearchFiles will use the EarlybirdConfig, the provided file list, decompressed zip files tempoary paths to send found secrets to the Hit channel
func SearchFiles(cfg *cfgReader.EarlybirdConfig, files []File, compressPaths []string, hits chan<- Hit) {
	//Delete tmp file directory when we're done
	defer DeleteFiles(compressPaths)

	//Create our channels and mutex
	var jobMutex = &sync.Mutex{}
	jobs := make(chan WorkJob)
	wg := new(sync.WaitGroup)

	//Create our worker pool
	scanPool(cfg, wg, jobMutex, jobs, hits)

	//Scan the file names
	nameScanner(cfg, files, hits)

	//Create work from file content for the scanPool
	contentJobWriter(cfg, files, jobMutex, jobs)

	//Close our channels
	close(jobs)
	wg.Wait()
	close(hits)
}

//scanPool searches incoming jobs for secrets and write findings to hits channel
func scanPool(cfg *cfgReader.EarlybirdConfig, wg *sync.WaitGroup, jobMutex *sync.Mutex, jobs chan WorkJob, hits chan<- Hit) {
	//Create duplicate map
	dupeMap := make(map[string]bool) //HASH:true
	for w := 1; w <= 100; w++ {
		wg.Add(1)
		go func(w int) {
			for j := range jobs {
				if IsIgnoreAnnotation(cfg, j.WorkLine.LineValue) {
					j.WorkLine.LineValue = ""
				}

				// Scan the line based on common password rules
				hitFound, tmpHits := scanLine(j.WorkLine, j.FileLines, cfg)
				if cfg.Suppress {
					for i := range tmpHits {
						tmpHits[i].MatchValue = maskValue(tmpHits[i].MatchValue)
						tmpHits[i].LineValue = maskValue(tmpHits[i].LineValue)
					}
				}
				if hitFound {
					for _, hit := range tmpHits {
						jobMutex.Lock() // put a mutex on it to avoid collisions/misses
						if !hitUnique(dupeMap, hit) {
							jobMutex.Unlock()
							continue
						}
						jobMutex.Unlock()
						hits <- hit //Push hits to channel

						cfg.FailScan = determineScanFail(cfg, &hit)
					}
				}
			}
			defer wg.Done()
		}(w)
	}
}

// determine if we should fail scan based on severity and confidence
func determineScanFail(cfg *cfgReader.EarlybirdConfig, hit *Hit) bool {
	if hit.Severity == infoLevelSeverity {
		return false
	}

	return hit.SeverityID <= cfg.SeverityFailLevel && hit.ConfidenceID <= cfg.ConfidenceFailLevel
}

//contentJobWriter creates work based off file content for scanning
func contentJobWriter(cfg *cfgReader.EarlybirdConfig, files []File, jobMutex *sync.Mutex, jobs chan WorkJob) {
	var e error
	// Loop through each File
	for _, searchFile := range files {
		//FileOS refers to the file object that's open, not the file object which contains the name and path
		if searchFile.Path == "buffer" || searchFile.Name == "buffer" {
			for _, workline := range searchFile.Lines {
				jobs <- WorkJob{
					WorkLine:  workline,
					FileLines: searchFile.Lines,
				}
			}
		} else {
			//Don't do file read/scan on files we know will trigger the filename scan -- Don't open compressed files either
			if !isExcludedFileType(cfg, searchFile.Name) && len(CompressPattern.FindStringSubmatch(searchFile.Name)) <= 0 {
				fileOS, err := os.Open(searchFile.Path) //Open file path
				if err != nil {
					fileOS, err = os.Open(searchFile.Name) //If file path open fails, try file name
					if err != nil {
						log.Fatal("Can't open file", err)
					}
				}
				var work []WorkJob
				var job WorkJob
				job.FileLines = searchFile.Lines

				//Search line by line
				reader := bufio.NewReader(fileOS)
				job.WorkLine.LineValue, e = readln(reader)
				for e == nil {
					job.WorkLine.LineNum = job.WorkLine.LineNum + 1
					job.WorkLine.FileName = jobFileName(cfg.Gitrepo, searchFile.Name)
					job.WorkLine.FilePath = searchFile.Path
					job.FileLines = append(job.FileLines, job.WorkLine)

					//Add our split up jobs to the work array
					work = append(work, splitJob(job, cfg.WorkLength)...)
					//Search next line to break out of loop
					job.WorkLine.LineValue, e = readln(reader)
					if e != nil && e != io.EOF {
						log.Println("Error reading file:", e)
					}
				}
				//Push our work to the jobs channel
				for _, job := range work {
					jobs <- job
				}
				fileOS.Close()
			}
		}
	}
}

//nameScanner scans file names for sensitive values
func nameScanner(cfg *cfgReader.EarlybirdConfig, files []File, hits chan<- Hit) {
	for _, file := range files {
		// Scan the filename based on the Filename rules
		hitFound, hit := scanName(file, CombinedRules, cfg.LevelMap, cfg.ShowSolutions)
		if hitFound {

			// Append the hit to our slice for return
			if i := cfg.LevelMap[hit.Severity]; i <= cfg.SeverityDisplayLevel {
				hits <- hit //push hit to channel
			}

			// If a hit severity is less than the failLevel, set failScan = true
			if i := cfg.LevelMap[hit.Severity]; i <= cfg.SeverityFailLevel {
				cfg.FailScan = true
			}

			// If a hit confidence is less than the failLevel, set failScan = true
			if i := cfg.LevelMap[hit.Confidence]; i <= cfg.ConfidenceFailLevel {
				cfg.FailScan = true
			}
		}
	}

}

//DeleteFiles removes files and folders in target path array
func DeleteFiles(paths []string) {
	for _, p := range paths {
		err := os.RemoveAll(p)
		if err != nil {
			log.Println("Failed to delete temporary file", err)
		}
	}
}

// Check if the file extension is something we know will trigger a hit on the filename scan (e.g., .pem, .p12, etc.
func isExcludedFileType(cfg *cfgReader.EarlybirdConfig, filename string) (excluded bool) {
	for _, ext := range cfg.ExtensionsToSkipScan {
		if strings.EqualFold(filepath.Ext(filename), ext) {
			return true
		}

		//filename ends in extension stripped of period, e.g., 'foobarmin.js'
		trimmedExt := string(ext[1:])
		if strings.HasSuffix(filename, trimmedExt) {
			return true
		}
	}
	return false
}

func hitUnique(dupeMap map[string]bool, hit Hit) bool {
	digest := sha1.New()
	_, err := digest.Write([]byte(hit.Filename + strconv.Itoa(hit.Line) + hit.MatchValue))
	if err != nil {
		log.Println("Failed to produce digest of hit", err)
	}
	hithash := string(digest.Sum(nil))
	//hash hit here
	if exists, ok := (dupeMap)[hithash]; ok && exists {
		//This is a duplicate
		return false
	}
	(dupeMap)[hithash] = true
	return true
}

// Take a line and run through the rules, looking for a hit
func scanLine(line Line, fileLines []Line, cfg *cfgReader.EarlybirdConfig) (isHit bool, hits []Hit) {
	for _, rule := range CombinedRules {
		var hit Hit
		//Skip rules that do not apply
		if rule.Searcharea == "filename" || cfg.SkipComments && rule.Category == "comment" {
			continue
		}

		patternMatch, matchValue := findHit(line.LineValue, rule.CompiledPattern)

		if !patternMatch {
			continue
		}

		//If we found a Regexp match, build a Hit
		hit.Code = rule.Code
		hit.Confidence = getLevelNameFromID(rule.Confidence, cfg.LevelMap)
		hit.ConfidenceID = rule.Confidence
		hit.Caption = rule.Caption
		hit.Category = rule.Category
		if cfg.ShowSolutions {
			hit.Solution = SolutionConfigs[rule.SolutionID].Text
		}
		hit.CWE = rule.CWE
		hit.Line = line.LineNum
		hit.LineValue = strings.TrimSpace(line.LineValue)
		hit.MatchValue = matchValue
		if line.FilePath != "buffer" {
			hit.Filename = removeTempPrefix(line.FilePath)
		} else {
			hit.Filename = line.FileName
		}
		hit.Time = time.Now().UTC().Format(time.RFC3339)
		hit.determineSeverity(cfg, &rule)

		// Apply labels to the hit if appropriate
		labelHit(&hit, fileLines)

		//Check if our hit has any false positives
		isStillHit := hit.postProcess(cfg, &rule)
		if isStillHit {
			isHit = true
			hits = append(hits, hit)
		}

	}
	return isHit, hits
}

// Take a filename and run through the rules, looking for a hit
func scanName(file File, rules []Rule, levelMap map[string]int, showSolutions bool) (isHit bool, hit Hit) {
	for _, rule := range rules {
		if rule.Searcharea == "body" { //Skip rules that do not apply
			continue
		}

		if file.Path == "buffer" {
			file.Path = file.Name
		}

		patternMatch, _ := findHit(file.Path, rule.CompiledPattern)

		// If we found a match to the Regexp pattern, build a Hit
		if patternMatch {
			hit.Code = rule.Code
			hit.Severity = getLevelNameFromID(rule.Severity, levelMap)
			hit.SeverityID = rule.Severity
			hit.Caption = rule.Caption
			hit.Category = rule.Category
			hit.CWE = rule.CWE
			hit.Confidence = getLevelNameFromID(rule.Confidence, levelMap)
			hit.ConfidenceID = rule.Confidence
			if showSolutions {
				hit.Solution = SolutionConfigs[rule.SolutionID].Text
			}
			hit.Line = 0
			hit.Filename = filepath.Base(file.Path)
			hit.MatchValue = file.Name
			hit.LineValue = file.Name
			hit.Time = time.Now().UTC().Format(time.RFC3339)
			return true, hit
		}
	}
	return false, hit
}

// readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

//IsIgnoreAnnotation Checks for ignore annotation
func IsIgnoreAnnotation(cfg *cfgReader.EarlybirdConfig, line string) bool {
	for _, annotation := range cfg.AnnotationsToSkipLine {
		if strings.Contains(line, annotation) {
			return true
		}
	}
	return false
}

// If we want to suppress a secret value from being displayed in the results, mask it with maskCharacter
func maskValue(input string) string {
	return strings.Repeat(maskCharacter, len(input))
}

func jobFileName(gitRepo, fileName string) string {
	if gitRepo != "" {
		return getFileURL(gitRepo, filepath.Base(fileName))
	}
	return fileName
}

//splitJob splits up the job into an array of jobs if too long otherwise returns a single job
func splitJob(inJob WorkJob, worklength int) (work []WorkJob) {
	//If line isn't too long, just push job to work channel
	if len(inJob.WorkLine.LineValue) <= worklength {
		return []WorkJob{inJob}
	}

	//For VERY long lines, split it up at WORK_LENGTH, creating another string that overlaps the split
	linesValues := splitSubN(inJob.WorkLine.LineValue, worklength)
	for _, value := range linesValues {
		outJob := WorkJob{
			WorkLine: Line{
				LineNum:   inJob.WorkLine.LineNum,
				FileName:  inJob.WorkLine.FileName,
				FilePath:  inJob.WorkLine.FilePath,
				LineValue: value,
			},
			FileLines: inJob.FileLines,
		}
		work = append(work, outJob)
	}
	return work
}

//splitSubN Create the overlap string when splitting long strings
func splitSubN(s string, n int) []string {
	runes := []rune(s)
	chunks := make([]string, 0, len(runes)/n)
	for start := 0; start < len(runes); start += n {
		if end := start + n; end < len(runes) {
			chunks = append(chunks, string(runes[start:end]))
		} else {
			chunks = append(chunks, string(runes[start:]))
		}
	}

	results := []string{}
	//subs contains all split strings
	//iterate over strings parsing
	toggle := true
	for _, sub := range chunks {
		var tmpString string
		if toggle { // Check if we should parse from the end of the string
			toggle = false
			results = append(results, sub) // Append split string
			if len(sub) >= overlapLength {
				tmpString = sub[len(sub)-overlapLength:]
			} else {
				if len(sub) > 0 {
					tmpString = sub[len(sub)-(len(sub)-1):]
				} else {
					tmpString = sub[0:]
				}
			}
			continue
		}
		// parse from the start of the string
		toggle = true
		if len(sub) > overlapLength {
			tmpString = tmpString + sub[0:overlapLength-1]
			results = append(results, tmpString) //Append overlapped data
			results = append(results, sub)       // Append split string
			tmpString = ""
		} else {
			results = append(results, sub) // Append split string
			break                          //stop if last element is too short
		}
	}
	return results
}

// From the configs in labels.json, apply labels to each hit as appropriate
func labelHit(hit *Hit, fileLines []Line) {
	rules, ok := Labels[hit.Code]
	if !ok {
		return
	}

	for _, rule := range rules.Labels {
		if !rule.Multiline {
			// Some labels get applied based on the actual hit, not the key context
			if len(rule.Keys) == 0 {
				hit.Labels = append(hit.Labels, rule.Label)
				continue
			}
			var matched bool
			for _, key := range rule.Keys {
				if substringExistsInString(hit.LineValue, key) {
					matched = true
					break
				}
			}
			if matched {
				hit.Labels = append(hit.Labels, rule.Label)
			}
			continue
		}

		criteriaMatched := 0
		for _, key := range rule.Keys {
			if substringExistsInLines(fileLines, key) {
				criteriaMatched++
			}
		}
		neededCriteria := len(rule.Keys)
		if criteriaMatched == neededCriteria {
			hit.Labels = append(hit.Labels, rule.Label)
		}
	}
}

func (hit *Hit) determineSeverity(cfg *cfgReader.EarlybirdConfig, rule *Rule) {
	// check if for the given category we need to adjust the severity based on user config
	for _, adjustedSeverityCategoryCfg := range cfg.AdjustedSeverityCategories {
		if adjustedSeverityCategoryCfg.Category == rule.Category {
			for _, p := range adjustedSeverityCategoryCfg.CompiledPatterns {
				var test string

				if adjustedSeverityCategoryCfg.UseFilename {
					test = hit.Filename
				} else if adjustedSeverityCategoryCfg.UseLineValue {
					test = hit.LineValue
				} else {
					test = hit.MatchValue
				}

				if p.Match([]byte(test)) {
					severityId := getIdFromLevelName(adjustedSeverityCategoryCfg.AdjustedDisplaySeverity, cfg.LevelMap)

					hit.Severity = getLevelNameFromID(severityId, cfg.LevelMap)
					hit.SeverityID = severityId
					return
				}
			}
		}
	}

	hit.Severity = getLevelNameFromID(rule.Severity, cfg.LevelMap)
	hit.SeverityID = rule.Severity
}

func (hit *Hit) postProcess(cfg *cfgReader.EarlybirdConfig, rule *Rule) (isHit bool) {
	fpHit := false
	if !cfg.IgnoreFPRules {
		fpHit = findFalsePositive(*hit)
	}
	switch {
	case fpHit:
		isHit = false
		// Check if a password is valid and weak.  Exclude if invalid, label as 'weak' if weak
	case rule.Postprocess == "password":
		// Skip account_token as password so that it can be reported under credit card
		SkipAccountToken := postprocess.SkipAccountTokenPassword(hit.LineValue)
		if SkipAccountToken {
			isHit = false
			break
		}
		// If it's a false positive return no match.
		Confidence, IsFalsePositive := postprocess.PasswordFalse(hit.MatchValue)
		if IsFalsePositive {
			isHit = false
			break
		}
		// Skip password as same key/value pair
		IsPasswordSameKeyValue := postprocess.SkipSameKeyValuePassword(hit.MatchValue, hit.LineValue)
		if IsPasswordSameKeyValue {
			isHit = false
			break
		}

		// Skip password if the value has unicode char in it
		passwordContainsUnicode := postprocess.SkipPasswordWithUnicode(hit.MatchValue)
		if passwordContainsUnicode {
			isHit = false
			break
		}

		hit.Confidence = getLevelNameFromID(Confidence, cfg.LevelMap)
		hit.ConfidenceID = Confidence

		if postprocess.PasswordWeak(hit.MatchValue) {
			hit.Caption = postprocess.WeakPswdCaption
			hit.Labels = append(hit.Labels, "weak password")
		}
		isHit = true

		// If a SSN hit doesn't meet certain criteria (e.g., all zeroes, certain test patterns, etc.), skip it
	case rule.Postprocess == "ssn":
		if postprocess.ValidSSN(hit.MatchValue) {
			isHit = true
		}

		// Verify credit card hits against a mod10 check
	case rule.Postprocess == "mod10":
		// If the match passed a Luhn/mod-10 check, build a Hit
		if postprocess.IsCard(hit.MatchValue) {
			isHit = true
		}

		// Calculate the entropy of a string and make sure it passes entropyThreshold
	case rule.Postprocess == "entropy":
		e := postprocess.Shannon(hit.MatchValue)

		// If the line's string entropy is high enough, build a Hit
		if e > entropyThreshold {
			isHit = true
		}

		// No additional validation needed
	default:
		isHit = true
	}
	return isHit
}

//removeTempPrefix removes the temp path prefix if it exists
func removeTempPrefix(path string) string {
	if strings.Contains(path, "ebzip") || strings.Contains(path, "ebgit") {
		if paths := tempPattern.FindStringSubmatch(path); len(paths) > 1 {
			path = paths[1]
		}
	}
	return path
}
