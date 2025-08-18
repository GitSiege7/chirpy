package main

import (
	"net/http"

	"github.com/GitSiege7/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}

	db_chirp, err := cfg.queries.GetChirpByID(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	if db_chirp.UserID != user_id {
		respondWithError(w, 403, "Invalid ownership")
		return
	}

	err = cfg.queries.DeleteChirp(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	w.WriteHeader(204)
}
