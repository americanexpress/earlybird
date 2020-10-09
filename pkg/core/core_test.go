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
	"os"
	"testing"

	"github.com/americanexpress/earlybird/pkg/scan"
	"github.com/americanexpress/earlybird/pkg/utils"
)

var eb EarlybirdCfg

func init() {
	scan.Init(eb.Config)
}

//Program will exit with error if config init fails
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

	giturl := os.Getenv("giturl")
	if giturl == "" {
		t.Skip("Skipping GitClone. Git repository URL is required. Include ENV var: giturl")
	}

	var (
		RepoUser string
		Project  string
	)
	ptr := PTRGitConfig{
		Repo:     &giturl,
		RepoUser: &RepoUser,
		Project:  &Project,
	}

	eb.GitClone(ptr)

	//Delete temporary cloned repository directory
	utils.DeleteGit(giturl, eb.Config.SearchDir)
}
