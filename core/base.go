package core

import (
	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// BaseApp implements core.App and defines the structure of the whole application
type BaseApp struct {
	IsMigrationChecked bool
	DB                 *gorm.DB
	Echo               *echo.Echo
	CinemaService      services.CinemaService
	CityService        services.CityService
	MovieService       services.MovieService
	ShowService        services.ShowService
	UserService        services.UserService
	AuthService        services.AuthService
}

func ConfigureApp() *BaseApp {
	//create a database connection
	db := db.ConnectToDatabase()

	app := &BaseApp{
		IsMigrationChecked: false,
		Echo:               echo.New(),
		DB:                 db,
		CinemaService:      services.CinemaService{DB: db},
		CityService:        services.CityService{DB: db},
		MovieService:       services.MovieService{DB: db},
		ShowService:        services.ShowService{DB: db},
		UserService:        services.UserService{DB: db},
		AuthService:        services.AuthService{DB: db},
	}

	return app
}

func (app *BaseApp) Db() *gorm.DB {
	return app.DB
}

func (app *BaseApp) GetEcho() *echo.Echo {
	return app.Echo
}
