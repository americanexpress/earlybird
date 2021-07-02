package postprocess

import (
	"github.com/americanexpress/earlybird/pkg/utils"
	"regexp"
	"strings"
)

var (
	splitKeyValue = regexp.MustCompile(splitPswdRegex)
)

// SkipSameKeyValue reads the matched value from hit object to compare if the key and value are same.
// UseCase: In config/properties files where key/value pair are used
func SkipSameKeyValue(matchValue string, lineValue string) (result bool) {
	findingsArray := [2]string{matchValue, lineValue}
	for _, password := range findingsArray {
		// Separate password from full value as key/value pair
		splitFindings := splitKeyValue.Split(password, -1)
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
