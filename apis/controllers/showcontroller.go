package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/apis/core"
	"github.com/Wolechacho/ticketmaster-backend/infastructure/services"
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

// GetAvailableShowSeat godoc
// @Summary      Get the available seat for a particular show
// @Description  Get the available seat for a particular show
// @Tags         shows
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Id"
// @Param        CinemaHallRequest  body  services.CinemaHallRequest  true  "CinemaHallRequest"
// @Success      200  {object}  services.CinemaHallResponse
// @Failure      400  {object}  []string
// @Failure      404  {object}  []string
// @Router       /api/v1/shows/{id}/available-seat [get]
func (showController ShowController) GetAvailableShowSeatHandler(showContext echo.Context) error {
	request := new(services.GetAvailableSeatRequest)
	err := showContext.Bind(request)

	if err != nil {
		return showContext.JSON(http.StatusBadRequest, err.Error())
	}

	resp, errors := showController.App.ShowService.GetAvailableShowSeat(*request)

	if len(errors) != 0 {
		respErrors := []string{}
		for _, err := range errors {
			respErrors = append(respErrors, err.Error())
		}

		return showContext.JSON(resp.StatusCode, respErrors)
	}

	return showContext.JSON(http.StatusOK, resp)
}
