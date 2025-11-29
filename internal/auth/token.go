package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization scheme must be Bearer")
	}

	return parts[1], nil
}