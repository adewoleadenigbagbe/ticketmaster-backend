package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

func showRoutes(app *core.BaseApp, group *echo.Group) {
	showController := controllers.ShowController{App: app}
	group.POST("shows", showController.CreateShowHandler)
	group.GET("shows/user-location", showController.GetShowsByUserLocationHandler)
}
