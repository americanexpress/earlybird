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

package wildcard

import (
	"reflect"
	"testing"
)

func Test_initLookupTable(t *testing.T) {
	type args struct {
		row    int
		column int
	}
	tests := []struct {
		name string
		args args
		want [][]bool
	}{
		{
			name: "Compare init table to result",
			args: args{
				row:    2,
				column: 2,
			},
			want: [][]bool{{false, false}, {false, false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initLookupTable(tt.args.row, tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initLookupTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWildcardPatternMatch(t *testing.T) {
	type args struct {
		str     string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Match google wildcard",
			args: args{
				str:     "google",
				pattern: "g*gle",
			},
			want: true,
		},
		{
			name: "Match LICENSE",
			args: args{
				str:     "LICENSE",
				pattern: "LICENSE",
			},
			want: true,
		},
		{
			name: "Match LICENSE with wildcards: 1",
			args: args{
				str:     "LICENSE",
				pattern: "*LICENSE*",
			},
			want: true,
		},
		{
			name: "Match filepath with wildcards",
			args: args{
				str:     "bundle.min.js",
				pattern: "*.min.js",
			},
			want: true,
		},
		{
			name: "Match filepath extension",
			args: args{
				str:     "something.woff",
				pattern: "*.woff",
			},
			want: true,
		},
		{
			name: "Match file recursively",
			args: args{
				str:     "/some/deeply/nested/dir/something.woff",
				pattern: "*.woff",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PatternMatch(tt.args.str, tt.args.pattern); got != tt.want {
				t.Errorf("WildcardPatternMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
