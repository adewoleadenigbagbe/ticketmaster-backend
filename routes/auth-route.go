package routes

func authRoutes(app *core.BaseApp, group *echo.Group) {
	cinemaController := controllers.CinemaController{App: app}
	group.POST("cinemas", cinemaController.CreateCinemaHandler)
	group.POST("cinemas/:id/cinemahall", cinemaController.CreateCinemaHallHandler)
	group.POST("cinemas/:id/cinemahall/:cinemahallId/seat", cinemaController.CreateCinemaSeatHandler)
}
