package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
)

type BookingController struct {
	App *core.BaseApp
}

func (bookingController BookingController) ChargeBookingHandler(bookingContext echo.Context) error {
	var err error
	request := new(services.CreatePaymentRequest)

	err = bookingContext.Bind(request)
	if err != nil {
		return bookingContext.JSON(http.StatusBadRequest, err.Error())
	}
	return nil
}
