package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GitSiege7/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	db_user, err := cfg.queries.GetUserByEmail(r.Context(), dat.Email)
	if err != nil {
		err = respondWithError(w, 404, "User not found")
		if err != nil {
			fmt.Println("Failed to respond")
		}
		return
	}

	err = auth.CheckPasswordHash(dat.Password, db_user.HashedPassword)
	if err != nil {
		err = respondWithError(w, 401, "Unauthorized: incorrect password")
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

	err = respondWithJSON(w, 200, user)
	if err != nil {
		fmt.Println("Failed to respond")
	}
}
