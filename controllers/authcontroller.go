package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
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

	response, errResp := authController.App.AuthService.RegisterUser(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return authContext.JSON(errResp.StatusCode, errs)
	}

	return authContext.JSON(http.StatusOK, response)
}

// Auth godoc
// @Summary     SignIn user
// @Description    SignIn user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        SignInRequest  body  services.SignInRequest  true  "SignInRequest"
// @Success      200  {object}  services.SignInResponse
// @Failure      400  {object}  []string
// @Failure      404  {object}  []string
// @Router       /api/v1/auth/signin [post]
func (authController AuthController) SignInHandler(authContext echo.Context) error {
	var err error
	request := new(services.SignInRequest)
	err = authContext.Bind(request)
	if err != nil {
		return authContext.JSON(http.StatusBadRequest, err.Error())
	}

	resp, errResp := authController.App.AuthService.SignIn(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return authContext.JSON(errResp.StatusCode, errs)
	}
	return authContext.JSON(http.StatusOK, resp)
}

func (authController AuthController) SignOutHandler(authContext echo.Context) error {
	authContext.SetCookie(&http.Cookie{})
	return authContext.JSON(http.StatusOK, "success")
}
