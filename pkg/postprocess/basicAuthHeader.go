package postprocess

import (
	"encoding/base64"
	"strings"
)

// decodeBase64 decodes a Base64-encoded string
func decodeBase64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func IsBasicAuthHeader(rawText string) bool {

	// check if the rawText contains "Basic" (case-insensitive).  Raw text is expected to looks like this somewhat this Authorization: Basic dXNlbmFtZTpwYXNzd29yZA==
	if !strings.Contains(strings.ToLower(rawText), "basic") {
		return false
	}

	// Getting the start index of base64 encoded text from the header.
	Index := strings.Index(strings.ToLower(rawText), "basic")
	encryptedText := ""
	for i := Index + 5; i < len(rawText); i++ {
		if rawText[i] == ' ' || rawText[i] == '"' || rawText[i] == '`' || rawText[i] == '$' || rawText[i] == '\'' { // removing the trailing spaces and quotes if any.
			continue
		}
		encryptedText += string(rawText[i]) // appending character to encryptedText
	}

	decodedText, err := decodeBase64(encryptedText) // decoding the base64 encoded text.
	if err != nil || decodedText == "" {
		return false
	}
	colonIndex := strings.Index(decodedText, ":")
	return colonIndex != -1 && colonIndex != 0 && colonIndex != len(decodedText)-1 && len(decodedText) > 2 // checking if the decoded text contains ':' and has two parts (username and password).
}
