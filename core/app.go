package core

import (
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// App is an interface implemented by struct to instantiate
// dependencies for the application to run
type App interface {
	Db() *gorm.DB
	GetEcho() *echo.Echo
	GetNats() *nats.Conn
}
