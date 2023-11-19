package controllers

import (
	"net/http"
	"reflect"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/samber/lo"
)

type CinemaController struct {
	App *core.BaseApp
}

// CreateCinema godoc
// @Summary      Create a new cinema
// @Description   Create a new cinema
// @Tags         cinemas
// @Accept       json
// @Produce      json
// @Param        CreateCinemaRequest  body  services.CreateCinemaRequest  true  "CreateCinemaRequest"
// @Success      200  {object}  services.CreateCinemaResponse
// @Failure      400  {object}  []string
// @Failure      404  {object}  []string
// @Router       /api/v1/cinemas [post]
func (cinemaController CinemaController) CreateCinemaHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CreateCinemaRequest)

	err = cinemaContext.Bind(request)
	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := cinemaController.App.CinemaService.CreateCinema(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return cinemaContext.JSON(errResp.StatusCode, errs)
	}
	return cinemaContext.JSON(http.StatusOK, dataResp)
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
// @Router       /api/v1/cinemas/{id}/cinemahall [post]
func (cinemaController CinemaController) CreateCinemaHallHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CinemaHallRequest)
	err = cinemaContext.Bind(request)

	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := cinemaController.App.CinemaService.AddCinemaHall(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return cinemaContext.JSON(errResp.StatusCode, errs)
	}
	return cinemaContext.JSON(http.StatusOK, dataResp)
}

// CreateCinemaSeat godoc
// @Summary      Create a existing seats to halls
// @Description   Create a existing seats to halls
// @Tags         cinemas
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Id"
// @Param        cinemahallId  path  string  true  "CinemaHallId"
// @Param        CreateCinemaSeatRequest  body  services.CreateCinemaSeatRequest  true  "CreateCinemaSeatRequest"
// @Success      200  {object}  services.CreateCinemaSeatResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Router       /api/v1/cinemas/{id}/cinemahall/{cinemahallId}/seat [post]
func (cinemaController CinemaController) CreateCinemaSeatHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CreateCinemaSeatRequest)
	err = cinemaContext.Bind(request)

	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := cinemaController.App.CinemaService.AddCinemaSeat(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return cinemaContext.JSON(errResp.StatusCode, errs)
	}
	return cinemaContext.JSON(http.StatusOK, dataResp)
}

// CinemaRateRequest godoc
// @Summary      Add new rate for cinema
// @Description   Add new rate for cinema
// @Tags         cinemas
// @Accept       json
// @Produce      json
// @Param        CinemaId  path  string  true  "CinemaId"
// @Param        CinemaRateRequest  body  services.CinemaRateRequest  true  "CinemaRateRequest"
// @Success      200  {object}  services.CinemaRateResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Router       /api/v1/cinemas/{id}/rate [post]
func (cinemaController CinemaController) AddCinemaRateHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CinemaRateRequest)
	err = cinemaContext.Bind(request)

	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := cinemaController.App.CinemaService.AddCinemaRate(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return cinemaContext.JSON(errResp.StatusCode, errs)
	}
	return cinemaContext.JSON(http.StatusOK, dataResp)

}

// CinemaRateRequest godoc
// @Summary      Add new rate for cinema
// @Description   Add new rate for cinema
// @Tags         cinemas
// @Accept       json
// @Produce      json
// @Param        CinemaId  path  string  true  "CinemaId"
// @Param        CinemaRateRequest  body  services.CinemaRateRequest  true  "CinemaRateRequest"
// @Success      200  {object}  services.CinemaRateResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Router       /api/v1/cinemas/{id}/rate [post]
func (cinemaController CinemaController) AddCinemaRateHandler(cinemaContext echo.Context) error {
	var err error
	request := new(services.CinemaRateRequest)
	err = cinemaContext.Bind(request)

	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := cinemaController.App.CinemaService.AddCinemaRate(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return cinemaContext.JSON(errResp.StatusCode, errs)
	}
	return cinemaContext.JSON(http.StatusOK, dataResp)

}
