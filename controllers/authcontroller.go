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

// Auth godoc
// @Summary     Register new user
// @Description   Register new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        CreateUserRequest  body  services.CreateUserRequest  true  "CreateUserRequest"
// @Success      200  {object}  services.CreateUserResponse
// @Failure      400  {object}  []string
// @Failure      404  {object}  []string
// @Router       /api/v1/auth/register [post]
func (authController AuthController) RegisterHandler(authContext echo.Context) error {
	var err error
	request := new(services.CreateUserRequest)
	err = authContext.Bind(request)
	if err != nil {
		return authContext.JSON(http.StatusBadRequest, err.Error())
	}

	response, errs := authController.App.AuthService.RegisterUser(*request)

	if len(errs) > 0 {
		respErrors := []string{}
		for _, e := range errs {
			respErrors = append(respErrors, e.Error())
		}
		return authContext.JSON(response.StatusCode, respErrors)
	}
	return authContext.JSON(http.StatusOK, response)
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
