package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type CinemaSeat struct {
	Id         string `gorm:"column:Id"`
	SeatNumber int    `gorm:"column:SeatNumber"`
	//Type is an enum
	Type         int       `gorm:"column:Type"`
	CinemaHallId string    `gorm:"column:CinemaHallId"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaSeat) TableName() string {
	return "CinemaSeats"
}

func (cinemaSeat *CinemaSeat) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinemaSeat.Id) == 0 || cinemaSeat.Id == DEFAULT_UUID {
		cinemaSeat.Id = sequentialguid.New().String()
	}
	cinemaSeat.IsDeprecated = false

	return
}
