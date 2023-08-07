package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type City struct {
	Id           string         `gorm:"primaryKey;size:36;column:Id"`
	Name         string         `gorm:"not null;index;column:Name"`
	State        string         `gorm:"not null;column:State"`
	Zipcode      sql.NullString `gorm:"column:ZipCode"`
	IsDeprecated bool           `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time      `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time      `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (City) TableName() string {
	return "Cities"
}

func (city *City) BeforeCreate(tx *gorm.DB) (err error) {
	if len(city.Id) == 0 || city.Id == DEFAULT_UUID {
		city.Id = sequentialguid.New().String()
	}
	return
}
