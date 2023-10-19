package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
)

type CinemaController struct {
	App *core.BaseApp
}

func (cinemaController CinemaController) CreateCinemaHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CreateCinemaRequest)

	err = cinemaContext.Bind(request)
	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	response := cinemaController.App.CinemaService.CreateCinema(*request)
	return cinemaContext.JSON(http.StatusOK, response)
}

// CreateCinemaHall godoc
// @Summary      Create a halls to existing cinema
// @Description  Create a halls to existing cinema
// @Tags         cinemas
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Id"
// @Param        CinemaHallRequest  body  services.CinemaHallRequest  true  "CinemaHallRequest"
// @Success      200  {object}  services.CinemaHallResponse
// @Failure      400  {object}  []string
// @Failure      404  {object}  []string
// @Router       /api/v1/cinemas/{id}/cinemahall [get]
func (cinemaController CinemaController) CreateCinemaHallHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CinemaHallRequest)
	err = cinemaContext.Bind(request)

	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	response, fieldErrors := cinemaController.App.CinemaService.AddCinemaHall(*request)

	if len(fieldErrors) > 0 {
		errors := []string{}
		for _, err = range fieldErrors {
			errors = append(errors, err.Error())
		}
		return cinemaContext.JSON(http.StatusBadRequest, errors)
	}
	return cinemaContext.JSON(http.StatusOK, response)
}

func (cinemaController CinemaController) CreateCinemaSeatHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CreateCinemaSeatRequest)
	err = cinemaContext.Bind(request)

	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := cinemaController.App.CinemaService.AddCinemaSeat(*request)
	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}
	return cinemaContext.JSON(http.StatusOK, resp)
}
