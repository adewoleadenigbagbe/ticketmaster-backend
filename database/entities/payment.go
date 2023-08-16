package entities

import (
	"database/sql"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
)

type Payment struct {
	Id                  string    `gorm:"primaryKey;size:36;type:char(36)"`
	Amount              float64   `gorm:"not null"`
	PaymentDate         time.Time `gorm:"index;not null"`
	DiscountCouponId    sql.NullString
	RemoteTransactionId sql.NullString
	//PaymentMethod is enum
	PaymentMethod int    `gorm:"not null"`
	BookingId     string `gorm:"index;not null;type:char(36)"`
	IsDeprecated  bool
	CreatedOn     time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn    time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if len(payment.Id) == 0 || payment.Id == DEFAULT_UUID {
		payment.Id = sequentialguid.New().String()
	}
	return
}

func (payment Payment) GetId() string {
	return payment.Id
}

