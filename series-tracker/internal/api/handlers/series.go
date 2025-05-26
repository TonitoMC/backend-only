package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	domainErrors "series-tracker/internal/errors"

	"series-tracker/internal/models"
	"series-tracker/internal/services"

	"github.com/labstack/echo/v4"
)

// This file contains all of the handlers for the API, each one of them is commented thoroughly as
// to be able to generate swagger documentation based on comments. The objective of these functions
// is to take information from the request & pass it on to the respective service's function to execute
// operations. As a note, these functions are able to return errors & have them be returned in a friendly
// format to the user by implementing a custom error handler & middleware which takes care of the logging.

// SeriesHandler holds all the dependencies for the series handler
type SeriesHandler struct {
	service services.SeriesService
}

// NewSeriesHandler returns a new SeriesHandler with the given dependencies
func NewSeriesHandler(service services.SeriesService) *SeriesHandler {
	return &SeriesHandler{
		service: service,
	}
}

// GetAllSeries godoc
// @Summary 			Retrieve all series
// @Description 	Get a list of all series in the database
// @Tags 					series
// @Accept 				json
// @Produce 			json
// @Success 			200 	{array} 		models.Serie 				 "Success, array of all fetched series"
// @Failure 			500 	{object} 		models.ErrorResponse "Internal Server Error, unspecified internal errors"
// @Router 				/api/series 		 	[get]
func (h *SeriesHandler) GetAllSeries(c echo.Context) error {
	// Get serie via service
	seriesList, err := h.service.GetAllSeries()
	if err != nil {
		return err
	}

	// Return fetched list of series
	return c.JSON(http.StatusOK, seriesList)
}

// GetSerie godoc
// @Summary 			Retrieve a series by ID
// @Description 	Get details of a series using the provided ID
// @Tags 					series
// @Accept 				json
// @Produce 			json
// @Param 				id 		path 			int 		true 		"Series ID"
// @Success 			200 	{object} 		models.Serie 				 "Success, fetched series"
// @Failure 			400 	{object} 		models.ErrorResponse "Invalid input, eg. non-integer ID"
// @Failure 			404 	{object} 		models.ErrorResponse "Series not found, eg. an ID that doesn't exist in the database"
// @Failure 			500 	{object} 		models.ErrorResponse "Internal server error, unspecified internal errors"
// @Router 				/api/series/{id} 	[get]
func (h *SeriesHandler) GetSerie(c echo.Context) error {
	// Get URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("handler GetSerie: invalid id: %w", domainErrors.ErrInvalidInput)
	}

	// Get serie via service
	serie, err := h.service.GetSerieByID(id)
	if err != nil {
		return err
	}

	// Return fetched serie
	return c.JSON(http.StatusOK, serie)
}

// UpdateSerie godoc
// @Summary 			Update an existing series
// @Description 	Updates details of an existing series
// @Tags 					series
// @Accept 				json
// @Produce 			json
// @Param 				id 		path 			int 		true 		"Series ID"
// @Success 			200 	{object} 		models.Serie 				 "Success, updated series"
// @Failure 			400 	{object} 		models.ErrorResponse "Invalid input, eg. non-integer ID or invalid input fields"
// @Failure 			404 	{object} 		models.ErrorResponse "Series not found, eg. ID that doesn't exist in the database"
// @Failure 			409 	{object}  	models.ErrorResponse "Series conflict, eg. invalid values for episode or negative ranking"
// @Failure 			500 	{object} 		models.ErrorResponse "Internal server error, unspecified internal errors"
// @Router 				/api/series/{id} 	[put]
func (h *SeriesHandler) UpdateSerie(c echo.Context) error {
	// Get URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("handler UpdateSerie: invalid id: %w", domainErrors.ErrInvalidInput)
	}

	// Bind and validate request body
	var serie models.Serie
	if err := c.Bind(&serie); err != nil {
		return fmt.Errorf("handler UpdateSerie: failed to unwrap json: %w", domainErrors.ErrInvalidInput)
	}

	// Insert ID into struct
	serie.ID = id

	// Update series via service
	updatedSeries, err := h.service.UpdateSerie(serie)
	if err != nil {
		return err
	}

	// Return OK & updated data
	return c.JSON(http.StatusOK, updatedSeries)
}

// CreateSerie godoc
// @Summary      Create a new series
// @Description  Inserts a new series into the database, make sure the series object includes all the necessary fields.
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        body  body      models.Serie  true  "Series info"
// @Success      201   {object}  models.Serie "Created, newly created series"
// @Failure      400   {object}  models.ErrorResponse "Invalid input, eg. invalid input fields"
// @Failure 		 409 	 {object}  models.ErrorResponse "Series conflict, eg. negative ranking or duplicate title"
// @Failure      500   {object}  models.ErrorResponse "Internal server error, unspecified internal errors"
// @Router       /api/series [post]
func (h *SeriesHandler) CreateSerie(c echo.Context) error {
	// Bind and validate request body
	var serie models.Serie
	if err := c.Bind(&serie); err != nil {
		return fmt.Errorf("handler CreateSerie: failed to bind json: %w", domainErrors.ErrInvalidInput)
	}

	// Create series via service
	createdSerie, err := h.service.CreateSerie(serie)
	if err != nil {
		return err
	}

	// Returned created serie
	return c.JSON(http.StatusCreated, createdSerie)
}

