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

func TestValidSSN(t *testing.T) {
	tests := []struct {
		name string
		ssn  string
		want bool
	}{
		{
			"wrong digit count ssn",
			"00000-000-0000000",
			false,
		},
		{
			"....",
			`"""ajsdilaksjdklasjhdklasjdlaskjd--a nsmcbnsjkcbnasdkjcbnsdjkf- 90123 40921u3 4"`,
			false},
		{
			"empty!",
			"---",
			false},
		{
			name: "Test fake SSN",
			ssn:  "123-45-6789",
			want: true,
		},
		{
			name: "Test invalid SSN",
			ssn:  "666-000-0000",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidSSN(tt.ssn); got != tt.want {
				t.Errorf("ValidSSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
