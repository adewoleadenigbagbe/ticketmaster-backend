package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/labstack/echo/v4"
)

func cinemaRoutes(group *echo.Group) {
	cinemaController := controllers.CinemaController{}
	group.POST("cinema", cinemaController.CreateCinema)
}
