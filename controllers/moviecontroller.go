package controllers

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type MovieController struct {
}

func (movieController *MovieController) GetMovies(movieContext echo.Context) error {
	fmt.Println("Getting the first movie list")
	return nil
}
