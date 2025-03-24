package main

import (
	"log"
	"net/http"
)

type ServeMux struct {
}

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error could not start server: %v", err)
	}
}
