package repositories

import (
	"database/sql"
	"errors"

	"series-tracker/internal/models"
)

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
	query := `DELETE FROM series WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("mmgvo")
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
		return nil, err
	}
	defer rows.Close()

	// Scan results into Serie & append to Series slice
	for rows.Next() {
		var s models.Serie
		if err := rows.Scan(&s.ID, &s.Title, &s.Ranking, &s.Status, &s.CurrentEpisode, &s.TotalEpisodes); err != nil {
			return nil, err
		}
		series = append(series, s)
	}

	return series, nil
}

// CreateNewSeries inserts a new series into the database.
func (r *seriesRepository) CreateNewSerie(s models.Serie) (*models.Serie, error) {
	// Build query
	query := `INSERT INTO series (title, ranking, status, current_episode, total_episodes)
            VALUES ($1, $2, $3, $4, $5)`

	// Execute the query
	result, err := r.db.Exec(query, s.Title, s.Ranking, s.Status, s.CurrentEpisode, s.TotalEpisodes)
	if err != nil {
		return nil, err
	}

	// Update input struct's ID to match the DB
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	s.ID = int(id)

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
		return nil, err
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
		return nil, err
	}

	// Check rows affected to see if update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("aaa")
	}

	return &s, nil
}
