package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/szeri323/go-server/internal/auth"
	"github.com/szeri323/go-server/internal/database"
)

type Data struct {
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	api_key, err := auth.GetAPIKey(r.Header)
	if cfg.polka_key != api_key {
		respondWithError(w, http.StatusUnauthorized, "No permission")
		return
	}
	type request struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := request{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse body")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = cfg.database.UpdateUsersMembership(r.Context(), database.UpdateUsersMembershipParams{
		IsChirpyRed: true,
		ID:          params.Data.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
