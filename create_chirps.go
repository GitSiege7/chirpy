package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/GitSiege7/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	post := req{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		err = respondWithError(w, 500, "Failed to decode")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	if len(post.Body) > 140 {
		err := respondWithError(w, 400, "Chirp is too long")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	} else {
		words := strings.Split(post.Body, " ")
		for i := range words {
			if strings.ToLower(words[i]) == "kerfuffle" || strings.ToLower(words[i]) == "sharbert" || strings.ToLower(words[i]) == "fornax" {
				words[i] = "****"
			}
		}

		cleaned := strings.Join(words, " ")

		db_chirp, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleaned,
			UserID: post.UserID,
		})
		if err != nil {
			err = respondWithError(w, 500, "Failed to create chirp")
			if err != nil {
				fmt.Println("Failed to respond")
			}
			return
		}

		chirp := Chirp{
			ID:        db_chirp.ID,
			CreatedAt: db_chirp.CreatedAt,
			UpdatedAt: db_chirp.UpdatedAt,
			Body:      db_chirp.Body,
			UserID:    db_chirp.UserID,
		}

		err = respondWithJSON(w, 201, chirp)
		if err != nil {
			fmt.Println("Failed to respond")
		}
	}
}
