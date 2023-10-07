package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type Booking struct {
	Id            string    `gorm:"primaryKey;size:36;type:char(36)"`
	NumberOfSeats uint      `gorm:"not null"`
	BookDateTime  time.Time `gorm:"index; not null"`
	//Status is Enum in case of change
	Status       int    `gorm:"not null"`
	UserId       string `gorm:"index;not null;type:char(36)"`
	ShowId       string `gorm:"index;not null;type:char(36)"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (booking *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if len(booking.Id) == 0 || booking.Id == utilities.DEFAULT_UUID {
		booking.Id = sequentialguid.New().String()
	}

	return
}

func (booking Booking) GetId() string {
	return booking.Id
}
