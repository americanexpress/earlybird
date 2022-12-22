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

package postprocess

import "testing"

func TestIsCard(t *testing.T) {
	type args struct {
		cc string
	}
	tests := []struct {
		name string
		cc   string
		want bool
	}{
		{
			"way too short - this is handled by the pattern matching",
			"0",
			true,
		},
		{
			"way too long - this may be handled by the pattern matching",
			"37000000000000200000000000000000000000000011",
			true,
		},

		{
			"garbage in middle - this is handled by the pattern matching",
			"370000ajlsdklasdj000000002",
			true,
		},
		{
			"Test AMEX example credit card",
			"370000000000002",
			true,
		},
		{
			"Test bogus credit card",
			"100000000000000",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCard(tt.cc); got != tt.want {
				t.Errorf("IsCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isolateNumber(t *testing.T) {
	type args struct {
		cc string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Remove characters from CC number",
			args: args{
				cc: "1000-0000-0000-000",
			},
			want: "100000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isolateNumber(tt.args.cc); got != tt.want {
				t.Errorf("isolateNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
