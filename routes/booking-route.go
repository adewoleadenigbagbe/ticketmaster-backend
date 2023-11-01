package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

func bookingRoutes(app *core.BaseApp, group *echo.Group) {
	bookingController := controllers.BookingController{App: app}
	group.POST("booking/book-show", bookingController.BookShowHandler)
}
