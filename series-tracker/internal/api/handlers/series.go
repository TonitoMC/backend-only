package handlers

import (
	"net/http"
	"strconv"

	"series-tracker/internal/models"
	"series-tracker/internal/services"

	"github.com/labstack/echo/v4"
)

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
// @Success 			200 	{array} 		models.Serie
// @Failure 			400 	{object} 		map[string]string
// @Failure 			500 	{object} 		map[string]string
// @Router 				/api/series 		 	[get]
func (h *SeriesHandler) GetAllSeries(c echo.Context) error {
	// Get serie via service
	seriesList, err := h.service.GetAllSeries()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, seriesList)
}

// GetSerie godoc
// @Summary 			Retrieve a series by ID
// @Description 	Get details of a series using the provided ID
// @Tags 					series
// @Accept 				json
// @Produce 			json
// @Param 				id 		path 			int 		true 		"Series ID"
// @Success 			200 	{object} 		models.Serie
// @Failure 			400 	{object} 		map[string]string
// @Failure 			500 	{object} 		map[string]string
// @Router 				/api/series/{id} 	[get]
func (h *SeriesHandler) GetSerie(c echo.Context) error {
	// Get URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	// Get serie via service
	serie, err := h.service.GetSerieByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
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
// @Success 			200 	{object} 		models.Serie
// @Failure 			400 	{object} 		map[string]string
// @Failure 			500 	{object} 		map[string]string
// @Router 				/api/series/{id} 	[put]
func (h *SeriesHandler) UpdateSerie(c echo.Context) error {
	// Get URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	// Bind and validate request body
	var serie models.Serie
	if err := c.Bind(&serie); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	// Insert ID into struct
	serie.ID = id

	// Update series via service
	updatedSeries, err := h.service.CreateSerie(serie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
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
// @Success      201   {object}  models.Serie "Newly created series"
// @Failure      400   {object}  map[string]string "Bad request, e.g, invalid input"
// @Failure      500   {object}  map[string]string "Internal Server Error, e.g, database error"
// @Router       /api/series [post]
func (h *SeriesHandler) CreateSerie(c echo.Context) error {
	// Bind and validate request body
	var serie models.Serie
	if err := c.Bind(&serie); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	// Create series via service
	createdSerie, err := h.service.CreateSerie(serie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create series"})
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
// @Success      204   "No content"
// @Failure      400   {object}  map[string]string "Bad request, e.g, invalid input"
// @Failure      500   {object}  map[string]string "Internal Server Error, e.g, database error"
// @Router       /api/series/{id} [delete]
func (h *SeriesHandler) DeleteSerie(c echo.Context) error {
	// Get URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := h.service.DeleteSerie(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error"})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateSerieStatus godoc
// @Summary      Update series status
// @Description  Updates the status of the series with the specified ID.
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        id      path int              true "Series ID"
// @Success      200     {object} models.Serie  "Successfully updated series status"
// @Failure      400     {object} map[string]string "Invalid input or status value"
// @Failure      404     {object} map[string]string "Series not found"
// @Failure      500     {object} map[string]string "Internal server error"
// @Router       /api/series/{id}/status [patch]
func (h *SeriesHandler) UpdateSerieStatus(c echo.Context) error {
	// Get ID URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	// Unpack status, binded request to map since declaring struct / model
	// is a bit overkill
	var reqMap map[string]string
	if err := c.Bind(&reqMap); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "error"})
	}

	newStatus, exists := reqMap["status"]
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing status"})
	}

	updatedSeries, err := h.service.UpdateSerieStatus(id, newStatus)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error"})
	}
	return c.JSON(http.StatusOK, updatedSeries)
}

// IncrementEpisode godoc
// @Summary      Advance series episode count
// @Description  Increments the current episode number of a series by one
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        id      path int              true "Series ID"
// @Success      200     {object} models.Serie  "Successfully updated series status"
// @Failure      400     {object} map[string]string "Invalid input or status value"
// @Failure      404     {object} map[string]string "Series not found"
// @Failure      500     {object} map[string]string "Internal server error"
// @Router       /api/series/{id}/episode [patch]
func (h *SeriesHandler) IncrementEpisode(c echo.Context) error {
	// Get ID URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	updatedSeries, err := h.service.IncrementSerieEpisode(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error"})
	}
	return c.JSON(http.StatusOK, updatedSeries)
}

// UpvoteSerie godoc
// @Summary      Increase series score
// @Description  Increases the rating (score) of the series with the specified ID by one vote.
// @Tags         series
// @Accept       json
// @Produce      json
// @Param        id   path      int   true  "Series ID"
// @Success      200  {object}  models.Serie "Series successfully upvoted"
// @Failure      400  {object}  map[string]string "Invalid series ID"
// @Failure      404  {object}  map[string]string "Series not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/series/{id}/upvote [patch]
func (h *SeriesHandler) UpvoteSerie(c echo.Context) error {
	// Extract the series ID from the URL parameter.
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	// Call the service layer to upvote the series.
	updatedSerie, err := h.service.UpvoteSerie(id)
	if err != nil {
		// Optionally handle not found errors separately.
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to upvote series"})
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
// @Success      200  {object}  models.Serie "Series successfully downvoted"
// @Failure      400  {object}  map[string]string "Invalid series ID"
// @Failure      404  {object}  map[string]string "Series not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/series/{id}/downvote [patch]
func (h *SeriesHandler) DownvoteSerie(c echo.Context) error {
	// Extract the series ID from the URL parameter.
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	// Call the service layer to downvote the series.
	updatedSerie, err := h.service.DownvoteSerie(id)
	if err != nil {
		// Optionally handle not found errors separately.
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to downvote series"})
	}

	// Return the updated series.
	return c.JSON(http.StatusOK, updatedSerie)
}
