package entities

import (
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type Show struct {
	Id                 string                     `gorm:"column:Id"`
	Date               time.Time                  `gorm:"column:Date"`
	StartTime          int64                      `gorm:"column:StartTime"`
	EndTime            int64                      `gorm:"column:EndTime"`
	CinemaHallId       string                     `gorm:"column:CinemaHallId"`
	MovieId            string                     `gorm:"column:MovieId"`
	IsCancelled        bool                       `gorm:"column:IsCancelled"`
	CancellationReason utilities.Nullable[string] `gorm:"column:CancellationReason"`
	IsDeprecated       bool                       `gorm:"column:IsDeprecated"`
	CreatedOn          time.Time                  `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn         time.Time                  `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Show) TableName() string {
	return "Shows"
}

func (show *Show) BeforeCreate(tx *gorm.DB) (err error) {
	if len(show.Id) == 0 || show.Id == utilities.DEFAULT_UUID {
		show.Id = sequentialguid.New().String()
	}
	show.IsDeprecated = false
	return
}

func (show Show) GetId() string {
	return show.Id
}
