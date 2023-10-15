package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
)

type MovieController struct {
	App *core.BaseApp
}

func (movieController *MovieController) GetMoviesHandler(movieContext echo.Context) error {
	req := new(services.GetMoviesRequest)
	err := movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}
	resp, _ := movieController.App.MovieService.GetMovies(*req)

	return movieContext.JSON(http.StatusOK, resp)
}

func (movieController *MovieController) GetMovieByIdHandler(movieContext echo.Context) error {
	req := new(services.GetMovieByIdRequest)

	err := movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	resp, err := movieController.App.MovieService.GetMovieById(*req)

	if err != nil {
		return movieContext.JSON(resp.StatusCode, err.Error())
	}

	return movieContext.JSON(http.StatusOK, resp)
}

func (movieController MovieController) SearchMovieHandler(movieContext echo.Context) error {
	var err error
	req := new(services.GetSearchRequest)

	err = movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	resp, err := movieController.App.MovieService.SearchMovie(*req)

	if err != nil {
		return movieContext.JSON(resp.StatusCode, err.Error())
	}

	return movieContext.JSON(http.StatusOK, resp)
}
