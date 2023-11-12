package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type ShowController struct {
	App *core.BaseApp
}

func (showController *ShowController) CreateShowHandler(showContext echo.Context) error {
	var err error
	request := new(services.CreateShowRequest)
	err = showContext.Bind(request)
	if err != nil {
		return showContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := showController.App.ShowService.CreateShow(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return showContext.JSON(errResp.StatusCode, errs)
	}
	return showContext.JSON(http.StatusOK, dataResp)
}

func (showController ShowController) GetShowsByUserLocationHandler(showContext echo.Context) error {
	var err error
	request := new(services.GetShowsByLocationRequest)

	err = showContext.Bind(request)
	if err != nil {
		return showContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := showController.App.ShowService.GetShowsByUserLocation(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return showContext.JSON(errResp.StatusCode, errs)
	}
	return showContext.JSON(http.StatusOK, dataResp)
}
