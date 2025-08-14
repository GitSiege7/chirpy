package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		err := respondWithError(w, 403, "Reset only accessible in local development environment")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	cfg.fileserverHits.Store(0)

	err := cfg.queries.DeleteUsers(r.Context())
	if err != nil {
		err = respondWithError(w, 500, "Failed to delete users")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
