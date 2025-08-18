package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GitSiege7/chirpy/internal/auth"
	"github.com/GitSiege7/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

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

	db_user, err := cfg.queries.GetUserByEmail(r.Context(), dat.Email)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	err = auth.CheckPasswordHash(dat.Password, db_user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: incorrect password")
		return
	}

	jwt, err := auth.MakeJWT(db_user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, 500, "Failed to make JWT")
		return
	}

	ref, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Failed to make refresh token")
		return
	}

	_, err = cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: ref,
		UserID: uuid.NullUUID{
			UUID:  db_user.ID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, 500, "Failed to add refresh token to db")
		return
	}

	res := response{
		ID:           db_user.ID,
		CreatedAt:    db_user.CreatedAt,
		UpdatedAt:    db_user.UpdatedAt,
		Email:        db_user.Email,
		Token:        jwt,
		RefreshToken: ref,
	}

	err = respondWithJSON(w, 200, res)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
