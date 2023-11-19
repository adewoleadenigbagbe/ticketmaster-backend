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

// CreateShow godoc
// @Summary      Create a new show
// @Description  Create a new show
// @Tags         shows
// @Accept       json
// @Produce      json
// @Param        CreateShowRequest  body  services.CreateShowRequest  true  "CreateShowRequest"
// @Success      200  {object}  services.CreateShowResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  []string
// @Router       /api/v1/shows [post]
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

// GetShowsByUserLocation godoc
// @Summary      Get a list of user based on user current location
// @Description  Get a list of user based on user current location
// @Tags         shows
// @Accept       json
// @Produce      json
// @Param        GetShowsByLocationRequest  body  services.GetShowsByLocationRequest  true  "GetShowsByLocationRequest"
// @Success      200  {object}  services.GetShowsByLocationResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  []string
// @Router       /api/v1/shows/user-location [get]
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
