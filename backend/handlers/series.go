package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend-only/models"

	"github.com/go-chi/chi/v5"
)

// Handler function for GET /api/series
func GetSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the DB
		rows, err := db.Query("SELECT id, title, ranking, status, current_episode, total_episodes FROM series")
		if err != nil {
			http.Error(w, "failed to fetch series", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Create series models from the response
		var series []models.Serie
		for rows.Next() {
			var s models.Serie
			if err := rows.Scan(&s.ID, &s.Title, &s.Ranking, &s.Status, &s.CurrentEpisode, &s.TotalEpisodes); err != nil {
				http.Error(w, "failed to scan row", http.StatusInternalServerError)
				return
			}
			series = append(series, s)
		}

		// Encode the series
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(series)
	}
}

// Handler function for POST /api/series
func PostSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serie models.Serie

		// Decode JSON
		if err := json.NewDecoder(r.Body).Decode(&serie); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Insert into DB
		query := `INSERT INTO series (title, ranking, status, current_episode, total_episodes)
		          VALUES ($1, $2, $3, $4, $5)`
		result, err := db.Exec(query, serie.Title, serie.Ranking, serie.Status, serie.CurrentEpisode, serie.TotalEpisodes)
		if err != nil {
			http.Error(w, "Failed to insert series", http.StatusInternalServerError)
			return
		}

		// Get the new ID and return it in the response
		id, _ := result.LastInsertId()
		serie.ID = int(id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(serie)
	}
}

// Handler function GET /api/series/{id}
func GetSeriesByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get series ID parameter from request
		id := chi.URLParam(r, "seriesID")

		var serie models.Serie

		// Select from DB
		query := `SELECT id, title, ranking, status, current_episode, total_episodes FROM series WHERE id = $1`
		err := db.QueryRow(query, id).Scan(
			&serie.ID,
			&serie.Title,
			&serie.Ranking,
			&serie.Status,
			&serie.CurrentEpisode,
			&serie.TotalEpisodes,
		)

		if err == sql.ErrNoRows {
			http.Error(w, "Series not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to get series", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serie)
	}
}

// Handler function PUT /api/series/{id}
func PutSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get series ID parameter from request
		id := chi.URLParam(r, "seriesID")

		var serie models.Serie

		// Decode series from request
		if err := json.NewDecoder(r.Body).Decode(&serie); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Update series in DB
		query := `UPDATE series SET title = $1, ranking = $2, status = $3, current_episode = $4, total_episodes = $5
							WHERE id = $6`
		if _, err := db.Exec(query, serie.Title, serie.Ranking, serie.Status, serie.CurrentEpisode, serie.TotalEpisodes, id); err != nil {
			http.Error(w, "Failed to update series", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serie)
	}
}

// Handler function for DELETE /api/series/{id}
func DeleteSeriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "seriesID")

		// Delete series in DB
		query := `DELETE FROM series WHERE id = $1`
		if _, err := db.Exec(query, id); err != nil {
			http.Error(w, "Failed to delete series", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
