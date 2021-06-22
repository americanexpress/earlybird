/*
 * Copyright 2021 American Express
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
	"regexp"
)

var (
	notNum = regexp.MustCompile(`[^0-9]+`)
	delta  = []int{0, 1, 2, 3, 4, -4, -3, -2, -1, 0}
)

//IsCard Run a mod10 check on a potential card number
func IsCard(cc string) bool {
	cc = isolateNumber(cc)
	checksum := 0
	bOdd := false
	card := []byte(cc)
	for i := len(card) - 1; i > -1; i-- {
		cn := int(card[i]) - '0'
		checksum += cn
		if bOdd {
			checksum += delta[cn]
		}
		bOdd = !bOdd
	}
	return checksum%10 == 0
}

// Strip off any non numerical garbage from the card number before processing
func isolateNumber(cc string) string {
	processedCC := notNum.ReplaceAllString(cc, "")
	return processedCC
}
