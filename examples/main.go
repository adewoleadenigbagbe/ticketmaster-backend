package main

import (
	"fmt"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
)

func main() {
	guid := sequentialguid.New()
	fmt.Println(guid)
}
