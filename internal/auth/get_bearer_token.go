package auth

import (
	"fmt"
	"net/http"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("no bearer token found")
	}

	return bearer[7:], nil
}
