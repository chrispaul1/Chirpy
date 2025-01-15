package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("authorization header was not found")
	}

	tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		return "", fmt.Errorf("token could not be found")
	}

	return tokenString, nil
}
