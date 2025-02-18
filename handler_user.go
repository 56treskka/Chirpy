package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/56treskka/Chirpy/internal/auth"
	"github.com/56treskka/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing password", err)
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: password,
	})
	if err != nil {
		respondWithError(w, 500, "Error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request body", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing password", err)
		return
	}

	user, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: password,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, 500, "Error updating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}
