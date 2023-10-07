package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/labstack/echo/v4"
)

func userRoutes(group *echo.Group) {
	userController := controllers.UserController{}
	group.GET("users", userController.CreateUser)
}
