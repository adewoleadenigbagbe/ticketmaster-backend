package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	middlewares "github.com/Wolechacho/ticketmaster-backend/middleware"
	"github.com/labstack/echo/v4"
)

func showRoutes(app *core.BaseApp, group *echo.Group) {
	showController := controllers.ShowController{App: app}
	group.POST("shows", showController.CreateShowHandler, middlewares.AuthorizeAdmin)
	//group.GET("shows/user-location", showController.GetShowsByUserLocationHandler, middlewares.AuthorizeUser)
	group.GET("shows/user-location", showController.GetShowsByUserLocationHandler)
}
