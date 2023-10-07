package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type City struct {
	Id           string    `gorm:"primaryKey;size:36;column:Id;type:char(36)"`
	Name         string    `gorm:"not null;index;column:Name;type:varchar(255)"`
	State        string    `gorm:"not null;column:State;type:varchar(255)"`
	Zipcode      string    `gorm:"not null;column:ZipCode;type:varchar(255)"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (City) TableName() string {
	return "Cities"
}

func (city *City) BeforeCreate(tx *gorm.DB) (err error) {
	if len(city.Id) == 0 || city.Id == utilities.DEFAULT_UUID {
		city.Id = sequentialguid.New().String()
	}
	return
}

func (city City) GetId() string {
	return city.Id
}
