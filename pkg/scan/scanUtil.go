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
	"regexp"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

// Look for a regexp pattern hit in a string
func findHit(target string, CompiledPattern *regexp.Regexp) (isHit bool, retMatch string) {
	if target != "" {
		matchValue := CompiledPattern.FindString(target)
		if matchValue != "" {
			return true, prepareMatchValue(matchValue)
		}
	}
	return false, ""
}

//substringExistsInLines Search for a regexp pattern occurring anywhere in a file
func substringExistsInLines(fileLines []Line, str string) bool {
	reg := regexp.MustCompile("(?i)" + str)
	for _, line := range fileLines {
		if reg.MatchString(line.LineValue) {
			return true
		}
	}
	return false
}

//substringExistsInString check if sub exists in string
func substringExistsInString(str string, substr string) bool {
	m := search.New(language.English, search.IgnoreCase)
	start, _ := m.IndexString(str, substr)
	return start >= 0
}

// If it makes sense to strip leading/trailing characters for readability, let's do it
func prepareMatchValue(matchValue string) (isolatedValue string) {
	origMatchValue := matchValue
	if len(matchValue) > 0 && shouldStrip(matchValue[0]) {
		matchValue = matchValue[1:]
	}
	if len(matchValue) > 0 && shouldStrip(matchValue[len(matchValue)-1]) {
		matchValue = matchValue[:len(matchValue)-1]
	}

	// If the string still contains quotes/ticks, we may not have wanted to strip the leading/trailing quotes/ticks
	if stringContainsQuotesTicks(matchValue) {
		return origMatchValue
	}
	return matchValue
}

// We shouldn't strip leading/trailing quotes/ticks/etc if the string still has quotes/ticks
func stringContainsQuotesTicks(s string) bool {
	if strings.Contains(s, "\"") || strings.Contains(s, "'") {
		return true
	}
	return false
}

// Is the character in the strippable list
func shouldStrip(b byte) bool {
	switch b {
	case '"', '\'', '\\', '/', '=', ',', '|':
		return true
	default:
		return false
	}

}

// Verify that a hit is a false positive
func findFalsePositive(hit Hit) (isFP bool) {
	var scan bool
	if rules, ok := FalsePositiveRules[hit.Code]; ok { //Veriy a false positive rule exists for this hit code
		for _, rule := range rules.FalsePositives {
			scan = true
			if len(rule.FileExtensions) > 0 { //Check if this rule only applies to certain files
				scan = false
				for _, fileExtension := range rule.FileExtensions { //Cycle through file extensions to verify rule applies to hit
					if strings.Contains(hit.Filename, fileExtension) {
						scan = true //Trigger a value scan if the file name matches
					}
				}
			}
			if scan {
				valueToScan := hit.MatchValue
				if rule.UseFullLine {
					valueToScan = hit.LineValue
				}
				if rule.CompiledPattern.MatchString(valueToScan) {
					return true
				}
			}
		}
	}
	return false
}

// Translate the severity level from int value to string value
func getLevelNameFromID(level int, levelMap map[string]int) string {
	levelName := "low"
	for key, value := range levelMap {
		if value == level {
			levelName = key
		}
	}
	return levelName
}

func getZIPURL(url string) string {
	splits := strings.SplitAfter(url, ".zip")
	if len(splits) <= 2 {
		return splits[0]
	}
	return url
}

func getFileURL(giturl, filepath string) (fileurl string) {
	//strip .git at end if link if exists
	fileurl = strings.Replace(giturl, ".git", "", 1)
	fileurl = fileurl + "/blob/master/" + filepath

	//Strip file paths within zip file
	if strings.Contains(fileurl, ".zip") {
		fileurl = getZIPURL(fileurl)
	}

	//preappend filepath
	fileurl = filepath + ":" + fileurl
	return fileurl
}
