package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type User struct {
	Id           string         `gorm:"primaryKey;size:36;type:char(36)"`
	FirstName    string         `gorm:"not null;type:varchar(255)"`
	LastName     string         `gorm:"not null;type:varchar(255)"`
	Email        string         `gorm:"not null;type:varchar(255)"`
	PhoneNumber  sql.NullString `gorm:"type:varchar(20)"`
	Password     string         `gorm:"not null"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(user.Id) == 0 || user.Id == DEFAULT_UUID {
		user.Id = sequentialguid.New().String()
	}
	return
}

func (user User) GetId() string {
	return user.Id
}
