package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.queries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}

	data := make([]Chirp, len(chirps))

	for i, elem := range chirps {
		data[i] = Chirp{
			ID:        elem.ID,
			CreatedAt: elem.CreatedAt,
			UpdatedAt: elem.UpdatedAt,
			Body:      elem.Body,
			UserID:    elem.UserID.UUID,
		}
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirps, err := cfg.queries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}

	for _, elem := range chirps {
		if elem.ID.String() == chirpID {
			respondWithJSON(w, http.StatusOK, Chirp{
				ID:        elem.ID,
				CreatedAt: elem.CreatedAt,
				UpdatedAt: elem.UpdatedAt,
				Body:      elem.Body,
				UserID:    elem.UserID.UUID,
			})
			return
		}
	}
	respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
}
