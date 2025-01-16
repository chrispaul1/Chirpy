package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIkey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("authorization header was not found")
	}

	apikey, ok := strings.CutPrefix(authHeader, "ApiKey ")
	if !ok {
		return "", fmt.Errorf("apikey could not be found")
	}

	return apikey, nil
}
