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

func (userController UserController) AddRoleHandler(userContext echo.Context) error {
	var err error
	request := new(services.CreateRoleRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	resp, errors := userController.App.UserService.AddRole(*request)

	if len(errors) > 0 {
		respErrors := []string{}
		for _, er := range errors {
			respErrors = append(respErrors, er.Error())
		}
		return userContext.JSON(resp.StatusCode, respErrors)
	}

	return userContext.JSON(http.StatusOK, resp)
}
