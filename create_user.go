package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GitSiege7/chirpy/internal/auth"
	"github.com/GitSiege7/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var dat req

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dat)
	if err != nil {
		respondWithError(w, 500, "Failed to decode request")
		return
	}

	hash, err := auth.HashPassword(dat.Password)
	if err != nil {
		respondWithError(w, 500, "Failed to hash password")
		return
	}

	db_user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          dat.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to create user")
		return
	}

	user := User{
		ID:          db_user.ID,
		CreatedAt:   db_user.CreatedAt,
		UpdatedAt:   db_user.UpdatedAt,
		Email:       db_user.Email,
		IsChirpyRed: db_user.IsChirpyRed.Bool,
	}

	err = respondWithJSON(w, 201, user)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
