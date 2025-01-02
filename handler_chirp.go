package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/56treskka/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: uuid.NullUUID{
			UUID:  params.User_ID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.UUID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	banWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	cleaned := removeWords(body, banWords)

	return cleaned, nil
}

func removeWords(body string, banWords map[string]struct{}) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		if _, ok := banWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
