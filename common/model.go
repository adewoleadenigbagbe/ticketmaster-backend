package common

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/enums"
)

type BookingMessage struct {
	UserId          string
	ShowId          string
	CinemaSeatIds   []string
	Status          enums.ShowSeatStatus
	BookingDateTime time.Time
	ExpiryDateTime  time.Time
}
