/*
 * Copyright 2021 American Express
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
	"testing"
)

func TestBBCommitAPIURL(t *testing.T) {
	type args struct {
		host       string
		source     string
		project    string
		repository string
		commitid   string
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
	}{
		{
			name: "format bitbucket API URL",
			args: args{
				host:       "example.com",
				source:     "stash",
				project:    "test-project",
				repository: "earlybird",
				commitid:   "123456789",
			},
			wantURL: "https://example.com/stash/rest/api/latest/projects/test-project/repos/earlybird/commits/123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotURL := BBCommitAPIURL(tt.args.host, tt.args.source, tt.args.project, tt.args.repository, tt.args.commitid); gotURL != tt.wantURL {
				t.Errorf("BBCommitAPIURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func TestBBCommitURL(t *testing.T) {
	type args struct {
		host       string
		source     string
		project    string
		repository string
		commitid   string
	}
	tests := []struct {
		name    string
		args    args
		wantURL string
	}{
		{
			name: "format bitbucket commit URL",
			args: args{
				host:       "example.com",
				source:     "stash",
				project:    "test-project",
				repository: "earlybird",
				commitid:   "123456789",
			},
			wantURL: "https://example.com/stash/projects/test-project/repos/earlybird/commits/123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotURL := BBCommitURL(tt.args.host, tt.args.source, tt.args.project, tt.args.repository, tt.args.commitid); gotURL != tt.wantURL {
				t.Errorf("BBCommitURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}
