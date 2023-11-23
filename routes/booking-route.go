package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/controllers"
	"github.com/Wolechacho/ticketmaster-backend/core"
	middlewares "github.com/Wolechacho/ticketmaster-backend/middleware"
	"github.com/labstack/echo/v4"
)

func bookingRoutes(app *core.BaseApp, group *echo.Group) {
	bookingController := controllers.BookingController{App: app}
	group.GET("booking/generate-pdf", bookingController.GenerateInvoiceHandler, middlewares.AuthorizeUser)
}
