package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/apis/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

func authRoutes(app *core.BaseApp, group *echo.Group) {
	authController := controllers.AuthController{App: app}
	group.POST("auth/register", authController.RegisterHandler)
	group.POST("auth/signin", authController.SignInHandler)
	group.POST("auth/signout", authController.SignOutHandler)
}
