package routes

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes - Configuration for all incoming routes
func RegisterAllRoutes(echo *echo.Echo) {
	//movie routes
	movieRoutes(echo)
}
