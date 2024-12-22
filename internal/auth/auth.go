package auth

import  (
	"net/http"
	"errors"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("no auth header include")
	}

	token := strings.Split(header, " ")
	if len(token) < 2 || token[0] != "Bearer" {
		return "", errors.New("malformed auth header")
	}

	return token[1], nil
}
