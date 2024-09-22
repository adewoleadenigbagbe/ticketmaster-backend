package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/apis/core"
	"github.com/Wolechacho/ticketmaster-backend/infastructure/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type UserController struct {
	App *core.BaseApp
}

// CreateUserRole godoc
// @Summary      Create a new role
// @Description   Create a new role
// @Tags         userole
// @Accept       json
// @Produce      json
// @Param        CreateRoleRequest  body  services.CreateRoleRequest  true  "CreateRoleRequest"
// @Success      200  {object}  services.CreateRoleResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  []string
// @Router       /api/v1/user/add-role [post]
func (userController UserController) AddRoleHandler(userContext echo.Context) error {
	var err error
	request := new(services.CreateRoleRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := userController.App.UserService.AddRole(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return userContext.JSON(errResp.StatusCode, errs)
	}
	return userContext.JSON(http.StatusOK, dataResp)
}

// UpdateUserLocation godoc
// @Summary      Add a new location for user
// @Description   Add a new location for user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        UserLocationRequest  body  services.UserLocationRequest  true  "UserLocationRequest"
// @Success      200  {object}  services.UserLocationResponse
// @Failure      400  {object}  string
// @Failure      404  {object}  []string
// @Router       /api/v1/user/{id}/add-location [post]
func (userController UserController) UpdateUserLocationHandler(userContext echo.Context) error {
	var err error
	request := new(services.UserLocationRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	dataResp, errResp := userController.App.UserService.UpdateUserLocation(*request)
	if !reflect.ValueOf(errResp).IsZero() {
		errs := lo.Map(errResp.Errors, func(er error, index int) string {
			return er.Error()
		})
		return userContext.JSON(errResp.StatusCode, errs)
	}
	return userContext.JSON(http.StatusOK, dataResp)
}
