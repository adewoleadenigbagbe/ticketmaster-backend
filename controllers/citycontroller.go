package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
)

type CityController struct {
	App *core.BaseApp
}

// GetCity godoc
// @Summary      Get city by ID
// @Description  Get a particular city by ID
// @Tags         cities
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Id"
// @Success      200  {object}  services.CityModelResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Router       /api/v1/city/{id} [get]
func (cityController *CityController) GetCityByIdHandler(cityContext echo.Context) error {
	var err error
	req := new(services.GetCityByIdRequest)
	err = cityContext.Bind(req)
	if err != nil {
		return cityContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	resp, err := cityController.App.CityService.GetCityById(*req)
	if err != nil {
		return cityContext.JSON(resp.StatusCode, err.Error())
	}
	return cityContext.JSON(http.StatusOK, resp.City)
}
