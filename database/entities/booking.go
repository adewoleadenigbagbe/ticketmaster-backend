package entities

import (
	"time"
)

const (
	DEFAULT_TIME = "0001-01-01T00:00:00Z"
	DEFAULT_UUID = "00000000-0000-0000-0000-000000000000"
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
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}
