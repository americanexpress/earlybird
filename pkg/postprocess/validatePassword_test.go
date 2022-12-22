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
		name: "Skip passwords with a dot",
		args: args{
			testPD: "password: ignore.me",
		},
		wantConfidence: 3,
		wantIgnore:     true,
	},
	{
		name: "Skip passwords with two equals",
		args: args{
			testPD: "password: ignoreme==please",
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
	{
		name: "Verify = delimited values",
		args: args{
			testPD: "my.property=propertyEqualDelimitedPassword",
		},
		wantConfidence: 2,
		wantIgnore:     false,
	},
	{
		name: "Verify : delimited values",
		args: args{
			testPD: "my.property:propertyColonDelimitedPassword",
		},
		wantConfidence: 2,
		wantIgnore:     false,
	},
	{
		name: "Do not skip, real finding, ensure whitespace is permitted around delimited values",
		args: args{
			testPD: "my.property    =     propertySpacesAroundDelimited",
		},
		wantConfidence: 2,
		wantIgnore:     false,
	},
	{
		name: "Do not skip, real finding, ensure yml files are handled",
		args: args{
			testPD: "my.property: sampleYmlPassword",
		},
		wantConfidence: 2,
		wantIgnore:     false,
	},
	{
		name: "Do not skip, real finding, ensure json files are handled",
		args: args{
			testPD: "\"my.property\": \"sample%3YmlPassword\"",
		},
		wantConfidence: 3,
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

var testSkipUnicodeInPasswords = []struct {
	name       string
	args       args
	wantIgnore bool
}{
	{
		name: "Skip password with unicode which is not ASCII",
		args: args{
			testPD: `"password": "\u0049\u0044\u306e\u78ba\u8a8d\u3001\u30d1\u30b9\u30ef\u30fc\u30c9\u306e\u5909\u66f4"`,
		},
		wantIgnore: true,
	},
	{
		name: "Skip passwords that has non ASCII char",
		args: args{
			testPD: `"password": "VeryStrong$$\u306e\u78ba"`,
		},
		wantIgnore: true,
	},
	{
		name: "Do not skip password with unicode that convert to valid of ASCII.",
		args: args{
			testPD: `"password": "VeryStrong$$\u0049\u0044"`,
		},
		wantIgnore: false,
	},
	{
		name: "Do not skip, real password finding",
		args: args{
			testPD: "password: VeryStrong857!@$^&*#",
		},
		wantIgnore: false,
	},
	{
		name: "Do not skip, real secret finding",
		args: args{
			testPD: "secret: VeryStrong857#",
		},
		wantIgnore: false,
	},
	{
		name: "Skip passwords that has non ASCII char and is invalid string due to unicode being in CAPS",
		args: args{
			testPD: `"password"= "Informationsb\U00e4rare"`,
		},
		wantIgnore: true,
	},
}

func TestSkipUnicodeInPasswords(t *testing.T) {
	for _, tt := range testSkipUnicodeInPasswords {
		t.Run(tt.name, func(t *testing.T) {
			gotIgnore := SkipPasswordWithUnicode(tt.args.testPD)
			if gotIgnore != tt.wantIgnore {
				t.Errorf("SkipUnicodeInPasswords() gotIgnore = %v, want %v", gotIgnore, tt.wantIgnore)
			}
		})
	}
}

var testSkipHTMLEntitiesInPasswords = []struct {
	name       string
	args       args
	wantIgnore bool
}{
	{
		name: "Skip password with HTML entities",
		args: args{
			testPD: `"password": "L&#246;senord"`,
		},
		wantIgnore: true,
	},
	{
		name: "Skip secret with HTML entities",
		args: args{
			testPD: `"secret": "&#12497;&#12473;&#12527;&#12540;&#12489"`,
		},
		wantIgnore: true,
	},
	{
		name: "Do not skip password without HTML entities",
		args: args{
			testPD: `"password": "VeryStrong$$\u0049\u0044"`,
		},
		wantIgnore: false,
	},
	{
		name: "Do not skip, real password finding",
		args: args{
			testPD: "password: VeryStrong857!@$^&*#",
		},
		wantIgnore: false,
	},
	{
		name: "Do not skip, real secret finding",
		args: args{
			testPD: "secret: VeryStrong857#",
		},
		wantIgnore: false,
	},
}

func TestSkipHTMLEntitiesInPasswords(t *testing.T) {
	for _, tt := range testSkipHTMLEntitiesInPasswords {
		t.Run(tt.name, func(t *testing.T) {
			gotIgnore := SkipPasswordWithHTMLEntities(tt.args.testPD)
			if gotIgnore != tt.wantIgnore {
				t.Errorf("SkipUnicodeInPasswords() gotIgnore = %v, want %v", gotIgnore, tt.wantIgnore)
			}
		})
	}
}
