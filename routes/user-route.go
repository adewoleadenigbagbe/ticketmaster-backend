package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

func userRoutes(app *core.BaseApp, group *echo.Group) {
	userController := controllers.UserController{App: app}
	group.GET("users", userController.CreateUserHandler)
}
