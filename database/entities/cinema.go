package entities

import "time"

type Cinema struct {
	Id                string    `gorm:"primaryKey;size:36;column:Id;type:char(36)"`
	Name              string    `gorm:"index;not null;column:Name;type:varchar(255)"`
	TotalCinemalHalls int       `gorm:"not null;column:TotalCinemalHalls"`
	CityId            string    `gorm:"index;not null;column:CityId;type:char(36)"`
	IsDeprecated      bool      `gorm:"column:IsDeprecated"`
	CreatedOn         time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn        time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Cinema) TableName() string {
	return "Cinemas"
}
