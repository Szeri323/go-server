package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/szeri323/go-server/internal/auth"
	"github.com/szeri323/go-server/internal/database"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid Authentication header")
		return
	}
	userID, err := cfg.database.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No user with this refresh token")
		return
	}
	jwt, err := auth.MakeJWT(userID, cfg.secret, time.Second*60*60)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't make a JWT")
		return
	}
	respondWithJSON(w, http.StatusOK, response{Token: jwt})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid Authentication header")
		return
	}
	fmt.Println("token:")
	fmt.Println(token)
	err = cfg.database.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: time.Now(),
		Token:     token,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token doesn't exists")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
