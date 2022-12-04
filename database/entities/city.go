package entities

import "database/sql"

type City struct {
	Id           string `gorm:"primaryKey;size:36"`
	Name         string `gorm:"not null"`
	State        string `gorm:"not null"`
	Zipcode      sql.NullString
	IsDeprecated bool
}
