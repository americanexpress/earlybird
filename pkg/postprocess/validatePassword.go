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
	"github.com/americanexpress/earlybird/pkg/utils"
	"regexp"
	"strconv"
	"strings"
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
		confidence = 2 // Mark confidence high if passes all the following
	}

	return confidence, false
}

// SkipSameKeyValuePassword reads the matched value from hit object to compare if the key and value are same for passwords or secrets.
// UseCase: In config/properties files where key/value pair are used to inject passwords.
func SkipSameKeyValuePassword(matchValue string, lineValue string) (result bool) {
	findingsArray := [2]string{matchValue, lineValue}
	for _, password := range findingsArray {
		// Separate password from full value as key/value pair
		splitFindings := splitPswdPattern.Split(password, -1)
		// Check if length = 2 for true key/value pair.
		if len(splitFindings) == 2 {
			// Trim the spaces and remove the Quotes(single or double) if they exist and lower case it for equality check.
			trimKey := strings.ToLower(strings.Trim(strings.TrimSpace(splitFindings[0]), "\"'"))
			trimValue := strings.ToLower(strings.Trim(strings.TrimSpace(splitFindings[1]), "\"'"))

			// Check if key and value are equal. Only validates the alphanumeric part.
			// Might cause us to skip some real findings but has a very low possibility.
			if strings.EqualFold(utils.GetAlphaNumericValues(trimKey), utils.GetAlphaNumericValues(trimValue)) {
				// return early if its same key/value
				return true
			}
		}
	}
	return false
}

// SkipPasswordWithUnicode returns true if the password value contains a non ASCII character.
// UseCase: Localized content contains unicode char for different languages which cannot be passwords in real world.
func SkipPasswordWithUnicode(password string) bool {
	passwords := splitPswdPattern.Split(password, -1)
	// Check if length = 2 for true key/value pair.
	if len(passwords) == 2 {
		// convert the unicode chars into string to compare with original password value
		// If the password value is not a string with quotes then this would never return a value
		passwordStringValue, err := strconv.Unquote(strings.TrimSpace(passwords[1]))
		if err != nil {
			// string with capital unicode errors out, invalid- "Information\U00e4", valid- "Information\u00e4"
			// Convert the value to lowercase to make sure the unicode is really invalid.
			passwordStringValue, err = strconv.Unquote(strings.TrimSpace(strings.ToLower(passwords[1])))
			if err != nil {
				return false
			}
		}

		for i, c := range passwordStringValue {
			// The rune at the index should match the string rune at the same index for ASCII values.
			if string(passwordStringValue[i]) != string(c) {
				// Skips early as soon as it finds a non ASCII rune while iterating the string.
				return true
			}
		}
	}
	return false
}
