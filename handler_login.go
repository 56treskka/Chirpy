package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/56treskka/Chirpy/internal/auth"
	"github.com/56treskka/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	expirationTime := time.Hour
	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	refreshTokenStr, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refersh token string", err)
		return
	}

	refreshToken, err := cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshTokenStr,
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refersh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}
