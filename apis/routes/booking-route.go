package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/apis/controllers"
	middlewares "github.com/Wolechacho/ticketmaster-backend/apis/middleware"
	"github.com/Wolechacho/ticketmaster-backend/apis/core"
	"github.com/labstack/echo/v4"
)

func bookingRoutes(app *core.BaseApp, group *echo.Group) {
	bookingController := controllers.BookingController{App: app}
	group.POST("booking/book-show", bookingController.BookShowHandler)
	group.POST("booking/payment", bookingController.ChargeBookingHandler)
	group.GET("booking/generate-pdf", bookingController.GenerateInvoiceHandler, middlewares.AuthorizeUser)
}
