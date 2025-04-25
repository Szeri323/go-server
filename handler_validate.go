package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	body := request{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding parameters:: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(body.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	slice := strings.Split(body.Body, " ")
	for i := 0; i < len(slice); i++ {
		if strings.ToLower(slice[i]) == "kerfuffle" || strings.ToLower(slice[i]) == "sharbert" || strings.ToLower(slice[i]) == "fornax" {
			slice[i] = "****"
		}
		fmt.Println(slice[i])
	}
	type response struct {
		Valid        bool   `json:"valid"`
		Cleaned_body string `json:"cleaned_body"`
	}
	respBody := response{
		Valid:        true,
		Cleaned_body: strings.Join(slice, " "),
	}
	res, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}
