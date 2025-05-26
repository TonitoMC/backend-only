package api

import (
	"series-tracker/internal/api/handlers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// This file configures the routing for our application, it links specific HTTP routes
// to their respective handler

// Struct for the router
type RouterConfig struct {
	SeriesHandler *handlers.SeriesHandler
}

// Route setup
func SetupRoutes(e *echo.Echo, config *RouterConfig) {
	// GET /api/series for getting all series in the database
	e.GET("api/series", config.SeriesHandler.GetAllSeries)
	// GET /api/series/:id for getting a specific series by ID
	e.GET("api/series/:id", config.SeriesHandler.GetSerie)
	// PUT /api/series/:id for updating a specific series
	e.PUT("api/series/:id", config.SeriesHandler.UpdateSerie)
	// POST /api/series for creating a new series
	e.POST("api/series", config.SeriesHandler.CreateSerie)
	// DELETE /api/series/:id for deleting a series
	e.DELETE("api/series/:id", config.SeriesHandler.DeleteSerie)
	// PATCH /api/series/:id/status for updating a series' status
	e.PATCH("api/series/:id/status", config.SeriesHandler.UpdateSerieStatus)
	// PATCH /api/series/:id/episode for incrementing the series' current episode by 1
	e.PATCH("api/series/:id/episode", config.SeriesHandler.IncrementEpisode)
	// PATCH /api/series/:id/upvote to increment the series' ranking score by 1
	e.PATCH("api/series/:id/upvote", config.SeriesHandler.UpvoteSerie)
	// PATCH /api/series/:id/downvote to decrease the series' ranking score by 1
	e.PATCH("api/series/:id/downvote", config.SeriesHandler.DownvoteSerie)
	// Handler for swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
