package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
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
// @Failure      404  {object}  []string
// @Router       /api/v1/city/{id} [get]
func (cityController *CityController) GetCityByIdHandler(cityContext echo.Context) error {
	var err error
	req := new(services.GetCityByIdRequest)
	err = cityContext.Bind(req)
	if err != nil {
		return cityContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	dataResp, errResp := cityController.App.CityService.GetCityById(*req)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return cityContext.JSON(errResp.StatusCode, errs)
	}

	return cityContext.JSON(http.StatusOK, dataResp)
}
