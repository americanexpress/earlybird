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
	"bufio"
	"crypto/sha1"
	"reflect"
	"strconv"
	"strings"
	"testing"

	cfgReader "github.com/americanexpress/earlybird/pkg/config"
	"github.com/americanexpress/earlybird/pkg/utils"
)

var cfg = cfgReader.EarlybirdConfig{
	ConfigDir:              utils.GetConfigDir(),
	IgnoreFile:             utils.GetConfigDir() + ".ge_ignore",
	SeverityDisplayLevel:   4,
	SeverityFailLevel:      4,
	ConfidenceDisplayLevel: 4,
	ConfidenceFailLevel:    4,
	MaxFileSize:            10240000,
	WorkLength:             2500,
	EnabledModules:         []string{"ccnumber", "common", "content", "filename", "entropy"},
}

func init() {
	Init(cfg)
}

func TestScanFiles(t *testing.T) {
	CombinedRules = loadRuleConfigs(2, 4, utils.GetConfigDir()+"content.json")
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
			go SearchFiles(tt.args.cfg, tt.args.files, tt.args.compressPaths, hits)
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
		MatchValue: "password = '123'",
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
	CombinedRules = loadRuleConfigs(4, 4, utils.GetConfigDir()+"content.json")
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
				rules: loadRuleConfigs(2, 2, utils.GetConfigDir()+"filename.json"),
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
			gotIsHit, gotHit := scanName(tt.args.file, tt.args.rules, cfg.LevelMap, cfg.ShowSolutions)
			if !gotIsHit {
				t.Errorf("scanName() gotIsHit = %v, want true", gotIsHit)
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
		MatchValue: "test_password = '123'",
	}
	labelHit(hit, []Line{})
	if len(hit.Labels) != 0 {
		t.Errorf("LabelHit() = [], failed to label hit")
	}
}
