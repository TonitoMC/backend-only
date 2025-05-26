package main

import (
	"log"

	"series-tracker/internal/api"
	"series-tracker/internal/api/handlers"
	"series-tracker/internal/database"
	"series-tracker/internal/repositories"
	"series-tracker/internal/services"

	_ "series-tracker/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// This project is a REST API for managing TV series. It provides endpoints for creating,
// reading, updating and deleting series data, as well as additional operations like upvoting,
// downvoting and episode tracking. Full instructions on how to utilize it can be found inside
// of the swagger documentation as well as the README.md in the root of the repository. This project
// implements a layered architecture with repositories, services and handlers for a clear separation
// of concerns, where each of them can be found in their respective modules.

func main() {
	// Initialize the database connection
	dbConn, err := database.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("FATAL: No db")
	}
	defer dbConn.Close()

	// Initialize the repository, service & handler for series
	seriesRepo := repositories.NewSeriesRepository(dbConn)
	seriesService := services.NewSeriesService(seriesRepo)
	seriesHandler := handlers.NewSeriesHandler(seriesService)

	// Configure the router
	routerConfig := &api.RouterConfig{
		SeriesHandler: seriesHandler,
	}

	// Echo setup
	e := echo.New()
	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS is set up to be able to receive requests from localhost / localhost:80,
	// this is the default port nginx is set up to run on & just making it more
	// accessible
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost",
			"http://localhost:80",
		},
		AllowMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	api.SetupRoutes(e, routerConfig)

	e.Logger.Fatal(e.Start(":8080"))
}
