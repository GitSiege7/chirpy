package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GitSiege7/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	type req struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		ExpiresInSecs int    `json:"expires_in_seconds,omitempty"`
	}
	var dat req

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dat)
	if err != nil {
		err = respondWithError(w, 500, "Failed to decode request")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	if dat.ExpiresInSecs == 0 || dat.ExpiresInSecs > 3600 {
		dat.ExpiresInSecs = 3600
	}

	db_user, err := cfg.queries.GetUserByEmail(r.Context(), dat.Email)
	if err != nil {
		err = respondWithError(w, 404, "User not found")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	err = auth.CheckPasswordHash(dat.Password, db_user.HashedPassword)
	if err != nil {
		err = respondWithError(w, 401, "Unauthorized: incorrect password")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	jwt, err := auth.MakeJWT(db_user.ID, cfg.secret, time.Duration(dat.ExpiresInSecs)*time.Second)
	if err != nil {
		err = respondWithError(w, 500, "Failed to make JWT")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	res := response{
		ID:        db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email:     db_user.Email,
		Token:     jwt,
	}

	err = respondWithJSON(w, 200, res)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
