package entities

type ShowSeat struct {
	Id           string `gorm:"primaryKey;size:36"`
	Status       int `gorm:"not null"`
	Price        float64 `gorm:"not null"`
	CinemaSeatId string `gorm:"index;not null"`
	ShowId       string `gorm:"index;not null"`
	BookingId    string `gorm:"index;not null"`
	IsDeprecated bool
}
