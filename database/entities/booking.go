package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

const (
	MIN_DATE     = "0001-01-01T00:00:00Z"
	MAX_DATE     = "9999-12-31T00:00:00Z"
	DEFAULT_UUID = "00000000-0000-0000-0000-000000000000"
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
	if len(booking.Id) == 0 || booking.Id == DEFAULT_UUID {
		booking.Id = sequentialguid.New().String()
	}

	return
}
