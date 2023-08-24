package entities

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type Address struct {
	Id           string                  `gorm:"primaryKey;size:36;type:char(36);column:Id"`
	AddressLine  string                  `gorm:"not null;type:mediumtext;colunm:Address"`
	EntityId     string                  `gorm:"index;not null;column:UserId;type:char(36)"`
	AddressType  enums.EntityAddressType `gorm:"index;not null;column:AddressType;type:int"`
	CityId       string                  `gorm:"index;not null;column:CityId;type:char(36)"`
	Coordinates  Coordinate              `gorm:"not null;column:Location"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (address *Address) BeforeCreate(tx *gorm.DB) (err error) {
	if len(address.Id) == 0 || address.Id == DEFAULT_UUID {
		address.Id = sequentialguid.New().String()
	}
	return
}
