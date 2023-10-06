package core

import (
	"fmt"
	"net/http"

	middlewares "github.com/Wolechacho/ticketmaster-backend/middleware"
	"github.com/Wolechacho/ticketmaster-backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func createDatabase() *gorm.DB {
	var err error
	dsn := "root:P@ssw0r1d@tcp(127.0.0.1:3306)/ticketmasterDB?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the Database")

	return db
}

func initializeApi() *echo.Echo {
	e := echo.New()

	// set migration middleware
	e.Use(middlewares.CheckMigrationCompatibility)

	e.Logger.SetLevel(log.INFO)

	// Define a route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//Register All Routes
	routes.RegisterAllRoutes(e)

	return e
}
