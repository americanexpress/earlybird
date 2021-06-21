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

package utils

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/howeyc/gopass"
)

var gitURLPattern = regexp.MustCompile(`([^/]*)(?:.git)`)

// Contains Does an array/slice contain a string
func Contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// PathMustExist exit if path is invalid
func PathMustExist(path string) {
	if fileExists, err := Exists(path); !fileExists {
		if err != nil {
			log.Fatal(errInvalidPath)
		}
	}
}

// GetConfigDir Determine the operating system and pull the path to the go-earlybird config directory
func GetConfigDir() (configDir string) {
	if strings.HasSuffix(os.Args[0], ".test") { // Return repository config directory when testing ../../config
		return ".." + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "config" + string(os.PathSeparator)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Home directory doesn't exist", err)
	}
	cwd := MustGetED()

	overrideDir := string(os.PathSeparator) + ebConfFileDir + string(os.PathSeparator)
	localOverrideDir := cwd + overrideDir
	localOverrideFileCheck := localOverrideDir + ebConfFileName
	if fe, _ := Exists(localOverrideFileCheck); fe {
		configDir = localOverrideDir
		log.Println("Using local config directory: ", localOverrideDir)
	} else {
		switch runtime.GOOS {
		case "windows":
			configDir += ebWinConfFileDir
		case "linux": // also can be specified to FreeBSD
			configDir += overrideDir
		case "darwin":
			configDir += overrideDir
		}
		configDir = userHomeDir + configDir
	}
	return configDir
}

// MustGetED Get the executable directory or exit
func MustGetED() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(ex)
}

// MustGetWD Get the CWD for the default target Directory or exit
func MustGetWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

// GetTargetType returns the file scan context
func GetTargetType(GitStagedFlag, GitTrackedFlag bool) (targetType string) {
	if GitStagedFlag {
		targetType = Staged
	} else if GitTrackedFlag {
		targetType = Tracked
	} else {
		targetType = All
	}
	return targetType
}

// GetEnabledModulesMap returns a map of module name to filename enabled by default or explicitly defined with CLI paramters
func GetEnabledModulesMap(enableFlags []string, availableModules map[string]string) (enabledModules map[string]string) {
	enabledModules = make(map[string]string)

	if len(enableFlags) == 0 {
		return availableModules
	}

	for _, moduleName := range enableFlags {
		enabledModules[moduleName] = availableModules[moduleName]
	}
	return enabledModules
}

// GetDisplayList Build the string to display an array in a human readable format
func GetDisplayList(levelNames []string) string {
	return "[ " + strings.Join(levelNames, " | ") + " ]"
}

// DeleteGit Check if we've cloned a git repo, if so delete it
func DeleteGit(ptrRepo string, path string) {
	if ptrRepo != "" {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println(errGitDelete, err)
		}
	}
}

// GetGitRepo Parse repository name from URL
func GetGitRepo(gitURL string) (repository string) {
	if strings.Contains(gitURL, "github.com/") {
		u, err := url.Parse(gitURL)
		if err != nil {
			return
		}
		repository = strings.TrimPrefix(u.Path, "/")
	} else {
		items := gitURLPattern.FindStringSubmatch(gitURL)
		if len(items) > 1 {
			repository = items[1]
		}
	}
	return repository
}

// GetBBProject Parse project name from bitbucket URL
func GetBBProject(bbURL string) (project string) {
	re := regexp.MustCompile(`(?:projects/)([^/]*)`)
	results := re.FindStringSubmatch(bbURL) // Match second capture group, 1 = project/XXX, 2 = XXX
	if len(results) < 1 {
		log.Println("Failed To Get BB Project from URL:", bbURL)
		os.Exit(1)
	} else {
		project = results[1]
	}
	return project
}

// ParseBBURL Parse the base URL, Path and project name from BB URL
func ParseBBURL(bbURL string) (baseurl, path, project string) {
	u, err := url.Parse(bbURL)
	if err != nil {
		log.Println("Failed to parse Bitbucket URL")
		return
	}
	baseurl = u.Scheme + "://" + u.Host               // Parse Base URL
	parts := strings.Split(bbURL, "/projects/")       // Get URL before /projects
	path = strings.Replace(parts[0], baseurl, "", -1) // Delete the base url leaving the path
	return baseurl, path, GetBBProject(bbURL)
}

// GetGitProject parse project from URL
func GetGitProject(gitURL string) (project string) {
	u, err := url.Parse(gitURL)
	if err != nil {
		return
	}
	return strings.TrimPrefix(u.Path, "/")
}

// GetGitURL Format GI URL and parse/prompt user password
func GetGitURL(ptrRepo, ptrRepoUser *string) (Password string) {
	// Parse Username from URL
	u, err := url.Parse(*ptrRepo)
	if err != nil {
		return
	}

	// Remove Username prefix
	*ptrRepo = strings.Replace(*ptrRepo, u.User.Username()+"@", "", 1)
	if *ptrRepoUser == "" {
		*ptrRepoUser = u.User.Username()
	}

	// Remove HTTP and HTTPS prefix
	*ptrRepo = strings.Replace(*ptrRepo, gitHTTP, "", 1)
	*ptrRepo = strings.Replace(*ptrRepo, gitHTTPS, "", 1)
	*ptrRepo = gitHTTPS + *ptrRepo
	if *ptrRepoUser != "" {
		log.Print(gitPasswdPrompt)
		RepoPass, err := gopass.GetPasswdMasked()
		if err != nil {
			log.Println(errGitPasswd, err)
			os.Exit(1)
		}
		// Format git URL with user and password
		return string(RepoPass)
	}
	return ""
}

// Exists Check to see if a path exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

// GetAlphaNumericValues returns the alphanumeric part of the input string
func GetAlphaNumericValues(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, input)
}
