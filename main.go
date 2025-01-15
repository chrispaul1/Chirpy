package main

import (
	"chrispaul1/chirpy/internal/config"
	"chrispaul1/chirpy/internal/database"
	"chrispaul1/chirpy/internal/handlers"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error, loading env file : ", err)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error, opening database connection", err)
	}

	dbQueries := database.New(db)
	_ = dbQueries

	newApiConfig := config.ApiConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       os.Getenv("PLATFORM"),
		JWT_SECRET:     os.Getenv("JWT_SECRET"),
	}

	if newApiConfig.JWT_SECRET == "" {
		err := errors.New("token secret missing, failed to generate JWT")
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	userHandler := handlers.NewUserHandler(&newApiConfig)
	mux.HandleFunc("/api/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusPermanentRedirect)
	})
	handle := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", userHandler.MiddlewareMetricsInc(handle))
	mux.HandleFunc("GET /api/healthz", okHandler)
	mux.HandleFunc("POST /admin/reset", userHandler.ResetHandler)
	mux.HandleFunc("GET /admin/metrics", userHandler.MetricsHandler)
	mux.HandleFunc("POST /api/users", userHandler.HandleUserRegistration)
	mux.HandleFunc("POST /api/chirps", userHandler.HandleChirps)
	mux.HandleFunc("GET /api/chirps", userHandler.HandleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", userHandler.HandleGetChirpByID)
	mux.HandleFunc("POST /api/login", userHandler.HandleUserLogin)
	mux.HandleFunc("POST /api/refresh", userHandler.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", userHandler.HandleRevoke)
	mux.HandleFunc("PUT /api/users", userHandler.HandleUserEmailAndPassUpdate)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func okHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

//goose -dir ./sql/schema postgres "Connect String" down
