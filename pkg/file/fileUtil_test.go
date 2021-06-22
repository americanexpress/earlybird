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

package file

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/americanexpress/earlybird/pkg/scan"
)

var userHomeDir string

func init() {
	var err error
	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatal("Home directory doesn't exist", err)
	}
	ignorePatterns = getIgnorePatterns(userHomeDir+string(os.PathSeparator), ".ge_ignore", false)
}

func TestGetFiles(t *testing.T) {
	searchDir := "test_data"
	ignoreFile := userHomeDir + string(os.PathSeparator) + ".ge_ignore"
	verbose := false
	maxFileSize := int64(1000000)

	fileContext, err := GetFiles(searchDir, ignoreFile, verbose, maxFileSize)
	if err != nil {
		t.Errorf("GetFiles() err = %v", err)
	}
	if len(fileContext.Files) == 0 {
		t.Errorf("GetFiles() found none, expected atleast one file")
	} else if fileContext.Files[0].Name == "test_data/sample.zip/sample.py" {
		t.Errorf("GetFiles() first file doesn't match example file")
	}

	if len(fileContext.IgnorePatterns) == 0 {
		t.Errorf("GetFiles() IgnorePatterns, got none")
	}
}

func TestGetFileSize(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantSize int64
	}{
		{
			name: "Check file size",
			args: args{
				path: "test_data/sample.zip",
			},
			wantSize: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSize, err := GetFileSize(tt.args.path)
			if gotSize < tt.wantSize { //Want value bigger than 100
				t.Errorf("GetFileSize() = %v, want %v", gotSize, tt.wantSize)
			}
			if err != nil {
				t.Errorf("GetFileSize() error = %v", err)
			}
		})
	}
}

func Test_getFileSizeOK(t *testing.T) {
	type args struct {
		path        string
		maxFileSize int64
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "Check if file smaller than max",
			args: args{
				path:        "test_data/sample.zip",
				maxFileSize: 100000,
			},
			wantResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := getFileSizeOK(tt.args.path, tt.args.maxFileSize); gotResult != tt.wantResult {
				t.Errorf("getFileSizeOK() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_getIgnorePatterns(t *testing.T) {
	if gotIgnorePatterns := getIgnorePatterns(userHomeDir+string(os.PathSeparator), ".ge_ignore", false); len(ignorePatterns) == 0 {
		t.Errorf("getIgnorePatterns() = %v, want multiple patterns", gotIgnorePatterns)
	}
}

func Test_isIgnoredFile(t *testing.T) {
	if os.Getenv("local") == "" {
		t.Skip("If test cases not running locally, skip cloning external repositories for CI/CD purposes.")
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check if file is ignored",
			args: args{
				fileName: "ignore.png",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIgnoredFile(tt.args.fileName); got != tt.want {
				t.Errorf("isIgnoredFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDirectory(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Check if path is a directory",
			args: args{
				path: "test_data",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isDirectory(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("isDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWD(t *testing.T) {
	got, err := GetWD()
	if err != nil {
		t.Errorf("GetWD() error = %v, wantErr nil", err)
		return
	}
	if !strings.Contains(got, "pkg/file") {
		t.Errorf("GetWD() = %v, want pkg/file directory", got)
	}

}

func TestIsEmpty(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Check if path is empty",
			args: args{
				path: "test_data",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsEmpty(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsEmpty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Check if path exists",
			args: args{
				path: "test_data",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Check if path exists",
			args: args{
				path: "dont_exist",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exists(tt.args.path)
			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCompressedFiles(t *testing.T) {
	type args struct {
		files []scan.File
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if path exists",
			args: args{
				files: []scan.File{
					{
						Path: "test_data/sample.zip",
						Name: "sample.zip",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewfiles, gotCompresspaths, err := GetCompressedFiles(tt.args.files)
			if err != nil {
				t.Errorf("GetCompressedFiles() err = %v", err)
			}
			if len(gotNewfiles) == 0 {
				t.Errorf("GetCompressedFiles() gotNewfiles = %v, want multiple files", gotNewfiles)
			}
			if len(gotCompresspaths) == 0 {
				t.Errorf("GetCompressedFiles() gotCompresspaths = %v, want multiple paths", gotCompresspaths)
			}
		})
	}
}

func TestUncompress(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name          string
		args          args
		wantFilenames []string
		wantErr       bool
	}{
		{
			name: "Check if path exists",
			args: args{
				src:  "test_data/sample.zip",
				dest: "test_data",
			},
			wantFilenames: []string{"test_data/sample.py"},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilenames, err := Uncompress(tt.args.src, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uncompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFilenames, tt.wantFilenames) {
				t.Errorf("Uncompress() = %v, want %v", gotFilenames, tt.wantFilenames)
			}
		})
	}
	//Delete left over file
	os.Remove("test_data/sample.py")
}
