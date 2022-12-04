package entities

type CinemaHall struct {
	Id           string `gorm:"primaryKey;size:36"`
	Name         string `gorm:"not null"`
	TotalHall    int    `gorm:"not null"`
	CinemaId     string `gorm:"index;not null"`
	IsDeprecated bool
}
