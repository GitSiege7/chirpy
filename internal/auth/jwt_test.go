package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	UUID := uuid.New()
	tokenSecret := "supersecretstring"

	jwt, err := MakeJWT(UUID, tokenSecret, 5*time.Second)
	if err != nil {
		fmt.Println("Failed to make JWT")
		t.Fail()
		return
	}

	UUID2, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		fmt.Println("Failed to validate JWT")
		t.Fail()
		return
	}

	if UUID != UUID2 {
		t.Fail()
	}
}
