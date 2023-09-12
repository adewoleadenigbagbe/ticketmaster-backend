package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type Show struct {
	Id                 string         `gorm:"column:Id"`
	Date               time.Time      `gorm:"column:Date"`
	StartTime          int64          `gorm:"column:StartTime"`
	EndTime            int64          `gorm:"column:EndTime"`
	CinemaHallId       string         `gorm:"column:CinemaHallId"`
	MovieId            string         `gorm:"column:MovieId"`
	IsCancelled        bool           `gorm:"column:IsCancelled"`
	CancellationReason sql.NullString `gorm:"column:CancellationReason"`
	IsDeprecated       bool           `gorm:"column:IsDeprecated"`
	CreatedOn          sql.NullTime   `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn         sql.NullTime   `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Show) TableName() string {
	return "Shows"
}

func (show *Show) BeforeCreate(tx *gorm.DB) (err error) {
	if len(show.Id) == 0 || show.Id == DEFAULT_UUID {
		show.Id = sequentialguid.New().String()
	}
	show.IsDeprecated = false
	return
}
