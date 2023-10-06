package core

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// App is an interface implemented by struct to instantiate
// dependencies for the application to run
type App interface {
	DB() *gorm.DB

	GetEcho() *echo.Echo
}
