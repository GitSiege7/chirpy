package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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
			err = respondWithError(w, 500, "Failed to marshal error message")
			if err != nil {
				fmt.Println("Failed to respond")
			}
			return
		}

		err = respondWithJSON(w, 500, dat)
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	if len(chirp.Body) > 140 {
		resp := fail{}
		resp.Error = "Chirp is too long"

		err := respondWithJSON(w, 400, resp)
		if err != nil {
			err = respondWithError(w, 500, "Failed to marshal error message")
			if err != nil {
				fmt.Println("Failed to respond")
			}
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
			err = respondWithError(w, 500, "Failed to marshal pass message")
			if err != nil {
				fmt.Println("Failed to respond")
			}
			return
		}
	}
}
