package postprocess

import (
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// IsValidJWT checks if the JWT token is valid and not expired.
func IsValidJWT(rawText string) bool {
	// Extract the JWT token from the raw text.
	jwtToken := ""
	for _, ch := range rawText {
		if ch != ' ' && ch != '"' && ch != ':' && ch != '=' && ch != '\'' {
			jwtToken += string(ch)
		}
	}

	// Token is encrypted.
	if len(strings.Split(jwtToken, ".")) == 5 {
		return false
	}

	// Parse the JWT token.
	token, _, err := new(jwt.Parser).ParseUnverified(jwtToken, jwt.MapClaims{})
	if err != nil {
		return true
	}

	// Extract the claims and check the expiration time.
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			// Checking if the token is expired.
			expirationTime := time.Unix(int64(exp), 0)
			if expirationTime.Before(time.Now()) {
				return false
			}
		}
	}
	return true // Expiration claim is missing or invalid or token is valid or not expired.
}
