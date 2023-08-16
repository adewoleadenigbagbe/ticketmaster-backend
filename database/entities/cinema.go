package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type Cinema struct {
	Id                string    `gorm:"primaryKey;size:36;column:Id"`
	Name              string    `gorm:"index;not null;column:Name"`
	TotalCinemalHalls int       `gorm:"not null;column:TotalCinemalHalls"`
	CityId            string    `gorm:"index;not null1;column:CityId"`
	IsDeprecated      bool      `gorm:"column:IsDeprecated"`
	CreatedOn         time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn        time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Cinema) TableName() string {
	return "Cinemas"
}

func (cinema *Cinema) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinema.Id) == 0 || cinema.Id == DEFAULT_UUID {
		cinema.Id = sequentialguid.New().String()
	}

	return
}
