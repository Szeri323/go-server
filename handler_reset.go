package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) clearUsersTable(r *http.Request) {
	cfg.database.TruncateUsersTable(r.Context())
}

func (cfg *apiConfig) clearChirpsTable(r *http.Request) {
	cfg.database.TruncateChirpsTable(r.Context())
}

func (cfg *apiConfig) clearRefreshTokenTable(r *http.Request) {
	cfg.database.TruncateRefreshTokensTable(r.Context())
}

func (cfg *apiConfig) resetHits() {
	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.resetHits()
	platform := os.Getenv("PLATFORM")
	if platform == "dev" {
		cfg.clearUsersTable(r)
		cfg.clearChirpsTable(r)
		cfg.clearRefreshTokenTable(r)
		respondWithPlaintext(w, http.StatusOK, cfg.printHits())
	} else {
		respondWithPlaintext(w, http.StatusForbidden, []byte("Forbidden"))
	}

}
