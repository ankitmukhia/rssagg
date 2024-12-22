package auth

import (
	"net/http"
	"errors"
	"strings"
)

// extracts and API key from the header of an HTTP req
// Authorization: Bearer {api_key/token}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("undefined auth header")
	}

	token := strings.Split(val, " ")
	if len(token) != 2 {
		return "", errors.New("malformed token")
	}

	if token[0] != "Bearer" {
		return "", errors.New("malformed token bearer")
	}
	return token[1], nil
}
