package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/core"
	_ "github.com/Wolechacho/ticketmaster-backend/docs"
	middlewares "github.com/Wolechacho/ticketmaster-backend/middleware"
	"github.com/Wolechacho/ticketmaster-backend/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title TicketMaster Endpoints
// @version 1.0
// @description Contains all the endpoint for the ticketmaster app
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8185
// @BasePath /
// Consumes:
//      - application/json
//   Produces:
//   - application/json
// @schemes http
func main() {
	//load env variables
	err := godotenv.Load("..//.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//configure application
	app, err := core.ConfigureApp()
	if err != nil {
		log.Fatal(err)
	}

	// set migration middleware
	mc := middlewares.MigrationChanges{App: app}
	app.Echo.Use(mc.CheckMigrationCompatibility)

	app.Echo.Logger.SetLevel(log.INFO)

	// Define a route
	app.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	app.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	//Register All Routes
	routes.RegisterAllRoutes(app)

	// Start server
	go func() {
		if err := app.Echo.Start(":8185"); err != nil && err != http.ErrServerClosed {
			app.Echo.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Echo.Shutdown(ctx); err != nil {
		app.Echo.Logger.Fatal(err)
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
