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

import (
	"strconv"
	"strings"
)

// ValidSSN checks if a SSN meets standard
func ValidSSN(ssn string) bool {
	ssn = strings.Trim(ssn, "\"'\n ")
	groups := strings.Split(ssn, "-")
	if len(groups) != 3 {
		return false
	}
	if first, _ := strconv.Atoi(groups[0]); first == 666 || first <= 0 || first > 999 {
		return false
	} else if second, _ := strconv.Atoi(groups[1]); second <= 0 || second > 99 {
		return false
	} else if third, _ := strconv.Atoi(groups[2]); third <= 0 || third > 9999 {
		return false
	}
	return true

}
