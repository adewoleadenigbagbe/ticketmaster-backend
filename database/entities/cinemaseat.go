package entities

type CinemaSeat struct {
	Id         string `gorm:"primaryKey;size:36"`
	SeatNumber int    `gorm:"not null"`
	//Type is an enum
	Type         int    `gorm:"not null"`
	CinemaHallId string `gorm:"index;not null"`
	IsDeprecated bool
}
