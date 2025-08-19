package auth

import (
	"fmt"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth_header := headers.Get("Authorization")
	if auth_header == "" {
		return "", fmt.Errorf("invalid authorization")
	}

	return auth_header[7:], nil
}
