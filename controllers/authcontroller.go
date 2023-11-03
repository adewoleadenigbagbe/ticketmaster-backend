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

func (authController AuthController) RegisterHandler(authContext echo.Context) error {
	var err error
	request := new(services.CreateUserRequest)
	err = authContext.Bind(request)
	if err != nil {
		return authContext.JSON(http.StatusBadRequest, err.Error())
	}

	response := authController.App.AuthService.RegisterUser(*request)
	return authContext.JSON(response.StatusCode, response)
}

func (authController AuthController) SignInHandler(authContext echo.Context) error {
	var err error
	request := new(services.SignInRequest)
	err = authContext.Bind(request)
	if err != nil {
		return authContext.JSON(http.StatusBadRequest, err.Error())
	}

	resp, errors := authController.App.AuthService.SignIn(*request)
	if len(errors) > 0 {
		messageErrors := []string{}
		for _, er := range errors {
			messageErrors = append(messageErrors, er.Error())
		}
		return authContext.JSON(resp.StatusCode, messageErrors)
	}
	return authContext.JSON(http.StatusOK, resp)
}
