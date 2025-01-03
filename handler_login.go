package main

import (
	"encoding/json"
	"net/http"

	"github.com/56treskka/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding parameters", err)
		return
	}

	user, err := cfg.queries.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	checkErr := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if checkErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", checkErr)
	}

	respondWithJSON(w, http.StatusOK, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
