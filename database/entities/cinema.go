package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type Cinema struct {
	Id                string       `gorm:"column:Id"`
	Name              string       `gorm:"column:Name"`
	TotalCinemalHalls int          `gorm:"column:TotalCinemalHalls"`
	CityId            string       `gorm:"column:CityId"`
	CinemaHalls       []CinemaHall `gorm:"foreignKey:CinemaId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsDeprecated      bool         `gorm:"column:IsDeprecated"`
	CreatedOn         time.Time    `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn        time.Time    `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Cinema) TableName() string {
	return "Cinemas"
}

func (cinema *Cinema) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinema.Id) == 0 || cinema.Id == utilities.DEFAULT_UUID {
		cinema.Id = sequentialguid.New().String()
	}

	cinema.IsDeprecated = false
	return
}

func (cinema Cinema) GetId() string {
	return cinema.Id
}
