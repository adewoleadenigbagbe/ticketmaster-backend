package controllers

import (
	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	App *core.BaseApp
}

func (authController AuthController) RegisterHandler(authContext echo.Context) error {
	return nil
}
