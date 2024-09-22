package common

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/shared/enums"
)

type BookingMessage struct {
	UserId          string
	ShowId          string
	CinemaSeatIds   []string
	Status          enums.ShowSeatStatus
	BookingDateTime time.Time
	ExpiryDateTime  time.Time
}

type SeatAvailableMessage struct {
	ShowId        string
	CinemaSeatIds []string
}
