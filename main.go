package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	file_server := http.FileServer(http.Dir("/home/siege/workspace/Go/chirpy/"))

	mux.Handle("/", file_server)

	server.ListenAndServe()
}
