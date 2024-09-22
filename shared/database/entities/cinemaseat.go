package entities

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/shared/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/shared/helpers"
	"github.com/Wolechacho/ticketmaster-backend/shared/helpers/utilities"
	"gorm.io/gorm"
)

type CinemaSeat struct {
	Id           string         `gorm:"column:Id"`
	SeatNumber   int            `gorm:"column:SeatNumber"`
	Type         enums.SeatType `gorm:"column:Type"`
	CinemaHallId string         `gorm:"column:CinemaHallId"`
	IsDeprecated bool           `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time      `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time      `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaSeat) TableName() string {
	return "CinemaSeats"
}

func (cinemaSeat *CinemaSeat) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinemaSeat.Id) == 0 || cinemaSeat.Id == utilities.DEFAULT_UUID {
		cinemaSeat.Id = sequentialguid.New().String()
	}
	cinemaSeat.IsDeprecated = false

	return
}

func (cinemaSeat CinemaSeat) GetId() string {
	return cinemaSeat.Id
}
