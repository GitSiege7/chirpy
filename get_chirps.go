package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	db_chirps, err := cfg.queries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to get chirps")
		return
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
