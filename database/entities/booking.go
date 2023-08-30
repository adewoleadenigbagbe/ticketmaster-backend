package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

const (
	MIN_DATE     = "0001-01-01T00:00:00Z"
	MAX_DATE     = "9999-12-31T23:59:59Z"
	DEFAULT_UUID = "00000000-0000-0000-0000-000000000000"
)

type Booking struct {
	Id            string    `gorm:"column:Id"`
	NumberOfSeats uint      `gorm:"column:NumberOfSeats"`
	BookDateTime  time.Time `gorm:"column:BookDateTime"`
	//Status is Enum in case of change
	Status       int       `gorm:"column:Status"`
	UserId       string    `gorm:"column:UserId"`
	ShowId       string    `gorm:"column:ShowId"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (booking *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if len(booking.Id) == 0 || booking.Id == DEFAULT_UUID {
		booking.Id = sequentialguid.New().String()
	}

	booking.IsDeprecated = false
	return
}
