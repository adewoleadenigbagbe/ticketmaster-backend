package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/labstack/echo/v4"
)

func movieRoutes(e *echo.Echo) {
	movieController := controllers.MovieController{}
	e.GET("/api/v1/movies", movieController.GetMovies)
}
