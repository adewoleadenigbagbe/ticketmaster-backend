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
	ReleaseDate  time.Time `gorm:"not null"`
	Duration     sql.NullInt32
	Genre        int     `gorm:"not null"`
	Popularity   float32 `gorm:"not null"`
	VoteCount    int     `gorm:"not null"`
	IsDeprecated bool
}

func (Movie) TableName() string {
	return "Movies"
}
