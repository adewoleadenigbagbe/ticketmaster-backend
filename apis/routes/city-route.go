package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	middlewares "github.com/Wolechacho/ticketmaster-backend/middleware"
	"github.com/labstack/echo/v4"
)

func cityRoutes(app *core.BaseApp, group *echo.Group) {
	cityController := controllers.CityController{App: app}
	group.GET("city/:id", cityController.GetCityByIdHandler, middlewares.AuthorizeAdmin)
}
