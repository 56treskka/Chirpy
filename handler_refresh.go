package main

import (
	"net/http"
	"time"

	"github.com/56treskka/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	user, err := cfg.queries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token is invalid", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(time.Hour))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}
	err = cfg.queries.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
