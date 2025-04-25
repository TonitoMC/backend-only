package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend-only/models"
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
