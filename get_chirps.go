package main

import (
	"fmt"
	"net/http"

	"github.com/GitSiege7/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	author_id_str := r.URL.Query().Get("author_id")

	var db_chirps []database.Chirp
	var err error
	if author_id_str == "" {
		db_chirps, err = cfg.queries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, "Failed to get chirps")
			return
		}
	} else {
		author_id, err := uuid.Parse(author_id_str)
		if err != nil {
			respondWithError(w, 400, "Invalid author id")
			return
		}

		db_chirps, err = cfg.queries.GetChirpsByUser(r.Context(), author_id)
		if err != nil {
			respondWithError(w, 404, "Author not found")
			return
		}
	}

	chirps := []Chirp{}
	for _, chirp := range db_chirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	err = respondWithJSON(w, 200, chirps)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
