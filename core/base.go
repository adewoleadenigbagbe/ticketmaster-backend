package core

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// BaseApp implements core.App and defines the structure of the whole application
type BaseApp struct {
	db   *gorm.DB
	echo *echo.Echo
}

func NewBaseApp() *BaseApp {
	app := &BaseApp{
		db:   createDatabase(),
		echo: initializeApi(),
	}
	return app
}

func (app *BaseApp) DB() *gorm.DB {
	return app.db
}

func (app *BaseApp) GetEcho() *echo.Echo {
	return app.echo
}
