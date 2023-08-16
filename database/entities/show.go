package entities

import (
	"database/sql"
	"time"
)

type Tabler interface {
	TableName() string
}

type Show struct {
	Id                 string       `gorm:"primaryKey;size:36;column:Id;type:char(36)"`
	Date               time.Time    `gorm:"index;not null;column:Date"`
	StartTime          int64        `gorm:"not null;column:StartTime"`
	EndTime            int64        `gorm:"not null;column:EndTime"`
	CinemaHallId       string       `gorm:"index;not null;size:36;column:CinemalHallId;type:char(36)"`
	MovieId            string       `gorm:"index;not null;size:36;column:MovieId;type:char(36)"`
	IsCancelled        bool         `gorm:"column:IsCancelled"`
	CancellationReason string       `gorm:"column:CancellationReason;type:mediumtext"`
	IsDeprecated       bool         `gorm:"column:IsDeprecated"`
	CreatedOn          sql.NullTime `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn         sql.NullTime `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Show) TableName() string {
	return "Shows"
}
