package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type User struct {
	Id           string         `gorm:"column:Id"`
	FirstName    string         `gorm:"column:FirstName"`
	LastName     string         `gorm:"column:LastName"`
	RoleId       string         `gorm:"column:RoleId"`
	Email        string         `gorm:"column:Email"`
	PhoneNumber  sql.NullString `gorm:"column:PhoneNumber"`
	Password     string         `gorm:"column:Password"`
	IsDeprecated bool           `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time      `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time      `gorm:"column:ModifiedOn;autoUpdateTime"`
	UserRole     UserRole       `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(user.Id) == 0 || user.Id == utilities.DEFAULT_UUID {
		user.Id = sequentialguid.New().String()
	}
	user.IsDeprecated = false
	return
}

func (user User) GetId() string {
	return user.Id
}
