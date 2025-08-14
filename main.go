package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/GitSiege7/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
}

func respondWithJSON(w http.ResponseWriter, code int, dat interface{}) error {
	resp, err := json.Marshal(dat)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func handlerReadiness(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerCount(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()

	fmt.Fprintf(w, "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", count)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func handlerValid(w http.ResponseWriter, r *http.Request) {
	type fail struct {
		Error string `json:"error"`
	}
	type pass struct {
		Cleaned string `json:"cleaned_body"`
	}

	type req struct {
		Body string `json:"body"`
	}
	chirp := req{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		resp := fail{}
		resp.Error = "Failed to decode"

		dat, err := json.Marshal(resp)
		if err != nil {
			respondWithError(w, 500, "Failed to marshal error message")
			return
		}

		respondWithJSON(w, 500, dat)
		return
	}

	if len(chirp.Body) > 140 {
		resp := fail{}
		resp.Error = "Chirp is too long"

		err := respondWithJSON(w, 400, resp)
		if err != nil {
			respondWithError(w, 500, "Failed to marshal error message")
			return
		}
	} else {
		resp := pass{}

		words := strings.Split(chirp.Body, " ")
		for i := range words {
			if strings.ToLower(words[i]) == "kerfuffle" || strings.ToLower(words[i]) == "sharbert" || strings.ToLower(words[i]) == "fornax" {
				words[i] = "****"
			}
		}

		resp.Cleaned = strings.Join(words, " ")

		err := respondWithJSON(w, 200, resp)
		if err != nil {
			respondWithError(w, 500, "Failed to marshal pass message")
			return
		}
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Failed to open db: %v", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	cfg := &apiConfig{queries: dbQueries}
	cfg.fileserverHits.Store(0)

	filepathRoot := "/home/siege/workspace/Go/chirpy/"

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerCount)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValid)

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	server.ListenAndServe()
}
