package controllers

import (
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/services"
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
