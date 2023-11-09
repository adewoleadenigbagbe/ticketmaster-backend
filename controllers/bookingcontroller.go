package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type BookingController struct {
	App *core.BaseApp
}

func (bookingController BookingController) BookShowHandler(bookingContext echo.Context) error {
	var err error
	request := new(services.BookRequest)

	err = bookingContext.Bind(request)
	if err != nil {
		return bookingContext.JSON(http.StatusBadRequest, err.Error())
	}

	resp, errors := bookingController.App.BookService.BookShow(*request)
	if len(errors) > 0 {
		errs := lo.Map(errors, func(e error, index int) string {
			return e.Error()
		})
		return bookingContext.JSON(resp.StatusCode, errs)
	}

	return bookingContext.JSON(http.StatusOK, resp)
}
