package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type CinemaHall struct {
	Id           string    `gorm:"column:Id"`
	Name         string    `gorm:"column:Name"`
	TotalSeat    int       `gorm:"column:TotalSeat"`
	CinemaId     string    `gorm:"column:CinemaId"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaHall) TableName() string {
	return "CinemaHalls"
}

func (cinemaHall *CinemaHall) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinemaHall.Id) == 0 || cinemaHall.Id == DEFAULT_UUID {
		cinemaHall.Id = sequentialguid.New().String()
	}
	cinemaHall.IsDeprecated = false

	return
}
