package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	_ "github.com/Wolechacho/ticketmaster-backend/docs"
	middlewares "github.com/Wolechacho/ticketmaster-backend/middleware"
	"github.com/Wolechacho/ticketmaster-backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8185
// @BasePath /
// @schemes http
func main() {

	//create a database connection
	db.ConnectToDatabase()

	// Create a new Echo instance
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

	// Start the server
	e.Start(":8185")

	// Start server
	go func() {
		if err := e.Start(":8185"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
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
