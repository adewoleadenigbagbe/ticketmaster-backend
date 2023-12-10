package common

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/enums"
)

type BookingMessage struct {
	UserId          string
	ShowId          string
	BookingId       string
	CinemaSeatIds   []string
	Status          enums.ShowSeatStatus
	BookingDateTime time.Time
	ExpiryDateTime  time.Time
}

type SeatAvailableMessage struct {
	ShowId        string
	CinemaSeatIds []string
}

type SeatAvailableDTO struct {
	ShowId       string
	CinemaSeatId string
}

type ShowSeatDTO struct {
	Id           string
	Status       enums.ShowSeatStatus
	CinemaSeatId string
	ShowId       string
	SeatNumber   int
	UserId       string
}

type BookingDTO struct {
	UserId       string
	ShowId       string
	BookingId    string
	CinemaSeatId string
	Status       enums.ShowSeatStatus
}
