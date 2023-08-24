package entities

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Coordinate struct {
	Longitude int
	Latitude  int
}

func (coord Coordinate) GormDataType() string {
	return "geometry"
}

func (coord Coordinate) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", coord.Longitude, coord.Latitude)},
	}
}

// Scan implements the sql.Scanner interface
func (coord *Coordinate) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	return nil
}
