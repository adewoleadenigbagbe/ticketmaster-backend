package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
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

	resp, errors := showController.App.ShowService.CreateShow(*request)
	respErrors := []string{}
	if len(errors) != 0 {
		for _, er := range errors {
			respErrors = append(respErrors, er.Error())
		}
		return showContext.JSON(http.StatusBadRequest, respErrors)
	}

	return showContext.JSON(http.StatusOK, resp)
}

func (showController ShowController) GetShowsByUserLocationHandler(showContext echo.Context) error {
	var err error
	request := new(services.GetShowsByLocationRequest)

	err = showContext.Bind(request)
	if err != nil {
		return showContext.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := showController.App.ShowService.GetShowsByUserLocation(*request)
	if err != nil {
		return showContext.JSON(resp.StatusCode, err.Error())
	}
	return showContext.JSON(http.StatusOK, resp)
}
