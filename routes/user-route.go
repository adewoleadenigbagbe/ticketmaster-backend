package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	middlewares "github.com/Wolechacho/ticketmaster-backend/middleware"
	"github.com/labstack/echo/v4"
)

func userRoutes(app *core.BaseApp, group *echo.Group) {
	userController := controllers.UserController{App: app}
	group.POST("user/add-role", userController.AddRoleHandler, middlewares.AuthorizeAdmin)
	group.POST("user/:id/add-location", userController.UpdateUserLocationHandler, middlewares.AuthorizeUser)
}
