/*
 * Copyright 2023 American Express
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

package writers

import (
	"testing"

	"github.com/americanexpress/earlybird/pkg/scan"
)

func TestWriteConsole(t *testing.T) {
	HitChan := make(chan scan.Hit)
	go func() {
		HitChan <- scan.Hit{
			Code:       3003,
			Line:       1,
			Filename:   "sample.py",
			MatchValue: "tomcat_password = '123'",
		}
	}()
	type args struct {
		hits         chan scan.Hit
		fileName     string
		showFullLine bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Write a hit to console",
			args: args{
				hits:         HitChan,
				showFullLine: true,
			},
			wantErr: false,
		},
	}
	for _, myTest := range tests {
		tt := myTest
		go func() {
			t.Run(tt.name, func(t *testing.T) {
				if err := WriteConsole(tt.args.hits, tt.args.fileName, tt.args.showFullLine); (err != nil) != tt.wantErr {
					t.Errorf("WriteConsole() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}()
	}
}

func Test_displayCWE(t *testing.T) {
	type args struct {
		cwe []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "Format CWE output",
			args: args{
				cwe: []string{"CWE-1", "CWE-2", "CWE3"},
			},
			wantResult: "CWE-1/CWE-2/CWE3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := displayCWE(tt.args.cwe); gotResult != tt.wantResult {
				t.Errorf("displayCWE() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_printableASCII(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Strip non ASCII bytes from string",
			args: args{
				str: `Hello中国!`,
			},
			want: "Hello!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := printableASCII(tt.args.str); got != tt.want {
				t.Errorf("printableASCII() = %v, want %v", got, tt.want)
			}
		})
	}
}
