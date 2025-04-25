// Package main implements a REST API for managing series using Go and Chi,
// it's developed as a backend for an existing frontend found at
// https://github.com/denn1s/series-tracker
package main

import (
	"log"
	"net/http"

	"backend-only/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

// Main function handles database & routing setup
func main() {
	// Connect to the DB
	db, err := dbConnect()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	log.Println("Successful DB connection")

	// Router setup
	r := chi.NewRouter()

	// CORS setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:80"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// GET /api/series handler
	r.Get("/api/series", handlers.GetSeriesHandler(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}
