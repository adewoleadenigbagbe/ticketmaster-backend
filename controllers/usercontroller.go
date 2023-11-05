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

// CreateUserRole godoc
// @Summary      Create a new role
// @Description   Create a new role
// @Tags         userole
// @Accept       json
// @Produce      json
// @Param        CreateRoleRequest  body  services.CreateRoleRequest  true  "CreateRoleRequest"
// @Success      200  {object}  services.CreateRoleResponse
// @Failure      400  {object}  []string
// @Failure      404  {object}  []string
// @Router       /api/v1/user/add-role [post]
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
