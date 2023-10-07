package entities

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Coordinate struct {
	Longitude float32
	Latitude  float32
}

func (coord Coordinate) GormDataType() string {
	return "geometry"
}

func (coord Coordinate) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%.6f %.6f)", coord.Longitude, coord.Latitude)},
	}
}

// Scan implements the sql.Scanner interface
func (coord *Coordinate) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	switch v := v.(type) {
	case []uint8:
		var longitude float64
		var latitude float64
		buf := bytes.NewReader(v[9:17])
		err := binary.Read(buf, binary.LittleEndian, &longitude)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(v[17:25])
		err = binary.Read(buf, binary.LittleEndian, &latitude)
		if err != nil {
			return err
		}

		coord.Latitude = float32(latitude)
		coord.Longitude = float32(longitude)
	default:
		return errors.New("incompatible type for Coordinates")
	}
	return nil
}
