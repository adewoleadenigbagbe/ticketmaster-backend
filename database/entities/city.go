package entities

import (
	"context"
	"fmt"
	"time"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Location struct {
	Longitude int
	Latitude  int
}

func (loc Location) GormDataType() string {
	return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.Longitude, loc.Latitude)},
	}
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	return nil
}

type City struct {
	Id           string    `gorm:"primaryKey;size:36;column:Id;type:char(36)"`
	Name         string    `gorm:"not null;index;column:Name;type:varchar(255)"`
	State        string    `gorm:"not null;column:State;type:varchar(255)"`
	Location     Location       `gorm:"not null;column:Location"`
	Zipcode      string    `gorm:"not null;column:ZipCode;type:varchar(255)"`
	IsDeprecated bool      `gorm:"column:IsDeprecated"`
	CreatedOn    time.Time `gorm:"column:CreatedOn;autoCreateTime"`
	ModifiedOn   time.Time `gorm:"column:ModifiedOn;autoUpdateTime"`
}

func (City) TableName() string {
	return "Cities"
}

func (city *City) BeforeCreate(tx *gorm.DB) (err error) {
	if len(city.Id) == 0 || city.Id == DEFAULT_UUID {
		city.Id = sequentialguid.New().String()
	}
	return
}
