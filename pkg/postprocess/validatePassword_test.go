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

package postprocess

import "testing"

type args struct {
	testPD string
}

var tests = []struct {
	name           string
	args           args
	wantConfidence int
	wantIgnore     bool
}{
	{
		name: "Skip passwords too short",
		args: args{
			testPD: "fo",
		},
		wantConfidence: 3,
		wantIgnore:     true,
	},
	{
		name: "Skip variables",
		args: args{
			testPD: "$variable",
		},
		wantConfidence: 3,
		wantIgnore:     true,
	},
	{
		name: "Skip functions",
		args: args{
			testPD: "func()",
		},
		wantConfidence: 3,
		wantIgnore:     true,
	},
	{
		name: "Skip passwords with spaces and no quotes",
		args: args{
			testPD: "ignore me please",
		},
		wantConfidence: 3,
		wantIgnore:     true,
	},
	{
		name: "Do not skip, real finding",
		args: args{
			testPD: "VeryStrong857#",
		},
		wantConfidence: 2,
		wantIgnore:     false,
	},
}

func TestPasswordFalse(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfidence, gotIgnore := PasswordFalse(tt.args.testPD)
			if gotConfidence != tt.wantConfidence {
				t.Errorf("PasswordFalse() gotConfidence = %v, want %v", gotConfidence, tt.wantConfidence)
			}
			if gotIgnore != tt.wantIgnore {
				t.Errorf("PasswordFalse() gotIgnore = %v, want %v", gotIgnore, tt.wantIgnore)
			}
		})
	}
}
