package repositories

import (
	"database/sql"
	"fmt"

	domainErrors "series-tracker/internal/errors"

	"series-tracker/internal/models"
)

// This file defines and implements the series repository, this is the layer that's in charge
// of executing basic operations and communicating with the database. The methods implemented
// are kept quite simple, as this avoids creating an overly specific repository with non-reusable
// functions.

// SeriesRepository defines all the methods to be implemented for series data access
type SeriesRepository interface {
	// GetAllSeries returns a list of all series from the database
	GetAllSeries() ([]models.Serie, error)
	// CreateNewSerie inserts a new series into the database
	CreateNewSerie(models.Serie) (*models.Serie, error)
	// GetSerieByID finds a series by its ID in the database
	GetSerieByID(id int) (*models.Serie, error)
	// UpdateSerie updates a series with all values detailed in a Serie struct based on its ID
	UpdateSerie(models.Serie) (*models.Serie, error)
	// DeleteSerie deletes a series by its ID
	DeleteSerie(id int) error
}

// seriesRepository holds all the dependencies for the repository
type seriesRepository struct {
	db *sql.DB
}

// NewSeriesRepository creates a new SeriesRepository with the given DB connection
func NewSeriesRepository(dbConn *sql.DB) SeriesRepository {
	return &seriesRepository{
		db: dbConn,
	}
}

// DeleteSerie deletes a series by its ID.
func (r *seriesRepository) DeleteSerie(id int) error {
	// Build the query
	query := `DELETE FROM series WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("repository DeleteSerie: failed to execute query: %w", err)
	}

	// Verify affected rows to make sure effect took place
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository DeleteSerie: failed to verify rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domainErrors.ErrSeriesNotFound
	}

	return nil
}

// GetAllSerie returns a list of all series from the database.
func (r *seriesRepository) GetAllSeries() ([]models.Serie, error) {
	// Create return slice
	series := []models.Serie{}

	// Query the DB
	rows, err := r.db.Query("SELECT id, title, ranking, status, current_episode, total_episodes FROM series")
	if err != nil {
		return nil, fmt.Errorf("repository GetAllSerie: failed to execute query: %w", err)
	}
	defer rows.Close()

	// Scan results into Serie & append to Series slice
	for rows.Next() {
		var s models.Serie
		if err := rows.Scan(&s.ID, &s.Title, &s.Ranking, &s.Status, &s.CurrentEpisode, &s.TotalEpisodes); err != nil {
			return nil, fmt.Errorf("repository GetAllSeries: failed to scan row: %w", err)
		}
		series = append(series, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository GetAllSeries: row iteration error: %w", err)
	}

	return series, nil
}

// CreateNewSeries inserts a new series into the database.
func (r *seriesRepository) CreateNewSerie(s models.Serie) (*models.Serie, error) {
	// Build query
	query := `INSERT INTO series (title, ranking, status, current_episode, total_episodes)
            VALUES ($1, $2, $3, $4, $5)
  					RETURNING id`

	// Execute the query
	err := r.db.QueryRow(query, s.Title, s.Ranking, s.Status, s.CurrentEpisode, s.TotalEpisodes).Scan(&s.ID)
	if err != nil {
		return nil, fmt.Errorf("repository CreateNewSeries: failed to execute query: %w", domainErrors.ErrSeriesConflict)
	}

	return &s, nil
}

// GetSerieByID finds a Serie by its ID in the database.
func (r *seriesRepository) GetSerieByID(id int) (*models.Serie, error) {
	// Create series struct for response
	var serie models.Serie

	// Build the query
	query := `SELECT id, title, ranking, status, current_episode, total_episodes
            FROM series
            WHERE id = $1`

	// Execute the query & scan into Serie struct
	if err := r.db.QueryRow(query, id).Scan(
		&serie.ID,
		&serie.Title,
		&serie.Ranking,
		&serie.Status,
		&serie.CurrentEpisode,
		&serie.TotalEpisodes,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainErrors.ErrSeriesNotFound
		}
		return nil, fmt.Errorf("repository GetSerieByID: failed to query / scan: %w", err)
	}

	return &serie, nil
}

// UpdateSerie updates a serie with all values detailed in a Serie struct based on its ID
func (r *seriesRepository) UpdateSerie(s models.Serie) (*models.Serie, error) {
	// Build the query
	query := `UPDATE series 
            SET title = $1, ranking = $2, status = $3, current_episode = $4, total_episodes = $5
            WHERE id = $6`

	// Execute the query
	result, err := r.db.Exec(query, s.Title, s.Ranking, s.Status, s.CurrentEpisode, s.TotalEpisodes, s.ID)
	if err != nil {
		return nil, fmt.Errorf("repository UpdateSerie: failed to execute update query: %w", err)
	}

	// Check rows affected to see if update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("repository UpdateSerie: failed to verify affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return nil, domainErrors.ErrSeriesNotFound
	}

	return &s, nil
}
