package entities

import (
	"time"
)

type Show struct {
	Id           string    `gorm:"primaryKey;size:36"`
	Date         time.Time `gorm:"index;not null"`
	StartTime    time.Time `gorm:"not null"`
	EndTime      time.Time `gorm:"not null"`
	CinemaHallId string    `gorm:"index;not null"`
	MovieId      string    `gorm:"index;not null"`
	IsDeprecated bool
}
