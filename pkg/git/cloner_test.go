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

package git

import (
	"os"
	"testing"
)

func TestReposPerProject(t *testing.T) {
	if os.Getenv("local") == "" {
		t.Skip("If test cases not running locally, skip project API call for CI/CD purposes.")
	}

	if os.Getenv("projecturl") == "" || os.Getenv("gituser") == "" || os.Getenv("gitpassword") == "" {
		t.Skip("Skipping ReposPerProject. Authentication needed. Include ENV vars: projecturl, gituser, gitpassword")
	}

	if gotScanRepos := ReposPerProject(os.Getenv("projecturl"), os.Getenv("gituser"), os.Getenv("gitpassword")); len(gotScanRepos) == 0 {
		t.Errorf("ReposPerProject() = %v, want multiple repository names", gotScanRepos)
	}
}

func TestCloneGitRepos(t *testing.T) {
	if os.Getenv("local") == "" {
		t.Skip("If test cases not running locally, skip cloning external repositories for CI/CD purposes.")
	}

	giturl := os.Getenv("giturl")
	if giturl == "" {
		t.Skip("Skipping CloneGitRepo. Git repository URL is required. Include ENV var: giturl")
	}

	SearchDir, err := CloneGitRepos([]string{giturl}, "", "", true)
	if err != nil {
		t.Errorf("Failed to clone repository: %s", giturl)
	}

	//Delete temporary cloned repository directory
	err = os.RemoveAll(SearchDir)
	if err != nil {
		t.Errorf("Failed to delete git dir: %s", err)
	}
}
