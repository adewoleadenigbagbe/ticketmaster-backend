package entities

import (
	"time"
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
