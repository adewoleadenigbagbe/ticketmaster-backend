package entities

import (
	"database/sql"
	"time"
)

type Movie struct {
	Id           string `gorm:"primaryKey;size:36"`
	Title        string `gorm:"not null"`
	Description  sql.NullString
	Language     string    `gorm:"not null"`
	Duration     time.Time `gorm:"not null"`
	ReleaseDate  time.Time `gorm:"index;not null"`
	Country      string    `gorm:"not null"`
	Genre        string    `gorm:"not null"`
	IsDeprecated bool
}
