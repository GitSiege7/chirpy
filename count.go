package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerCount(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()

	fmt.Fprintf(w, "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", count)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")
}
