package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// dbConnect function manages the database connection. Connects, pings
// and returns DB instance on success.
func dbConnect() (*sql.DB, error) {
	// Load environment variables for connection
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Create connection string with environment variables
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Connect to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Ping the database
	err = db.Ping()
	if err != nil {
		db.Close() // Close if unsuccessful
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Return the DB instance
	return db, nil
}
