package routes

import (
	"github.com/Wolechacho/ticketmaster-backend/core"
)

// RegisterRoutes - Configuration for all incoming routes
func RegisterAllRoutes(app *core.BaseApp) {
	group := app.Echo.Group("/api/v1/")
	authRoutes(app , group)
	movieRoutes(app, group)
	showRoutes(app, group)
	cityRoutes(app, group)
	cinemaRoutes(app, group)
	userRoutes(app, group)
	bookingRoutes(app, group)
}
