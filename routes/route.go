package routes

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes - Configuration for all incoming routes
func RegisterAllRoutes(echo *echo.Echo) {
	//movie routes
	group := echo.Group("/api/v1/")
	movieRoutes(group)
	cityRoutes(group)
}
