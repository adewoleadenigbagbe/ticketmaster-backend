package entities

import "time"

type CinemaHall struct {
	Id           string    `gorm:"primaryKey;size:36;column:Id;type:char(36)"`
	Name         string    `gorm:"index;not null;column:Name;type:varchar(255)"`
	TotalSeat    int       `gorm:"not null;column:TotalSeat"`
	CinemaId     string    `gorm:"index;not null;column:CinemalId;type:char(36)"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaHall) TableName() string {
	return "CinemaHalls"
}