// DeleteSerie 	 godoc
// @Summary      Remove an existing series from the database
// @Description  Inserts a new series into the database, make sure the series object includes all the necessary fields.
// @Tags         series
// @Accept       json
// @Produce      json
// @Param 				id 		path 			int 		true 		"Series ID"
// @Success      204   "No Content"
// @Failure      400   {object}  models.ErrorResponse "Invalid input, eg. non-integer ID"
// @Failure  		 404 	 {object}  models.ErrorResponse "Series not found, eg. ID that doesn't exist in the databas"
// @Failure      500   {object}  models.ErrorResponse "Internal server error, eg. database error"
// @Router       /api/series/{id} [delete]
func (h *SeriesHandler) DeleteSerie(c echo.Context) error {
	// Get URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("handler DeleteSerie: invalid id: %w", domainErrors.ErrInvalidInput)
	}

	// Delete serie via service
	if err := h.service.DeleteSerie(id); err != nil {
		return err
	}

	// Return no content response
	return c.NoContent(http.StatusNoContent)
}

// UpdateSerieStatus godoc
// @Summary      Update series status
// @Description  Updates the status of the series with the specified ID.
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        id      path int              true "Series ID"
// @Success      200     {object} models.Serie  "Success, updated series"
// @Failure      400     {object} models.ErrorResponse "Invalid input, eg. non-integer ID or invalid status"
// @Failure      404     {object} models.ErrorResponse "Series not found, eg. ID that doesn't exist in the database"
// @Failure      500     {object} models.ErrorResponse "Internal server error, unspecified internal errors"
// @Router       /api/series/{id}/status [patch]
func (h *SeriesHandler) UpdateSerieStatus(c echo.Context) error {
	// Get ID URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("handler updateSerieStatus: invalid id: %w", domainErrors.ErrInvalidInput)
	}

	// Unpack status, binded request to map since declaring struct / model
	// is a bit overkill
	var reqMap map[string]string
	if err := c.Bind(&reqMap); err != nil {
		return fmt.Errorf("handler UpdateSerieStatus: failed to bind json: %w", domainErrors.ErrInvalidInput)
	}

	newStatus, exists := reqMap["status"]
	if !exists {
		return fmt.Errorf("handler UpdateSerieStatus: failed to bind json: %w", domainErrors.ErrInvalidInput)
	}

	// Update serie status via service
	updatedSeries, err := h.service.UpdateSerieStatus(id, newStatus)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error"})
	}

	// Return updated series
	return c.JSON(http.StatusOK, updatedSeries)
}

// IncrementEpisode godoc
// @Summary      Advance series episode count
// @Description  Increments the current episode number of a series by one
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        id      path int              true "Series ID"
// @Success      200     {object} models.Serie  "Success, updated series"
// @Failure      400     {object} models.ErrorResponse "Invalid input, eg. non-integer ID"
// @Failure      404     {object} models.ErrorResponse "Series not found, eg. ID that doesn't exist in the database"
// @Failure 		 409 		 {object} models.ErrorResponse "Series conflict, eg. trying to increment episode past last one"
// @Failure      500     {object} models.ErrorResponse "Internal server error, unspecified internal errors"
// @Router       /api/series/{id}/episode [patch]
func (h *SeriesHandler) IncrementEpisode(c echo.Context) error {
	// Get ID URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("handler IncrementEpisode: invalid id: %w", domainErrors.ErrInvalidInput)
	}

	// Update series via service
	updatedSeries, err := h.service.IncrementSerieEpisode(id)
	if err != nil {
		return err
	}

	// Return updated series
	return c.JSON(http.StatusOK, updatedSeries)
}

// UpvoteSerie godoc
// @Summary      Increase series score
// @Description  Increases the rating (score) of the series with the specified ID by one vote.
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        id   path      int   true  "Series ID"
// @Success      200  {object}  models.Serie "Success, updated series"
// @Failure      400  {object}  models.ErrorResponse "Invalid input, eg. non-integer ID"
// @Failure      404  {object}  models.ErrorResponse "Series not found, eg. ID that doesn't exist in the database"
// @Failure      500  {object}  models.ErrorResponse "Internal server error, unspecified internal errors"
// @Router       /api/series/{id}/upvote [patch]
func (h *SeriesHandler) UpvoteSerie(c echo.Context) error {
	// Extract the series ID from the URL parameter.
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("handler UpvoteSerie: invalid id: %w", domainErrors.ErrInvalidInput)
	}

	// Call the service layer to upvote the series.
	updatedSerie, err := h.service.UpvoteSerie(id)
	if err != nil {
		return err
	}

	// Return the updated series.
	return c.JSON(http.StatusOK, updatedSerie)
}

// DownvoteSerie godoc
// @Summary      Decrease series score
// @Description  Decreases the rating (score) of the series with the specified ID by one vote.
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        id   path      int   true  "Series ID"
// @Success      200  {object}  models.Serie "Success, updated series"
// @Failure      400  {object}  models.ErrorResponse "Invalid input, eg. non-integer ID"
// @Failure      404  {object}  models.ErrorResponse "Series not found, eg. ID that doesn't exist in the database"
// @Failure      409  {object}  models.ErrorResponse "Series conflict, eg. downvoting past 0"
// @Failure      500  {object}  models.ErrorResponse "Internal server error, unspecified internal errors"
// @Router       /api/series/{id}/downvote [patch]
func (h *SeriesHandler) DownvoteSerie(c echo.Context) error {
	// Extract the series ID from the URL parameter.
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fmt.Errorf("handler DownvoteSerie: invalid id: %w", domainErrors.ErrInvalidInput)
	}

	// Call the service layer to downvote the series.
	updatedSerie, err := h.service.DownvoteSerie(id)
	if err != nil {
		return err
	}

	// Return the updated series.
	return c.JSON(http.StatusOK, updatedSerie)
}
