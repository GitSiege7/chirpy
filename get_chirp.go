package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirp_id := r.PathValue("chirpID")

	id, err := uuid.Parse(chirp_id)
	if err != nil {
		respondWithError(w, 400, "Invalid ID")
		return
	}

	db_chirp, err := cfg.queries.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	chirp := Chirp{
		ID:        db_chirp.ID,
		CreatedAt: db_chirp.CreatedAt,
		UpdatedAt: db_chirp.UpdatedAt,
		Body:      db_chirp.Body,
		UserID:    db_chirp.UserID,
	}

	err = respondWithJSON(w, 200, chirp)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
