package entities

import "time"

type ShowSeat struct {
	Id           string  `gorm:"primaryKey;size:36;type:char(36)"`
	Status       int     `gorm:"not null"`
	Price        float64 `gorm:"not null"`
	CinemaSeatId string  `gorm:"index;not null;type:char(36)"`
	ShowId       string  `gorm:"index;not null;type:char(36)"`
	BookingId    string  `gorm:"index;not null;type:char(36)"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}
