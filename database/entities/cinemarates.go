package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type CinemaRate struct {
	Id         string          `gorm:"column:Id"`
	CinemaId   string          `gorm:"column:CinemaId"`
	BaseFee    float32         `gorm:"column:BaseFee"`
	IsActive   bool            `gorm:"column:IsActive"`
	Discount   sql.NullFloat64 `gorm:"column:Discount"`
	IsSpecials sql.NullBool    `gorm:"column:IsSpecials"`
	CreatedOn  time.Time       `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn time.Time       `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (CinemaRate) TableName() string {
	return "CinemaRates"
}

func (cinemaRate *CinemaRate) BeforeCreate(tx *gorm.DB) (err error) {
	if len(cinemaRate.Id) == 0 || cinemaRate.Id == utilities.DEFAULT_UUID {
		cinemaRate.Id = sequentialguid.New().String()
	}

	return
}

func (cinemaRate CinemaRate) GetId() string {
	return cinemaRate.Id
}
