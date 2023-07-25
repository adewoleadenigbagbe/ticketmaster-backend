package entities

import (
	"database/sql"
	"time"
)

type City struct {
	Id           string         `gorm:"primaryKey;size:36;column:Id"`
	Name         string         `gorm:"not null;index;column:Name"`
	State        string         `gorm:"not null;column:State"`
	Zipcode      sql.NullString `gorm:"column:ZipCode"`
	IsDeprecated bool           `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time      `gorm:"column:CreatedOn"`
	ModifiedOn   time.Time      `gorm:"column:ModifiedOn"`
}

func (City) TableName() string {
	return "Cities"
}
