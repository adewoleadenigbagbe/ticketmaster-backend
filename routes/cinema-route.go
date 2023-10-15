package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

func cinemaRoutes(app *core.BaseApp, group *echo.Group) {
	cinemaController := controllers.CinemaController{App: app}
	group.POST("cinemas", cinemaController.CreateCinemaHandler)
	group.POST("cinemas/:id/cinemaHall", cinemaController.CreateCinemaHallHandler)
}