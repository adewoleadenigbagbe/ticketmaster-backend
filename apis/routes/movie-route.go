package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/apis/controllers"
	"github.com/Wolechacho/ticketmaster-backend/apis/core"
	"github.com/labstack/echo/v4"
)

func movieRoutes(app *core.BaseApp, group *echo.Group) {
	movieController := controllers.MovieController{App: app}
	group.GET("movies", movieController.GetMoviesHandler)
	group.GET("movies/:id", movieController.GetMovieByIdHandler)
	group.GET("movies/search", movieController.SearchMovieHandler)
}
