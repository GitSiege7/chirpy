package auth

import (
	"fmt"
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	UUID := uuid.New()

	jwt, err := MakeJWT(UUID, "IzstEIf0C5YVYQ539dqglVtD+ZMZN4GYcOP4r6AgHb/q+zrBW3lhtvQdgTz850Vt0OSCiFgGOC23tevYbT14ug==", 5*time.Minute)
	if err != nil {
		fmt.Println("Failed to make jwt")
		t.Fail()
		return
	}

	header.Set("Authorization", "Bearer "+jwt)

	got, err := GetBearerToken(header)
	if err != nil {
		fmt.Println("Failed to get bearer token")
	}

	if got != jwt {
		fmt.Printf("%v... != %v...\n", got[0:20], jwt[0:20])
		fmt.Println("Mismatch outputs")
		t.Fail()
	}
}
