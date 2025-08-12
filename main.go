package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handlerReadiness(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerCount(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()

	fmt.Fprintf(w, "Hits: %v", count)
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	cfg := &apiConfig{}
	cfg.fileserverHits.Store(0)

	filepathRoot := "/home/siege/workspace/Go/chirpy/"

	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", cfg.handlerCount)
	mux.HandleFunc("/reset", cfg.handlerReset)

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	server.ListenAndServe()
}
