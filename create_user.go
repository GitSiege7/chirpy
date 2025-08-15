package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email string `json:"email"`
	}
	var dat req

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dat)
	if err != nil {
		err = respondWithError(w, 500, "Failed to decode request")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	db_user, err := cfg.queries.CreateUser(r.Context(), dat.Email)
	if err != nil {
		err = respondWithError(w, 500, "Failed to create user")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	user := User{
		ID:        db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email:     db_user.Email,
	}

	err = respondWithJSON(w, 201, user)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
