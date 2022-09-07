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
	"bufio"
	"crypto/sha1"
	cfgReader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/utils"
	"path"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var cfg = cfgReader.EarlybirdConfig{
	ConfigDir:               path.Join(utils.MustGetWD(), utils.GetConfigDir()),
	LabelsConfigDir:         path.Join(utils.MustGetWD(), utils.GetConfigDir(), "labels"),
	FalsePositivesConfigDir: path.Join(utils.MustGetWD(), utils.GetConfigDir(), "falsepositives"),
	RulesConfigDir:          path.Join(utils.MustGetWD(), utils.GetConfigDir(), "rules"),
	IgnoreFile:              path.Join(utils.MustGetWD(), "../../.ge_ignore"),
	SeverityDisplayLevel:    4,
	SeverityFailLevel:       4,
	ConfidenceDisplayLevel:  4,
	ConfidenceFailLevel:     4,
	MaxFileSize:             10240000,
	WorkLength:              2500,
	EnabledModulesMap:       map[string]string{"content": "content.yaml", "password-secret": "password-secret.yaml"},
	LevelMap: map[string]int{
		"critical": 1,
		"high":     2,
		"info":     5,
		"low":      4,
		"medium":   3,
	},
	WorkerCount: 100,
}

func init() {
	err := cfgReader.LoadConfig(&cfgReader.Settings, path.Join(utils.MustGetWD(), utils.GetConfigDir(), "earlybird.json"))
	if err != nil {
		panic(err)
	}

	cfg.AdjustedSeverityCategories = []cfgReader.AdjustedSeverityCategory{
		{
			Category:                "password-secret",
			AdjustedDisplaySeverity: "medium",
			Patterns: []string{
				"(?i)/lowEnv/",
				"(?i)/test/",
				"(?i)/tests/",
				"(?i)/__tests__/",
				"(?i)lowEnv\\.(yaml|yml|properties|js|json)",
			},
			UseFilename: true,
		},
	}

	Init(cfg)
}

func TestScanFiles(t *testing.T) {
	files := []File{
		{
			Name: "file.py",
			Path: "buffer",
			Lines: []Line{
				{
					LineValue: `password = "SecretValue1673"`,
					LineNum:   1,
				},
			},
		},
	}
	hits := make(chan Hit)

	type args struct {
		cfg           *cfgReader.EarlybirdConfig
		files         []File
		compressPaths []string
		convertPaths  []string
		wantCode      int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				cfg:      &cfg,
				files:    files,
				wantCode: 3001,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SearchFiles(tt.args.cfg, tt.args.files, tt.args.compressPaths, tt.args.convertPaths, hits)
			for i := range hits {
				if i.Code != tt.args.wantCode {
					t.Errorf("ScanFiles() found code %v, want code %v", i.Code, tt.args.wantCode)
				}
			}
		})
	}
}

