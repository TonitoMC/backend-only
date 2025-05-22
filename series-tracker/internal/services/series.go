package services

import (
	"errors"

	"series-tracker/internal/models"
	"series-tracker/internal/repositories"
)

// Set of valid statuses for the serie, used in various
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
		return nil, err
	}

	// Return result
	return serie, nil
}

// CreateSerie creates a new series
func (s *seriesService) CreateSerie(serie models.Serie) (*models.Serie, error) {
	// Create series in the repository
	createdSerie, err := s.seriesRepo.CreateNewSerie(serie)
	if err != nil {
		return nil, err
	}

	return createdSerie, nil
}

// UpdateSerie updates a series with all values detailed in the struct based on the ID
func (s *seriesService) UpdateSerie(serie models.Serie) (*models.Serie, error) {
	updatedSerie, err := s.seriesRepo.UpdateSerie(serie)
	if err != nil {
		return nil, err
	}

	return updatedSerie, nil
}

// DeleteSerie deletes a serie by its ID
func (s *seriesService) DeleteSerie(id int) error {
	if err := s.seriesRepo.DeleteSerie(id); err != nil {
		return err
	}
	return nil
}

// UpdateSerieStatus updates the status of a serie by updating the information & updating via repository
func (s *seriesService) UpdateSerieStatus(id int, status string) (*models.Serie, error) {
	// Check validity of given status
	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, err
	}

	// Set the status to the updated one
	serie.Status = status

	// Call repository to update
	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, err
	}
	return updatedSerie, nil
}

// UpvoteSerie updates the ranking of a serie incrementing by one
func (s *seriesService) UpvoteSerie(id int) (*models.Serie, error) {
	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, err
	}

	// Increment ranking score by 1
	serie.Ranking += 1

	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, err
	}
	return updatedSerie, nil
}

// DownvoteSerie updates the ranking of a serie decreasing by one
func (s *seriesService) DownvoteSerie(id int) (*models.Serie, error) {
	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, err
	}

	if serie.Ranking <= 0 {
		return nil, errors.New("series can't be downvoted further")
	}

	// Decrease value by one
	serie.Ranking -= 1

	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, err
	}
	return updatedSerie, nil
}

// IncrementSerieEpisode incrementes the current episode by one
func (s *seriesService) IncrementSerieEpisode(id int) (*models.Serie, error) {
	// Get series information from repository
	serie, err := s.seriesRepo.GetSerieByID(id)
	if err != nil {
		return nil, err
	}

	if serie.CurrentEpisode >= serie.TotalEpisodes {
		return nil, errors.New("series hit max episodes")
	}

	// Increment value by one
	serie.CurrentEpisode += 1

	updatedSerie, err := s.seriesRepo.UpdateSerie(*serie)
	if err != nil {
		return nil, err
	}
	return updatedSerie, nil
}
