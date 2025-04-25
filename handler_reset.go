package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) clearUsersTable(r *http.Request) {
	cfg.database.TruncateUsersTable(r.Context())
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
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}
}
