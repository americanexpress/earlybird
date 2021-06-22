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

package scan

import (
	"regexp"
	"testing"
)

func Test_findHit(t *testing.T) {
	type args struct {
		target          string
		CompiledPattern *regexp.Regexp
	}
	tests := []struct {
		name         string
		args         args
		wantIsHit    bool
		wantRetMatch string
	}{
		{
			name: "Check for hit in compress file name",
			args: args{
				target:          "compressed.zip",
				CompiledPattern: CompressPattern,
			},
			wantIsHit:    true,
			wantRetMatch: ".zip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsHit, gotRetMatch := findHit(tt.args.target, tt.args.CompiledPattern)
			if gotIsHit != tt.wantIsHit {
				t.Errorf("findHit() gotIsHit = %v, want %v", gotIsHit, tt.wantIsHit)
			}
			if gotRetMatch != tt.wantRetMatch {
				t.Errorf("findHit() gotRetMatch = %v, want %v", gotRetMatch, tt.wantRetMatch)
			}
		})
	}
}

func Test_substringExistsInLines(t *testing.T) {
	type args struct {
		fileLines []Line
		str       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Find secret word hideme in lines",
			args: args{
				fileLines: []Line{
					{LineValue: "Nothing to see here"},
					{LineValue: "Not sure what you're looking hideme for"},
					{LineValue: "This line shouldn't be searched"},
				},
				str: "hideme",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := substringExistsInLines(tt.args.fileLines, tt.args.str); got != tt.want {
				t.Errorf("substringExistsInLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_substringExistsInString(t *testing.T) {
	type args struct {
		str    string
		substr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Find foo in foo bar",
			args: args{
				str:    "foo bar",
				substr: "foo",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := substringExistsInString(tt.args.str, tt.args.substr); got != tt.want {
				t.Errorf("substringExistsInString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareMatchValue(t *testing.T) {
	tests := []struct {
		name              string
		matchValue        string
		wantIsolatedValue string
	}{
		{
			name:              "Prep matched value",
			matchValue:        `"test"`,
			wantIsolatedValue: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsolatedValue := prepareMatchValue(tt.matchValue); gotIsolatedValue != tt.wantIsolatedValue {
				t.Errorf("prepareMatchValue() = %v, want %v", gotIsolatedValue, tt.wantIsolatedValue)
			}
		})
	}
}

func Test_stringContainsQuotesTicks(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "Check string for quote ticks",
			s:    `"this is quoted"`,
			want: true,
		},
		{
			name: "Check if string doesn't have quote ticks",
			s:    `this is NOT quoted`,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringContainsQuotesTicks(tt.s); got != tt.want {
				t.Errorf("stringContainsQuotesTicks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shouldStrip(t *testing.T) {
	tests := []struct {
		name string
		s    byte
		want bool
	}{
		{
			name: "Check that , strippable",
			s:    ',',
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldStrip(tt.s); got != tt.want {
				t.Errorf("shouldStrip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLevelNameFromID(t *testing.T) {
	type args struct {
		level    int
		levelMap map[string]int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ID to Name",
			args: args{
				level: 3,
				levelMap: map[string]int{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			want: "three",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLevelNameFromID(tt.args.level, tt.args.levelMap); got != tt.want {
				t.Errorf("getLevelNameFromID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getZIPURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Get zip URL",
			url:  "https://github.com/test/sample.zip/ignorethis",
			want: "https://github.com/test/sample.zip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getZIPURL(tt.url); got != tt.want {
				t.Errorf("getZIPURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFileURL(t *testing.T) {
	type args struct {
		giturl   string
		filepath string
	}
	tests := []struct {
		name        string
		args        args
		wantFileurl string
	}{
		{
			name: "Get file link from github link and path",
			args: args{
				giturl:   "https://github.com/test",
				filepath: "test_data/sample.zip",
			},
			wantFileurl: "test_data/sample.zip:https://github.com/test/blob/master/test_data/sample.zip",
		},
		{
			name: "Get file link from bitbucket link and path",
			args: args{
				giturl:   "https://bitbucket.com/stash/scm/project/test",
				filepath: "test_data/sample.zip",
			},
			wantFileurl: "test_data/sample.zip:https://bitbucket.com/stash/projects/project/repos/test/browse/test_data/sample.zip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFileurl := getFileURL(tt.args.giturl, tt.args.filepath); gotFileurl != tt.wantFileurl {
				t.Errorf("getFileURL() = %v, want %v", gotFileurl, tt.wantFileurl)
			}
		})
	}
}
