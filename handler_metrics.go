package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) printHits() []byte {
	hits := cfg.fileserverHits.Load()
	hitStr := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
	return []byte(hitStr)
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(cfg.printHits())
}
