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
	"unicode"
)

// PasswordWeak If the password doesn't meet minimum requirements, call it weak
func PasswordWeak(password string) (weak bool) {
	var number, special, upper bool

	//Seperate password from full value
	passwords := pswdPattern.FindStringSubmatch(password)
	if len(passwords) > 2 {
		password = passwords[2]
	}

	for _, s := range password {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
		}
	}

	//Password is weak if it doesn't contain atleast one upper case letter, special character, number and atleast 7 letters
	if !(number && special && upper && len(password) > strongPswdLen) {
		return true
	}

	//Password is weak if it contains sequential numbers
	return containsSequenceOfNumbers(password)
}

func containsSequenceOfNumbers(s string) bool {
	for i := 0; i < len(s)-3; i++ {
		a, b, c := s[i], s[i+1], s[i+2]
		if a < '0' || b < '0' || c < '0' || a > '9' || b > '9' || c > '9' {
			continue
		}

		if b == a+1 && c == a+2 { // ascending
			return true
		}
		if b == a-1 && c == a-2 { // descending
			return true
		}
	}
	return false
}
