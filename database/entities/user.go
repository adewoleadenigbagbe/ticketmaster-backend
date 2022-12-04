package entities

import (
	"database/sql"
)

type User struct {
	Id           string `gorm:"primaryKey;size:36"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	PhoneNumber  sql.NullString
	Password     string `gorm:"not null"`
	IsDeprecated bool
}
