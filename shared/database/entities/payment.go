package entities

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/shared/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/shared/helpers"
	"github.com/Wolechacho/ticketmaster-backend/shared/helpers/utilities"
	"gorm.io/gorm"
)

type Payment struct {
	Id                       string                     `gorm:"column:Id"`
	Amount                   float64                    `gorm:"column:Amount"`
	PaymentDate              time.Time                  `gorm:"column:PaymentDate"`
	DiscountCouponId         utilities.Nullable[string] `gorm:"column:DiscountCouponId"`
	RemoteTransactionId      utilities.Nullable[string] `gorm:"column:RemoteTransactionId"`
	PaymentMethod            enums.PaymentType          `gorm:"column:PaymentMethod"`
	BookingId                string                     `gorm:"column:BookingId"`
	ProviderExtraInformation utilities.Nullable[string] `gorm:"column:ProviderExtraInformation"`
	IsDeprecated             bool                       `gorm:"column:IsDeprecated"`
	CreatedOn                time.Time                  `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn               time.Time                  `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (Payment) TableName() string {
	return "Payments"
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if len(payment.Id) == 0 || payment.Id == utilities.DEFAULT_UUID {
		payment.Id = sequentialguid.New().String()
	}
	payment.IsDeprecated = false

	return
}

func (payment Payment) GetId() string {
	return payment.Id
}
