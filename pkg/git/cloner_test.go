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

var FakeRepo = "https://github.com/carnal0wnage/fake_commited_secrets"

func TestReposPerProject(t *testing.T) {
	if os.Getenv("gituser") == "" && os.Getenv("gitpassword") == "" {
		t.Skip("Skipping ReposPerProject. Authentication needed. Include ENV vars: gituser, gitpassword")
	}

	if gotScanRepos := ReposPerProject("https://github.com/americanexpress", os.Getenv("gituser"), os.Getenv("gitpassword")); len(gotScanRepos) == 0 {
		t.Errorf("ReposPerProject() = %v, want multiple repository names", gotScanRepos)
	}
}

func TestCloneGitRepos(t *testing.T) {
	if os.Getenv("local") == "" {
		t.Skip("If test cases not running locally, skip cloning external repositories for CI/CD purposes.")
	}

	SearchDir, err := CloneGitRepos([]string{FakeRepo}, "", "", true)
	if err != nil {
		t.Errorf("Failed to clone repository: %s", FakeRepo)
	}

	//Delete temporary cloned repository directory
	err = os.RemoveAll(SearchDir)
	if err != nil {
		t.Errorf("Failed to delete git dir: %s", err)
	}
}
