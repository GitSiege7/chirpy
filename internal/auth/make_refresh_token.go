package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("Failed rand.read")
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
