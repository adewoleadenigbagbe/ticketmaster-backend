package entities

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type Address struct {
	Id           string                  `gorm:"column:Id"`
	AddressLine  string                  `gorm:"column:AddressLine"`
	EntityId     string                  `gorm:"column:EntityId"`
	CityId       string                  `gorm:"column:CityId"`
	AddressType  enums.EntityAddressType `gorm:"column:AddressType"`
	Coordinates  Coordinate              `gorm:"column:Coordinates"`
	IsCurrent    bool                    `gorm:"column:IsCurrent"`
	IsDeprecated bool                    `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time               `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time               `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Address) TableName() string {
	return "Addresses"
}

func (address *Address) BeforeCreate(tx *gorm.DB) (err error) {
	if len(address.Id) == 0 || address.Id == utilities.DEFAULT_UUID {
		address.Id = sequentialguid.New().String()
	}

	address.IsDeprecated = false
	return
}

func (address Address) GetId() string {
	return address.Id
}
