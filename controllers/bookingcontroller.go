package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type BookingController struct {
	App *core.BaseApp
}

func (bookingController BookingController) GenerateInvoiceHandler(bookingContext echo.Context) error {
	var err error
	request := new(services.GeneratePdfRequest)
	err = bookingContext.Bind(request)
	if err != nil {
		return bookingContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := bookingController.App.BookingService.GenerateInvoicePDF(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return bookingContext.JSON(errResp.StatusCode, errs)
	}

	bookingContext.Response().Header().Set("Content-Disposition", "attachment; filename=ticketinvoice.pdf")
	bookingContext.Response().Header().Set("Content-Type", "application/pdf")
	bookingContext.Response().WriteHeader(http.StatusOK)
	bookingContext.Response().Write(dataResp.PdfBytes)
	return nil
}
