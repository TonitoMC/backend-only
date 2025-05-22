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

func main() {
	dbConn, err := database.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("FATAL: No db")
	}
	defer dbConn.Close()

	seriesRepo := repositories.NewSeriesRepository(dbConn)
	seriesService := services.NewSeriesService(seriesRepo)
	seriesHandler := handlers.NewSeriesHandler(seriesService)

	routerConfig := &api.RouterConfig{
		SeriesHandler: seriesHandler,
	}

	e := echo.New()
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

	e.Logger.Print("hola")
}
