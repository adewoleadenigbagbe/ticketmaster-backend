package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/labstack/echo/v4"
)

func cityRoutes(group *echo.Group) {
	cityController := controllers.CityController{}
	group.GET("city/:id", cityController.GetCityById)
}
