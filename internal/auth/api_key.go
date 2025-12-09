package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	if !strings.EqualFold(parts[0], "ApiKey") {
		return "", errors.New("authorization scheme must be ApiKey")
	}

	return parts[1], nil
}