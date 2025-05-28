package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/szeri323/go-server/internal/auth"
	"github.com/szeri323/go-server/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := request{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters.")
		return
	}

	user, err := cfg.database.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user.")
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't authenticate user.")
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 60 {
		params.ExpiresInSeconds = 360
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Second*time.Duration(params.ExpiresInSeconds))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token.")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a refresh token.")
		return
	}
	_, err = cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Second * 60 * 60 * 24 * 60),
		RevokedAt: sql.NullTime{Time: time.Time{}, Valid: false},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't insert the refresh token record in database.")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Name,
			Token:        token,
			RefreshToken: refreshToken,
			IsChirpyRed:  user.IsChirpyRed,
		}})
}

func validateExpiresInSecondsHeader() {

}
