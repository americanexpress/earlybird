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
	"regexp"
	"strings"
)

var (
	pswdPattern = regexp.MustCompile(pswdRegex)
)

//PasswordFalse Make sure we don't report strings as passwords if they contain indicators that it's a variable assignment,
// or some other disqualifier based on line context
func PasswordFalse(password string) (confidence int, ignore bool) {
	var quoted int
	confidence = 3 // Mark default confidence medium
	//Seperate password from full value
	passwords := pswdPattern.FindStringSubmatch(password)
	if len(passwords) > 2 {
		password = passwords[2]
	}

	//Remove space if present
	password = strings.TrimSpace(password)

	//Remove double quotes if present
	if len(password) > 0 && password[0] == '"' {
		password = password[1:]
		quoted++
	}
	if len(password) > 0 && password[len(password)-1] == '"' {
		password = password[:len(password)-1]
		quoted++
	}

	if len(password) < pswdMinLen {
		return confidence, true
	}

	if quoted != 2 {
		//Check if contains dot for object label references
		if strings.Contains(password, ".") {
			return confidence, true
		}

		//Check if starts with $ could be a variable reference
		if password[0] == '$' {
			return confidence, true
		}

		//If begins with = could be comparison between two values
		if strings.Contains(password, "==") {
			return confidence, true
		}

		//Password can't contain spaces if it's not quoted unless there's a comma
		if strings.Contains(password, " ") && !strings.Contains(password, ",") && !strings.Contains(password, "=") {
			return confidence, true
		}

		//Function references contain () unquoted
		if strings.Contains(password, "(") && strings.Contains(password, ")") {
			return confidence, true
		}
		confidence = 2 // Mark confidence high if passes all the following
	}

	return confidence, false
}
