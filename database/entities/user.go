package entities

import (
	"database/sql"
	"time"
)

type User struct {
	Id           string `gorm:"primaryKey;size:36"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	PhoneNumber  sql.NullString
	Password     string `gorm:"not null"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}
