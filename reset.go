package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Reset only accessible in local development environment")
		return
	}

	cfg.fileserverHits.Store(0)

	err := cfg.queries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to delete users")
		return
	}

	w.WriteHeader(http.StatusOK)
}
