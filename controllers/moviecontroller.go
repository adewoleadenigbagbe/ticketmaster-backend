package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type MovieController struct {
	App *core.BaseApp
}

func (movieController MovieController) GetMoviesHandler(movieContext echo.Context) error {
	req := new(services.GetMoviesRequest)
	err := movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}
	resp, _ := movieController.App.MovieService.GetMovies(*req)

	return movieContext.JSON(http.StatusOK, resp)
}

// GetMovieByID godoc
// @Summary      Get movie by ID
// @Description  Get a particular movie by ID
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Id"
// @Success      200  {object}  services.MovieDataResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  []string
// @Router       /api/v1/movies/{id} [get]
func (movieController MovieController) GetMovieByIdHandler(movieContext echo.Context) error {
	req := new(services.GetMovieByIdRequest)

	err := movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	dataResp, errResp := movieController.App.MovieService.GetMovieById(*req)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return movieContext.JSON(errResp.StatusCode, errs)
	}
	return movieContext.JSON(http.StatusOK, dataResp)
}

// SearchMovie godoc
// @Summary     Search for movie
// @Description  Search for movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        GetSearchRequest  body  services.GetSearchRequest  true  "GetSearchRequest"
// @Success      200  {object}  services.GetSearchResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  []string
// @Router       /api/v1/movies/search [get]
func (movieController MovieController) SearchMovieHandler(movieContext echo.Context) error {
	var err error
	req := new(services.GetSearchRequest)

	err = movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	dataResp, errResp := movieController.App.MovieService.SearchMovie(*req)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return movieContext.JSON(errResp.StatusCode, errs)
	}
	return movieContext.JSON(http.StatusOK, dataResp)
}
