package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/szeri323/go-server/internal/database"
)

// handler func(http.ResponseWriter, *http.Request)
type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error could not connect to database: %v", err)
	}
	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()

	var apiCfg apiConfig
	apiCfg.database = dbQueries
	fileServer := http.FileServer(http.Dir(filepathRoot))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("GET /api/healthz", apiCfg.handlerHealthz)
	serveMux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidate)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error could not start server: %v", err)
	}

}
