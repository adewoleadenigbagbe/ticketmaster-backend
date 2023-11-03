package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	App *core.BaseApp
}

func (authController AuthController) RegisterHandler(userContext echo.Context) error {
	var err error
	request := new(services.CreateUserRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	response := authController.App.AuthService.RegisterUser(*request)
	return userContext.JSON(response.StatusCode, response)
}
