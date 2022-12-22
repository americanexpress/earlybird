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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/americanexpress/earlybird/pkg/scan"
)

func TestWriteJSON(t *testing.T) {
	start := time.Now()
	Hits := []scan.Hit{
		{
			Code:       3003,
			Line:       1,
			Filename:   "sample.py",
			MatchValue: "tomcat_password = '123'",
		},
	}
	type args struct {
		report   scan.Report
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Write Hits to JSON",
			args: args{
				report: scan.Report{
					Hits:      Hits,
					HitCount:  len(Hits),
					Version:   "Test 1.0",
					StartTime: start.UTC().Format(time.RFC3339),
					EndTime:   time.Now().UTC().Format(time.RFC3339),
					Duration:  fmt.Sprintf("%d ms", time.Since(start)/time.Millisecond),
				},
			},
			wantErr: false,
		},
	}
	for _, myTest := range tests {
		tt := myTest
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteJSON(tt.args.report, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var js map[string]interface{}
			if json.Unmarshal([]byte(got), &js) != nil {
				t.Errorf("WriteJSON() = %v, want output in JSON format", got)
			}
		})
	}
}
