package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type Payment struct {
	Id                  string         `gorm:"column:Id"`
	Amount              float64        `gorm:"column:Amount"`
	PaymentDate         time.Time      `gorm:"column:PaymentDate"`
	DiscountCouponId    sql.NullString `gorm:"column:DiscountCouponId"`
	RemoteTransactionId sql.NullString `gorm:"column:RemoteTransactionId"`
	//PaymentMethod is enum
	PaymentMethod int       `gorm:"column:PaymentMethod"`
	BookingId     string    `gorm:"column:BookingId"`
	IsDeprecated  bool      `gorm:"column:IsDeprecated"`
	CreatedOn     time.Time `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn    time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if len(payment.Id) == 0 || payment.Id == DEFAULT_UUID {
		payment.Id = sequentialguid.New().String()
	}
	payment.IsDeprecated = false

	return
}
