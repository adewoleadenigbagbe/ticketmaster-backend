package entities

import (
	"database/sql"
	"time"
)

type Movie struct {
	Id           string `gorm:"primaryKey;size:36;type:char(36)"`
	Title        string `gorm:"not null;type:mediumtext"`
	Description  sql.NullString
	Language     string    `gorm:"not null;type:char(10)"`
	ReleaseDate  time.Time `gorm:"not null"`
	Duration     sql.NullInt32
	Genre        int     `gorm:"not null"`
	Popularity   float32 `gorm:"not null"`
	VoteCount    int     `gorm:"not null"`
	IsDeprecated bool
	CreatedOn    time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Movie) TableName() string {
	return "Movies"
}
