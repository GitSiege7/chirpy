package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GitSiege7/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := bearer[7:]

	ref, err := cfg.queries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Auth refresh token not found")
		return
	}

	if ref.RevokedAt.Valid {
		respondWithError(w, 401, "Auth refresh token revoked")
		return
	}

	if time.Now().Compare(ref.ExpiresAt.Time) == 1 {
		respondWithError(w, 401, "Auth refresh token expired")
		return
	}

	null_uuid, err := cfg.queries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, "Failed to get user by refresh token")
		return
	}
	UUID := null_uuid.UUID

	type response struct {
		Token string `json:"token"`
	}
	var res response

	jwt, err := auth.MakeJWT(UUID, cfg.secret)
	if err != nil {
		respondWithError(w, 500, "Failed to make jwt")
		return
	}
	res.Token = jwt

	err = respondWithJSON(w, 200, res)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
