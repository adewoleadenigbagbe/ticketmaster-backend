package entities

import (
	"database/sql"
	"time"
)

type Payment struct {
	Id                  string    `gorm:"primaryKey;size:36"`
	Amount              float64   `gorm:"not null"`
	PaymentDate         time.Time `gorm:"index;not null"`
	DiscountCouponId    sql.NullString
	RemoteTransactionId sql.NullString
	//PaymentMethod is enum
	PaymentMethod int    `gorm:"not null"`
	BookingId     string `gorm:"index;not null"`
	IsDeprecated  bool
	CreatedOn     time.Time `gorm:"index;column:CreatedOn;autoCreateTime"`
	ModifiedOn    time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}
