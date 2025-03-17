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
	"golang.org/x/net/html"
	"regexp"
	"strings"
	"unicode"
)

var (
	pswdPattern      = regexp.MustCompile(pswdRegex)
	splitPswdPattern = regexp.MustCompile(splitPswdRegex)
)

// PasswordFalse Make sure we don't report strings as passwords if they contain indicators that it's a variable assignment,
// or some other disqualifier based on line context
func PasswordFalse(password string) (confidence int, ignore bool) {
	var quoted int
	confidence = 3 // Mark default confidence medium
	// Separate password from full value
	passwords := pswdPattern.FindStringSubmatch(password)
	// Note: the conditional >2 will never evaluate to true with the "current" pswdRegex.
	// However, leaving to ensure we do not break any other flows and leaving to core development team to identify if can be removed later
	if len(passwords) > 2 {
		password = passwords[2]
	} else if len(passwords) > 1 {
		password = passwords[1]
	}

	// Remove space if present
	password = strings.TrimSpace(password)

	// Remove double quotes if present
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
		// Check if contains dot for object label references
		if strings.Contains(password, ".") {
			return confidence, true
		}

		// Check if starts with $ could be a variable reference
		if password[0] == '$' {
			return confidence, true
		}

		// If begins with = could be comparison between two values
		if strings.Contains(password, "==") {
			return confidence, true
		}

		// Password can't contain spaces if it's not quoted unless there's a comma
		if strings.Contains(password, " ") && !strings.Contains(password, ",") && !strings.Contains(password, "=") {
			return confidence, true
		}

		// Function references contain () unquoted
		if strings.Contains(password, "(") && strings.Contains(password, ")") {
			return confidence, true
		}
	}
	confidence = 2 // Mark confidence high if passes all the following

	return confidence, false
}

// SkipPasswordWithUnicode returns true if the password value contains a non ASCII character.
// UseCase: Localized content contains unicode char for different languages which cannot be passwords in real world.
func SkipPasswordWithUnicode(password string) bool {
	for _, c := range password {
		// The rune c is a Unicode code point, and we check if it is greater than 127 to identify non-ASCII characters.
		if c > unicode.MaxASCII {
			// Skips early as soon as it finds a non ASCII rune while iterating the string.
			return true
		}
	}

	return false
}

// SkipPasswordWithHTMLEntities returns true if the password value contains a HTML entities.
// UseCase: Html entities are used as contents in localized files. ex - "L&#246;senord" is "Lösenord"
func SkipPasswordWithHTMLEntities(password string) bool {
	passwords := splitPswdPattern.Split(password, -1)
	// Check if length = 2 for true key/value pair.
	if len(passwords) == 2 {
		// Unescape the html entities into string to compare with original password value
		passwordStringValue := html.UnescapeString(strings.TrimSpace(passwords[1]))
		// ex - "L&#246;senord" is "Lösenord" if after html.UnescapeString they are not equal then it has html entities.
		if passwordStringValue != strings.TrimSpace(passwords[1]) {
			return true
		}
	}
	return false
}
