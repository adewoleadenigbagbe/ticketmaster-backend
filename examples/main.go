package main

import (
	"log"

	"github.com/Wolechacho/ticketmaster-backend"
)

func main() {
	app := ticketmaster.New()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