func Test_isExcludedFileType(t *testing.T) {
	cfg := cfgReader.EarlybirdConfig{
		ExtensionsToSkipScan: []string{".jpg"},
	}
	type args struct {
		cfg      *cfgReader.EarlybirdConfig
		filename string
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		{
			name: "skip image",
			args: args{
				cfg:      &cfg,
				filename: "photo.jpg",
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := isExcludedFileType(tt.args.cfg, tt.args.filename); gotOk != tt.wantOk {
				t.Errorf("isExcludedFileType() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_hitUnique(t *testing.T) {
	hit := Hit{
		Code:       3003,
		Line:       1,
		Filename:   "sample.py",
		MatchValue: "tomcat_password = '123'",
	}

	digest := sha1.New()
	_, err := digest.Write([]byte(hit.Filename + strconv.Itoa(hit.Line) + hit.MatchValue))
	if err != nil {
		t.Errorf("Failed to produce digest of hit: %v", err)
	}
	hithash := string(digest.Sum(nil))
	dupeMap := make(map[string]bool)
	(dupeMap)[hithash] = true

	type args struct {
		dupeMap map[string]bool
		hit     Hit
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check if hit is unique",
			args: args{
				hit:     hit,
				dupeMap: dupeMap,
			},
			want: false,
		},
		{
			name: "Check if hit is unique",
			args: args{
				hit: Hit{
					Code:       1111,
					Line:       1,
					Filename:   "different.py",
					MatchValue: "password = '123'",
				},
				dupeMap: dupeMap,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hitUnique(tt.args.dupeMap, tt.args.hit); got != tt.want {
				t.Errorf("hitUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanLine(t *testing.T) {
	var fileLines []Line
	type args struct {
		line      Line
		fileLines []Line
		rules     []Rule
	}
	tests := []struct {
		name      string
		args      args
		wantIsHit bool
	}{
		{
			name: "Find secret in line",
			args: args{
				line: Line{
					LineValue: "password=TrueFinding7842!",
				},
				fileLines: fileLines,
			},
			wantIsHit: true,
		},
		{
			name: "Find fall as a password",
			args: args{
				line: Line{
					LineValue: "password = fall123",
				},
				fileLines: fileLines,
			},
			wantIsHit: true,
		},
		{
			name: "Ignore fall in a sentence",
			args: args{
				line: Line{
					LineValue: "using fall in a general sentence should not error",
				},
				fileLines: fileLines,
			},
			wantIsHit: false,
		},
		{
			name: "Find twitter API key as a password",
			args: args{
				line: Line{
					LineValue: `twitterApiSecret:"111aAa222bBb333cCc444dDd555eEe666fFf777"`,
				},
				fileLines: fileLines,
			},
			wantIsHit: true,
		},
		{
			name: "Ignore potential twitter API key separated by too many characters",
			args: args{
				line: Line{
					LineValue: `twitter="twitter";//This LineValue emulates extremely long one-liner code files that can cause false positives "111aAa222bBb333cCc444dDd555eEe666fFf777"`,
				},
				fileLines: fileLines,
			},
			wantIsHit: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsHit, _ := scanLine(tt.args.line, tt.args.fileLines, &cfg)
			if gotIsHit != tt.wantIsHit {
				t.Errorf("scanLine() gotIsHit = %v, want %v", gotIsHit, tt.wantIsHit)
			}
		})
	}
}

func Test_scanName(t *testing.T) {
	type args struct {
		file  File
		rules []Rule
	}
	tests := []struct {
		name      string
		args      args
		wantIsHit bool
		wantHit   Hit
	}{
		{
			name: "Check if keyfile matches",
			args: args{
				rules: loadRuleConfigs(cfg, "filename", "filename.yaml"),
				file: File{
					Name: "findme.pem",
					Path: "/bad/findme.pem",
				},
			},
			wantIsHit: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsHit, gotHit := scanName(tt.args.file, tt.args.rules, &cfg)

			if gotIsHit != tt.wantIsHit {
				t.Errorf("scanName() gotIsHit = %v, want %v", gotIsHit, tt.wantIsHit)
			}

			if gotHit.Code == 0 {
				t.Errorf("scanName() gotIsHit = %v, want non empty hit", gotHit)
			}
		})
	}
}

func Test_readln(t *testing.T) {
	w := strings.NewReader("test\n")
	rbuf := bufio.NewReader(w)
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test line reader",
			args: args{
				r: rbuf,
			},
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readln(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readln() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readln() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIgnoreAnnotation(t *testing.T) {
	cfg := cfgReader.EarlybirdConfig{
		AnnotationsToSkipLine: []string{"#SKIPME"},
	}
	type args struct {
		cfg  *cfgReader.EarlybirdConfig
		line string
	}
	tests := []struct {
		name    string
		args    args
		wantRet bool
	}{
		{
			name: "Skip annotated line",
			args: args{
				cfg:  &cfg,
				line: "This is a false positive #SKIPME",
			},
			wantRet: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := IsIgnoreAnnotation(tt.args.cfg, tt.args.line); gotRet != tt.wantRet {
				t.Errorf("IsIgnoreAnnotation() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func Test_maskValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Supress value",
			input: "foobar",
			want:  "******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskValue(tt.input); got != tt.want {
				t.Errorf("maskValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitJob(t *testing.T) {
	type args struct {
		job        WorkJob
		worklength int
	}
	tests := []struct {
		name     string
		args     args
		wantWork []WorkJob
	}{
		{
			name: "Split up work",
			args: args{
				job: WorkJob{
					WorkLine: Line{
						LineValue: "foobar",
					},
				},
				worklength: 3,
			},
			wantWork: []WorkJob{
				{
					WorkLine: Line{
						LineValue: "foo",
					},
				},
				{
					WorkLine: Line{
						LineValue: "bar",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotjobs := splitJob(tt.args.job, tt.args.worklength); !reflect.DeepEqual(gotjobs, tt.wantWork) {
				t.Errorf("splitJob() = %v, want %v", gotjobs, tt.wantWork)
			}
		})
	}
}

func Test_splitSubN(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Split foo bar",
			args: args{
				s: "foobar",
				n: 3,
			},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitSubN(tt.args.s, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitSubN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_labelHit(t *testing.T) {
	hit := &Hit{
		Code:       3003,
		Line:       1,
		Filename:   "sample.py",
		MatchValue: "tomcat_password = '123'",
	}
	labelHit(hit, []Line{})
	if len(hit.Labels) != 0 {
		t.Errorf("LabelHit() = [], failed to label hit")
	}
}

func Test_determineSeverity(t *testing.T) {
	rule := Rule{
		Category: "password-secret",
		Code:     3001,
		Severity: 2,
	}

	tests := []struct {
		name               string
		hit                *Hit
		expectedSeverity   string
		expectedSeverityId int
	}{
		{
			name: "it should return normal severity findings",
			hit: &Hit{
				Code:       3001,
				Line:       1,
				Filename:   "bar/foo.js",
				MatchValue: "password = 'aReallyBadPassword'",
			},
			expectedSeverity:   "high",
			expectedSeverityId: 2,
		},
		{
			name: "it should reduce severity findings based on test dir",
			hit: &Hit{
				Code:       3001,
				Line:       1,
				Filename:   "a/b/__tests__/foo.js",
				MatchValue: "password = 'aReallyBadPassword'",
			},
			expectedSeverity:   "medium",
			expectedSeverityId: 3,
		},
		{
			name: "it should reduce severity findings based on lowEnv dir",
			hit: &Hit{
				Code:       3001,
				Line:       1,
				Filename:   "root/of/repo/lowEnv/foo.js",
				MatchValue: "password = 'aReallyBadPassword'",
			},
			expectedSeverity:   "medium",
			expectedSeverityId: 3,
		},
		{
			name: "it should reduce severity findings for files with lowEnv in name",
			hit: &Hit{
				Code:       3001,
				Line:       1,
				Filename:   "bar/config-lowEnv.js",
				MatchValue: "password = 'aReallyBadPassword'",
			},
			expectedSeverity:   "medium",
			expectedSeverityId: 3,
		},
		{
			name: "it should reduce severity findings for files with lowEnv in name #2",
			hit: &Hit{
				Code:       3001,
				Line:       1,
				Filename:   "bar/latest/lowEnv.js",
				MatchValue: "password = 'aReallyBadPassword'",
			},
			expectedSeverity:   "medium",
			expectedSeverityId: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hit.determineSeverity(&cfg, &rule)

			if tt.hit.SeverityID != tt.expectedSeverityId && tt.hit.Severity != tt.expectedSeverity {
				t.Errorf("hit.Severity = %s, want %s, hit.SeverityId = %d, want %d", tt.hit.Severity, tt.expectedSeverity, tt.hit.SeverityID, tt.expectedSeverityId)
			}
		})
	}
}

func Test_determineScanFail(t *testing.T) {
	levelMap := map[string]int{
		"critical": 1,
		"high":     2,
		"info":     5,
		"low":      4,
		"medium":   3,
	}
	tests := []struct {
		name           string
		hit            *Hit
		cfg            *cfgReader.EarlybirdConfig
		shouldFailScan bool
	}{
		{
			name: "hit severity < fail severity and hit confidence < fail confidence",
			hit: &Hit{
				SeverityID:   1, // critical
				ConfidenceID: 1, // critical
			},
			cfg: &cfgReader.EarlybirdConfig{
				SeverityFailLevel:   2, // high
				ConfidenceFailLevel: 2, // high
				LevelMap:            levelMap,
			},
			shouldFailScan: true,
		},
		{
			name: "hit severity == fail severity and hit confidence == fail confidence",
			hit: &Hit{
				SeverityID:   2, // high
				ConfidenceID: 2, // high
			},
			cfg: &cfgReader.EarlybirdConfig{
				SeverityFailLevel:   2, // high
				ConfidenceFailLevel: 2, // high
				LevelMap:            levelMap,
			},
			shouldFailScan: true,
		},
		{
			name: "hit severity < fail severity and hit confidence > fail confidence",
			hit: &Hit{
				SeverityID:   1, // high
				ConfidenceID: 3, // medium
			},
			cfg: &cfgReader.EarlybirdConfig{
				SeverityFailLevel:   2, // high
				ConfidenceFailLevel: 2, // high
				LevelMap:            levelMap,
			},
			shouldFailScan: false,
		},
		{
			name: "hit severity > fail severity and hit confidence < fail confidence",
			hit: &Hit{
				SeverityID:   3, // medium
				ConfidenceID: 1, // critical
			},
			cfg: &cfgReader.EarlybirdConfig{
				SeverityFailLevel:   2, // high
				ConfidenceFailLevel: 2, // high
				LevelMap:            levelMap,
			},
			shouldFailScan: false,
		},
		{
			name: "hit severity > fail severity and hit confidence > fail confidence",
			hit: &Hit{
				SeverityID:   3, // medium
				ConfidenceID: 3, // medium
			},
			cfg: &cfgReader.EarlybirdConfig{
				SeverityFailLevel:   2, // high
				ConfidenceFailLevel: 2, // high
				LevelMap:            levelMap,
			},
			shouldFailScan: false,
		},
		{
			name: "hit severity is info level",
			hit: &Hit{
				SeverityID:   5, // info
				ConfidenceID: 1, // critical
			},
			cfg: &cfgReader.EarlybirdConfig{
				SeverityFailLevel:   2, // high
				ConfidenceFailLevel: 2, // high
				LevelMap:            levelMap,
			},
			shouldFailScan: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineScanFail(tt.cfg, tt.hit); got != tt.shouldFailScan {
				t.Errorf("determineScanFail() = %t, want %t", got, tt.shouldFailScan)
			}
		})
	}
}
