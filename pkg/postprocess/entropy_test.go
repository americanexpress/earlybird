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

// Compute Shannon Entropy of a byte stream
// H = - Σ P(x) * log P(x)
package postprocess

import "testing"

func TestShannon(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Test complex string for Shannon complexity",
			args: args{
				s: "G$#G%^J&%$J52165251h6$%FTH%$H!",
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "Test single character for Shannon complexity",
			args: args{
				s: "0",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Shannon(tt.args.s)
			if !(got >= float64(tt.want)) {
				t.Errorf("Shannon() = %v >= want %v", got, tt.want)
			}
		})
	}
}
