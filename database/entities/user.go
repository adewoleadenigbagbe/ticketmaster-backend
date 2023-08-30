package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type User struct {
	Id           string         `gorm:"column:Id"`
	FirstName    string         `gorm:"column:FirstName"`
	LastName     string         `gorm:"column:LastName"`
	Email        string         `gorm:"column:Email"`
	PhoneNumber  sql.NullString `gorm:"column:PhoneNumber"`
	Password     string         `gorm:"column:Password"`
	IsDeprecated bool           `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time      `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time      `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(user.Id) == 0 || user.Id == DEFAULT_UUID {
		user.Id = sequentialguid.New().String()
	}
	user.IsDeprecated = false
	return
}
