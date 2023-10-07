package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type Movie struct {
	Id           string         `gorm:"column:Id"`
	Title        string         `gorm:"column:Title"`
	Description  sql.NullString `gorm:"column:Description"`
	Language     string         `gorm:"column:Language"`
	ReleaseDate  time.Time      `gorm:"column:ReleaseDate"`
	Duration     sql.NullInt32  `gorm:"column:Duration"`
	Genre        int            `gorm:"column:Genre"`
	Popularity   float32        `gorm:"column:Popularity"`
	VoteCount    int            `gorm:"column:VoteCount"`
	IsDeprecated bool           `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time      `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time      `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Movie) TableName() string {
	return "Movies"
}

func (movie *Movie) BeforeCreate(tx *gorm.DB) (err error) {
	if len(movie.Id) == 0 || movie.Id == utilities.DEFAULT_UUID {
		movie.Id = sequentialguid.New().String()
	}

	movie.IsDeprecated = false
	return
}

func (movie Movie) GetId() string {
	return movie.Id
}
