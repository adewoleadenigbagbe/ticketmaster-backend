package main

import (
	"os"

	"github.com/Wolechacho/ticketmaster-backend/cmd"
)

func main() {
	ticketmaster := cmd.NewTicketMaster()
	err := ticketmaster.Start()
	if err != nil {
		os.Exit(1)
	}
}
