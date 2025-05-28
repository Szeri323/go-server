package main

import "net/http"

func (cfg *apiConfig) handlerHealthz(w http.ResponseWriter, r *http.Request) {
	respondWithPlaintext(w, http.StatusOK, []byte("OK"))
}
