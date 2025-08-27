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

	if !strings.Contains(strings.ToLower(rawText), "authorization") || !strings.Contains(strings.ToLower(rawText), "basic") {
		return false
	}

	Index := strings.Index(strings.ToLower(rawText), "basic")
	encryptedText := ""
	for i := Index + 5; i < len(rawText); i++ {
		if rawText[i] == ' ' || rawText[i] == '"' || rawText[i] == '`' || rawText[i] == '$' || rawText[i] == '\'' {
			continue
		}
		encryptedText += string(rawText[i])
	}

	decodedText, err := decodeBase64(encryptedText)
	if err != nil || decodedText == "" {
		return false
	}
	return strings.Contains(decodedText, ":") && len(strings.Split(decodedText, ":")) == 2 && len(decodedText) > 2
}
