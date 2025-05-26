package services

import (
	"fmt"

	domainErrors "series-tracker/internal/errors"

	"series-tracker/internal/models"
	"series-tracker/internal/repositories"
)

// This file defines and implements the series service, this is the layer
// where our business logic lies and we coordinate our HTTP handlers with
// the repository for linking requests to our data storage. The methods here
// correspond 1:1 to an HTTP handler

// Set of valid statuses for the serie, used in various points
var validStatuses = map[string]bool{
	"Watching":      true,
	"Plan to Watch": true,
	"Dropped":       true,
	"Completed":     true,
}

// SeriesService defines all the methods to be implemented for series management
type SeriesService interface {
	// GetSerieByID returns a series by its ID
	GetSerieByID(id int) (*models.Serie, error)
	// GetAllSeries returns a list of all series
	GetAllSeries() ([]models.Serie, error)
	// CreateSerie creates a new series
	CreateSerie(models.Serie) (*models.Serie, error)
	// UpdateSerie updates a series with all values detailed in a Serie struct based on its ID
	UpdateSerie(models.Serie) (*models.Serie, error)
	// DeleteSerie deletes a series by its ID
	DeleteSerie(id int) error
	// UpdateSerieStatus updates the status of a series by its ID
	UpdateSerieStatus(id int, status string) (*models.Serie, error)
	// UpvoteSerie increases the ranking score of a series by 1
	UpvoteSerie(id int) (*models.Serie, error)
	// DownvoteSerie decreases the ranking score of a series by 1
	DownvoteSerie(id int) (*models.Serie, error)
	// IncrementSerieEpisode increases the current episode of a series by 1
	IncrementSerieEpisode(id int) (*models.Serie, error)
}

// seriesService holds all the dependencies for the service
type seriesService struct {
	seriesRepo repositories.SeriesRepository
}

// NewSeriesService returns a seriesService with the given dependencies
func NewSeriesService(seriesRepo repositories.SeriesRepository) SeriesService {
	return &seriesService{
		seriesRepo: seriesRepo,
	}
}

// GetAllSeries returns a list of all series
func (s *seriesService) GetAllSeries() ([]models.Serie, error) {
	// Get series by ID from repository
	series, err := s.seriesRepo.GetAllSeries()
	if err != nil {
		return nil, err
	}

	return series, nil
}

// GetSerieByID returns a series by itsd ID
func (s *seriesService) GetSerieByID(id int) (*models.Serie, error) {
	// Get the series from the repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, fmt.Errorf("service GetSerieByID: failed to get series with id %d: %w", id, err)
	}

	// Return result
	return serie, nil
}

// CreateSerie creates a new series
func (s *seriesService) CreateSerie(serie models.Serie) (*models.Serie, error) {
	// Validate the input
	if !validateSerie(serie) {
		return nil, fmt.Errorf("service CreateSerie: failed to validate series: %w", domainErrors.ErrSeriesConflict)
	}
	// Create series in the repository
	createdSerie, err := s.seriesRepo.CreateNewSerie(serie)
	if err != nil {
		return nil, fmt.Errorf("service CreateSerie: failed to create series: %w", err)
	}

	return createdSerie, nil
}

// UpdateSerie updates a series with all values detailed in the struct based on the ID
func (s *seriesService) UpdateSerie(serie models.Serie) (*models.Serie, error) {
	// Validate the input
	if !validateSerie(serie) {
		return nil, fmt.Errorf("service UpdateSerie: failed to validate series: %w", domainErrors.ErrSeriesConflict)
	}

	updatedSerie, err := s.seriesRepo.UpdateSerie(serie)
	if err != nil {
		return nil, fmt.Errorf("service UpdateSerie: failed to update series with id %d: %w", serie.ID, err)
	}

	return updatedSerie, nil
}

// DeleteSerie deletes a serie by its ID
func (s *seriesService) DeleteSerie(id int) error {
	if err := s.seriesRepo.DeleteSerie(id); err != nil {
		return fmt.Errorf("service DeleteSerie: failed to delete series with id %d: %w", id, err)
	}
	return nil
}

// UpdateSerieStatus updates the status of a serie by updating the information & updating via repository
func (s *seriesService) UpdateSerieStatus(id int, status string) (*models.Serie, error) {
	// Check validity of given status
	if !validStatuses[status] {
		return nil, fmt.Errorf("service UpdateSerieStatus: failed to update serie with id %d invalid status %s: %w", id, status, domainErrors.ErrInvalidInput)
	}

	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, fmt.Errorf("service UpdateSerieStatus: failed to update serie with id %d and status %s: %w", id, status, err)
	}

	// Set the status to the updated one
	serie.Status = status

	// Call repository to update
	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, fmt.Errorf("service UpdateSerieStatus: failed to update serie with id %d and status %s: %w", id, status, err)
	}

	return updatedSerie, nil
}

// UpvoteSerie updates the ranking of a serie incrementing by one
func (s *seriesService) UpvoteSerie(id int) (*models.Serie, error) {
	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, fmt.Errorf("service UpvoteSerie: failed to upvote serie with id %d: %w", id, err)
	}

	// Increment ranking score by 1
	serie.Ranking += 1

	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, fmt.Errorf("service UpvoteSerie: failed to upvote serie with id %d: %w", id, err)
	}
	return updatedSerie, nil
}

// DownvoteSerie updates the ranking of a serie decreasing by one
func (s *seriesService) DownvoteSerie(id int) (*models.Serie, error) {
	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, fmt.Errorf("service DownvoteSerie: failed to downvote serie with id %d: %w", id, err)
	}

	if serie.Ranking <= 0 {
		return nil, fmt.Errorf("service DownvoteSerie: failed to downvote serie with id %d - ranking too low: %w", id, domainErrors.ErrSeriesConflict)
	}

	// Decrease value by one
	serie.Ranking -= 1

	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, fmt.Errorf("service DownvoteSerie: failed to downvote serie with id %d: %w", id, err)
	}
	return updatedSerie, nil
}

// IncrementSerieEpisode incrementes the current episode by one
func (s *seriesService) IncrementSerieEpisode(id int) (*models.Serie, error) {
	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, fmt.Errorf("service IncrementEpisode: failed to increment serie episode for id %d: %w", id, err)
	}

	if serie.CurrentEpisode >= serie.TotalEpisodes {
		return nil, fmt.Errorf("service IncrementEpisode: failed to increment serie episode for id %d - reached last episode: %w", id, domainErrors.ErrSeriesConflict)
	}

	// Increment value by one
	serie.CurrentEpisode += 1

	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, fmt.Errorf("service IncrementEpisode: failed to increment serie episode for id %d: %w", id, err)
	}
	return updatedSerie, nil
}

// Utility function to validate serie validity for inserts and udpates
func validateSerie(serie models.Serie) bool {
	if serie.Title == "" || serie.Ranking <= 0 || serie.CurrentEpisode <= 0 || serie.TotalEpisodes <= 0 || !validStatuses[serie.Status] {
		return false
	}
	return true
}
