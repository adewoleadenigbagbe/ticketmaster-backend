package entities

import (
	"time"

	"github.com/Wolechacho/ticketmaster-backend/enums"
)

type UserRole struct {
	Id          string     `gorm:"column:Id"`
	Name        string     `gorm:"column:Name"`
	Role        enums.Role `gorm:"column:Role"`
	Description string     `gorm:"column:Description"`
	CreatedOn   time.Time  `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn  time.Time  `gorm:"column:ModifiedOn;autoUpdateTime"`
}
