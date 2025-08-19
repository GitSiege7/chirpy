package main

import (
	"encoding/json"
	"net/http"

	"github.com/GitSiege7/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	var req request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 401, "Failed to decode request")
		return
	}

	api_key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid api key")
		return
	}

	if api_key != cfg.apiKey {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	id, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, 401, "Invalid user id")
		return
	}

	err = cfg.queries.UpgradeUser(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}

	w.WriteHeader(204)
}
