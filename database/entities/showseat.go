package entities

import (
	"database/sql"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type ShowSeat struct {
	Id           string               `gorm:"column:Id"`
	Status       enums.ShowSeatStatus `gorm:"column:Status"`
	Price        float64              `gorm:"column:Price"`
	CinemaSeatId string               `gorm:"column:CinemaSeatId"`
	ShowId       string               `gorm:"column:ShowId"`
	BookingId    sql.NullString       `gorm:"column:BookingId"`
	IsDeprecated bool                 `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time            `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time            `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (ShowSeat) TableName() string {
	return "ShowSeats"
}

func (showSeat *ShowSeat) BeforeCreate(tx *gorm.DB) (err error) {
	if len(showSeat.Id) == 0 || showSeat.Id == utilities.DEFAULT_UUID {
		showSeat.Id = sequentialguid.New().String()
	}
	showSeat.IsDeprecated = false
	return
}

func (showSeat ShowSeat) GetId() string {
	return showSeat.Id
}
