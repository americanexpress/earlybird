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

import "strings"

func initLookupTable(row, column int) [][]bool {
	lookup := make([][]bool, row)
	for i := range lookup {
		lookup[i] = make([]bool, column)
	}
	return lookup
}

//PatternMatch Function that matches input str with given wildcard pattern
func PatternMatch(str, pattern string) bool {
	s := []rune(strings.ToLower(str))
	p := []rune(pattern)

	// empty pattern can only match with empty string
	if len(p) == 0 {
		return len(s) == 0
	}

	// lookup table for storing results of subproblems
	// zero value of lookup is false
	lookup := initLookupTable(len(s)+1, len(p)+1)

	// empty pattern can match with empty string
	lookup[0][0] = true

	// Only '*' can match with empty string
	for j := 1; j < len(p)+1; j++ {
		if p[j-1] == '*' {
			lookup[0][j] = lookup[0][j-1]
		}
	}

	// fill the table in bottom-up fashion
	for i := 1; i < len(s)+1; i++ {
		for j := 1; j < len(p)+1; j++ {
			if p[j-1] == '*' {
				// Two cases if we see a '*'
				// a) We ignore ‘*’ character and move
				//    to next  character in the pattern,
				//     i.e., ‘*’ indicates an empty sequence.
				// b) '*' character matches with ith
				//     character in input
				lookup[i][j] = lookup[i][j-1] || lookup[i-1][j]

			} else if p[j-1] == '?' || s[i-1] == p[j-1] {
				// Current characters are considered as
				// matching in two cases
				// (a) current character of pattern is '?'
				// (b) characters actually match
				lookup[i][j] = lookup[i-1][j-1]

			} else {
				// If characters don't match
				lookup[i][j] = false
			}
		}
	}

	return lookup[len(s)][len(p)]
}
