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

func TestPasswordWeak(t *testing.T) {
	type args struct {
		testPD string
	}
	tests := []struct {
		name     string
		args     args
		wantWeak bool
	}{
		{
			name: "Test a strong password",
			args: args{
				testPD: "VeryStrong8161#",
			},
			wantWeak: false,
		},
		{
			name: "Test a very short password",
			args: args{
				testPD: "foo",
			},
			wantWeak: true,
		},
		{
			name: "Test another weak password",
			args: args{
				testPD: "test123",
			},
			wantWeak: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWeak := PasswordWeak(tt.args.testPD); gotWeak != tt.wantWeak {
				t.Errorf("PasswordWeak() = %v, want %v", gotWeak, tt.wantWeak)
			}
		})
	}
}
