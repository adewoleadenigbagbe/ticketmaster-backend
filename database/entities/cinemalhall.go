package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type CinemaHall struct {
	Id           string       `gorm:"column:Id"`
	Name         string       `gorm:"column:Name"`
	TotalSeat    int          `gorm:"column:TotalSeat"`
	CinemaId     string       `gorm:"column:CinemaId"`
	CinemaSeats  []CinemaSeat `gorm:"foreignKey:CinemaHallId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsDeprecated bool         `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time    `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time    `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaHall) TableName() string {
	return "CinemaHalls"
}

func (cinemaHall *CinemaHall) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinemaHall.Id) == 0 || cinemaHall.Id == utilities.DEFAULT_UUID {
		cinemaHall.Id = sequentialguid.New().String()
	}
	cinemaHall.IsDeprecated = false

	return
}

func (cinemaHall CinemaHall) GetId() string {
	return cinemaHall.Id
}
