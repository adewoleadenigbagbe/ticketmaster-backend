package entities

import "time"

type CinemaHall struct {
	Id           string    `gorm:"primaryKey;size:36;column:Id"`
	Name         string    `gorm:"index;not null;column:Name"`
	TotalSeat    int       `gorm:"not null;column:TotalSeat"`
	CinemaId     string    `gorm:"index;not null;column:CinemalId"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"index;column:CreatedOn"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn"`
}
