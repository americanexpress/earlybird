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
	"io/ioutil"
	"os"
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		haystack []string
		needle   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check for bar in foo bar",
			args: args{
				haystack: []string{"foo", "bar"},
				needle:   "bar",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathMustExist(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current directory")
	}

	PathMustExist(currentDir) //Will exit program causing an error if invalid path
}

func TestGetConfigDir(t *testing.T) {
	gotConfigDir := GetConfigDir()
	if gotConfigDir == "" {
		t.Errorf("GetConfigDir() = %v, want non nil value", gotConfigDir)
	}

	exists, _ := Exists(gotConfigDir)
	if !exists {
		t.Errorf("GetConfigDir() = %v does not exist", gotConfigDir)
	}
}

func TestMustGetED(t *testing.T) {
	gotED := MustGetED()
	if gotED == "" {
		t.Errorf("GetED() = %v, want non nil value", gotED)
	}

	exists, _ := Exists(gotED)
	if !exists {
		t.Errorf("MustGetED()= %v does not exist", gotED)
	}
}

func TestMustGetWD(t *testing.T) {
	gotWD := MustGetWD()
	if gotWD == "" {
		t.Errorf("GetWD() = %v, want non nil value", gotWD)
	}

	exists, _ := Exists(gotWD)
	if !exists {
		t.Errorf("GetWD()= %v does not exist", gotWD)
	}
}

func TestGetDisplayList(t *testing.T) {
	type args struct {
		levelNames []string
	}
	tests := []struct {
		name             string
		args             args
		wantLevelOptions string
	}{
		{
			name: "Format level map into string",
			args: args{
				levelNames: []string{"low", "medium", "high"},
			},
			wantLevelOptions: "[ low | medium | high ]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLevelOptions := GetDisplayList(tt.args.levelNames); gotLevelOptions != tt.wantLevelOptions {
				t.Errorf("GetDisplayListy() = %v, want %v", gotLevelOptions, tt.wantLevelOptions)
			}
		})
	}
}

func TestDeleteGit(t *testing.T) {
	type args struct {
		ptrRepo string
		path    string
	}

	testRepo := "test-repo"
	tmpDir, err := ioutil.TempDir("", "ebgit")
	if err != nil {
		t.Errorf("Failed to create temporoary directory: %v", err)
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Delete tmp directory",
			args: args{
				ptrRepo: testRepo,
				path:    tmpDir,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteGit(tt.args.ptrRepo, tt.args.path)
		})
	}
}

func TestGetGitRepo(t *testing.T) {
	tests := []struct {
		name           string
		gitURL         string
		wantRepository string
	}{
		{
			name:           "Parse repository name from git repo",
			gitURL:         "https://github.com/americanexpress/earlybird",
			wantRepository: "americanexpress/earlybird",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRepository := GetGitRepo(tt.gitURL); gotRepository != tt.wantRepository {
				t.Errorf("GetGitRepo() = %v, want %v", gotRepository, tt.wantRepository)
			}
		})
	}
}

func TestGetBBProject(t *testing.T) {
	tests := []struct {
		name        string
		bbURL       string
		wantProject string
	}{
		{
			name:        "Parse project from bitbucket URL",
			bbURL:       "https://example.com/stash/projects/TEST123/repos/test-repo/browse",
			wantProject: "TEST123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotProject := GetBBProject(tt.bbURL); gotProject != tt.wantProject {
				t.Errorf("GetBBProject() = %v, want %v", gotProject, tt.wantProject)
			}
		})
	}
}

func TestParseBBURL(t *testing.T) {
	tests := []struct {
		name        string
		bbURL       string
		wantBaseurl string
		wantPath    string
		wantProject string
	}{
		{
			name:        "Parse bitbucket URL format into separate values",
			bbURL:       "https://example.com/stash/projects/TEST123/repos/test-repo/browse",
			wantBaseurl: "https://example.com",
			wantPath:    "/stash",
			wantProject: "TEST123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBaseurl, gotPath, gotProject := ParseBBURL(tt.bbURL)
			if gotBaseurl != tt.wantBaseurl {
				t.Errorf("ParseBBURL() gotBaseurl = %v, want %v", gotBaseurl, tt.wantBaseurl)
			}
			if gotPath != tt.wantPath {
				t.Errorf("ParseBBURL() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotProject != tt.wantProject {
				t.Errorf("ParseBBURL() gotProject = %v, want %v", gotProject, tt.wantProject)
			}
		})
	}
}

func TestGetGitProject(t *testing.T) {
	testURL := "https://example.com/test-project"
	tests := []struct {
		name        string
		gitURL      string
		wantProject string
	}{
		{
			name:        "Parse project from git URL",
			gitURL:      testURL,
			wantProject: "test-project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotProject := GetGitProject(tt.gitURL); gotProject != tt.wantProject {
				t.Errorf("GetGitProject() = %v, want %v", gotProject, tt.wantProject)
			}
		})
	}
}

func TestGetGitURL(t *testing.T) {
	testURL := "https://example.com/repository"
	var emptyUser string
	type args struct {
		ptrRepo     *string
		ptrRepoUser *string
	}
	tests := []struct {
		name         string
		args         args
		wantPassword string
	}{
		{
			name: "Format git url and return blank password",
			args: args{
				ptrRepo:     &testURL,
				ptrRepoUser: &emptyUser,
			},
			wantPassword: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPassword := GetGitURL(tt.args.ptrRepo, tt.args.ptrRepoUser); gotPassword != tt.wantPassword {
				t.Errorf("GetGitURL() = %v, want %v", gotPassword, tt.wantPassword)
			}
		})
	}
}

func TestExists(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current directory")
	}
	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{
			name:    "Check if current directory exists",
			path:    currentDir,
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exists(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAlphaNumericValues(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "Remove Underscore from string",
			value: "COUCHBASE_PASSWORD",
			want:  "COUCHBASEPASSWORD",
		},
		{
			name:  "Remove all special chars from the string",
			value: "MONGO$##@@)(!~_PASSWORD##$$%%_123",
			want:  "MONGOPASSWORD123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAlphaNumericValues(tt.value)
			if got != tt.want {
				t.Errorf("GetAlphaNumericValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
