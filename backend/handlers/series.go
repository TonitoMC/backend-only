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

// Handler function for Patch /api/series/:id/status
func PatchStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "seriesID")
		var status models.Status

		// Unpack status
		if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Update in DB
		query := `UPDATE series SET status = $1 WHERE id = $2`
		if _, err := db.Exec(query, status.Status, id); err != nil {
			http.Error(w, "Failed to update status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}

func PatchEpisodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "seriesID")

		// Update in DB
		query := `UPDATE series SET current_episode = current_episode + 1 WHERE id = $1`
		if _, err := db.Exec(query, id); err != nil {
			http.Error(w, "Failed to increment episode", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		// Frontend actually does nothing with the response except read it, we HAVE
		// to send a JSON but we can but whatever we want in the content
		json.NewEncoder(w).Encode(map[string]string{"tralalero": "tralala"})
	}
}

func PatchUpvoteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "seriesID")

		// Update in DB
		query := `UPDATE series SET ranking = ranking + 1 WHERE id = $1`
		if _, err := db.Exec(query, id); err != nil {
			http.Error(w, "Failed to upvote episode", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		// Again, frontend expects a JSON response but does nothing with it.
		json.NewEncoder(w).Encode(map[string]string{"brr brr": "patapim"})
	}
}

func PatchDownvoteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "seriesID")

		query := `UPDATE series SET ranking = ranking - 1 WHERE id = $1`
		if _, err := db.Exec(query, id); err != nil {
			http.Error(w, "Failed to downvote episode", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		// Frontend expects JSON but does nothing with it
		json.NewEncoder(w).Encode(map[string]string{"bombardiro": "crocodilo"})
	}
}
