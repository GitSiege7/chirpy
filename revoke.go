package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := bearer[7:]

	err := cfg.queries.SetRevoked(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, "Failed to set as revoked")
		return
	}

	w.WriteHeader(204)
}
