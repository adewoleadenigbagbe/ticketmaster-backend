package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type City struct {
	Id           string     `gorm:"column:Id"`
	Name         string     `gorm:"column:Name"`
	State        string     `gorm:"column:State"`
	Coordinates  Coordinate `gorm:"column:Coordinates"`
	Zipcode      string     `gorm:"column:ZipCode"`
	IsDeprecated bool       `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time  `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time  `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (City) TableName() string {
	return "Cities"
}

func (city *City) BeforeCreate(tx *gorm.DB) (err error) {
	if len(city.Id) == 0 || city.Id == DEFAULT_UUID {
		city.Id = sequentialguid.New().String()
	}
	city.IsDeprecated = false

	return
}
