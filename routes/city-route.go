package routes

import (
	"fmt"

	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

func cityRoutes(app *core.BaseApp, group *echo.Group) {
	cityController := controllers.CityController{App: app}
	group.GET("city/:id", cityController.GetCityByIdHandler)
	fmt.Println(group)
}
