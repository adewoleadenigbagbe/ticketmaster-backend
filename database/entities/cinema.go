package entities

type Cinema struct {
	Id                string `gorm:"primaryKey;size:36"`
	Name              string `gorm:"not null"`
	TotalCinemalHalls int    `gorm:"not null"`
	CityId            string `gorm:"index;not null1"`
	IsDeprecated      bool
}
