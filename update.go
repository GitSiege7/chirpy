package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GitSiege7/chirpy/internal/auth"
	"github.com/GitSiege7/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 400, "Failed to decode request")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, 500, "Failed to hash password")
		return
	}

	db_user, err := cfg.queries.UpdateCredentials(r.Context(), database.UpdateCredentialsParams{
		ID:             id,
		Email:          req.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		respondWithError(w, 400, "Failed to update user credentials")
		return
	}

	user := User{
		ID:          db_user.ID,
		CreatedAt:   db_user.CreatedAt,
		UpdatedAt:   db_user.UpdatedAt,
		Email:       db_user.Email,
		IsChirpyRed: db_user.IsChirpyRed.Bool,
	}

	err = respondWithJSON(w, 200, user)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
