package entities

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/shared/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/shared/helpers"
	"github.com/Wolechacho/ticketmaster-backend/shared/helpers/utilities"
	"gorm.io/gorm"
)

type UserRole struct {
	Id          string     `gorm:"column:Id"`
	Name        string     `gorm:"column:Name"`
	Role        enums.Role `gorm:"column:Role"`
	Description string     `gorm:"column:Description"`
	CreatedOn   time.Time  `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn  time.Time  `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (userRole *UserRole) BeforeCreate(tx *gorm.DB) (err error) {
	if len(userRole.Id) == 0 || userRole.Id == utilities.DEFAULT_UUID {
		userRole.Id = sequentialguid.New().String()
	}

	return
}

func (UserRole) TableName() string {
	return "UserRoles"
}

func (userRole UserRole) GetId() string {
	return userRole.Id
}
