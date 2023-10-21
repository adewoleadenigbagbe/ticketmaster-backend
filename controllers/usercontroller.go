package controllers

import (
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	App *core.BaseApp
}

func (userController UserController) CreateUserHandler(userContext echo.Context) error {
	var err error
	request := new(services.CreateUserRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	response := userController.App.UserService.CreateUser(*request)
	return userContext.JSON(response.StatusCode, response)
}
