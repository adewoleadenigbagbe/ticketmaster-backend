package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type UserLocation struct {
	Id           string `gorm:"primaryKey;size:36;type:char(36);column:Id"`
	Address      string `gorm:"not null;type:mediumtext;colunm:Address"`
	UserId       string `gorm:"index;not null;column:UserId;type:char(36)"`
	CityId       string `gorm:"index;not null;column:CityId;type:char(36)"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (userLocation *UserLocation) BeforeCreate(tx *gorm.DB) (err error) {
	if len(userLocation.Id) == 0 || userLocation.Id == DEFAULT_UUID {
		userLocation.Id = sequentialguid.New().String()
	}
	return
}
