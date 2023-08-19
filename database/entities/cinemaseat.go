package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type CinemaSeat struct {
	Id         string `gorm:"primaryKey;size:36;column:Id;type:char(36)"`
	SeatNumber int    `gorm:"not null;column:SeatNumber"`
	//Type is an enum
	Type         int       `gorm:"not null;column:Type"`
	CinemaHallId string    `gorm:"index;not null;column:CinemaHallId;type:char(36)"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaSeat) TableName() string {
	return "CinemaSeats"
}

func (cinemaSeat *CinemaSeat) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinemaSeat.Id) == 0 || cinemaSeat.Id == DEFAULT_UUID {
		cinemaSeat.Id = sequentialguid.New().String()
	}
	return
}
