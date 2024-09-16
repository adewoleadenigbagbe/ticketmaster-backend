package core

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Wolechacho/ticketmaster-backend/infastructure/services"
	"github.com/Wolechacho/ticketmaster-backend/infastructure/tools"
	db "github.com/Wolechacho/ticketmaster-backend/shared/database"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

const (
	DbFile = "database.json"
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

func ConfigureApp() (*BaseApp, error) {
	//connect to nats
	natUrl := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natUrl)
	if err != nil {
		log.Fatal(err)
	}

	//create a database connection
	dbConfigPath, err := filepath.Abs(DbFile)
	if err != nil {
		return nil, err
	}

	db, err := db.ConnectToDatabase(dbConfigPath)
	if err != nil {
		return nil, err
	}
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
		BookService:        services.BookService{DB: db, Nc: nc, PDFService: tools.PDFService{}},
		AuthService:        services.AuthService{DB: db},
	}

	return app, nil
}

func ConnectToNats() (*nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	return nc, nil
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
