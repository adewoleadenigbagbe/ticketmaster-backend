package core

import (
	"fmt"
	"log"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/services"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// BaseApp implements core.App and defines the structure of the whole application
type BaseApp struct {
	IsMigrationChecked bool
	DB                 *gorm.DB
	Nats               *nats.Conn
	Echo               *echo.Echo
	CinemaService      services.CinemaService
	CityService        services.CityService
	MovieService       services.MovieService
	ShowService        services.ShowService
	UserService        services.UserService
	BookService        services.BookService
	AuthService        services.AuthService
}

func ConfigureApp() *BaseApp {
	//create a database connection
	db := db.ConnectToDatabase()
	nc := ConnectToNats()

	app := &BaseApp{
		IsMigrationChecked: false,
		Echo:               echo.New(),
		DB:                 db,
		Nats:               nc,
		CinemaService:      services.CinemaService{DB: db},
		CityService:        services.CityService{DB: db},
		MovieService:       services.MovieService{DB: db},
		ShowService:        services.ShowService{DB: db},
		UserService:        services.UserService{DB: db},
		BookService:        services.BookService{DB: db, Nc: nc},
		AuthService:        services.AuthService{DB: db},
	}

	return app
}

func ConnectToNats() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Nats connected")
	return nc
}

func (app *BaseApp) Db() *gorm.DB {
	return app.DB
}

func (app *BaseApp) GetEcho() *echo.Echo {
	return app.Echo
}

func (app *BaseApp) GetNats() *nats.Conn {
	return app.Nats
}
