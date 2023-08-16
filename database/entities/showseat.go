package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type ShowSeat struct {
	Id           string  `gorm:"primaryKey;size:36;type:char(36)"`
	Status       int     `gorm:"not null"`
	Price        float64 `gorm:"not null"`
	CinemaSeatId string  `gorm:"index;not null;type:char(36)"`
	ShowId       string  `gorm:"index;not null;type:char(36)"`
	BookingId    string  `gorm:"index;not null;type:char(36)"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (showSeat *ShowSeat) BeforeCreate(tx *gorm.DB) (err error) {
	if len(showSeat.Id) == 0 || showSeat.Id == DEFAULT_UUID {
		showSeat.Id = sequentialguid.New().String()
	}
	return
}

func (showSeat ShowSeat) GetId() string {
	return showSeat.Id
}
