package controllers

import (
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
)

type CityController struct {
	App *core.BaseApp
}

func (cityController *CityController) GetCityById(cityContext echo.Context) error {
	var err error
	req := new(services.GetCityByIdRequest)
	fmt.Println("Did we get here GetCityById")
	err = cityContext.Bind(req)
	if err != nil {
		return cityContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	resp, err := cityController.App.CityService.GetCityById(*req)
	if err != nil {
		return cityContext.JSON(http.StatusNotFound, err.Error())
	}
	return cityContext.JSON(http.StatusOK, resp)
}
