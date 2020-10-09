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

import (
	"math"
)

//Shannon is an algorithm used to calculate the complexity of the string
func Shannon(s string) float64 {
	// count as integers to maintain precision case we have a very large (>10**24 byte) string.
	var freq [256]int
	var totalInt int
	for _, b := range []byte(s) {
		freq[b]++
		totalInt++
	}

	var entropy float64
	total := float64(totalInt)
	for _, count := range freq {
		if count == 0 {
			continue
		}
		count := float64(count)
		pval := count / float64(total)
		pinv := total / float64(count)
		entropy += pval * math.Log2(pinv)
	}
	return entropy
}
