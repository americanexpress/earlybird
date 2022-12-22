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

func TestWriteCSV(t *testing.T) {
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
		hits     chan scan.Hit
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Convert hit from channel to CSV",
			args: args{
				hits: HitChan,
			},
			wantErr: false,
		},
	}
	for _, myTest := range tests {
		tt := myTest
		go func() {
			t.Run(tt.name, func(t *testing.T) {
				if err := WriteCSV(tt.args.hits, tt.args.fileName); (err != nil) != tt.wantErr {
					t.Errorf("WriteCSV() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}()
	}
}
