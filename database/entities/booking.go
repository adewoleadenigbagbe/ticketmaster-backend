package entities

import (
	"time"
)

type Booking struct {
	Id            string    `gorm:"primaryKey;size:36"`
	NumberOfSeats uint      `gorm:"not null"`
	BookDateTime  time.Time `gorm:"index; not null"`
	//Status is Enum in case of change
	Status       int    `gorm:"not null"`
	UserId       string `gorm:"index;not null"`
	ShowId       string `gorm:"index;not null"`
	IsDeprecated bool
}
