package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("Authorization header not found")
	}

	bearerToken := strings.TrimSpace(strings.ReplaceAll(auth, "Bearer", ""))
	return bearerToken, nil
}
