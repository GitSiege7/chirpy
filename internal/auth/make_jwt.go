package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(3600) * time.Second)),
		Subject:   fmt.Sprintf("%v", userID),
	})

	jwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return jwt, nil
}
