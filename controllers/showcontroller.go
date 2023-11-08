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
