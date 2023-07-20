package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/labstack/echo/v4"
)

func movieRoutes(group *echo.Group) {
	movieController := controllers.MovieController{}
	group.GET("movies", movieController.GetMovies)
	group.GET("movies/:id", movieController.GetMovieById)
}
