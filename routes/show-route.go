package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/labstack/echo/v4"
)

func showRoutes(group *echo.Group) {
	showController := controllers.ShowController{}
	group.POST("shows", showController.CreateShow)
	group.GET("shows/user-location", showController.GetShowsByUserLocation)
}
