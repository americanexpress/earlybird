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
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/howeyc/gopass"
)

//Contains Does an array/slice contain a string
func Contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

//PathMustExist exit if path is invalid
func PathMustExist(path string) {
	if fileExists, err := Exists(path); !fileExists {
		if err != nil {
			log.Fatal(errInvalidPath)
		}
	}
}

//GetConfigDir Determine the operating system and pull the path to the go-earlybird config directory
func GetConfigDir() (configDir string) {
	if strings.HasSuffix(os.Args[0], ".test") { //Return repository config directory when testing ../../config
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
		fmt.Println("Using local config directory: ", localOverrideDir)
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

//MustGetED Get the executable directory or exit
func MustGetED() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(ex)
}

//MustGetWD Get the CWD for the default target Directory or exit
func MustGetWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

//GetTargetType returns the file scan context
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

//GetEnabledModules returns a list of modules enabled by default or explicitly defined with CLI paramters
func GetEnabledModules(enableFlags, availableModules []string) (enabledModules []string) {
	if len(enableFlags) == 0 {
		return availableModules
	}

	for _, moduleName := range enableFlags {
		if Contains(availableModules, moduleName) {
			enabledModules = append(enabledModules, moduleName)
		}
	}
	return enabledModules
}

//GetDisplayList Build the string to display an array in a human readable format
func GetDisplayList(levelNames []string) string {
	return "[ " + strings.Join(levelNames, " | ") + " ]"
}

//DeleteGit Check if we've cloned a git repo, if so delete it
func DeleteGit(ptrRepo string, path string) {
	if ptrRepo != "" {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println(errGitDelete, err)
		}
	}
}

//GetGitRepo Parse repository name from URL
func GetGitRepo(gitURL string) (repository string) {
	u, err := url.Parse(gitURL)
	if err != nil {
		return
	}
	repository = strings.TrimPrefix(u.Path, "/")
	//Trim suffix .git if exists
	repository = strings.TrimSuffix(repository, filepath.Ext(repository))
	return repository
}

//GetGitProject parse project from URL
func GetGitProject(gitURL string) (project string) {
	u, err := url.Parse(gitURL)
	if err != nil {
		return
	}
	return strings.TrimPrefix(u.Path, "/")
}

//GetGitURL Format GI URL and parse/prompt user password
func GetGitURL(ptrRepo, ptrRepoUser *string) (Password string) {
	//Parse Username from URL
	u, err := url.Parse(*ptrRepo)
	if err != nil {
		return
	}

	//Remove Username prefix
	*ptrRepo = strings.Replace(*ptrRepo, u.User.Username()+"@", "", 1)
	if *ptrRepoUser == "" {
		*ptrRepoUser = u.User.Username()
	}

	//Remove HTTP and HTTPS prefix
	*ptrRepo = strings.Replace(*ptrRepo, gitHTTP, "", 1)
	*ptrRepo = strings.Replace(*ptrRepo, gitHTTPS, "", 1)
	*ptrRepo = gitHTTPS + *ptrRepo
	if *ptrRepoUser != "" {
		fmt.Print(gitPasswdPrompt)
		RepoPass, err := gopass.GetPasswdMasked()
		if err != nil {
			fmt.Println(errGitPasswd, err)
			os.Exit(1)
		}
		//Format git URL with user and password
		return string(RepoPass)
	}
	return ""
}

//Exists Check to see if a path exists
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
