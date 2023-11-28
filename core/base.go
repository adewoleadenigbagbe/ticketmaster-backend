package core

import (
	"os"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
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
	Logger             zerolog.Logger
}

func ConfigureApp() *BaseApp {
	//create a database connection
	db := db.ConnectToDatabase()

	app := &BaseApp{
		IsMigrationChecked: false,
		Echo:               echo.New(),
		DB:                 db,
		CinemaService:      services.CinemaService{DB: db, Logger: zerolog.New(os.Stdout).With().Timestamp().Logger()},
		CityService:        services.CityService{DB: db, Logger: zerolog.New(os.Stdout).With().Timestamp().Logger()},
		MovieService:       services.MovieService{DB: db, Logger: zerolog.New(os.Stdout).With().Timestamp().Logger()},
		ShowService:        services.ShowService{DB: db, Logger: zerolog.New(os.Stdout).With().Timestamp().Logger()},
		UserService:        services.UserService{DB: db, Logger: zerolog.New(os.Stdout).With().Timestamp().Logger()},
		AuthService:        services.AuthService{DB: db, Logger: zerolog.New(os.Stdout).With().Timestamp().Logger()},
	}

	return app
}

func (app *BaseApp) Db() *gorm.DB {
	return app.DB
}

func (app *BaseApp) GetEcho() *echo.Echo {
	return app.Echo
}
