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
	Id                 string       `gorm:"primaryKey;size:36;column:Id"`
	Date               time.Time    `gorm:"index;not null;column:Date"`
	StartTime          int64        `gorm:"not null;column:StartTime"`
	EndTime            int64        `gorm:"not null;column:EndTime"`
	CinemaHallId       string       `gorm:"index;not null;size:36;column:CinemalHallId"`
	MovieId            string       `gorm:"index;not null;size:36;column:MovieHallId"`
	IsCancelled        bool         `gorm:"column:IsCancelled"`
	CancellationReason string       `gorm:"column:CancellationReason"`
	IsDeprecated       bool         `gorm:"column:IsDeprecated"`
	CreatedOn          sql.NullTime `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn         sql.NullTime `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Show) TableName() string {
	return "Shows"
}

func (show *Show) BeforeCreate(tx *gorm.DB) (err error) {
	show.Id = sequentialguid.New().String()
	return nil
}
