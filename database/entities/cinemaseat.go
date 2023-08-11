package entities

import "time"

type CinemaSeat struct {
	Id         string `gorm:"primaryKey;size:36;column:Id;type:char(36)"`
	SeatNumber int    `gorm:"not null;column:SeatNumber"`
	//Type is an enum
	Type         int       `gorm:"not null;column:Type"`
	CinemaHallId string    `gorm:"index;not null;column:CinemaHallId;type:char(36)"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaSeat) TableName() string {
	return "CinemaSeats"
}
