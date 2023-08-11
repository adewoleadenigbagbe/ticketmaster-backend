package entities

import (
	"database/sql"
	"time"
)

type User struct {
	Id           string         `gorm:"primaryKey;size:36;type:char(36)"`
	Name         string         `gorm:"not null;type:mediumtext"`
	Email        string         `gorm:"not null;type:varchar(255)"`
	PhoneNumber  sql.NullString `gorm:"type:varchar(20)"`
	Password     string         `gorm:"not null"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}
