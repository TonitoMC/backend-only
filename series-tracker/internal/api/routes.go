package api

import (
	"series-tracker/internal/api/handlers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type RouterConfig struct {
	SeriesHandler *handlers.SeriesHandler
}

func SetupRoutes(e *echo.Echo, config *RouterConfig) {
	e.GET("api/series", config.SeriesHandler.GetAllSeries)
	e.GET("api/series/:id", config.SeriesHandler.GetSerie)
	e.PUT("api/series/:id", config.SeriesHandler.UpdateSerie)
	e.POST("api/series", config.SeriesHandler.CreateSerie)
	e.DELETE("api/series/:id", config.SeriesHandler.DeleteSerie)
	e.PATCH("api/series/:id/status", config.SeriesHandler.UpdateSerieStatus)
	e.PATCH("api/series/:id/episode", config.SeriesHandler.IncrementEpisode)
	e.PATCH("api/series/:id/upvote", config.SeriesHandler.UpvoteSerie)
	e.PATCH("api/series/:id/downvote", config.SeriesHandler.DownvoteSerie)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
